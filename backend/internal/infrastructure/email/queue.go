package email

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/trickreport/backend/internal/email"
)

// queuedEmail holds the data needed to send a single email.
type queuedEmail struct {
	to      string
	subject string
	body    string
}

// EmailQueue is a channel-based email queue with retry and dead-letter
// logging. It decouples email sending from the request path.
//
// TODO(wire): wire EmailQueue into wire.go's NewServices and start its
// Process loop in the server bootstrap.
type EmailQueue struct {
	queue chan queuedEmail
	wg    sync.WaitGroup
}

// NewEmailQueue creates a new EmailQueue with the given buffer size.
func NewEmailQueue(bufferSize int) *EmailQueue {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	return &EmailQueue{
		queue: make(chan queuedEmail, bufferSize),
	}
}

// Enqueue adds an email to the queue. It is non-blocking — if the queue is
// full the email is logged as dropped (dead-letter) and not sent.
func (q *EmailQueue) Enqueue(to, subject, body string) {
	select {
	case q.queue <- queuedEmail{to: to, subject: subject, body: body}:
	default:
		log.Warn().
			Str("to", to).
			Str("subject", subject).
			Msg("Email queue full — message dead-lettered")
	}
}

// Process drains the queue and sends emails via the provided sender. It blocks
// until ctx is cancelled. Failed sends are retried with exponential backoff
// (1s, 2s, 4s) up to maxRetries times. Emails that still fail after all
// retries are logged as dead letters.
func (q *EmailQueue) Process(ctx context.Context, sender email.Sender) {
	backoffs := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	maxRetries := 3

	for {
		select {
		case <-ctx.Done():
			q.wg.Wait()
			return
		case msg := <-q.queue:
			q.wg.Add(1)
			func(msg queuedEmail) {
				defer q.wg.Done()
				q.sendWithRetry(ctx, sender, msg, backoffs, maxRetries)
			}(msg)
		}
	}
}

func (q *EmailQueue) sendWithRetry(ctx context.Context, sender email.Sender, msg queuedEmail, backoffs []time.Duration, maxRetries int) {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if err := sender.Send(msg.to, msg.subject, msg.body); err == nil {
			return
		} else {
			lastErr = err
		}

		if attempt < maxRetries {
			backoff := backoffs[attempt%len(backoffs)]
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}
		}
	}

	// Dead letter — all retries exhausted.
	log.Error().
		Err(lastErr).
		Str("to", msg.to).
		Str("subject", msg.subject).
		Int("attempts", maxRetries+1).
		Msg("Email dead-lettered after max retries")
}
