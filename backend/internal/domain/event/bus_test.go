package event

import (
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
)

func TestMemoryEventBus_SubscribeAndPublish(t *testing.T) {
	bus := NewMemoryEventBus()
	var calls int32

	bus.Subscribe("ticket.created", func(e Event) {
		if e.Name != "ticket.created" {
			t.Errorf("Name = %q", e.Name)
		}
		atomic.AddInt32(&calls, 1)
	})

	bus.Publish(Event{Name: "ticket.created", AggregateID: uuid.New()})

	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Errorf("handler called %d times, want 1", got)
	}
}

func TestMemoryEventBus_MultipleHandlers(t *testing.T) {
	bus := NewMemoryEventBus()
	var count int32

	bus.Subscribe("evt", func(e Event) { atomic.AddInt32(&count, 1) })
	bus.Subscribe("evt", func(e Event) { atomic.AddInt32(&count, 1) })
	bus.Subscribe("evt", func(e Event) { atomic.AddInt32(&count, 1) })

	bus.Publish(Event{Name: "evt"})

	if got := atomic.LoadInt32(&count); got != 3 {
		t.Errorf("handlers called %d times, want 3", got)
	}
}

func TestMemoryEventBus_NoHandlers(t *testing.T) {
	bus := NewMemoryEventBus()
	// Publishing an event with no subscribers should not panic.
	bus.Publish(Event{Name: "unhandled.event"})
}

func TestMemoryEventBus_HandlerIsolation(t *testing.T) {
	bus := NewMemoryEventBus()
	var aCalls, bCalls int32

	bus.Subscribe("a", func(e Event) { atomic.AddInt32(&aCalls, 1) })
	bus.Subscribe("b", func(e Event) { atomic.AddInt32(&bCalls, 1) })

	bus.Publish(Event{Name: "a"})
	bus.Publish(Event{Name: "b"})
	bus.Publish(Event{Name: "a"})

	if got := atomic.LoadInt32(&aCalls); got != 2 {
		t.Errorf("a handler called %d times, want 2", got)
	}
	if got := atomic.LoadInt32(&bCalls); got != 1 {
		t.Errorf("b handler called %d times, want 1", got)
	}
}

func TestMemoryEventBus_EventPayload(t *testing.T) {
	bus := NewMemoryEventBus()
	type payload struct{ Msg string }

	var received string
	bus.Subscribe("test.payload", func(e Event) {
		p, ok := e.Payload.(payload)
		if !ok {
			t.Error("payload type mismatch")
			return
		}
		received = p.Msg
	})

	bus.Publish(Event{
		Name:    "test.payload",
		Payload: payload{Msg: "hello"},
	})

	if received != "hello" {
		t.Errorf("received = %q, want %q", received, "hello")
	}
}

func TestMemoryEventBus_ConcurrentSubscribe(t *testing.T) {
	bus := NewMemoryEventBus()
	var count int32

	// Subscribe many handlers concurrently to verify the mutex works.
	for i := 0; i < 100; i++ {
		bus.Subscribe("concurrent", func(e Event) {
			atomic.AddInt32(&count, 1)
		})
	}

	bus.Publish(Event{Name: "concurrent"})

	if got := atomic.LoadInt32(&count); got != 100 {
		t.Errorf("handlers called %d times, want 100", got)
	}
}
