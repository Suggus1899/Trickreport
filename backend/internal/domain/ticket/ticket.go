package ticket

import (
	"fmt"
	"time"
)

// Ticket is the aggregate root for the ticket domain.
type Ticket struct {
	ID          string
	TenantID    string
	Title       string
	Description string
	Status      Status
	Priority    Priority
	Category    string
	CreatedBy   string
	AssignedTo  *string
	SLADeadline *time.Time
	SLABreached bool
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Derived fields (populated by joins, not persisted directly)
	CreatorName  string
	AssigneeName *string
}

// CanBeViewedBy returns true if the given user can view this ticket.
// End users can only see their own tickets; agents and admins see all.
func (t *Ticket) CanBeViewedBy(userID, role string) bool {
	if role == "admin" || role == "agent" {
		return true
	}
	return t.CreatedBy == userID
}

// CanStatusBeChangedBy returns true if the user role can change the status.
// End users can only close their own tickets.
func CanStatusBeChangedBy(role string, ticket *Ticket, userID string, target Status) bool {
	if role == "admin" || role == "agent" {
		return true
	}
	if role == "end_user" {
		return ticket.CreatedBy == userID && target == StatusClosed
	}
	return false
}

// ChangeStatus transitions the ticket to a new status.
// Returns ErrInvalidTransition if the transition is not allowed.
func (t *Ticket) ChangeStatus(newStatus Status) error {
	if t.Status == newStatus {
		return nil
	}
	if !t.Status.CanTransitionTo(newStatus) {
		return fmt.Errorf("%w: %s → %s", ErrInvalidTransition, t.Status, newStatus)
	}
	t.Status = newStatus
	t.UpdatedAt = time.Now()
	return nil
}

// Assign sets the assignee of the ticket.
func (t *Ticket) Assign(userID string) {
	t.AssignedTo = &userID
	t.UpdatedAt = time.Now()
}
