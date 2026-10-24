// Package itsm implements the Stage 3 ITSM application logic: service
// catalog, approvals, assets/CMDB, problems and changes, plus the ITSM
// report. Authorization is enforced here, in Go, not in the frontend:
// route middleware only narrows write endpoints to staff, while the
// approval decision rule (domain.CanDecide) and the change-approval gate
// are checked per record in this service.
package itsm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	domain "github.com/trickreport/backend/internal/domain/itsm"
)

// Repository is the persistence boundary for ITSM aggregates. Every
// method scopes by tenantID; implementations must return ErrNotFound
// for records outside the tenant.
type Repository interface {
	// Catalog
	ListCategories(ctx context.Context, tenantID uuid.UUID) ([]domain.CatalogCategory, error)
	CreateCategory(ctx context.Context, c *domain.CatalogCategory) error
	UpdateCategory(ctx context.Context, c *domain.CatalogCategory) error
	DeleteCategory(ctx context.Context, id, tenantID uuid.UUID) error
	ListServices(ctx context.Context, tenantID uuid.UUID, categoryID *uuid.UUID, activeOnly bool) ([]domain.CatalogService, error)
	GetService(ctx context.Context, id, tenantID uuid.UUID) (*domain.CatalogService, error)
	CreateService(ctx context.Context, s *domain.CatalogService) error
	UpdateService(ctx context.Context, s *domain.CatalogService) error
	DeleteService(ctx context.Context, id, tenantID uuid.UUID) error
	// Tickets (only what the catalog request and links need)
	CreateTicketFromService(ctx context.Context, tenantID, requesterID uuid.UUID, title, description, category, priority string) (uuid.UUID, error)
	TicketExists(ctx context.Context, id, tenantID uuid.UUID) (bool, error)
	UserExists(ctx context.Context, id, tenantID uuid.UUID) (bool, error)
	// Approvals
	ListApprovals(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string) ([]domain.Approval, error)
	GetApproval(ctx context.Context, id, tenantID uuid.UUID) (*domain.Approval, error)
	CreateApproval(ctx context.Context, a *domain.Approval) error
	DecideApproval(ctx context.Context, a *domain.Approval, actorID uuid.UUID) error
	ApprovalHistory(ctx context.Context, id, tenantID uuid.UUID) ([]domain.ApprovalEvent, error)
	HasPendingApproval(ctx context.Context, tenantID uuid.UUID, targetType string, targetID uuid.UUID) (bool, error)
	HasApprovedApproval(ctx context.Context, tenantID uuid.UUID, targetType string, targetID uuid.UUID) (bool, error)
	// Assets / CMDB
	ListCITypes(ctx context.Context, tenantID uuid.UUID) ([]domain.CIType, error)
	CreateCIType(ctx context.Context, ct *domain.CIType) error
	DeleteCIType(ctx context.Context, id, tenantID uuid.UUID) error
	ListAssets(ctx context.Context, tenantID uuid.UUID, ciTypeID *uuid.UUID, status string) ([]domain.Asset, error)
	GetAsset(ctx context.Context, id, tenantID uuid.UUID) (*domain.Asset, error)
	CreateAsset(ctx context.Context, a *domain.Asset) error
	UpdateAsset(ctx context.Context, a *domain.Asset) error
	DeleteAsset(ctx context.Context, id, tenantID uuid.UUID) error
	ListRelationships(ctx context.Context, tenantID, assetID uuid.UUID) ([]domain.CIRelationship, error)
	CreateRelationship(ctx context.Context, rel *domain.CIRelationship) error
	DeleteRelationship(ctx context.Context, id, tenantID uuid.UUID) error
	LinkTicketAsset(ctx context.Context, tenantID, ticketID, assetID uuid.UUID) error
	UnlinkTicketAsset(ctx context.Context, tenantID, ticketID, assetID uuid.UUID) error
	AssetsForTicket(ctx context.Context, tenantID, ticketID uuid.UUID) ([]domain.Asset, error)
	// Problems
	ListProblems(ctx context.Context, tenantID uuid.UUID, status string) ([]domain.Problem, error)
	GetProblem(ctx context.Context, id, tenantID uuid.UUID) (*domain.Problem, error)
	CreateProblem(ctx context.Context, p *domain.Problem) error
	UpdateProblem(ctx context.Context, p *domain.Problem) error
	LinkIncident(ctx context.Context, tenantID, problemID, ticketID uuid.UUID) error
	UnlinkIncident(ctx context.Context, tenantID, problemID, ticketID uuid.UUID) error
	IncidentsForProblem(ctx context.Context, tenantID, problemID uuid.UUID) ([]uuid.UUID, error)
	// Changes
	ListChanges(ctx context.Context, tenantID uuid.UUID, status string) ([]domain.Change, error)
	GetChange(ctx context.Context, id, tenantID uuid.UUID) (*domain.Change, error)
	CreateChange(ctx context.Context, c *domain.Change) error
	UpdateChange(ctx context.Context, c *domain.Change) error
	LinkChangeTicket(ctx context.Context, tenantID, changeID, ticketID uuid.UUID) error
	UnlinkChangeTicket(ctx context.Context, tenantID, changeID, ticketID uuid.UUID) error
	TicketsForChange(ctx context.Context, tenantID, changeID uuid.UUID) ([]uuid.UUID, error)
	// Reports
	Report(ctx context.Context, tenantID uuid.UUID) (*domain.Report, error)
}

