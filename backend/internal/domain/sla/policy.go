package sla

import "time"

// Policy defines SLA deadlines for a specific priority level.
type Policy struct {
	ID                    string
	TenantID              string
	Priority              string
	ResponseTimeMinutes   int
	ResolutionTimeMinutes int
	EscalationMinutes     int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
