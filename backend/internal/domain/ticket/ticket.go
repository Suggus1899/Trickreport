package ticket

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/domain/event"
)

// Ticket is the aggregate root for the ticket domain.
type Ticket struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Title       string
	Description string
	Status      Status
	Priority    Priority
	Category    string
	CreatedBy   uuid.UUID
	AssignedTo  *uuid.UUID
	SLADeadline *time.Time
	SLABreached bool
	Version     int
	CreatedAt   time.Time
	UpdatedAt   time.Time

	// Events holds domain events raised by this aggregate. They are cleared
	// after being published by the application layer.
	Events []event.Event

	// Derived fields (populated by joins, not persisted directly)
	CreatorName  string
	AssigneeName *string
}

// RaiseEvent appends a domain event to the aggregate's event list.
func (t *Ticket) RaiseEvent(name string, payload any) {
	t.Events = append(t.Events, event.Event{
		Name:        name,
		AggregateID: t.ID,
		Payload:     payload,
		OccurredAt:  time.Now(),
	})
}

// ClearEvents removes all pending domain events from the aggregate.
func (t *Ticket) ClearEvents() {
	t.Events = nil
}

// CanBeViewedBy returns true if the given user can view this ticket.
// End users can only see their own tickets; agents and admins see all.
func (t *Ticket) CanBeViewedBy(userID uuid.UUID, role string) bool {
	if role == "admin" || role == "agent" {
		return true
	}
	return t.CreatedBy == userID
}

// CanStatusBeChangedBy returns true if the user role can change the status.
// End users can only close their own tickets.
func CanStatusBeChangedBy(role string, ticket *Ticket, userID uuid.UUID, target Status) bool {
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
	oldStatus := t.Status
	t.Status = newStatus
	t.UpdatedAt = time.Now()
	t.RaiseEvent("ticket_status_changed", map[string]any{
		"old_status": string(oldStatus),
		"new_status": string(newStatus),
	})
	return nil
}

// Assign sets the assignee of the ticket.
func (t *Ticket) Assign(userID uuid.UUID) {
	t.AssignedTo = &userID
	t.UpdatedAt = time.Now()
	t.RaiseEvent("ticket_assigned", map[string]any{
		"assigned_to": userID.String(),
	})
}
