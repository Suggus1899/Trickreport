package automation

import (
	"time"

	"github.com/google/uuid"
)

// Rule is an automation rule that triggers actions on ticket events.
type Rule struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Description string
	TriggerType string
	Conditions  map[string]any
	Actions     []any
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
