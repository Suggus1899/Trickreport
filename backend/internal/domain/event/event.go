package event

import (
	"time"

	"github.com/google/uuid"
)

// Event represents a domain event raised by an aggregate root.
type Event struct {
	Name        string
	AggregateID uuid.UUID
	Payload     any
	OccurredAt  time.Time
}

// EventBus is the port for publishing domain events.
type EventBus interface {
	Publish(event Event)
}

// EventHandler handles a single domain event.
type EventHandler func(event Event)
