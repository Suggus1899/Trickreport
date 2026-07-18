package ticket

import "fmt"

// Priority represents the urgency of a ticket.
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// IsValid returns true if the priority is a recognized value.
func (p Priority) IsValid() bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return true
	}
	return false
}

// ParsePriority converts a string to a Priority, returning an error if invalid.
func ParsePriority(s string) (Priority, error) {
	pr := Priority(s)
	if !pr.IsValid() {
		return "", fmt.Errorf("invalid ticket priority: %q", s)
	}
	return pr, nil
}
