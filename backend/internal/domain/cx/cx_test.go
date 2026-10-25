package cx

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCustomerIsolation(t *testing.T) {
	a, b := uuid.New(), uuid.New()
	if !CanAccessCustomer("admin", nil, a) || !CanAccessCustomer("agent", nil, b) {
		t.Error("staff must see every customer of their tenant")
	}
	if !CanAccessCustomer("end_user", []uuid.UUID{a}, a) {
		t.Error("contact must see their own customer")
	}
	if CanAccessCustomer("end_user", []uuid.UUID{a}, b) {
		t.Error("contact must NOT see another customer")
	}
	if CanAccessCustomer("end_user", nil, a) {
		t.Error("contact with no customers must see nothing, not everything")
	}
}

func TestTakeConversationRace(t *testing.T) {
	alice, bob := uuid.New(), uuid.New()
	free := &Conversation{Status: StatusOpen}
	if !CanTakeConversation(free, alice) {
		t.Error("unassigned conversation must be takeable")
	}
	taken := &Conversation{Status: StatusInProgress, AssignedTo: &alice}
	if !CanTakeConversation(taken, alice) {
		t.Error("assignee retaking their own conversation must pass")
	}
	if CanTakeConversation(taken, bob) {
		t.Error("a second agent must NOT take a conversation assigned to another agent")
	}
}

func TestConversationTransitions(t *testing.T) {
	c := &Conversation{Status: StatusOpen}
	if !c.CanMoveTo(StatusResolved) || c.CanMoveTo(StatusOpen) {
		t.Error("open -> resolved allowed, open -> open is not a transition")
	}
	c.Status = StatusResolved
	if !c.CanMoveTo(StatusInProgress) || c.CanMoveTo(StatusWaitingCustomer) {
		t.Error("resolved reopens to in_progress only (plus closed)")
	}
	if !MessageReopens(StatusResolved, SenderCustomer) || MessageReopens(StatusOpen, SenderCustomer) {
		t.Error("only a customer reply to a resolved conversation reopens it")
	}
}

func TestSLADefaultsAndBreach(t *testing.T) {
	now := time.Now()
	c := &Conversation{Channel: ChannelChat, Priority: PriorityMedium, Status: StatusOpen}
	ApplySLA(c, nil, now)
	if c.SLAResponseMinutes != 15 || c.SLAResolutionMinutes != 240 {
		t.Errorf("chat defaults = %d/%d, want 15/240", c.SLAResponseMinutes, c.SLAResolutionMinutes)
	}
	c.Priority = PriorityCritical
	ApplySLA(c, nil, now)
	if c.SLAResponseMinutes != 15 || c.SLAResolutionMinutes != 120 {
		t.Errorf("critical chat SLA = %d/%d, want 15/120 (halved, response floored at 15)", c.SLAResponseMinutes, c.SLAResolutionMinutes)
	}
	if c.IsSLABreached(now) {
		t.Error("conversation inside its deadlines must not be breached")
	}
	if !c.IsSLABreached(now.Add(3 * time.Hour)) {
		t.Error("conversation past its resolution deadline must be breached")
	}
	c.Status = StatusResolved
	if c.IsSLABreached(now.Add(3 * time.Hour)) {
		t.Error("resolved conversation must not count as breached")
	}
}

func TestAutoReplyMatching(t *testing.T) {
	replies := []AutoReply{
		{TriggerKind: "keyword", Keyword: "horario", ResponseText: "Horario", Active: true},
		{TriggerKind: "channel", Channel: ChannelChat, ResponseText: "Chat", Active: true},
	}
	if got := MatchAutoReply(replies, ChannelEmail, "¿Cuál es vuestro HORARIO?"); got == nil || got.ResponseText != "Horario" {
		t.Error("keyword match must be case-insensitive")
	}
	if got := MatchAutoReply(replies, ChannelChat, "hola"); got == nil || got.ResponseText != "Chat" {
		t.Error("channel trigger must fire for any message on that channel")
	}
	if got := MatchAutoReply(replies, ChannelEmail, "hola"); got != nil {
		t.Error("no trigger must fire for an unrelated email")
	}
}

func TestSurveyValidation(t *testing.T) {
	if err := ValidateSurvey(&Survey{Rating: 5}); err != nil {
		t.Error("rating 5 must be valid")
	}
	if err := ValidateSurvey(&Survey{Rating: 0}); err == nil {
		t.Error("rating 0 must fail validation")
	}
	if err := ValidateSurvey(&Survey{Rating: 6}); err == nil {
		t.Error("rating 6 must fail validation")
	}
}
