package event

import (
	"sync"
)

// MemoryEventBus is a simple in-memory implementation of EventBus.
// Handlers are invoked synchronously on Publish.
type MemoryEventBus struct {
	mu       sync.RWMutex
	handlers map[string][]EventHandler
}

// NewMemoryEventBus creates a new MemoryEventBus.
func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe registers a handler for events with the given name.
func (b *MemoryEventBus) Subscribe(name string, handler EventHandler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[name] = append(b.handlers[name], handler)
}

// Publish dispatches the event to all registered handlers synchronously.
func (b *MemoryEventBus) Publish(event Event) {
	b.mu.RLock()
	handlers := b.handlers[event.Name]
	b.mu.RUnlock()

	for _, h := range handlers {
		h(event)
	}
}

// Compile-time assertion that MemoryEventBus implements EventBus.
var _ EventBus = (*MemoryEventBus)(nil)
