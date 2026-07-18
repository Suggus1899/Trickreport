package ticket

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from Status
		to   Status
		want bool
	}{
		{"open to in_progress", StatusOpen, StatusInProgress, true},
		{"open to closed", StatusOpen, StatusClosed, true},
		{"open to open", StatusOpen, StatusOpen, false},
		{"in_progress to open", StatusInProgress, StatusOpen, false},
		{"resolved to in_progress", StatusResolved, StatusInProgress, true},
		{"resolved to closed", StatusResolved, StatusClosed, true},
		{"closed to any", StatusClosed, StatusOpen, false},
		{"invalid from", Status("invalid"), StatusOpen, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("%q.CanTransitionTo(%q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestStatus_IsValid(t *testing.T) {
	for _, s := range []Status{StatusOpen, StatusInProgress, StatusWaitingClient, StatusResolved, StatusClosed} {
		t.Run(string(s), func(t *testing.T) {
			if !s.IsValid() {
				t.Errorf("expected %q to be valid", s)
			}
		})
	}

	t.Run("invalid", func(t *testing.T) {
		if Status("gone").IsValid() {
			t.Error("expected invalid status to be false")
		}
	})
}

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Status
		wantErr bool
	}{
		{"open", "open", StatusOpen, false},
		{"in_progress", "in_progress", StatusInProgress, false},
		{"invalid", "invalid", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStatus(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseStatus(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPriority_IsValid(t *testing.T) {
	for _, p := range []Priority{PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical} {
		t.Run(string(p), func(t *testing.T) {
			if !p.IsValid() {
				t.Errorf("expected %q to be valid", p)
			}
		})
	}
}

func TestParsePriority(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Priority
		wantErr bool
	}{
		{"critical", "critical", PriorityCritical, false},
		{"invalid", "urgent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePriority(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParsePriority(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParsePriority(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestTicket_CanBeViewedBy(t *testing.T) {
	creator := uuid.New()
	other := uuid.New()
	ticket := &Ticket{CreatedBy: creator}

	tests := []struct {
		name   string
		userID uuid.UUID
		role   string
		want   bool
	}{
		{"end user owns ticket", creator, "end_user", true},
		{"end user views other", other, "end_user", false},
		{"agent can view all", other, "agent", true},
		{"admin can view all", other, "admin", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ticket.CanBeViewedBy(tt.userID, tt.role); got != tt.want {
				t.Errorf("CanBeViewedBy(%v, %q) = %v, want %v", tt.userID, tt.role, got, tt.want)
			}
		})
	}
}

func TestCanStatusBeChangedBy(t *testing.T) {
	creator := uuid.New()
	ticket := &Ticket{CreatedBy: creator}

	tests := []struct {
		name   string
		role   string
		userID uuid.UUID
		target Status
		want   bool
	}{
		{"admin changes to resolved", "admin", uuid.New(), StatusResolved, true},
		{"agent changes to closed", "agent", uuid.New(), StatusClosed, true},
		{"end user closes own", "end_user", creator, StatusClosed, true},
		{"end user closes other", "end_user", uuid.New(), StatusClosed, false},
		{"end user reopens", "end_user", creator, StatusOpen, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanStatusBeChangedBy(tt.role, ticket, tt.userID, tt.target); got != tt.want {
				t.Errorf("CanStatusBeChangedBy(%q, %v, %q) = %v, want %v", tt.role, tt.userID, tt.target, got, tt.want)
			}
		})
	}
}

func TestTicket_ChangeStatus(t *testing.T) {
	ticket := &Ticket{Status: StatusOpen}

	t.Run("valid transition", func(t *testing.T) {
		if err := ticket.ChangeStatus(StatusInProgress); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ticket.Status != StatusInProgress {
			t.Errorf("got %q, want %q", ticket.Status, StatusInProgress)
		}
	})

	t.Run("same status is no-op", func(t *testing.T) {
		if err := ticket.ChangeStatus(StatusInProgress); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("invalid transition returns wrapped error", func(t *testing.T) {
		ticket := &Ticket{Status: StatusClosed}
		err := ticket.ChangeStatus(StatusOpen)
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrInvalidTransition) {
			t.Errorf("expected ErrInvalidTransition, got %v", err)
		}
	})
}

func TestTicket_Assign(t *testing.T) {
	ticket := &Ticket{}
	agentID := uuid.New()
	ticket.Assign(agentID)

	if ticket.AssignedTo == nil || *ticket.AssignedTo != agentID {
		t.Errorf("expected assignee %v, got %v", agentID, ticket.AssignedTo)
	}
}
