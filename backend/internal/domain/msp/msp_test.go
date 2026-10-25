package msp

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestClientIsolation(t *testing.T) {
	a, b := uuid.New(), uuid.New()
	if !CanAccessClient("admin", nil, a) || !CanAccessClient("agent", nil, b) {
		t.Error("staff must see every client of their tenant")
	}
	if !CanAccessClient("end_user", []uuid.UUID{a}, a) {
		t.Error("contact must see their own client")
	}
	if CanAccessClient("end_user", []uuid.UUID{a}, b) {
		t.Error("contact must NOT see another client")
	}
	if CanAccessClient("end_user", nil, a) {
		t.Error("contact with no clients must see nothing, not everything")
	}
}

func TestContractActiveWindowAndSLA(t *testing.T) {
	now := time.Now()
	start := now.Add(-24 * time.Hour)
	end := now.Add(24 * time.Hour)
	c := &Contract{
		Status: ContractActive, StartsAt: start, EndsAt: &end,
		SLAResponseMinutes: 240, SLAResolutionMinutes: 1440,
		SLAOverrides: map[string]SLAOverride{"critical": {ResponseMinutes: 30, ResolutionMinutes: 240}},
	}
	if !c.IsActiveAt(now) {
		t.Error("contract inside its window must be active")
	}
	if c.IsActiveAt(now.Add(48 * time.Hour)) {
		t.Error("contract past its end must not be active")
	}
	if c.IsActiveAt(start.Add(-time.Hour)) {
		t.Error("contract before its start must not be active")
	}
	resp, res := c.ResolveSLA("critical")
	if resp != 30 || res != 240 {
		t.Errorf("critical override = %d/%d, want 30/240", resp, res)
	}
	resp, res = c.ResolveSLA("medium")
	if resp != 240 || res != 1440 {
		t.Errorf("medium base SLA = %d/%d, want 240/1440", resp, res)
	}
	suspended := *c
	suspended.Status = ContractSuspended
	if suspended.IsActiveAt(now) {
		t.Error("suspended contract must not resolve an SLA")
	}
	if err := ValidateContract(&Contract{Name: "x", Status: ContractActive, StartsAt: end, EndsAt: &start}); err == nil {
		t.Error("ends_at before starts_at must fail validation")
	}
}

func TestTimeEntryCost(t *testing.T) {
	if got := TimeEntryCost(30, 6000); got != 3000 {
		t.Errorf("30 min at 6000 = %d, want 3000", got)
	}
	if got := TimeEntryCost(45, 5000); got != 3750 {
		t.Errorf("45 min at 5000 = %d, want 3750", got)
	}
	if got := TimeEntryCost(7, 1000); got != 117 {
		t.Errorf("7 min at 1000 = %d, want 117 (half-up rounding)", got)
	}
	if got := TimeEntryCost(0, 5000); got != 0 {
		t.Errorf("zero minutes must cost 0, got %d", got)
	}
	if err := ValidateTimeEntry(&TimeEntry{Minutes: 0}); err == nil {
		t.Error("zero minutes must fail validation")
	}
	if err := ValidateTimeEntry(&TimeEntry{Minutes: 1441}); err == nil {
		t.Error("more than a day in one entry must fail validation")
	}
	if err := ValidateExpense(&Expense{AmountCents: 0}); err == nil {
		t.Error("zero expense must fail validation")
	}
}

func TestProjectProgress(t *testing.T) {
	if got := ProjectProgress(nil); got != 0 {
		t.Errorf("no tasks = %d%%, want 0", got)
	}
	tasks := []ProjectTask{{Status: "done"}, {Status: "doing"}, {Status: "done"}, {Status: "todo"}}
	if got := ProjectProgress(tasks); got != 50 {
		t.Errorf("2 of 4 done = %d%%, want 50", got)
	}
	p := &Project{Status: "planned"}
	if !p.CanMoveTo("active") || p.CanMoveTo("done") {
		t.Error("planned may go to active, never straight to done")
	}
	task := &ProjectTask{Status: "todo"}
	if !task.CanMoveTo("doing") || task.CanMoveTo("done") {
		t.Error("todo may go to doing, never straight to done")
	}
}

func TestAlertDedup(t *testing.T) {
	if !ShouldCreateTicketForAlert(false, "open") {
		t.Error("first open alert for a resource must create a ticket")
	}
	if ShouldCreateTicketForAlert(true, "open") {
		t.Error("an already-open alert for the same resource must NOT create a second ticket")
	}
	if ShouldCreateTicketForAlert(false, "resolved") {
		t.Error("a resolved alert must not create a ticket")
	}
}

func TestIntervalsOverlap(t *testing.T) {
	base := time.Now()
	if !IntervalsOverlap(base, base.Add(time.Hour), base.Add(30*time.Minute), base.Add(90*time.Minute)) {
		t.Error("overlapping windows must be detected")
	}
	if IntervalsOverlap(base, base.Add(time.Hour), base.Add(time.Hour), base.Add(2*time.Hour)) {
		t.Error("back-to-back windows must not overlap")
	}
	v := &Intervention{Status: "scheduled", Title: "Visit", StartsAt: base, EndsAt: base.Add(time.Hour)}
	if !v.CanMoveTo("in_progress") || v.CanMoveTo("done") {
		t.Error("scheduled may go to in_progress, never straight to done")
	}
}

func TestRemoteSessionRules(t *testing.T) {
	tech := uuid.New()
	other := uuid.New()
	s := &RemoteSession{Status: SessionAuthorized, TechnicianUserID: tech}
	if !CanStartSession(s, tech, "end_user") {
		t.Error("owning technician may start an authorized session")
	}
	if CanStartSession(s, other, "end_user") {
		t.Error("another user may not start the session")
	}
	if CanStartSession(&RemoteSession{Status: SessionRequested, TechnicianUserID: tech}, tech, "end_user") {
		t.Error("an unauthorized session must never start")
	}
	if !CanAuthorizeSession("admin") || !CanAuthorizeSession("agent") {
		t.Error("staff must authorize sessions")
	}
	if CanAuthorizeSession("end_user") {
		t.Error("a client contact must never authorize a remote session")
	}
}

func TestDraftTotalAndProfit(t *testing.T) {
	items := []BillingItem{{Kind: "time", AmountCents: 3000}, {Kind: "expense", AmountCents: 500}, {Kind: "fee", AmountCents: 10000}}
	if got := DraftTotal(items); got != 13500 {
		t.Errorf("draft total = %d, want 13500", got)
	}
	p := &Profitability{TimeRevenue: 10000, BilledDrafts: 2000, TimeCost: 4000, Expenses: 500}
	p.ComputeProfit()
	if p.Profit != 7500 {
		t.Errorf("profit = %d, want 7500", p.Profit)
	}
}
