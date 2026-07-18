package ticket

import (
	"time"

	"github.com/google/uuid"
)

// Notification represents a user-facing notification persisted in the database.
type Notification struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	UserID    uuid.UUID  `json:"user_id"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Type      string     `json:"type"` // info, ticket_created, ticket_updated, new_comment, sla_breach
	RefID     *uuid.UUID `json:"ref_id,omitempty"`
	RefType   string     `json:"ref_type,omitempty"` // ticket, comment
	Read      bool       `json:"read"`
	CreatedAt time.Time  `json:"created_at"`
}
