package realtime

import "github.com/google/uuid"

// EventType represents the type of a real-time event.
type EventType string

const (
	EventTicketCreated EventType = "TICKET_CREATED"
	EventTicketUpdated EventType = "TICKET_UPDATED"
	EventNewComment    EventType = "NEW_COMMENT"
)

// Message is the JSON payload sent to the client.
type Message struct {
	Type     EventType `json:"type"`
	TenantID uuid.UUID `json:"tenant_id,omitempty"`
	Data     any       `json:"data"`
}

// BroadcastPayload is used internally by the Hub to route messages.
type BroadcastPayload struct {
	TenantID uuid.UUID
	Message  []byte
}
