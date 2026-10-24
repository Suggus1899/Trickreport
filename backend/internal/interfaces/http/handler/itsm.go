package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appITSM "github.com/trickreport/backend/internal/application/itsm"
	domain "github.com/trickreport/backend/internal/domain/itsm"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/interfaces/http/validator"
)

// ITSMHandler serves the Stage 3 modules: service catalog, approvals,
// assets/CMDB, problems, changes and ITSM reports. Write endpoints for
// staff-only resources are additionally guarded by RequireRole in
// server.go; the approval decision rule itself is enforced in the
// application service, so hiding a button can never grant a decision.
type ITSMHandler struct {
	svc *appITSM.Service
}

// NewITSMHandler creates an ITSMHandler.
func NewITSMHandler(svc *appITSM.Service) *ITSMHandler { return &ITSMHandler{svc: svc} }

func itsmErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		response.Error(w, http.StatusNotFound, "not found")
	case errors.Is(err, domain.ErrValidation):
		response.Error(w, http.StatusBadRequest, "validation failed")
	case errors.Is(err, domain.ErrForbidden):
		response.Error(w, http.StatusForbidden, "insufficient permissions")
	case errors.Is(err, domain.ErrInvalidTransition):
		response.Error(w, http.StatusConflict, "invalid status transition")
	case errors.Is(err, domain.ErrConflict):
		response.Error(w, http.StatusConflict, "conflict")
	default:
		response.Error(w, http.StatusInternalServerError, "internal error")
	}
}

func itsmTenant(r *http.Request) uuid.UUID { return middleware.TenantFromContext(r.Context()) }

func itsmCaller(r *http.Request) (uuid.UUID, string) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		return uuid.Nil, ""
	}
	return claims.UserID, claims.Role
}

func itsmID(r *http.Request, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, name))
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func itsmDecode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	if err := validator.Validate.Struct(dst); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return false
	}
	return true
}

