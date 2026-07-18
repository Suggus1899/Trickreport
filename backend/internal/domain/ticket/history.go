package ticket

import "time"

// HistoryEntry is an audit record of a change made to a ticket.
type HistoryEntry struct {
	ID        string
	TicketID  string
	UserID    string
	Field     string
	OldValue  *string
	NewValue  *string
	CreatedAt time.Time

	// Derived
	UserName string
}
