package email

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type queueTestStore struct {
	message   *QueuedEmail
	delivered bool
}

func (s *queueTestStore) Enqueue(context.Context, string, string, string) error { return nil }
func (s *queueTestStore) Claim(context.Context) (*QueuedEmail, error)           { return s.message, nil }
func (s *queueTestStore) Finish(_ context.Context, _ *QueuedEmail, delivered bool) error {
	s.delivered = delivered
	return nil
}
func (s *queueTestStore) Counts(context.Context) (int64, int64, error) { return 0, 0, nil }

type cancelingContextSender struct {
	cancel context.CancelFunc
	called bool
}

func (s *cancelingContextSender) Send(string, string, string) error {
	return errors.New("legacy Send unexpectedly called")
}
func (s *cancelingContextSender) SendContext(ctx context.Context, _, _, _ string) error {
	s.called = true
	s.cancel()
	return ctx.Err()
}

func TestEmailQueuePropagatesCancellationToSender(t *testing.T) {
	store := &queueTestStore{message: &QueuedEmail{ID: uuid.New(), To: "to@example.test"}}
	q := NewEmailQueue(store)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sender := &cancelingContextSender{cancel: cancel}

	if err := q.deliverNext(ctx, sender); err != nil {
		t.Fatalf("deliverNext returned error: %v", err)
	}
	if !sender.called || store.delivered {
		t.Fatalf("sender called=%v, delivered=%v; want called and retained for retry", sender.called, store.delivered)
	}
}