func itsmOptUUID(s string) *uuid.UUID {
	if s == "" {
		return nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &id
}

func itsmOptTime(s string) (*time.Time, bool) {
	if s == "" {
		return nil, true
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, false
	}
	return &t, true
}

// ── DTOs ────────────────────────────────────────────────────────────

type categoryDTO struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toCategoryDTO(c *domain.CatalogCategory) categoryDTO {
	return categoryDTO{ID: c.ID, TenantID: c.TenantID, Name: c.Name, Description: c.Description, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}

type serviceDTO struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         uuid.UUID  `json:"tenant_id"`
	CategoryID       *uuid.UUID `json:"category_id"`
	CategoryName     string     `json:"category_name"`
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	RequiresApproval bool       `json:"requires_approval"`
	Active           bool       `json:"active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func toServiceDTO(s *domain.CatalogService) serviceDTO {
	return serviceDTO{ID: s.ID, TenantID: s.TenantID, CategoryID: s.CategoryID, CategoryName: s.CategoryName, Name: s.Name, Description: s.Description, RequiresApproval: s.RequiresApproval, Active: s.Active, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt}
}

type approvalDTO struct {
	ID          uuid.UUID  `json:"id"`
	TenantID    uuid.UUID  `json:"tenant_id"`
	TargetType  string     `json:"target_type"`
	TargetID    uuid.UUID  `json:"target_id"`
	RequesterID uuid.UUID  `json:"requester_id"`
	ApproverID  *uuid.UUID `json:"approver_id"`
	Status      string     `json:"status"`
	Comment     string     `json:"comment"`
	DecidedBy   *uuid.UUID `json:"decided_by"`
	DecidedAt   *time.Time `json:"decided_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func toApprovalDTO(a *domain.Approval) approvalDTO {
	return approvalDTO{ID: a.ID, TenantID: a.TenantID, TargetType: a.TargetType, TargetID: a.TargetID, RequesterID: a.RequesterID, ApproverID: a.ApproverID, Status: a.Status, Comment: a.Comment, DecidedBy: a.DecidedBy, DecidedAt: a.DecidedAt, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt}
}

type approvalEventDTO struct {
	ID         uuid.UUID `json:"id"`
	ApprovalID uuid.UUID `json:"approval_id"`
	ActorID    uuid.UUID `json:"actor_id"`
	FromStatus string    `json:"from_status"`
	ToStatus   string    `json:"to_status"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

type ciTypeDTO struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type assetDTO struct {
	ID           uuid.UUID      `json:"id"`
	TenantID     uuid.UUID      `json:"tenant_id"`
	CITypeID     uuid.UUID      `json:"ci_type_id"`
	CITypeName   string         `json:"ci_type_name"`
	Name         string         `json:"name"`
	SerialNumber string         `json:"serial_number"`
	Status       string         `json:"status"`
	Attributes   map[string]any `json:"attributes"`
	OwnerID      *uuid.UUID     `json:"owner_id"`
	OwnerName    string         `json:"owner_name"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func toAssetDTO(a *domain.Asset) assetDTO {
	return assetDTO{ID: a.ID, TenantID: a.TenantID, CITypeID: a.CITypeID, CITypeName: a.CITypeName, Name: a.Name, SerialNumber: a.SerialNumber, Status: a.Status, Attributes: a.Attributes, OwnerID: a.OwnerID, OwnerName: a.OwnerName, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt}
}

type relationshipDTO struct {
	ID               uuid.UUID `json:"id"`
	TenantID         uuid.UUID `json:"tenant_id"`
	SourceCIID       uuid.UUID `json:"source_ci_id"`
	TargetCIID       uuid.UUID `json:"target_ci_id"`
	RelationshipType string    `json:"relationship_type"`
	CreatedAt        time.Time `json:"created_at"`
}

type problemDTO struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	RootCause   string    `json:"root_cause"`
	Workaround  string    `json:"workaround"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toProblemDTO(p *domain.Problem) problemDTO {
	return problemDTO{ID: p.ID, TenantID: p.TenantID, Title: p.Title, Description: p.Description, Status: p.Status, RootCause: p.RootCause, Workaround: p.Workaround, CreatedBy: p.CreatedBy, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

type changeDTO struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      uuid.UUID  `json:"tenant_id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	ProblemID     *uuid.UUID `json:"problem_id"`
	Status        string     `json:"status"`
	Risk          string     `json:"risk"`
	PlannedStart  *time.Time `json:"planned_start"`
	PlannedEnd    *time.Time `json:"planned_end"`
	ImplementedAt *time.Time `json:"implemented_at"`
	ClosedAt      *time.Time `json:"closed_at"`
	CreatedBy     uuid.UUID  `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func toChangeDTO(c *domain.Change) changeDTO {
	return changeDTO{ID: c.ID, TenantID: c.TenantID, Title: c.Title, Description: c.Description, ProblemID: c.ProblemID, Status: c.Status, Risk: c.Risk, PlannedStart: c.PlannedStart, PlannedEnd: c.PlannedEnd, ImplementedAt: c.ImplementedAt, ClosedAt: c.ClosedAt, CreatedBy: c.CreatedBy, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}

// ── Catalog categories ──────────────────────────────────────────────

func (h *ITSMHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.svc.ListCategories(r.Context(), itsmTenant(r))
	if err != nil {
		itsmErr(w, err)
		return
	}
	out := make([]categoryDTO, len(cats))
	for i := range cats {
		out[i] = toCategoryDTO(&cats[i])
	}
	response.JSON(w, http.StatusOK, out)
}

type categoryReq struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=10000"`
}

func (h *ITSMHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req categoryReq
	if !itsmDecode(w, r, &req) {
		return
	}
	c, err := h.svc.CreateCategory(r.Context(), itsmTenant(r), req.Name, req.Description)
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toCategoryDTO(c))
}

func (h *ITSMHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req categoryReq
	if !itsmDecode(w, r, &req) {
		return
	}
	c, err := h.svc.UpdateCategory(r.Context(), itsmTenant(r), id, req.Name, req.Description)
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toCategoryDTO(c))
}

func (h *ITSMHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteCategory(r.Context(), itsmTenant(r), id); err != nil {
		itsmErr(w, err)
		return
	}
	response.NoContent(w)
}

// ── Catalog services ────────────────────────────────────────────────

func (h *ITSMHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	svcs, err := h.svc.ListServices(r.Context(), itsmTenant(r), itsmOptUUID(r.URL.Query().Get("category_id")), r.URL.Query().Get("active_only") == "true")
	if err != nil {
		itsmErr(w, err)
		return
	}
	out := make([]serviceDTO, len(svcs))
	for i := range svcs {
		out[i] = toServiceDTO(&svcs[i])
	}
	response.JSON(w, http.StatusOK, out)
}

