package ticket

import (
	"time"

	"github.com/google/uuid"
)

// HistoryEntry is an audit record of a change made to a ticket.
type HistoryEntry struct {
	ID        uuid.UUID
	TicketID  uuid.UUID
	UserID    uuid.UUID
	Field     string
	OldValue  *string
	NewValue  *string
	CreatedAt time.Time

	// Derived
	UserName string
}
