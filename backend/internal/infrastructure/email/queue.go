package email

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/trickreport/backend/internal/email"
	"sync"
	"sync/atomic"
	"time"
)

type QueuedEmail struct {
	ID                uuid.UUID
	To, Subject, Body string
	Attempts          int
}
type OutboxRepository interface {
	Enqueue(context.Context, string, string, string) error
	Claim(context.Context) (*QueuedEmail, error)
	Finish(context.Context, *QueuedEmail, bool) error
	Counts(context.Context) (pending, failed int64, err error)
}
type EmailQueue struct {
	store       OutboxRepository
	lastSuccess atomic.Int64
	done        chan struct{}
}

func NewEmailQueue(store OutboxRepository) *EmailQueue {
	return &EmailQueue{store: store, done: make(chan struct{})}
}
func (q *EmailQueue) Send(to, subject, body string) error {
	return q.SendContext(context.Background(), to, subject, body)
}
func (q *EmailQueue) SendContext(ctx context.Context, to, subject, body string) error {
	if q.store == nil {
		return errors.New("email outbox is not configured")
	}
	return q.store.Enqueue(ctx, to, subject, body)
}
func (q *EmailQueue) LastSuccessUnix() int64 { return q.lastSuccess.Load() }
func (q *EmailQueue) Done() <-chan struct{}  { return q.done }
func (q *EmailQueue) Process(ctx context.Context, sender email.Sender) {
	defer close(q.done)
	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				if ctx.Err() != nil {
					return
				}
				if err := q.deliverNext(ctx, sender); err != nil {
					log.Error().Msg("email outbox cycle failed")
				}
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
				}
			}
		}()
	}
	wg.Wait()
}
func (q *EmailQueue) deliverNext(ctx context.Context, sender email.Sender) error {
	msg, err := q.store.Claim(ctx)
	if err != nil {
		return err
	}
	q.lastSuccess.Store(time.Now().Unix())
	if msg == nil {
		return nil
	}
	if contextSender, ok := sender.(email.ContextSender); ok {
		err = contextSender.SendContext(ctx, msg.To, msg.Subject, msg.Body)
	} else {
		err = sender.Send(msg.To, msg.Subject, msg.Body)
	}
	if err != nil {
		log.Warn().Str("message_id", msg.ID.String()).Int("attempt", msg.Attempts).Msg("email delivery failed; retained for retry")
	}
	return q.store.Finish(ctx, msg, err == nil)
}