type serviceReq struct {
	CategoryID       *uuid.UUID `json:"category_id"`
	Name             string     `json:"name" validate:"required,min=1,max=255"`
	Description      string     `json:"description" validate:"max=10000"`
	RequiresApproval bool       `json:"requires_approval"`
	Active           bool       `json:"active"`
}

func (h *ITSMHandler) CreateService(w http.ResponseWriter, r *http.Request) {
	var req serviceReq
	if !itsmDecode(w, r, &req) {
		return
	}
	svc, err := h.svc.CreateService(r.Context(), itsmTenant(r), &domain.CatalogService{CategoryID: req.CategoryID, Name: req.Name, Description: req.Description, RequiresApproval: req.RequiresApproval, Active: req.Active})
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toServiceDTO(svc))
}

func (h *ITSMHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req serviceReq
	if !itsmDecode(w, r, &req) {
		return
	}
	svc, err := h.svc.UpdateService(r.Context(), itsmTenant(r), &domain.CatalogService{ID: id, CategoryID: req.CategoryID, Name: req.Name, Description: req.Description, RequiresApproval: req.RequiresApproval, Active: req.Active})
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toServiceDTO(svc))
}

func (h *ITSMHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteService(r.Context(), itsmTenant(r), id); err != nil {
		itsmErr(w, err)
		return
	}
	response.NoContent(w)
}

type catalogRequestReq struct {
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description string     `json:"description" validate:"max=10000"`
	Priority    string     `json:"priority" validate:"omitempty,oneof=low medium high critical"`
	ApproverID  *uuid.UUID `json:"approver_id"`
}

// RequestService creates a ticket from a catalog service and, when the
// service requires it, a pending approval for that ticket.
func (h *ITSMHandler) RequestService(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req catalogRequestReq
	if !itsmDecode(w, r, &req) {
		return
	}
	callerID, _ := itsmCaller(r)
	ticketID, err := h.svc.RequestService(r.Context(), itsmTenant(r), callerID, id, appITSM.CatalogRequest{Title: req.Title, Description: req.Description, Priority: req.Priority, ApproverID: req.ApproverID})
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"ticket_id": ticketID})
}

// ── Approvals ───────────────────────────────────────────────────────

func (h *ITSMHandler) ListApprovals(w http.ResponseWriter, r *http.Request) {
	callerID, role := itsmCaller(r)
	apprs, err := h.svc.ListApprovals(r.Context(), itsmTenant(r), callerID, role)
	if err != nil {
		itsmErr(w, err)
		return
	}
	out := make([]approvalDTO, len(apprs))
	for i := range apprs {
		out[i] = toApprovalDTO(&apprs[i])
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *ITSMHandler) ApprovalHistory(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	events, err := h.svc.ApprovalHistory(r.Context(), itsmTenant(r), id)
	if err != nil {
		itsmErr(w, err)
		return
	}
	out := make([]approvalEventDTO, len(events))
	for i := range events {
		out[i] = approvalEventDTO{ID: events[i].ID, ApprovalID: events[i].ApprovalID, ActorID: events[i].ActorID, FromStatus: events[i].FromStatus, ToStatus: events[i].ToStatus, Comment: events[i].Comment, CreatedAt: events[i].CreatedAt}
	}
	response.JSON(w, http.StatusOK, out)
}

type decideReq struct {
	Decision string `json:"decision" validate:"required,oneof=approved rejected"`
	Comment  string `json:"comment" validate:"max=10000"`
}

func (h *ITSMHandler) DecideApproval(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req decideReq
	if !itsmDecode(w, r, &req) {
		return
	}
	callerID, role := itsmCaller(r)
	appr, err := h.svc.DecideApproval(r.Context(), itsmTenant(r), id, callerID, role, req.Decision, req.Comment)
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toApprovalDTO(appr))
}

// ── CI types ────────────────────────────────────────────────────────

