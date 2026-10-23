package realtime

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	rt "github.com/trickreport/backend/internal/realtime"
)

type testRecipientResolver struct{ users []uuid.UUID }

func (r testRecipientResolver) EventRecipients(context.Context, uuid.UUID, uuid.UUID, bool) ([]uuid.UUID, error) {
	return r.users, nil
}

func TestHubAdapter_BroadcastEvent(t *testing.T) {
	hub := rt.NewHub()
	// Don't start Run — we'll read the targeted payload directly.
	tenantID := uuid.New()
	ticketID := uuid.New()
	recipientID := uuid.New()
	adapter := NewHubAdapter(hub, testRecipientResolver{users: []uuid.UUID{recipientID}})

	// Read from the Broadcast channel in a goroutine since BroadcastEvent
	// sends to an unbuffered channel and would block without a reader.
	type result struct {
		payload rt.BroadcastPayload
	}
	resCh := make(chan result, 1)
	go func() {
		payload := <-hub.Broadcast
		resCh <- result{payload: payload}
	}()

	adapter.BroadcastEvent(tenantID, "TICKET_CREATED", map[string]any{"id": ticketID, "title": "Printer unavailable"})

	select {
	case r := <-resCh:
		if r.payload.TenantID != tenantID {
			t.Errorf("tenantID = %v, want %v", r.payload.TenantID, tenantID)
		}
		if len(r.payload.Message) == 0 {
			t.Error("expected non-empty message")
		}
		if !r.payload.Scoped || len(r.payload.UserIDs) != 1 || r.payload.UserIDs[0] != recipientID {
			t.Errorf("event was not restricted to resolved recipients: %+v", r.payload)
		}
	case <-time.After(time.Second):
		t.Fatal("expected to receive broadcast via adapter")
	}
}

func TestNewHubAdapter(t *testing.T) {
	hub := rt.NewHub()
	adapter := NewHubAdapter(hub)
	if adapter == nil {
		t.Fatal("expected non-nil adapter")
	}
}
