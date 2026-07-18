package sla

import (
	"time"

	"github.com/google/uuid"
)

// Policy defines SLA deadlines for a specific priority level.
type Policy struct {
	ID                    uuid.UUID
	TenantID              uuid.UUID
	Priority              string
	ResponseTimeMinutes   int
	ResolutionTimeMinutes int
	EscalationMinutes     int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