func (h *ITSMHandler) ListCITypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.svc.ListCITypes(r.Context(), itsmTenant(r))
	if err != nil {
		itsmErr(w, err)
		return
	}
	out := make([]ciTypeDTO, len(types))
	for i := range types {
		out[i] = ciTypeDTO{ID: types[i].ID, TenantID: types[i].TenantID, Name: types[i].Name, Description: types[i].Description, CreatedAt: types[i].CreatedAt, UpdatedAt: types[i].UpdatedAt}
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *ITSMHandler) CreateCIType(w http.ResponseWriter, r *http.Request) {
	var req categoryReq
	if !itsmDecode(w, r, &req) {
		return
	}
	ct, err := h.svc.CreateCIType(r.Context(), itsmTenant(r), req.Name, req.Description)
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, ciTypeDTO{ID: ct.ID, TenantID: ct.TenantID, Name: ct.Name, Description: ct.Description, CreatedAt: ct.CreatedAt, UpdatedAt: ct.UpdatedAt})
}

func (h *ITSMHandler) DeleteCIType(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteCIType(r.Context(), itsmTenant(r), id); err != nil {
		itsmErr(w, err)
		return
	}
	response.NoContent(w)
}

// ── Assets ──────────────────────────────────────────────────────────

func (h *ITSMHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
	assets, err := h.svc.ListAssets(r.Context(), itsmTenant(r), itsmOptUUID(r.URL.Query().Get("ci_type_id")), r.URL.Query().Get("status"))
	if err != nil {
		itsmErr(w, err)
		return
	}
	out := make([]assetDTO, len(assets))
	for i := range assets {
		out[i] = toAssetDTO(&assets[i])
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *ITSMHandler) GetAsset(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	a, err := h.svc.GetAsset(r.Context(), itsmTenant(r), id)
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toAssetDTO(a))
}

type assetReq struct {
	CITypeID     uuid.UUID      `json:"ci_type_id" validate:"required"`
	Name         string         `json:"name" validate:"required,min=1,max=255"`
	SerialNumber string         `json:"serial_number" validate:"max=255"`
	Status       string         `json:"status" validate:"omitempty,oneof=active maintenance retired stock"`
	Attributes   map[string]any `json:"attributes"`
	OwnerID      *uuid.UUID     `json:"owner_id"`
}

func (h *ITSMHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var req assetReq
	if !itsmDecode(w, r, &req) {
		return
	}
	a, err := h.svc.CreateAsset(r.Context(), itsmTenant(r), &domain.Asset{CITypeID: req.CITypeID, Name: req.Name, SerialNumber: req.SerialNumber, Status: req.Status, Attributes: req.Attributes, OwnerID: req.OwnerID})
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toAssetDTO(a))
}

func (h *ITSMHandler) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req assetReq
	if !itsmDecode(w, r, &req) {
		return
	}
	a, err := h.svc.UpdateAsset(r.Context(), itsmTenant(r), &domain.Asset{ID: id, CITypeID: req.CITypeID, Name: req.Name, SerialNumber: req.SerialNumber, Status: req.Status, Attributes: req.Attributes, OwnerID: req.OwnerID})
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toAssetDTO(a))
}

func (h *ITSMHandler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteAsset(r.Context(), itsmTenant(r), id); err != nil {
		itsmErr(w, err)
		return
	}
	response.NoContent(w)
}

// ── CMDB relationships ──────────────────────────────────────────────

func (h *ITSMHandler) ListRelationships(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	rels, err := h.svc.ListRelationships(r.Context(), itsmTenant(r), id)
	if err != nil {
		itsmErr(w, err)
		return
	}
	out := make([]relationshipDTO, len(rels))
	for i := range rels {
		out[i] = relationshipDTO{ID: rels[i].ID, TenantID: rels[i].TenantID, SourceCIID: rels[i].SourceCIID, TargetCIID: rels[i].TargetCIID, RelationshipType: rels[i].RelationshipType, CreatedAt: rels[i].CreatedAt}
	}
	response.JSON(w, http.StatusOK, out)
}

type relationshipReq struct {
	SourceCIID       uuid.UUID `json:"source_ci_id" validate:"required"`
	TargetCIID       uuid.UUID `json:"target_ci_id" validate:"required"`
	RelationshipType string    `json:"relationship_type" validate:"required,oneof=depends_on hosts uses"`
}

func (h *ITSMHandler) CreateRelationship(w http.ResponseWriter, r *http.Request) {
	var req relationshipReq
	if !itsmDecode(w, r, &req) {
		return
	}
	rel, err := h.svc.CreateRelationship(r.Context(), itsmTenant(r), &domain.CIRelationship{SourceCIID: req.SourceCIID, TargetCIID: req.TargetCIID, RelationshipType: req.RelationshipType})
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, relationshipDTO{ID: rel.ID, TenantID: rel.TenantID, SourceCIID: rel.SourceCIID, TargetCIID: rel.TargetCIID, RelationshipType: rel.RelationshipType, CreatedAt: rel.CreatedAt})
}

