// Package itsm holds the pure domain model for Stage 3 (ITSM modules:
// service catalog, approvals, assets/CMDB, problems and changes).
// It has no database or HTTP dependency so its rules are unit-testable
// without PostgreSQL. Every repository query in the infrastructure layer
// scopes by tenant_id; these rules are the second, in-process check.
package itsm

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrNotFound is returned when a record does not exist or does not
	// belong to the caller's tenant.
	ErrNotFound = errors.New("itsm record not found")
	// ErrValidation is returned when input fails domain validation.
	ErrValidation = errors.New("itsm validation error")
	// ErrForbidden is returned when the caller may not perform the action.
	ErrForbidden = errors.New("itsm forbidden")
	// ErrInvalidTransition is returned for a disallowed status change.
	ErrInvalidTransition = errors.New("itsm invalid status transition")
	// ErrConflict is returned for duplicate links or names.
	ErrConflict = errors.New("itsm conflict")
)

// ── Catalog ─────────────────────────────────────────────────────────

type CatalogCategory struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CatalogService struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	CategoryID       *uuid.UUID
	CategoryName     string // derived, for display
	Name             string
	Description      string
	RequiresApproval bool
	Active           bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func ValidateCategoryName(name string) error {
	if strings.TrimSpace(name) == "" || len(name) > 255 {
		return ErrValidation
	}
	return nil
}

func ValidateService(s *CatalogService) error {
	if err := ValidateCategoryName(s.Name); err != nil {
		return err
	}
	if len(s.Description) > 10000 {
		return ErrValidation
	}
	return nil
}

// ── Approvals ───────────────────────────────────────────────────────

const (
	ApprovalTargetTicket  = "ticket"
	ApprovalTargetCatalog = "catalog_service"
	ApprovalTargetChange  = "change"

	ApprovalPending  = "pending"
	ApprovalApproved = "approved"
	ApprovalRejected = "rejected"
)

func ValidApprovalTarget(target string) bool {
	switch target {
	case ApprovalTargetTicket, ApprovalTargetCatalog, ApprovalTargetChange:
		return true
	}
	return false
}

type Approval struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	TargetType  string
	TargetID    uuid.UUID
	RequesterID uuid.UUID
	ApproverID  *uuid.UUID
	Status      string
	Comment     string
	DecidedBy   *uuid.UUID
	DecidedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ApprovalEvent struct {
	ID         uuid.UUID
	ApprovalID uuid.UUID
	ActorID    uuid.UUID
	FromStatus string
	ToStatus   string
	Comment    string
	CreatedAt  time.Time
}

// CanDecide reports whether caller may decide this approval. The rule is
// enforced in the application service, never only in the UI: the named
// approver, any admin, or — when no approver was assigned — any agent or
// admin. A requester who is neither approver nor staff cannot approve
// their own request.
func CanDecide(a *Approval, callerID uuid.UUID, callerRole string) bool {
	if a.Status != ApprovalPending {
		return false
	}
	if callerRole == "admin" {
		return true
	}
	if a.ApproverID != nil && *a.ApproverID == callerID {
		return true
	}
	if a.ApproverID == nil && callerRole == "agent" {
		return true
	}
	return false
}

// ── Assets / CMDB ───────────────────────────────────────────────────

const (
	RelationshipDependsOn = "depends_on"
	RelationshipHosts     = "hosts"
	RelationshipUses      = "uses"
)

func ValidRelationshipType(t string) bool {
	switch t {
	case RelationshipDependsOn, RelationshipHosts, RelationshipUses:
		return true
	}
	return false
}

type CIType struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Asset struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	CITypeID     uuid.UUID
	CITypeName   string // derived, for display
	Name         string
	SerialNumber string
	Status       string // active | maintenance | retired | stock
	Attributes   map[string]any
	OwnerID      *uuid.UUID
	OwnerName    string // derived, for display
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func ValidAssetStatus(s string) bool {
	switch s {
	case "active", "maintenance", "retired", "stock":
		return true
	}
	return false
}

