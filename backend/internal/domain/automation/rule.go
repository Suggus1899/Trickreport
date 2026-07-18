package automation

import "time"

// Rule is an automation rule that triggers actions on ticket events.
type Rule struct {
	ID          string
	TenantID    string
	Name        string
	Description string
	TriggerType string
	Conditions  map[string]any
	Actions     []any
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