func (h *ITSMHandler) DeleteRelationship(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteRelationship(r.Context(), itsmTenant(r), id); err != nil {
		itsmErr(w, err)
		return
	}
	response.NoContent(w)
}

// ── Ticket ↔ asset links ────────────────────────────────────────────

type linkReq struct {
	TicketID uuid.UUID `json:"ticket_id" validate:"required"`
	AssetID  uuid.UUID `json:"asset_id" validate:"required"`
}

func (h *ITSMHandler) LinkTicketAsset(w http.ResponseWriter, r *http.Request) {
	var req linkReq
	if !itsmDecode(w, r, &req) {
		return
	}
	if err := h.svc.LinkTicketAsset(r.Context(), itsmTenant(r), req.TicketID, req.AssetID); err != nil {
		itsmErr(w, err)
		return
	}
	response.NoContent(w)
}

func (h *ITSMHandler) UnlinkTicketAsset(w http.ResponseWriter, r *http.Request) {
	ticketID, ok := itsmID(r, "ticketId")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid ticket id")
		return
	}
	assetID, ok := itsmID(r, "assetId")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid asset id")
		return
	}
	if err := h.svc.UnlinkTicketAsset(r.Context(), itsmTenant(r), ticketID, assetID); err != nil {
		itsmErr(w, err)
		return
	}
	response.NoContent(w)
}

func (h *ITSMHandler) AssetsForTicket(w http.ResponseWriter, r *http.Request) {
	ticketID, ok := itsmID(r, "ticketId")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid ticket id")
		return
	}
	assets, err := h.svc.AssetsForTicket(r.Context(), itsmTenant(r), ticketID)
	if err != nil {
		itsmErr(w, err)
		return
	}
	out := make([]assetDTO, len(assets))
	for i := range assets {
		out[i] = toAssetDTO(&assets[i])
	}
	response.JSON(w, http.StatusOK, out)
}

// ── Problems ────────────────────────────────────────────────────────

func (h *ITSMHandler) ListProblems(w http.ResponseWriter, r *http.Request) {
	problems, err := h.svc.ListProblems(r.Context(), itsmTenant(r), r.URL.Query().Get("status"))
	if err != nil {
		itsmErr(w, err)
		return
	}
	out := make([]problemDTO, len(problems))
	for i := range problems {
		out[i] = toProblemDTO(&problems[i])
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *ITSMHandler) GetProblem(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	p, err := h.svc.GetProblem(r.Context(), itsmTenant(r), id)
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toProblemDTO(p))
}

type problemReq struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=10000"`
	RootCause   string `json:"root_cause" validate:"max=10000"`
	Workaround  string `json:"workaround" validate:"max=10000"`
}

func (h *ITSMHandler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	var req problemReq
	if !itsmDecode(w, r, &req) {
		return
	}
	callerID, _ := itsmCaller(r)
	p, err := h.svc.CreateProblem(r.Context(), itsmTenant(r), callerID, &domain.Problem{Title: req.Title, Description: req.Description, RootCause: req.RootCause, Workaround: req.Workaround})
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toProblemDTO(p))
}

type problemPatchReq struct {
	Status     *string `json:"status" validate:"omitempty,oneof=open investigating resolved closed"`
	RootCause  *string `json:"root_cause" validate:"omitempty,max=10000"`
	Workaround *string `json:"workaround" validate:"omitempty,max=10000"`
}

func (h *ITSMHandler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req problemPatchReq
	if !itsmDecode(w, r, &req) {
		return
	}
	status, changeStatus := "", false
	if req.Status != nil {
		status, changeStatus = *req.Status, true
	}
	rootCause, workaround := "", ""
	if req.RootCause != nil {
		rootCause = *req.RootCause
	}
	if req.Workaround != nil {
		workaround = *req.Workaround
	}
	p, err := h.svc.UpdateProblem(r.Context(), itsmTenant(r), id, status, rootCause, workaround, changeStatus)
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toProblemDTO(p))
}

type ticketLinkReq struct {
	TicketID uuid.UUID `json:"ticket_id" validate:"required"`
}

