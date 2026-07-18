package realtime

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewHub(t *testing.T) {
	h := NewHub()
	if h == nil {
		t.Fatal("expected non-nil hub")
	}
	if h.clients == nil {
		t.Error("expected initialized clients map")
	}
	if h.Broadcast == nil || h.Register == nil || h.Unregister == nil {
		t.Error("expected initialized channels")
	}
}

func TestHub_RegisterAndUnregister(t *testing.T) {
	h := NewHub()
	go h.Run()
	defer func() {
		// Stop the hub by sending to unregister (Run loop is infinite; we just let it be GC'd)
	}()

	tenantID := uuid.New()
	client := &Client{
		Hub:      h,
		TenantID: tenantID,
		UserID:   uuid.New(),
		send:     make(chan []byte, 256),
	}

	h.Register <- client
	time.Sleep(10 * time.Millisecond)

	h.mu.RLock()
	clients, ok := h.clients[tenantID]
	h.mu.RUnlock()
	if !ok || !clients[client] {
		t.Error("expected client to be registered")
	}

	h.Unregister <- client
	time.Sleep(10 * time.Millisecond)

	h.mu.RLock()
	_, stillExists := h.clients[tenantID]
	h.mu.RUnlock()
	if stillExists {
		t.Error("expected tenant map to be cleaned up after unregister")
	}
}

func TestHub_Broadcast(t *testing.T) {
	h := NewHub()
	go h.Run()

	tenantID := uuid.New()
	client := &Client{
		Hub:      h,
		TenantID: tenantID,
		UserID:   uuid.New(),
		send:     make(chan []byte, 256),
	}

	h.Register <- client
	time.Sleep(10 * time.Millisecond)

	h.Broadcast <- BroadcastPayload{
		TenantID: tenantID,
		Message:  []byte("hello"),
	}

	select {
	case msg := <-client.send:
		if string(msg) != "hello" {
			t.Errorf("got %q, want hello", string(msg))
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for broadcast message")
	}
}

func TestHub_BroadcastEvent(t *testing.T) {
	h := NewHub()
	go h.Run()

	tenantID := uuid.New()
	client := &Client{
		Hub:      h,
		TenantID: tenantID,
		UserID:   uuid.New(),
		send:     make(chan []byte, 256),
	}

	h.Register <- client
	time.Sleep(10 * time.Millisecond)

	h.BroadcastEvent(tenantID, EventTicketCreated, map[string]string{"id": "123"})

	select {
	case msg := <-client.send:
		if len(msg) == 0 {
			t.Error("expected non-empty message")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for broadcast event")
	}
}

func TestHub_Broadcast_DifferentTenant(t *testing.T) {
	h := NewHub()
	go h.Run()

	tenant1 := uuid.New()
	tenant2 := uuid.New()
	client := &Client{
		Hub:      h,
		TenantID: tenant1,
		UserID:   uuid.New(),
		send:     make(chan []byte, 256),
	}

	h.Register <- client
	time.Sleep(10 * time.Millisecond)

	// Broadcast to a different tenant — client should not receive it
	h.Broadcast <- BroadcastPayload{
		TenantID: tenant2,
		Message:  []byte("should-not-receive"),
	}

	select {
	case msg := <-client.send:
		t.Errorf("should not have received message from different tenant, got %q", string(msg))
	case <-time.After(50 * time.Millisecond):
		// Expected — no message received
	}
}

func TestHub_Broadcast_FullSendChannel(t *testing.T) {
	h := NewHub()
	go h.Run()

	tenantID := uuid.New()
	client := &Client{
		Hub:      h,
		TenantID: tenantID,
		UserID:   uuid.New(),
		send:     make(chan []byte, 1), // small buffer
	}

	h.Register <- client
	time.Sleep(10 * time.Millisecond)

	// Fill the send channel
	client.send <- []byte("msg1")

	// This broadcast should hit the default case and remove the client
	h.Broadcast <- BroadcastPayload{
		TenantID: tenantID,
		Message:  []byte("msg2"),
	}
	time.Sleep(10 * time.Millisecond)

	h.mu.RLock()
	clients := h.clients[tenantID]
	h.mu.RUnlock()
	if len(clients) != 0 {
		t.Errorf("expected client to be removed after full send channel, got %d clients", len(clients))
	}
}

func TestHub_BroadcastEvent_MarshalError(t *testing.T) {
	h := NewHub()
	// Don't start Run — we're testing the error path, not the broadcast path.

	tenantID := uuid.New()
	// A channel type cannot be marshaled to JSON, causing json.Marshal to fail.
	h.BroadcastEvent(tenantID, EventTicketCreated, make(chan int))

	// If we reach here without panicking, the error path was exercised.
	// The function should have logged the error and returned without sending.
}

func TestHub_Run_UnregisterUnknownClient(t *testing.T) {
	h := NewHub()
	go h.Run()

	tenantID := uuid.New()
	client := &Client{
		Hub:      h,
		TenantID: tenantID,
		UserID:   uuid.New(),
		send:     make(chan []byte, 256),
	}

	// Unregister a client that was never registered — should not panic
	h.Unregister <- client
	time.Sleep(10 * time.Millisecond)

	// Verify no tenant map was created
	h.mu.RLock()
	_, exists := h.clients[tenantID]
	h.mu.RUnlock()
	if exists {
		t.Error("expected no tenant map for unregistered client")
	}
}
