package realtime

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestScopedDeliveryDoesNotReachAnotherUserInSameTenant(t *testing.T) {
	h := NewHub()
	go h.Run()
	tenant := uuid.New()
	owner := uuid.New()
	a := &Client{TenantID: tenant, UserID: owner, send: make(chan []byte, 1)}
	b := &Client{TenantID: tenant, UserID: uuid.New(), send: make(chan []byte, 1)}
	h.Register <- a
	h.Register <- b
	sender, ok := any(h).(interface {
		SendToUsers(uuid.UUID, []uuid.UUID, EventType, any)
	})
	if !ok {
		t.Fatal("hub supports only tenant-wide delivery")
	}
	sender.SendToUsers(tenant, []uuid.UUID{owner}, EventNewComment, map[string]string{"content": "private"})
	select {
	case <-a.send:
	case <-time.After(time.Second):
		t.Fatal("owner missed event")
	}
	select {
	case <-b.send:
		t.Fatal("another user's private event leaked")
	default:
	}
	h.Unregister <- a
	h.Unregister <- b
}