func (h *ITSMHandler) LinkIncident(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req ticketLinkReq
	if !itsmDecode(w, r, &req) {
		return
	}
	if err := h.svc.LinkIncident(r.Context(), itsmTenant(r), id, req.TicketID); err != nil {
		itsmErr(w, err)
		return
	}
	response.NoContent(w)
}

func (h *ITSMHandler) UnlinkIncident(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	ticketID, ok := itsmID(r, "ticketId")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid ticket id")
		return
	}
	if err := h.svc.UnlinkIncident(r.Context(), itsmTenant(r), id, ticketID); err != nil {
		itsmErr(w, err)
		return
	}
	response.NoContent(w)
}

func (h *ITSMHandler) IncidentsForProblem(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	ids, err := h.svc.IncidentsForProblem(r.Context(), itsmTenant(r), id)
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, ids)
}

// ── Changes ─────────────────────────────────────────────────────────

func (h *ITSMHandler) ListChanges(w http.ResponseWriter, r *http.Request) {
	changes, err := h.svc.ListChanges(r.Context(), itsmTenant(r), r.URL.Query().Get("status"))
	if err != nil {
		itsmErr(w, err)
		return
	}
	out := make([]changeDTO, len(changes))
	for i := range changes {
		out[i] = toChangeDTO(&changes[i])
	}
	response.JSON(w, http.StatusOK, out)
}

func (h *ITSMHandler) GetChange(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	c, err := h.svc.GetChange(r.Context(), itsmTenant(r), id)
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toChangeDTO(c))
}

type changeReq struct {
	Title        string     `json:"title" validate:"required,min=1,max=255"`
	Description  string     `json:"description" validate:"max=10000"`
	ProblemID    *uuid.UUID `json:"problem_id"`
	Risk         string     `json:"risk" validate:"omitempty,oneof=low medium high critical"`
	PlannedStart string     `json:"planned_start"`
	PlannedEnd   string     `json:"planned_end"`
}

func (h *ITSMHandler) CreateChange(w http.ResponseWriter, r *http.Request) {
	var req changeReq
	if !itsmDecode(w, r, &req) {
		return
	}
	start, okStart := itsmOptTime(req.PlannedStart)
	end, okEnd := itsmOptTime(req.PlannedEnd)
	if !okStart || !okEnd {
		response.Error(w, http.StatusBadRequest, "planned dates must be RFC 3339")
		return
	}
	callerID, _ := itsmCaller(r)
	c, err := h.svc.CreateChange(r.Context(), itsmTenant(r), callerID, &domain.Change{Title: req.Title, Description: req.Description, ProblemID: req.ProblemID, Risk: req.Risk, PlannedStart: start, PlannedEnd: end})
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toChangeDTO(c))
}

type changeStatusReq struct {
	Status string `json:"status" validate:"required,oneof=draft planned awaiting_approval approved in_progress completed closed rejected"`
}

func (h *ITSMHandler) UpdateChangeStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req changeStatusReq
	if !itsmDecode(w, r, &req) {
		return
	}
	c, err := h.svc.UpdateChangeStatus(r.Context(), itsmTenant(r), id, req.Status)
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toChangeDTO(c))
}

func (h *ITSMHandler) LinkChangeTicket(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req ticketLinkReq
	if !itsmDecode(w, r, &req) {
		return
	}
	if err := h.svc.LinkChangeTicket(r.Context(), itsmTenant(r), id, req.TicketID); err != nil {
		itsmErr(w, err)
		return
	}
	response.NoContent(w)
}

func (h *ITSMHandler) UnlinkChangeTicket(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	ticketID, ok := itsmID(r, "ticketId")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid ticket id")
		return
	}
	if err := h.svc.UnlinkChangeTicket(r.Context(), itsmTenant(r), id, ticketID); err != nil {
		itsmErr(w, err)
		return
	}
	response.NoContent(w)
}

func (h *ITSMHandler) TicketsForChange(w http.ResponseWriter, r *http.Request) {
	id, ok := itsmID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	ids, err := h.svc.TicketsForChange(r.Context(), itsmTenant(r), id)
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, ids)
}

// ── ITSM report ─────────────────────────────────────────────────────

func (h *ITSMHandler) Report(w http.ResponseWriter, r *http.Request) {
	rep, err := h.svc.Report(r.Context(), itsmTenant(r))
	if err != nil {
		itsmErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, rep)
}
