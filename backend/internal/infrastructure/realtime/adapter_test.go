package realtime

import (
	"testing"
	"time"

	"github.com/google/uuid"
	rt "github.com/trickreport/backend/internal/realtime"
)

func TestHubAdapter_BroadcastEvent(t *testing.T) {
	hub := rt.NewHub()
	// Don't start Run — we'll read from the Broadcast channel directly.

	adapter := NewHubAdapter(hub)
	tenantID := uuid.New()

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

	adapter.BroadcastEvent(tenantID, "TICKET_CREATED", map[string]string{"id": "1"})

	select {
	case r := <-resCh:
		if r.payload.TenantID != tenantID {
			t.Errorf("tenantID = %v, want %v", r.payload.TenantID, tenantID)
		}
		if len(r.payload.Message) == 0 {
			t.Error("expected non-empty message")
		}
	case <-time.After(100 * time.Millisecond):
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
