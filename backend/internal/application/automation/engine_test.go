package automation

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	domainautomation "github.com/trickreport/backend/internal/domain/automation"
	"github.com/rs/zerolog"
)

// --- Mock TicketExecutor ---

type mockExecutor struct {
	snapshot       *TicketSnapshot
	snapshotErr    error
	setPriorityErr error
	setStatusErr   error
	assignErr      error
	commentErr     error
	tagErr         error
	gotPriority    string
	gotStatus      string
	gotAssignee    uuid.UUID
	gotComment     string
	gotInternal    bool
	gotTag         string
	calls          int
}

func (m *mockExecutor) GetTicket(ctx context.Context, tenantID, ticketID uuid.UUID) (*TicketSnapshot, error) {
	return m.snapshot, m.snapshotErr
}

func (m *mockExecutor) SetPriority(ctx context.Context, tenantID, ticketID uuid.UUID, priority string) error {
	m.calls++
	m.gotPriority = priority
	return m.setPriorityErr
}

func (m *mockExecutor) SetStatus(ctx context.Context, tenantID, ticketID uuid.UUID, status string) error {
	m.calls++
	m.gotStatus = status
	return m.setStatusErr
}

func (m *mockExecutor) AssignTo(ctx context.Context, tenantID, ticketID, userID uuid.UUID) error {
	m.calls++
	m.gotAssignee = userID
	return m.assignErr
}

func (m *mockExecutor) AddComment(ctx context.Context, tenantID, ticketID uuid.UUID, content string, isInternal bool) error {
	m.calls++
	m.gotComment = content
	m.gotInternal = isInternal
	return m.commentErr
}

func (m *mockExecutor) AddTag(ctx context.Context, tenantID, ticketID uuid.UUID, tag string) error {
	m.calls++
	m.gotTag = tag
	return m.tagErr
}

func newEngine(repo Repository, exec TicketExecutor) *Engine {
	return NewEngine(repo, exec, zerolog.Nop())
}

// --- triggerFromEvent tests ---

