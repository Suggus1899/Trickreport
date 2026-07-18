package ticket

import "fmt"

// Status represents the lifecycle state of a ticket.
type Status string

const (
	StatusOpen          Status = "open"
	StatusInProgress    Status = "in_progress"
	StatusWaitingClient Status = "waiting_client"
	StatusResolved      Status = "resolved"
	StatusClosed        Status = "closed"
)

// validTransitions defines which status transitions are allowed.
var validTransitions = map[Status][]Status{
	StatusOpen:          {StatusInProgress, StatusWaitingClient, StatusResolved, StatusClosed},
	StatusInProgress:    {StatusWaitingClient, StatusResolved, StatusClosed},
	StatusWaitingClient: {StatusInProgress, StatusResolved, StatusClosed},
	StatusResolved:      {StatusClosed, StatusInProgress},
	StatusClosed:        {},
}

// CanTransitionTo returns true if the current status can transition to the target.
func (s Status) CanTransitionTo(target Status) bool {
	allowed, ok := validTransitions[s]
	if !ok {
		return false
	}
	for _, a := range allowed {
		if a == target {
			return true
		}
	}
	return false
}

// IsValid returns true if the status is a recognized value.
func (s Status) IsValid() bool {
	switch s {
	case StatusOpen, StatusInProgress, StatusWaitingClient, StatusResolved, StatusClosed:
		return true
	}
	return false
}

// ParseStatus converts a string to a Status, returning an error if invalid.
func ParseStatus(s string) (Status, error) {
	st := Status(s)
	if !st.IsValid() {
		return "", fmt.Errorf("invalid ticket status: %q", s)
	}
	return st, nil
}
