package itsm

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCanDecide(t *testing.T) {
	approver := uuid.New()
	other := uuid.New()
	pending := &Approval{Status: ApprovalPending, ApproverID: &approver}
	unassigned := &Approval{Status: ApprovalPending}
	decided := &Approval{Status: ApprovalApproved, ApproverID: &approver}

	tests := []struct {
		name   string
		appr   *Approval
		caller uuid.UUID
		role   string
		want   bool
	}{
		{"named approver decides", pending, approver, "agent", true},
		{"non-approver agent cannot decide named approval", pending, other, "agent", false},
		{"non-approver end_user cannot decide", pending, other, "end_user", false},
		{"admin can decide named approval", pending, other, "admin", true},
		{"agent can decide unassigned approval", unassigned, other, "agent", true},
		{"end_user cannot decide unassigned approval", unassigned, other, "end_user", false},
		{"decided approval is closed to everyone", decided, approver, "admin", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanDecide(tt.appr, tt.caller, tt.role); got != tt.want {
				t.Errorf("CanDecide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProblemTransitions(t *testing.T) {
	p := &Problem{Status: "open"}
	if !p.CanMoveTo("investigating") {
		t.Error("open -> investigating should be allowed")
	}
	if p.CanMoveTo("open") {
		t.Error("open -> open should not be allowed")
	}
	p.Status = "closed"
	if p.CanMoveTo("resolved") {
		t.Error("closed -> resolved should not be allowed")
	}
	if !p.CanMoveTo("investigating") {
		t.Error("closed -> investigating (reopen) should be allowed")
	}
}

func TestChangeTransitionsAndApprovalGate(t *testing.T) {
	c := &Change{Status: "draft"}
	if c.CanMoveTo("approved") {
		t.Error("draft -> approved must skip the approval step")
	}
	if !c.CanMoveTo("planned") {
		t.Error("draft -> planned should be allowed")
	}
	c.Status = "completed"
	if c.CanMoveTo("in_progress") {
		t.Error("completed -> in_progress should not be allowed")
	}
	if !c.CanMoveTo("closed") {
		t.Error("completed -> closed should be allowed")
	}
}

func TestValidateChangeWindow(t *testing.T) {
	now := time.Now()
	later := now.Add(-time.Hour)
	c := &Change{Title: "Upgrade", Status: "draft", Risk: "low", PlannedStart: &now, PlannedEnd: &later}
	if err := ValidateChange(c); err == nil {
		t.Error("planned_end before planned_start should fail validation")
	}
	c.Risk = "unknown"
	c.PlannedEnd = nil
	if err := ValidateChange(c); err == nil {
		t.Error("unknown risk should fail validation")
	}
}

func TestValidateRelationship(t *testing.T) {
	a, b := uuid.New(), uuid.New()
	if err := ValidateRelationship(&CIRelationship{SourceCIID: a, TargetCIID: b, RelationshipType: "hosts"}); err != nil {
		t.Errorf("valid relationship rejected: %v", err)
	}
	if err := ValidateRelationship(&CIRelationship{SourceCIID: a, TargetCIID: a, RelationshipType: "hosts"}); err == nil {
		t.Error("self-relationship should fail validation")
	}
	if err := ValidateRelationship(&CIRelationship{SourceCIID: a, TargetCIID: b, RelationshipType: "contains"}); err == nil {
		t.Error("unknown relationship type should fail validation")
	}
}

func TestValidateAsset(t *testing.T) {
	if err := ValidateAsset(&Asset{Name: "srv-01", Status: "active"}); err != nil {
		t.Errorf("valid asset rejected: %v", err)
	}
	if err := ValidateAsset(&Asset{Name: "", Status: "active"}); err == nil {
		t.Error("empty asset name should fail validation")
	}
	if err := ValidateAsset(&Asset{Name: "srv-01", Status: "lost"}); err == nil {
		t.Error("unknown asset status should fail validation")
	}
}

func TestValidApprovalTarget(t *testing.T) {
	for _, target := range []string{"ticket", "catalog_service", "change"} {
		if !ValidApprovalTarget(target) {
			t.Errorf("target %q should be valid", target)
		}
	}
	if ValidApprovalTarget("problem") {
		t.Error("problems are not approval targets")
	}
}