func TestTriggerFromEvent(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"ticket_created", "ticket_created"},
		{"ticket_status_changed", "status_changed"},
		{"ticket_priority_changed", "priority_changed"},
		{"sla_breach", "sla_breach"},
		{"custom_event", "custom_event"},
	}
	for _, tt := range tests {
		if got := triggerFromEvent(tt.in); got != tt.want {
			t.Errorf("triggerFromEvent(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// --- Evaluate tests ---

func TestEngine_Evaluate_NilEngine(t *testing.T) {
	var e *Engine
	if err := e.Evaluate(context.Background(), uuid.New(), Event{Type: "ticket_created"}); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestEngine_Evaluate_NilRepo(t *testing.T) {
	e := &Engine{}
	if err := e.Evaluate(context.Background(), uuid.New(), Event{Type: "ticket_created"}); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestEngine_Evaluate_NilExecutor(t *testing.T) {
	e := &Engine{repo: &mockRepo{}}
	if err := e.Evaluate(context.Background(), uuid.New(), Event{Type: "ticket_created"}); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestEngine_Evaluate_RepoError(t *testing.T) {
	repo := &mockRepo{err: errors.New("db")}
	e := newEngine(repo, &mockExecutor{})

	err := e.Evaluate(context.Background(), uuid.New(), Event{Type: "ticket_created"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEngine_Evaluate_NoMatchingRules(t *testing.T) {
	repo := &mockRepo{rules: []domainautomation.Rule{
		{ID: uuid.New(), IsActive: true, TriggerType: "status_changed"},
	}}
	exec := &mockExecutor{}
	e := newEngine(repo, exec)

	if err := e.Evaluate(context.Background(), uuid.New(), Event{Type: "ticket_created"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.calls != 0 {
		t.Errorf("expected 0 executor calls, got %d", exec.calls)
	}
}

func TestEngine_Evaluate_InactiveRuleSkipped(t *testing.T) {
	repo := &mockRepo{rules: []domainautomation.Rule{
		{ID: uuid.New(), IsActive: false, TriggerType: "ticket_created", Actions: []any{
			map[string]any{"type": "set_priority", "priority": "high"},
		}},
	}}
	exec := &mockExecutor{}
	e := newEngine(repo, exec)

	if err := e.Evaluate(context.Background(), uuid.New(), Event{Type: "ticket_created"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.calls != 0 {
		t.Errorf("expected 0 calls, got %d", exec.calls)
	}
}

func TestEngine_Evaluate_ConditionsNotMatched(t *testing.T) {
	ticketID := uuid.New()
	tenant := uuid.New()
	repo := &mockRepo{rules: []domainautomation.Rule{
		{ID: uuid.New(), IsActive: true, TriggerType: "ticket_created",
			Conditions: map[string]any{"priority": "high"},
			Actions:    []any{map[string]any{"type": "set_status", "status": "in_progress"}},
		},
	}}
	exec := &mockExecutor{snapshot: &TicketSnapshot{ID: ticketID, Priority: "low"}}
	e := newEngine(repo, exec)

	if err := e.Evaluate(context.Background(), tenant, Event{Type: "ticket_created", TicketID: ticketID}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.calls != 0 {
		t.Errorf("expected 0 calls, got %d", exec.calls)
	}
}

func TestEngine_Evaluate_SnapshotError(t *testing.T) {
	ticketID := uuid.New()
	repo := &mockRepo{rules: []domainautomation.Rule{
		{ID: uuid.New(), IsActive: true, TriggerType: "ticket_created", Actions: []any{
			map[string]any{"type": "set_status", "status": "in_progress"},
		}},
	}}
	exec := &mockExecutor{snapshotErr: errors.New("not found")}
	e := newEngine(repo, exec)

	// Should still execute actions even if snapshot fails (conditions empty -> match)
	if err := e.Evaluate(context.Background(), uuid.New(), Event{Type: "ticket_created", TicketID: ticketID}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.calls != 1 {
		t.Errorf("expected 1 call, got %d", exec.calls)
	}
}

func TestEngine_Evaluate_NoTicketID(t *testing.T) {
	repo := &mockRepo{rules: []domainautomation.Rule{
		{ID: uuid.New(), IsActive: true, TriggerType: "ticket_created", Actions: []any{
			map[string]any{"type": "set_status", "status": "in_progress"},
		}},
	}}
	exec := &mockExecutor{}
	e := newEngine(repo, exec)

	if err := e.Evaluate(context.Background(), uuid.New(), Event{Type: "ticket_created"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.calls != 1 {
		t.Errorf("expected 1 call, got %d", exec.calls)
	}
}

func TestEngine_Evaluate_ActionErrorContinues(t *testing.T) {
	repo := &mockRepo{rules: []domainautomation.Rule{
		{ID: uuid.New(), IsActive: true, TriggerType: "ticket_created", Actions: []any{
			map[string]any{"type": "set_priority", "priority": "high"},
			map[string]any{"type": "set_status", "status": "in_progress"},
		}},
	}}
	exec := &mockExecutor{setPriorityErr: errors.New("fail")}
	e := newEngine(repo, exec)

	if err := e.Evaluate(context.Background(), uuid.New(), Event{Type: "ticket_created"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.calls != 2 {
		t.Errorf("expected 2 calls (continues after error), got %d", exec.calls)
	}
}

// --- matchConditions tests ---

func TestEngine_MatchConditions(t *testing.T) {
	e := newEngine(&mockRepo{}, &mockExecutor{})
	snap := &TicketSnapshot{Status: "open", Priority: "high", Category: "bug", Title: "crash"}
	event := Event{OldValue: "open", NewValue: "closed"}

	tests := []struct {
		name       string
		conditions map[string]any
		snapshot   *TicketSnapshot
		want       bool
	}{
		{"empty conditions match", map[string]any{}, snap, true},
		{"status matches", map[string]any{"status": "open"}, snap, true},
		{"status case-insensitive", map[string]any{"status": "OPEN"}, snap, true},
		{"status no match", map[string]any{"status": "closed"}, snap, false},
		{"priority matches", map[string]any{"priority": "high"}, snap, true},
		{"category matches", map[string]any{"category": "bug"}, snap, true},
		{"title matches", map[string]any{"title": "crash"}, snap, true},
		{"old_value matches", map[string]any{"old_value": "open"}, snap, true},
		{"new_value matches", map[string]any{"new_value": "closed"}, snap, true},
		{"oldvalue alt key", map[string]any{"oldvalue": "open"}, snap, true},
		{"newvalue alt key", map[string]any{"newvalue": "closed"}, snap, true},
		{"nil snapshot status", map[string]any{"status": "open"}, nil, false},
		{"nil snapshot old_value still works", map[string]any{"old_value": "open"}, nil, true},
		{"unknown field fails safe", map[string]any{"unknown": "x"}, snap, false},
		{"multiple all match", map[string]any{"status": "open", "priority": "high"}, snap, true},
		{"multiple one fails", map[string]any{"status": "open", "priority": "low"}, snap, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.matchConditions(tt.conditions, tt.snapshot, event); got != tt.want {
				t.Errorf("matchConditions = %v, want %v", got, tt.want)
			}
		})
	}
}

// --- executeActions / executeAction tests ---

func TestEngine_ExecuteActions_NonMapActionSkipped(t *testing.T) {
	e := newEngine(&mockRepo{}, &mockExecutor{})
	rule := domainautomation.Rule{ID: uuid.New(), Actions: []any{"not-a-map"}}
	if err := e.executeActions(context.Background(), uuid.New(), rule, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEngine_ExecuteActions_NoTypeSkipped(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	rule := domainautomation.Rule{ID: uuid.New(), Actions: []any{map[string]any{"foo": "bar"}}}
	if err := e.executeActions(context.Background(), uuid.New(), rule, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.calls != 0 {
		t.Errorf("expected 0 calls, got %d", exec.calls)
	}
}

func TestEngine_ExecuteAction_SetPriority(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	ticketID := uuid.New()

	if err := e.executeAction(context.Background(), uuid.New(), "set_priority",
		map[string]any{"priority": "high"}, Event{TicketID: ticketID}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.gotPriority != "high" {
		t.Errorf("priority = %q", exec.gotPriority)
	}
}

func TestEngine_ExecuteAction_SetPriority_ValueFallback(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	if err := e.executeAction(context.Background(), uuid.New(), "set_priority",
		map[string]any{"value": "critical"}, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.gotPriority != "critical" {
		t.Errorf("priority = %q", exec.gotPriority)
	}
}

func TestEngine_ExecuteAction_SetPriority_Missing(t *testing.T) {
	e := newEngine(&mockRepo{}, &mockExecutor{})
	err := e.executeAction(context.Background(), uuid.New(), "set_priority", map[string]any{}, Event{}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestEngine_ExecuteAction_SetStatus(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	if err := e.executeAction(context.Background(), uuid.New(), "set_status",
		map[string]any{"status": "in_progress"}, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.gotStatus != "in_progress" {
		t.Errorf("status = %q", exec.gotStatus)
	}
}

func TestEngine_ExecuteAction_SetStatus_ValueFallback(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	if err := e.executeAction(context.Background(), uuid.New(), "set_status",
		map[string]any{"value": "closed"}, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.gotStatus != "closed" {
		t.Errorf("status = %q", exec.gotStatus)
	}
}

func TestEngine_ExecuteAction_SetStatus_Missing(t *testing.T) {
	e := newEngine(&mockRepo{}, &mockExecutor{})
	if err := e.executeAction(context.Background(), uuid.New(), "set_status", map[string]any{}, Event{}, nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestEngine_ExecuteAction_AssignTo(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	uid := uuid.New()
	if err := e.executeAction(context.Background(), uuid.New(), "assign_to",
		map[string]any{"user_id": uid.String()}, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.gotAssignee != uid {
		t.Errorf("assignee = %v", exec.gotAssignee)
	}
}

func TestEngine_ExecuteAction_AssignTo_Missing(t *testing.T) {
	e := newEngine(&mockRepo{}, &mockExecutor{})
	if err := e.executeAction(context.Background(), uuid.New(), "assign_to", map[string]any{}, Event{}, nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestEngine_ExecuteAction_AssignTo_InvalidUUID(t *testing.T) {
	e := newEngine(&mockRepo{}, &mockExecutor{})
	if err := e.executeAction(context.Background(), uuid.New(), "assign_to",
		map[string]any{"user_id": "not-a-uuid"}, Event{}, nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestEngine_ExecuteAction_AddComment(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	if err := e.executeAction(context.Background(), uuid.New(), "add_comment",
		map[string]any{"content": "hello", "is_internal": false}, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.gotComment != "hello" || exec.gotInternal != false {
		t.Errorf("comment=%q internal=%v", exec.gotComment, exec.gotInternal)
	}
}

func TestEngine_ExecuteAction_AddComment_DefaultInternal(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	if err := e.executeAction(context.Background(), uuid.New(), "add_comment",
		map[string]any{"content": "hi"}, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exec.gotInternal {
		t.Error("expected default isInternal=true")
	}
}

func TestEngine_ExecuteAction_AddComment_ValueFallback(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	if err := e.executeAction(context.Background(), uuid.New(), "add_comment",
		map[string]any{"value": "fallback"}, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.gotComment != "fallback" {
		t.Errorf("comment = %q", exec.gotComment)
	}
}

func TestEngine_ExecuteAction_AddComment_Missing(t *testing.T) {
	e := newEngine(&mockRepo{}, &mockExecutor{})
	if err := e.executeAction(context.Background(), uuid.New(), "add_comment", map[string]any{}, Event{}, nil); err == nil {
		t.Fatal("expected error")
	}
}

func TestEngine_ExecuteAction_AddTag(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	if err := e.executeAction(context.Background(), uuid.New(), "add_tag",
		map[string]any{"tag": "urgent"}, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.gotTag != "urgent" {
		t.Errorf("tag = %q", exec.gotTag)
	}
}

func TestEngine_ExecuteAction_AddTag_ValueFallback(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	if err := e.executeAction(context.Background(), uuid.New(), "add_tag",
		map[string]any{"value": "vip"}, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.gotTag != "vip" {
		t.Errorf("tag = %q", exec.gotTag)
	}
}

func TestEngine_ExecuteAction_AddTag_Missing(t *testing.T) {
	e := newEngine(&mockRepo{}, &mockExecutor{})
	if err := e.executeAction(context.Background(), uuid.New(), "add_tag", map[string]any{}, Event{}, nil); err == nil {
		t.Fatal("expected error")
	}
}

// --- Mock EmailSender ---

type mockEmailSender struct {
	err        error
	gotTo      string
	gotSubject string
	gotBody    string
	calls      int
}

func (m *mockEmailSender) Send(to, subject, body string) error {
	m.calls++
	m.gotTo = to
	m.gotSubject = subject
	m.gotBody = body
	return m.err
}

func TestEngine_ExecuteAction_SendEmail(t *testing.T) {
	exec := &mockExecutor{snapshot: &TicketSnapshot{Title: "Printer on fire", Status: "open", Priority: "critical"}}
	sender := &mockEmailSender{}
	e := newEngine(&mockRepo{}, exec)
	e.SetEmailSender(sender)

	err := e.executeAction(context.Background(), uuid.New(), "send_email",
		map[string]any{"to": "manager@example.com", "subject": "Escalated"}, Event{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sender.calls != 1 {
		t.Fatalf("expected 1 send, got %d", sender.calls)
	}
	if sender.gotTo != "manager@example.com" {
		t.Errorf("to = %q", sender.gotTo)
	}
	if sender.gotSubject != "Escalated" {
		t.Errorf("subject = %q", sender.gotSubject)
	}
	if sender.gotBody == "" {
		t.Error("expected a non-empty body composed from the ticket snapshot")
	}
}

func TestEngine_ExecuteAction_SendEmail_DefaultSubject(t *testing.T) {
	exec := &mockExecutor{}
	sender := &mockEmailSender{}
	e := newEngine(&mockRepo{}, exec)
	e.SetEmailSender(sender)

	if err := e.executeAction(context.Background(), uuid.New(), "send_email",
		map[string]any{"to": "a@b.com"}, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sender.gotSubject == "" {
		t.Error("expected a default subject when none is provided")
	}
}

func TestEngine_ExecuteAction_SendEmail_Missing(t *testing.T) {
	e := newEngine(&mockRepo{}, &mockExecutor{})
	e.SetEmailSender(&mockEmailSender{})
	if err := e.executeAction(context.Background(), uuid.New(), "send_email", map[string]any{}, Event{}, nil); err == nil {
		t.Fatal("expected error for missing 'to' field")
	}
}

func TestEngine_ExecuteAction_SendEmail_NoSenderConfigured(t *testing.T) {
	e := newEngine(&mockRepo{}, &mockExecutor{})
	// Deliberately not calling SetEmailSender — should degrade gracefully,
	// same as an unknown action type, not fail the whole rule.
	if err := e.executeAction(context.Background(), uuid.New(), "send_email",
		map[string]any{"to": "a@b.com"}, Event{}, nil); err != nil {
		t.Fatalf("expected graceful no-op, got error: %v", err)
	}
}

func TestEngine_ExecuteAction_UnknownType(t *testing.T) {
	exec := &mockExecutor{}
	e := newEngine(&mockRepo{}, exec)
	if err := e.executeAction(context.Background(), uuid.New(), "unknown_type", map[string]any{}, Event{}, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exec.calls != 0 {
		t.Errorf("expected 0 calls, got %d", exec.calls)
	}
}

func TestNewEngine(t *testing.T) {
	e := NewEngine(&mockRepo{}, &mockExecutor{}, zerolog.Nop())
	if e == nil {
		t.Fatal("expected non-nil engine")
	}
}
