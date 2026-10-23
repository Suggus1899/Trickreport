package automation

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	domain "github.com/trickreport/backend/internal/domain/automation"
)

func TestEngine_CanonicalConditionsANDAndLegacy(t *testing.T) {
	for _, conditions := range []map[string]any{
		{"priority": "high", "oldvalue": "open"},
		{"conditions": []any{map[string]any{"field": "priority", "op": "equals", "value": "high"}, map[string]any{"field": "old_value", "op": "equals", "value": "open"}}},
	} {
		for _, old := range []string{"open", "closed"} {
			exec := &mockExecutor{snapshot: &TicketSnapshot{Priority: "high"}}
			e := newEngine(&mockRepo{rules: []domain.Rule{{IsActive: true, TriggerType: "status_changed", Conditions: conditions, Actions: []any{map[string]any{"type": "set_priority", "priority": "critical"}}}}}, exec)
			if err := e.Evaluate(context.Background(), uuid.New(), Event{Type: "ticket_status_changed", TicketID: uuid.New(), OldValue: old}); err != nil {
				t.Fatal(err)
			}
			want := 0
			if old == "open" {
				want = 1
			}
			if exec.calls != want {
				t.Fatalf("conditions=%v old=%s calls=%d want=%d", conditions, old, exec.calls, want)
			}
		}
	}
}

func TestService_RejectsInvalidAutomationContracts(t *testing.T) {
	for _, input := range []CreateInput{
		{Name: "invalid", TriggerType: "invented"},
		{Name: "invalid", TriggerType: "ticket_created", Conditions: map[string]any{"unknown": "x"}},
		{Name: "invalid", TriggerType: "ticket_created", Conditions: map[string]any{"conditions": []any{map[string]any{"field": "status", "op": "regex", "value": "open"}}}},
		{Name: "invalid", TriggerType: "ticket_created", Actions: []any{map[string]any{"type": "add_tag", "tag": "x"}}},
		{Name: "invalid", TriggerType: "ticket_created", Actions: []any{map[string]any{"type": "set_status", "value": "impossible"}}},
	} {
		repo := &mockRepo{}
		if _, err := NewService(repo).Create(context.Background(), input); !errors.Is(err, domain.ErrValidation) {
			t.Errorf("input=%+v error=%v want validation", input, err)
		}
		if repo.createCalls != 0 {
			t.Error("invalid contract reached repository")
		}
	}
}

func TestEngine_ActionErrorsAreReturned(t *testing.T) {
	expected := errors.New("write failed")
	exec := &mockExecutor{setPriorityErr: expected}
	e := newEngine(&mockRepo{rules: []domain.Rule{{IsActive: true, TriggerType: "ticket_created", Actions: []any{map[string]any{"type": "set_priority", "priority": "high"}}}}}, exec)
	if err := e.Evaluate(context.Background(), uuid.New(), Event{Type: "ticket_created"}); !errors.Is(err, expected) {
		t.Fatalf("error=%v want %v", err, expected)
	}
}