// Service is the ITSM application service.
type Service struct {
	repo Repository
}

// NewService creates a Service from a Repository implementation.
func NewService(repo Repository) *Service { return &Service{repo: repo} }

// ── Catalog ─────────────────────────────────────────────────────────

func (s *Service) ListCategories(ctx context.Context, tenantID uuid.UUID) ([]domain.CatalogCategory, error) {
	return s.repo.ListCategories(ctx, tenantID)
}

func (s *Service) CreateCategory(ctx context.Context, tenantID uuid.UUID, name, description string) (*domain.CatalogCategory, error) {
	if err := domain.ValidateCategoryName(name); err != nil {
		return nil, err
	}
	c := &domain.CatalogCategory{TenantID: tenantID, Name: strings.TrimSpace(name), Description: description}
	if err := s.repo.CreateCategory(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) UpdateCategory(ctx context.Context, tenantID, id uuid.UUID, name, description string) (*domain.CatalogCategory, error) {
	if err := domain.ValidateCategoryName(name); err != nil {
		return nil, err
	}
	c := &domain.CatalogCategory{ID: id, TenantID: tenantID, Name: strings.TrimSpace(name), Description: description}
	if err := s.repo.UpdateCategory(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCategory(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteCategory(ctx, id, tenantID)
}

func (s *Service) ListServices(ctx context.Context, tenantID uuid.UUID, categoryID *uuid.UUID, activeOnly bool) ([]domain.CatalogService, error) {
	return s.repo.ListServices(ctx, tenantID, categoryID, activeOnly)
}

func (s *Service) CreateService(ctx context.Context, tenantID uuid.UUID, svc *domain.CatalogService) (*domain.CatalogService, error) {
	svc.TenantID = tenantID
	svc.Name = strings.TrimSpace(svc.Name)
	if err := domain.ValidateService(svc); err != nil {
		return nil, err
	}
	if err := s.repo.CreateService(ctx, svc); err != nil {
		return nil, mapRepoErr(err)
	}
	return svc, nil
}

func (s *Service) UpdateService(ctx context.Context, tenantID uuid.UUID, svc *domain.CatalogService) (*domain.CatalogService, error) {
	svc.TenantID = tenantID
	svc.Name = strings.TrimSpace(svc.Name)
	if err := domain.ValidateService(svc); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateService(ctx, svc); err != nil {
		return nil, mapRepoErr(err)
	}
	return svc, nil
}

func (s *Service) DeleteService(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteService(ctx, id, tenantID)
}

// CatalogRequest is the input for requesting a service from the catalog.
type CatalogRequest struct {
	Title       string
	Description string
	Priority    string
	ApproverID  *uuid.UUID
}

// RequestService creates a ticket from a catalog service. When the
// service requires approval, a pending approval is created alongside
// the ticket: the requested approver, or an unassigned approval that
// any agent/admin may decide (see domain.CanDecide).
func (s *Service) RequestService(ctx context.Context, tenantID, requesterID, serviceID uuid.UUID, req CatalogRequest) (ticketID uuid.UUID, err error) {
	if strings.TrimSpace(req.Title) == "" || len(req.Title) > 255 || len(req.Description) > 10000 {
		return uuid.Nil, domain.ErrValidation
	}
	svc, err := s.repo.GetService(ctx, serviceID, tenantID)
	if err != nil {
		return uuid.Nil, err
	}
	if !svc.Active {
		return uuid.Nil, fmt.Errorf("%w: service is not active", domain.ErrValidation)
	}
	if req.ApproverID != nil {
		ok, err := s.repo.UserExists(ctx, *req.ApproverID, tenantID)
		if err != nil {
			return uuid.Nil, err
		}
		if !ok {
			return uuid.Nil, domain.ErrNotFound
		}
	}
	category := svc.CategoryName
	if category == "" {
		category = "general"
	}
	if len(category) > 100 {
		category = category[:100]
	}
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}
	switch priority {
	case "low", "medium", "high", "critical":
	default:
		return uuid.Nil, domain.ErrValidation
	}
	ticketID, err = s.repo.CreateTicketFromService(ctx, tenantID, requesterID, strings.TrimSpace(req.Title), req.Description, category, priority)
	if err != nil {
		return uuid.Nil, err
	}
	if svc.RequiresApproval {
		appr := &domain.Approval{
			TenantID: tenantID, TargetType: domain.ApprovalTargetTicket, TargetID: ticketID,
			RequesterID: requesterID, ApproverID: req.ApproverID, Status: domain.ApprovalPending,
		}
		if err := s.repo.CreateApproval(ctx, appr); err != nil {
			return uuid.Nil, err
		}
	}
	return ticketID, nil
}

// ── Approvals ───────────────────────────────────────────────────────

func (s *Service) ListApprovals(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string) ([]domain.Approval, error) {
	return s.repo.ListApprovals(ctx, tenantID, viewerID, viewerRole)
}

func (s *Service) ApprovalHistory(ctx context.Context, tenantID, id uuid.UUID) ([]domain.ApprovalEvent, error) {
	if _, err := s.repo.GetApproval(ctx, id, tenantID); err != nil {
		return nil, err
	}
	return s.repo.ApprovalHistory(ctx, id, tenantID)
}

// DecideApproval enforces the decision rule in Go: only the named
// approver, an admin, or staff for an unassigned approval may decide,
// and only while the approval is pending.
func (s *Service) DecideApproval(ctx context.Context, tenantID, id, callerID uuid.UUID, callerRole, decision, comment string) (*domain.Approval, error) {
	if decision != domain.ApprovalApproved && decision != domain.ApprovalRejected {
		return nil, domain.ErrValidation
	}
	if len(comment) > 10000 {
		return nil, domain.ErrValidation
	}
	appr, err := s.repo.GetApproval(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	if !domain.CanDecide(appr, callerID, callerRole) {
		return nil, domain.ErrForbidden
	}
	appr.Status = decision
	appr.Comment = comment
	now := time.Now()
	appr.DecidedBy = &callerID
	appr.DecidedAt = &now
	if err := s.repo.DecideApproval(ctx, appr, callerID); err != nil {
		return nil, err
	}
	return appr, nil
}

// ── Assets / CMDB ───────────────────────────────────────────────────

func (s *Service) ListCITypes(ctx context.Context, tenantID uuid.UUID) ([]domain.CIType, error) {
	return s.repo.ListCITypes(ctx, tenantID)
}

func (s *Service) CreateCIType(ctx context.Context, tenantID uuid.UUID, name, description string) (*domain.CIType, error) {
	if err := domain.ValidateCategoryName(name); err != nil {
		return nil, err
	}
	ct := &domain.CIType{TenantID: tenantID, Name: strings.TrimSpace(name), Description: description}
	if err := s.repo.CreateCIType(ctx, ct); err != nil {
		return nil, mapRepoErr(err)
	}
	return ct, nil
}

func (s *Service) DeleteCIType(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteCIType(ctx, id, tenantID)
}

func (s *Service) ListAssets(ctx context.Context, tenantID uuid.UUID, ciTypeID *uuid.UUID, status string) ([]domain.Asset, error) {
	if status != "" && !domain.ValidAssetStatus(status) {
		return nil, domain.ErrValidation
	}
	return s.repo.ListAssets(ctx, tenantID, ciTypeID, status)
}

func (s *Service) GetAsset(ctx context.Context, tenantID, id uuid.UUID) (*domain.Asset, error) {
	return s.repo.GetAsset(ctx, id, tenantID)
}

func (s *Service) CreateAsset(ctx context.Context, tenantID uuid.UUID, a *domain.Asset) (*domain.Asset, error) {
	a.TenantID = tenantID
	if a.Status == "" {
		a.Status = "active"
	}
	if err := domain.ValidateAsset(a); err != nil {
		return nil, err
	}
	if err := s.repo.CreateAsset(ctx, a); err != nil {
		return nil, mapRepoErr(err)
	}
	return a, nil
}

func (s *Service) UpdateAsset(ctx context.Context, tenantID uuid.UUID, a *domain.Asset) (*domain.Asset, error) {
	a.TenantID = tenantID
	if err := domain.ValidateAsset(a); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateAsset(ctx, a); err != nil {
		return nil, mapRepoErr(err)
	}
	return a, nil
}

func (s *Service) DeleteAsset(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteAsset(ctx, id, tenantID)
}

func (s *Service) ListRelationships(ctx context.Context, tenantID, assetID uuid.UUID) ([]domain.CIRelationship, error) {
	if _, err := s.repo.GetAsset(ctx, assetID, tenantID); err != nil {
		return nil, err
	}
	return s.repo.ListRelationships(ctx, tenantID, assetID)
}

func (s *Service) CreateRelationship(ctx context.Context, tenantID uuid.UUID, rel *domain.CIRelationship) (*domain.CIRelationship, error) {
	rel.TenantID = tenantID
	if err := domain.ValidateRelationship(rel); err != nil {
		return nil, err
	}
	// Both endpoints must exist in this tenant; the repo enforces it
	// and returns ErrNotFound otherwise.
	if err := s.repo.CreateRelationship(ctx, rel); err != nil {
		return nil, mapRepoErr(err)
	}
	return rel, nil
}

func (s *Service) DeleteRelationship(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteRelationship(ctx, id, tenantID)
}

func (s *Service) LinkTicketAsset(ctx context.Context, tenantID, ticketID, assetID uuid.UUID) error {
	if err := s.checkTicketAndAsset(ctx, tenantID, ticketID, assetID); err != nil {
		return err
	}
	return mapRepoErr(s.repo.LinkTicketAsset(ctx, tenantID, ticketID, assetID))
}

func (s *Service) UnlinkTicketAsset(ctx context.Context, tenantID, ticketID, assetID uuid.UUID) error {
	return s.repo.UnlinkTicketAsset(ctx, tenantID, ticketID, assetID)
}

func (s *Service) AssetsForTicket(ctx context.Context, tenantID, ticketID uuid.UUID) ([]domain.Asset, error) {
	ok, err := s.repo.TicketExists(ctx, ticketID, tenantID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrNotFound
	}
	return s.repo.AssetsForTicket(ctx, tenantID, ticketID)
}

func (s *Service) checkTicketAndAsset(ctx context.Context, tenantID, ticketID, assetID uuid.UUID) error {
	ok, err := s.repo.TicketExists(ctx, ticketID, tenantID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrNotFound
	}
	if _, err := s.repo.GetAsset(ctx, assetID, tenantID); err != nil {
		return err
	}
	return nil
}

// ── Problems ────────────────────────────────────────────────────────

func (s *Service) ListProblems(ctx context.Context, tenantID uuid.UUID, status string) ([]domain.Problem, error) {
	if status != "" {
		if _, ok := domain.ProblemTransitions[status]; !ok {
			return nil, domain.ErrValidation
		}
	}
	return s.repo.ListProblems(ctx, tenantID, status)
}

func (s *Service) GetProblem(ctx context.Context, tenantID, id uuid.UUID) (*domain.Problem, error) {
	return s.repo.GetProblem(ctx, id, tenantID)
}

func (s *Service) CreateProblem(ctx context.Context, tenantID, creatorID uuid.UUID, p *domain.Problem) (*domain.Problem, error) {
	p.TenantID = tenantID
	p.CreatedBy = creatorID
	if p.Status == "" {
		p.Status = "open"
	}
	if err := domain.ValidateProblem(p); err != nil {
		return nil, err
	}
	if err := s.repo.CreateProblem(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

// UpdateProblemStatus moves a problem through its state machine and
// optionally records root cause / workaround.
func (s *Service) UpdateProblem(ctx context.Context, tenantID, id uuid.UUID, status, rootCause, workaround string, changeStatus bool) (*domain.Problem, error) {
	p, err := s.repo.GetProblem(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	if changeStatus {
		if !p.CanMoveTo(status) {
			return nil, domain.ErrInvalidTransition
		}
		p.Status = status
	}
	if rootCause != "" {
		p.RootCause = rootCause
	}
	if workaround != "" {
		p.Workaround = workaround
	}
	if err := domain.ValidateProblem(p); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateProblem(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) LinkIncident(ctx context.Context, tenantID, problemID, ticketID uuid.UUID) error {
	if _, err := s.repo.GetProblem(ctx, problemID, tenantID); err != nil {
		return err
	}
	ok, err := s.repo.TicketExists(ctx, ticketID, tenantID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrNotFound
	}
	return mapRepoErr(s.repo.LinkIncident(ctx, tenantID, problemID, ticketID))
}

func (s *Service) UnlinkIncident(ctx context.Context, tenantID, problemID, ticketID uuid.UUID) error {
	return s.repo.UnlinkIncident(ctx, tenantID, problemID, ticketID)
}

func (s *Service) IncidentsForProblem(ctx context.Context, tenantID, problemID uuid.UUID) ([]uuid.UUID, error) {
	if _, err := s.repo.GetProblem(ctx, problemID, tenantID); err != nil {
		return nil, err
	}
	return s.repo.IncidentsForProblem(ctx, tenantID, problemID)
}

// ── Changes ─────────────────────────────────────────────────────────

func (s *Service) ListChanges(ctx context.Context, tenantID uuid.UUID, status string) ([]domain.Change, error) {
	if status != "" {
		if _, ok := domain.ChangeTransitions[status]; !ok {
			return nil, domain.ErrValidation
		}
	}
	return s.repo.ListChanges(ctx, tenantID, status)
}

func (s *Service) GetChange(ctx context.Context, tenantID, id uuid.UUID) (*domain.Change, error) {
	return s.repo.GetChange(ctx, id, tenantID)
}

func (s *Service) CreateChange(ctx context.Context, tenantID, creatorID uuid.UUID, c *domain.Change) (*domain.Change, error) {
	c.TenantID = tenantID
	c.CreatedBy = creatorID
	if c.Status == "" {
		c.Status = "draft"
	}
	if c.Risk == "" {
		c.Risk = "medium"
	}
	if err := domain.ValidateChange(c); err != nil {
		return nil, err
	}
	if c.ProblemID != nil {
		if _, err := s.repo.GetProblem(ctx, *c.ProblemID, tenantID); err != nil {
			return nil, err
		}
	}
	if err := s.repo.CreateChange(ctx, c); err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

// UpdateChangeStatus moves a change through its state machine. Two
// gates are enforced here: moving to "approved" requires an approved
// approval for the change, and moving to "awaiting_approval" creates a
// pending (unassigned) approval when none exists.
func (s *Service) UpdateChangeStatus(ctx context.Context, tenantID, id uuid.UUID, next string) (*domain.Change, error) {
	c, err := s.repo.GetChange(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	if !c.CanMoveTo(next) {
		return nil, domain.ErrInvalidTransition
	}
	if next == "approved" {
		ok, err := s.repo.HasApprovedApproval(ctx, tenantID, domain.ApprovalTargetChange, id)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("%w: change requires an approved approval first", domain.ErrInvalidTransition)
		}
	}
	c.Status = next
	now := time.Now()
	if next == "in_progress" && c.ImplementedAt == nil {
		c.ImplementedAt = &now
	}
	if next == "closed" && c.ClosedAt == nil {
		c.ClosedAt = &now
	}
	if err := s.repo.UpdateChange(ctx, c); err != nil {
		return nil, err
	}
	if next == "awaiting_approval" {
		pending, err := s.repo.HasPendingApproval(ctx, tenantID, domain.ApprovalTargetChange, id)
		if err != nil {
			return nil, err
		}
		if !pending {
			appr := &domain.Approval{
				TenantID: tenantID, TargetType: domain.ApprovalTargetChange, TargetID: id,
				RequesterID: c.CreatedBy, Status: domain.ApprovalPending,
			}
			if err := s.repo.CreateApproval(ctx, appr); err != nil {
				return nil, err
			}
		}
	}
	return c, nil
}

func (s *Service) LinkChangeTicket(ctx context.Context, tenantID, changeID, ticketID uuid.UUID) error {
	if _, err := s.repo.GetChange(ctx, changeID, tenantID); err != nil {
		return err
	}
	ok, err := s.repo.TicketExists(ctx, ticketID, tenantID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrNotFound
	}
	return mapRepoErr(s.repo.LinkChangeTicket(ctx, tenantID, changeID, ticketID))
}

func (s *Service) UnlinkChangeTicket(ctx context.Context, tenantID, changeID, ticketID uuid.UUID) error {
	return s.repo.UnlinkChangeTicket(ctx, tenantID, changeID, ticketID)
}

func (s *Service) TicketsForChange(ctx context.Context, tenantID, changeID uuid.UUID) ([]uuid.UUID, error) {
	if _, err := s.repo.GetChange(ctx, changeID, tenantID); err != nil {
		return nil, err
	}
	return s.repo.TicketsForChange(ctx, tenantID, changeID)
}

// ── Reports ─────────────────────────────────────────────────────────

func (s *Service) Report(ctx context.Context, tenantID uuid.UUID) (*domain.Report, error) {
	return s.repo.Report(ctx, tenantID)
}

// mapRepoErr translates persistence conflicts into domain errors; other
// errors pass through unchanged.
func mapRepoErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrConflict) || errors.Is(err, domain.ErrValidation) {
		return err
	}
	return err
}
