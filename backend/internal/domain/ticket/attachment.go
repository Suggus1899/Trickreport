package ticket

import (
	"time"

	"github.com/google/uuid"
)

// Attachment is a file attached to a ticket.
type Attachment struct {
	ID          uuid.UUID
	TicketID    uuid.UUID
	UserID      uuid.UUID
	Filename    string
	ContentType string
	FileSize    int64
	CreatedAt   time.Time

	// Derived
	UserName string
}