func ValidateAsset(a *Asset) error {
	if strings.TrimSpace(a.Name) == "" || len(a.Name) > 255 {
		return ErrValidation
	}
	if len(a.SerialNumber) > 255 {
		return ErrValidation
	}
	if !ValidAssetStatus(a.Status) {
		return ErrValidation
	}
	return nil
}

type CIRelationship struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	SourceCIID       uuid.UUID
	TargetCIID       uuid.UUID
	RelationshipType string
	CreatedAt        time.Time
}

func ValidateRelationship(rel *CIRelationship) error {
	if !ValidRelationshipType(rel.RelationshipType) {
		return ErrValidation
	}
	if rel.SourceCIID == uuid.Nil || rel.TargetCIID == uuid.Nil {
		return ErrValidation
	}
	if rel.SourceCIID == rel.TargetCIID {
		return ErrValidation
	}
	return nil
}

// ── Problems ────────────────────────────────────────────────────────

var ProblemTransitions = map[string][]string{
	"open":          {"investigating", "resolved", "closed"},
	"investigating": {"resolved", "closed", "open"},
	"resolved":      {"closed", "investigating"},
	"closed":        {"investigating"},
}

type Problem struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Title       string
	Description string
	Status      string // open | investigating | resolved | closed
	RootCause   string
	Workaround  string
	CreatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Problem) CanMoveTo(next string) bool {
	for _, allowed := range ProblemTransitions[p.Status] {
		if allowed == next {
			return true
		}
	}
	return false
}

func ValidateProblem(p *Problem) error {
	if strings.TrimSpace(p.Title) == "" || len(p.Title) > 255 {
		return ErrValidation
	}
	if len(p.Description) > 10000 || len(p.RootCause) > 10000 || len(p.Workaround) > 10000 {
		return ErrValidation
	}
	if _, ok := ProblemTransitions[p.Status]; p.Status != "" && !ok {
		return ErrValidation
	}
	return nil
}

// ── Changes ─────────────────────────────────────────────────────────

var ChangeTransitions = map[string][]string{
	"draft":             {"planned", "rejected"},
	"planned":           {"awaiting_approval", "rejected"},
	"awaiting_approval": {"approved", "rejected", "planned"},
	"approved":          {"in_progress", "rejected"},
	"in_progress":       {"completed", "rejected"},
	"completed":         {"closed"},
	"closed":            {},
	"rejected":          {"draft"},
}

func ValidChangeRisk(risk string) bool {
	switch risk {
	case "low", "medium", "high", "critical":
		return true
	}
	return false
}

type Change struct {
	ID            uuid.UUID
	TenantID      uuid.UUID
	Title         string
	Description   string
	ProblemID     *uuid.UUID
	Status        string
	Risk          string
	PlannedStart  *time.Time
	PlannedEnd    *time.Time
	ImplementedAt *time.Time
	ClosedAt      *time.Time
	CreatedBy     uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (c *Change) CanMoveTo(next string) bool {
	for _, allowed := range ChangeTransitions[c.Status] {
		if allowed == next {
			return true
		}
	}
	return false
}

func ValidateChange(c *Change) error {
	if strings.TrimSpace(c.Title) == "" || len(c.Title) > 255 {
		return ErrValidation
	}
	if len(c.Description) > 10000 {
		return ErrValidation
	}
	if !ValidChangeRisk(c.Risk) {
		return ErrValidation
	}
	if _, ok := ChangeTransitions[c.Status]; c.Status != "" && !ok {
		return ErrValidation
	}
	if c.PlannedStart != nil && c.PlannedEnd != nil && c.PlannedEnd.Before(*c.PlannedStart) {
		return ErrValidation
	}
	return nil
}

// ── Reports ─────────────────────────────────────────────────────────

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type TypeCount struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type Report struct {
	ChangesByStatus   []StatusCount `json:"changes_by_status"`
	ProblemsByStatus  []StatusCount `json:"problems_by_status"`
	AssetsByType      []TypeCount   `json:"assets_by_type"`
	ApprovalsByStatus []StatusCount `json:"approvals_by_status"`
	PendingApprovals  int           `json:"pending_approvals"`
	ActiveServices    int           `json:"active_services"`
}
