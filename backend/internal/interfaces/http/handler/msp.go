package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appMSP "github.com/trickreport/backend/internal/application/msp"
	domain "github.com/trickreport/backend/internal/domain/msp"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/interfaces/http/validator"
)

// MSPHandler serves the Stage 4 modules: clients and contacts,
// contracts and SLA, technicians and planning, time and expenses,
// projects, billing drafts and profitability, client inventory,
// monitoring ingestion and remote sessions. Write endpoints for
// staff-only resources are additionally guarded by RequireRole in
// server.go; the per-client isolation rule itself is enforced in the
// application service, so a client contact can never read another
// client's data even if the UI showed it.
type MSPHandler struct {
	svc *appMSP.Service
}

// NewMSPHandler creates an MSPHandler.
func NewMSPHandler(svc *appMSP.Service) *MSPHandler { return &MSPHandler{svc: svc} }

func mspErr(w http.ResponseWriter, err error) {
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

func mspTenant(r *http.Request) uuid.UUID { return middleware.TenantFromContext(r.Context()) }

func mspCaller(r *http.Request) (uuid.UUID, string) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		return uuid.Nil, ""
	}
	return claims.UserID, claims.Role
}

func mspID(r *http.Request, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, name))
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func mspOptID(s string) *uuid.UUID {
	if s == "" {
		return nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil
	}
	return &id
}

func mspDecode(w http.ResponseWriter, r *http.Request, dst any) bool {
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

func mspOptDate(s string) (*time.Time, bool) {
	if s == "" {
		return nil, true
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, false
	}
	return &t, true
}

func mspDate(s string) (time.Time, bool) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// ── DTOs ────────────────────────────────────────────────────────────

type clientDTO struct {
	ID           uuid.UUID `json:"id"`
	TenantID     uuid.UUID `json:"tenant_id"`
	Name         string    `json:"name"`
	Domain       string    `json:"domain"`
	City         string    `json:"city"`
	Status       string    `json:"status"`
	ContactName  string    `json:"contact_name"`
	ContactEmail string    `json:"contact_email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func toClientDTO(c *domain.Client) clientDTO {
	return clientDTO{ID: c.ID, TenantID: c.TenantID, Name: c.Name, Domain: c.Domain, City: c.City, Status: c.Status, ContactName: c.ContactName, ContactEmail: c.ContactEmail, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}

type clientInput struct {
	Name         string `json:"name" validate:"required,max=255"`
	Domain       string `json:"domain" validate:"max=255"`
	City         string `json:"city" validate:"max=255"`
	Status       string `json:"status"`
	ContactName  string `json:"contact_name" validate:"max=255"`
	ContactEmail string `json:"contact_email" validate:"max=255"`
}

type contactDTO struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	ClientID  uuid.UUID  `json:"client_id"`
	UserID    *uuid.UUID `json:"user_id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
}

func toContactDTO(c *domain.ClientContact) contactDTO {
	return contactDTO{ID: c.ID, TenantID: c.TenantID, ClientID: c.ClientID, UserID: c.UserID, Name: c.Name, Email: c.Email, Role: c.Role, CreatedAt: c.CreatedAt}
}

type contactInput struct {
	UserID *uuid.UUID `json:"user_id"`
	Name   string     `json:"name" validate:"required,max=255"`
	Email  string     `json:"email" validate:"required,max=255"`
	Role   string     `json:"role"`
}

type contractDTO struct {
	ID                   uuid.UUID                     `json:"id"`
	TenantID             uuid.UUID                     `json:"tenant_id"`
	ClientID             uuid.UUID                     `json:"client_id"`
	Name                 string                        `json:"name"`
	Status               string                        `json:"status"`
	StartsAt             time.Time                     `json:"starts_at"`
	EndsAt               *time.Time                    `json:"ends_at"`
	Coverage             string                        `json:"coverage"`
	Currency             string                        `json:"currency"`
	MonthlyFeeCents      int64                         `json:"monthly_fee_cents"`
	HourlyRateCents      int64                         `json:"hourly_rate_cents"`
	IncludedHoursMonthly float64                       `json:"included_hours_monthly"`
	SLAResponseMinutes   int                           `json:"sla_response_minutes"`
	SLAResolutionMinutes int                           `json:"sla_resolution_minutes"`
	SLAOverrides         map[string]domain.SLAOverride `json:"sla_overrides"`
	CreatedAt            time.Time                     `json:"created_at"`
	UpdatedAt            time.Time                     `json:"updated_at"`
}

func toContractDTO(c *domain.Contract) contractDTO {
	return contractDTO{ID: c.ID, TenantID: c.TenantID, ClientID: c.ClientID, Name: c.Name, Status: c.Status, StartsAt: c.StartsAt, EndsAt: c.EndsAt, Coverage: c.Coverage, Currency: c.Currency, MonthlyFeeCents: c.MonthlyFeeCents, HourlyRateCents: c.HourlyRateCents, IncludedHoursMonthly: c.IncludedHoursMonthly, SLAResponseMinutes: c.SLAResponseMinutes, SLAResolutionMinutes: c.SLAResolutionMinutes, SLAOverrides: c.SLAOverrides, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}

type contractInput struct {
	ClientID             uuid.UUID                     `json:"client_id" validate:"required"`
	Name                 string                        `json:"name" validate:"required,max=255"`
	Status               string                        `json:"status"`
	StartsAt             time.Time                     `json:"starts_at" validate:"required"`
	EndsAt               *time.Time                    `json:"ends_at"`
	Coverage             string                        `json:"coverage" validate:"max=10000"`
	Currency             string                        `json:"currency"`
	MonthlyFeeCents      int64                         `json:"monthly_fee_cents"`
	HourlyRateCents      int64                         `json:"hourly_rate_cents"`
	IncludedHoursMonthly float64                       `json:"included_hours_monthly"`
	SLAResponseMinutes   int                           `json:"sla_response_minutes"`
	SLAResolutionMinutes int                           `json:"sla_resolution_minutes"`
	SLAOverrides         map[string]domain.SLAOverride `json:"sla_overrides"`
}

type mspTicketDTO struct {
	ID                   uuid.UUID  `json:"id"`
	Title                string     `json:"title"`
	Status               string     `json:"status"`
	Priority             string     `json:"priority"`
	Category             string     `json:"category"`
	CreatedAt            time.Time  `json:"created_at"`
	ClientID             uuid.UUID  `json:"client_id"`
	ClientName           string     `json:"client_name"`
	ContractID           *uuid.UUID `json:"contract_id"`
	SLAResponseMinutes   int        `json:"sla_response_minutes"`
	SLAResolutionMinutes int        `json:"sla_resolution_minutes"`
}

func toMSPTicketDTO(t appMSP.MSPTicket) mspTicketDTO {
	return mspTicketDTO{ID: t.ID, Title: t.Title, Status: t.Status, Priority: t.Priority, Category: t.Category, CreatedAt: t.CreatedAt, ClientID: t.ClientID, ClientName: t.ClientName, ContractID: t.ContractID, SLAResponseMinutes: t.SLAResponseMinutes, SLAResolutionMinutes: t.SLAResolutionMinutes}
}

type mspTicketInput struct {
	ClientID    uuid.UUID  `json:"client_id" validate:"required"`
	ContractID  *uuid.UUID `json:"contract_id"`
	Title       string     `json:"title" validate:"required,max=255"`
	Description string     `json:"description" validate:"max=10000"`
	Category    string     `json:"category"`
	Priority    string     `json:"priority"`
}

type ticketContextDTO struct {
	TicketID             uuid.UUID  `json:"ticket_id"`
	ClientID             uuid.UUID  `json:"client_id"`
	ClientName           string     `json:"client_name"`
	ContractID           *uuid.UUID `json:"contract_id"`
	SLAResponseMinutes   int        `json:"sla_response_minutes"`
	SLAResolutionMinutes int        `json:"sla_resolution_minutes"`
	CreatedAt            time.Time  `json:"created_at"`
}

type technicianDTO struct {
	ID              uuid.UUID `json:"id"`
	TenantID        uuid.UUID `json:"tenant_id"`
	UserID          uuid.UUID `json:"user_id"`
	UserName        string    `json:"user_name"`
	Skills          []string  `json:"skills"`
	HourlyRateCents int64     `json:"hourly_rate_cents"`
	Active          bool      `json:"active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func toTechnicianDTO(t *domain.Technician) technicianDTO {
	return technicianDTO{ID: t.ID, TenantID: t.TenantID, UserID: t.UserID, UserName: t.UserName, Skills: t.Skills, HourlyRateCents: t.HourlyRateCents, Active: t.Active, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}
}

type technicianInput struct {
	UserID          uuid.UUID `json:"user_id" validate:"required"`
	Skills          []string  `json:"skills"`
	HourlyRateCents int64     `json:"hourly_rate_cents"`
	Active          *bool     `json:"active"`
}

type availabilityDTO struct {
	ID           uuid.UUID `json:"id"`
	TechnicianID uuid.UUID `json:"technician_id"`
	DayOfWeek    int       `json:"day_of_week"`
	StartTime    string    `json:"start_time"`
	EndTime      string    `json:"end_time"`
}

func toAvailabilityDTO(a *domain.Availability) availabilityDTO {
	return availabilityDTO{ID: a.ID, TechnicianID: a.TechnicianID, DayOfWeek: a.DayOfWeek, StartTime: a.StartTime, EndTime: a.EndTime}
}

type availabilityInput struct {
	DayOfWeek int    `json:"day_of_week"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}

type interventionDTO struct {
	ID           uuid.UUID  `json:"id"`
	ClientID     uuid.UUID  `json:"client_id"`
	TechnicianID uuid.UUID  `json:"technician_id"`
	TicketID     *uuid.UUID `json:"ticket_id"`
	Title        string     `json:"title"`
	StartsAt     time.Time  `json:"starts_at"`
	EndsAt       time.Time  `json:"ends_at"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func toInterventionDTO(v *domain.Intervention) interventionDTO {
	return interventionDTO{ID: v.ID, ClientID: v.ClientID, TechnicianID: v.TechnicianID, TicketID: v.TicketID, Title: v.Title, StartsAt: v.StartsAt, EndsAt: v.EndsAt, Status: v.Status, CreatedAt: v.CreatedAt, UpdatedAt: v.UpdatedAt}
}

type interventionInput struct {
	ClientID     uuid.UUID  `json:"client_id" validate:"required"`
	TechnicianID uuid.UUID  `json:"technician_id" validate:"required"`
	TicketID     *uuid.UUID `json:"ticket_id"`
	Title        string     `json:"title" validate:"required,max=255"`
	StartsAt     time.Time  `json:"starts_at" validate:"required"`
	EndsAt       time.Time  `json:"ends_at" validate:"required"`
}

type timeEntryDTO struct {
	ID                      uuid.UUID  `json:"id"`
	ClientID                uuid.UUID  `json:"client_id"`
	TicketID                *uuid.UUID `json:"ticket_id"`
	InterventionID          *uuid.UUID `json:"intervention_id"`
	TechnicianID            uuid.UUID  `json:"technician_id"`
	Minutes                 int        `json:"minutes"`
	Description             string     `json:"description"`
	Billable                bool       `json:"billable"`
	Billed                  bool       `json:"billed"`
	HourlyRateCents         int64      `json:"hourly_rate_cents"`
	ContractHourlyRateCents *int64     `json:"contract_hourly_rate_cents"`
	CostCents               int64      `json:"cost_cents"`
	OccurredAt              time.Time  `json:"occurred_at"`
	CreatedAt               time.Time  `json:"created_at"`
}

func toTimeEntryDTO(e *domain.TimeEntry) timeEntryDTO {
	return timeEntryDTO{ID: e.ID, ClientID: e.ClientID, TicketID: e.TicketID, InterventionID: e.InterventionID, TechnicianID: e.TechnicianID, Minutes: e.Minutes, Description: e.Description, Billable: e.Billable, Billed: e.Billed, HourlyRateCents: e.HourlyRateCents, ContractHourlyRateCents: e.ContractHourlyRateCents, CostCents: e.CostCents, OccurredAt: e.OccurredAt, CreatedAt: e.CreatedAt}
}

type timeEntryInput struct {
	ClientID       uuid.UUID  `json:"client_id" validate:"required"`
	TicketID       *uuid.UUID `json:"ticket_id"`
	InterventionID *uuid.UUID `json:"intervention_id"`
	TechnicianID   uuid.UUID  `json:"technician_id" validate:"required"`
	Minutes        int        `json:"minutes" validate:"required"`
	Description    string     `json:"description" validate:"max=10000"`
	Billable       *bool      `json:"billable"`
	OccurredAt     string     `json:"occurred_at"`
}

type expenseDTO struct {
	ID          uuid.UUID  `json:"id"`
	ClientID    uuid.UUID  `json:"client_id"`
	TicketID    *uuid.UUID `json:"ticket_id"`
	ProjectID   *uuid.UUID `json:"project_id"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	AmountCents int64      `json:"amount_cents"`
	IncurredAt  time.Time  `json:"incurred_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func toExpenseDTO(e *domain.Expense) expenseDTO {
	return expenseDTO{ID: e.ID, ClientID: e.ClientID, TicketID: e.TicketID, ProjectID: e.ProjectID, Description: e.Description, Category: e.Category, AmountCents: e.AmountCents, IncurredAt: e.IncurredAt, CreatedAt: e.CreatedAt}
}

type expenseInput struct {
	ClientID    uuid.UUID  `json:"client_id" validate:"required"`
	TicketID    *uuid.UUID `json:"ticket_id"`
	ProjectID   *uuid.UUID `json:"project_id"`
	Description string     `json:"description" validate:"max=10000"`
	Category    string     `json:"category"`
	AmountCents int64      `json:"amount_cents" validate:"required"`
	IncurredAt  string     `json:"incurred_at"`
}

type projectDTO struct {
	ID          uuid.UUID  `json:"id"`
	ClientID    uuid.UUID  `json:"client_id"`
	Name        string     `json:"name"`
	Status      string     `json:"status"`
	BudgetCents int64      `json:"budget_cents"`
	StartsAt    *time.Time `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func toProjectDTO(p *domain.Project) projectDTO {
	return projectDTO{ID: p.ID, ClientID: p.ClientID, Name: p.Name, Status: p.Status, BudgetCents: p.BudgetCents, StartsAt: p.StartsAt, EndsAt: p.EndsAt, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt}
}

type projectInput struct {
	ClientID    uuid.UUID  `json:"client_id" validate:"required"`
	Name        string     `json:"name" validate:"required,max=255"`
	BudgetCents int64      `json:"budget_cents"`
	StartsAt    *time.Time `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at"`
}

type taskDTO struct {
	ID               uuid.UUID  `json:"id"`
	ProjectID        uuid.UUID  `json:"project_id"`
	Title            string     `json:"title"`
	Status           string     `json:"status"`
	TechnicianID     *uuid.UUID `json:"technician_id"`
	EstimatedMinutes int        `json:"estimated_minutes"`
	Position         int        `json:"position"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func toTaskDTO(t *domain.ProjectTask) taskDTO {
	return taskDTO{ID: t.ID, ProjectID: t.ProjectID, Title: t.Title, Status: t.Status, TechnicianID: t.TechnicianID, EstimatedMinutes: t.EstimatedMinutes, Position: t.Position, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}
}

type taskInput struct {
	Title            string     `json:"title" validate:"required,max=255"`
	TechnicianID     *uuid.UUID `json:"technician_id"`
	EstimatedMinutes int        `json:"estimated_minutes"`
	Position         int        `json:"position"`
}

type billingDraftDTO struct {
	ID          uuid.UUID  `json:"id"`
	ClientID    uuid.UUID  `json:"client_id"`
	ContractID  *uuid.UUID `json:"contract_id"`
	PeriodStart time.Time  `json:"period_start"`
	PeriodEnd   time.Time  `json:"period_end"`
	Status      string     `json:"status"`
	TotalCents  int64      `json:"total_cents"`
	Notes       string     `json:"notes"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func toBillingDraftDTO(d *domain.BillingDraft) billingDraftDTO {
	return billingDraftDTO{ID: d.ID, ClientID: d.ClientID, ContractID: d.ContractID, PeriodStart: d.PeriodStart, PeriodEnd: d.PeriodEnd, Status: d.Status, TotalCents: d.TotalCents, Notes: d.Notes, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt}
}

type billingDraftInput struct {
	ClientID    uuid.UUID `json:"client_id" validate:"required"`
	PeriodStart string    `json:"period_start" validate:"required"`
	PeriodEnd   string    `json:"period_end" validate:"required"`
	Notes       string    `json:"notes"`
}

type clientAssetDTO struct {
	AssetID    uuid.UUID `json:"asset_id"`
	ClientID   uuid.UUID `json:"client_id"`
	Name       string    `json:"name"`
	CITypeName string    `json:"ci_type_name"`
	Status     string    `json:"status"`
	Serial     string    `json:"serial_number"`
}

type discoveredDTO struct {
	ID              uuid.UUID      `json:"id"`
	ClientID        uuid.UUID      `json:"client_id"`
	ExternalID      string         `json:"external_id"`
	Name            string         `json:"name"`
	Kind            string         `json:"kind"`
	Payload         map[string]any `json:"payload"`
	Status          string         `json:"status"`
	ImportedAssetID *uuid.UUID     `json:"imported_asset_id"`
	DiscoveredAt    time.Time      `json:"discovered_at"`
}

func toDiscoveredDTO(d *domain.DiscoveredAsset) discoveredDTO {
	return discoveredDTO{ID: d.ID, ClientID: d.ClientID, ExternalID: d.ExternalID, Name: d.Name, Kind: d.Kind, Payload: d.Payload, Status: d.Status, ImportedAssetID: d.ImportedAssetID, DiscoveredAt: d.DiscoveredAt}
}

type discoveredInput struct {
	ExternalID string         `json:"external_id" validate:"required,max=255"`
	Name       string         `json:"name" validate:"required,max=255"`
	Kind       string         `json:"kind"`
	Payload    map[string]any `json:"payload"`
}

type monitoringSourceDTO struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Kind       string    `json:"kind"`
	BaseURL    string    `json:"base_url"`
	Configured bool      `json:"configured"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func toMonitoringSourceDTO(s *domain.MonitoringSource) monitoringSourceDTO {
	return monitoringSourceDTO{ID: s.ID, Name: s.Name, Kind: s.Kind, BaseURL: s.BaseURL, Configured: s.Configured, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt}
}

type monitoringSourceInput struct {
	Name         string `json:"name" validate:"required,max=255"`
	BaseURL      string `json:"base_url" validate:"max=500"`
	WebhookToken string `json:"webhook_token"`
}

type alertDTO struct {
	ID               uuid.UUID  `json:"id"`
	SourceID         uuid.UUID  `json:"source_id"`
	ClientID         uuid.UUID  `json:"client_id"`
	ExternalResource string     `json:"external_resource"`
	Severity         string     `json:"severity"`
	Status           string     `json:"status"`
	Title            string     `json:"title"`
	Message          string     `json:"message"`
	TicketID         *uuid.UUID `json:"ticket_id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func toAlertDTO(a *domain.MonitoringAlert) alertDTO {
	return alertDTO{ID: a.ID, SourceID: a.SourceID, ClientID: a.ClientID, ExternalResource: a.ExternalResource, Severity: a.Severity, Status: a.Status, Title: a.Title, Message: a.Message, TicketID: a.TicketID, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt}
}

type alertWebhookInput struct {
	ClientID         uuid.UUID `json:"client_id" validate:"required"`
	ExternalResource string    `json:"external_resource" validate:"required,max=255"`
	Severity         string    `json:"severity"`
	Status           string    `json:"status"`
	Title            string    `json:"title" validate:"max=255"`
	Message          string    `json:"message" validate:"max=10000"`
}

type remoteSessionDTO struct {
	ID                uuid.UUID  `json:"id"`
	ClientID          uuid.UUID  `json:"client_id"`
	TicketID          *uuid.UUID `json:"ticket_id"`
	TechnicianID      uuid.UUID  `json:"technician_id"`
	Provider          string     `json:"provider"`
	ExternalSessionID string     `json:"external_session_id"`
	Status            string     `json:"status"`
	AuthorizedBy      *uuid.UUID `json:"authorized_by"`
	AuditNote         string     `json:"audit_note"`
	StartedAt         *time.Time `json:"started_at"`
	EndedAt           *time.Time `json:"ended_at"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func toRemoteSessionDTO(s *domain.RemoteSession) remoteSessionDTO {
	return remoteSessionDTO{ID: s.ID, ClientID: s.ClientID, TicketID: s.TicketID, TechnicianID: s.TechnicianID, Provider: s.Provider, ExternalSessionID: s.ExternalSessionID, Status: s.Status, AuthorizedBy: s.AuthorizedBy, AuditNote: s.AuditNote, StartedAt: s.StartedAt, EndedAt: s.EndedAt, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt}
}

type remoteSessionInput struct {
	ClientID          uuid.UUID  `json:"client_id" validate:"required"`
	TicketID          *uuid.UUID `json:"ticket_id"`
	TechnicianID      uuid.UUID  `json:"technician_id" validate:"required"`
	ExternalSessionID string     `json:"external_session_id" validate:"max=255"`
	AuditNote         string     `json:"audit_note" validate:"max=10000"`
}

// ── Handlers: clients ───────────────────────────────────────────────

func (h *MSPHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	callerID, role := mspCaller(r)
	clients, err := h.svc.ListClients(r.Context(), mspTenant(r), callerID, role)
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]clientDTO, len(clients))
	for i := range clients {
		dtos[i] = toClientDTO(&clients[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	callerID, role := mspCaller(r)
	c, err := h.svc.GetClient(r.Context(), mspTenant(r), callerID, role, id)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toClientDTO(c))
}

func (h *MSPHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var in clientInput
	if !mspDecode(w, r, &in) {
		return
	}
	c, err := h.svc.CreateClient(r.Context(), mspTenant(r), &domain.Client{Name: in.Name, Domain: in.Domain, City: in.City, Status: in.Status, ContactName: in.ContactName, ContactEmail: in.ContactEmail})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toClientDTO(c))
}

func (h *MSPHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in clientInput
	if !mspDecode(w, r, &in) {
		return
	}
	c, err := h.svc.UpdateClient(r.Context(), mspTenant(r), &domain.Client{ID: id, Name: in.Name, Domain: in.Domain, City: in.City, Status: in.Status, ContactName: in.ContactName, ContactEmail: in.ContactEmail})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toClientDTO(c))
}

func (h *MSPHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.DeleteClient(r.Context(), mspTenant(r), id); err != nil {
		mspErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MSPHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
	clientID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	callerID, role := mspCaller(r)
	contacts, err := h.svc.ListContacts(r.Context(), mspTenant(r), callerID, role, clientID)
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]contactDTO, len(contacts))
	for i := range contacts {
		dtos[i] = toContactDTO(&contacts[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	clientID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in contactInput
	if !mspDecode(w, r, &in) {
		return
	}
	c, err := h.svc.CreateContact(r.Context(), mspTenant(r), clientID, &domain.ClientContact{UserID: in.UserID, Name: in.Name, Email: in.Email, Role: in.Role})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toContactDTO(c))
}

func (h *MSPHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	clientID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	contactID, ok := mspID(r, "contactId")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid contact id")
		return
	}
	if err := h.svc.DeleteContact(r.Context(), mspTenant(r), clientID, contactID); err != nil {
		mspErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ── Handlers: contracts and tickets ─────────────────────────────────

func (h *MSPHandler) ListContracts(w http.ResponseWriter, r *http.Request) {
	callerID, role := mspCaller(r)
	contracts, err := h.svc.ListContracts(r.Context(), mspTenant(r), callerID, role, mspOptID(r.URL.Query().Get("client_id")))
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]contractDTO, len(contracts))
	for i := range contracts {
		dtos[i] = toContractDTO(&contracts[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateContract(w http.ResponseWriter, r *http.Request) {
	var in contractInput
	if !mspDecode(w, r, &in) {
		return
	}
	c, err := h.svc.CreateContract(r.Context(), mspTenant(r), &domain.Contract{ClientID: in.ClientID, Name: in.Name, Status: in.Status, StartsAt: in.StartsAt, EndsAt: in.EndsAt, Coverage: in.Coverage, Currency: in.Currency, MonthlyFeeCents: in.MonthlyFeeCents, HourlyRateCents: in.HourlyRateCents, IncludedHoursMonthly: in.IncludedHoursMonthly, SLAResponseMinutes: in.SLAResponseMinutes, SLAResolutionMinutes: in.SLAResolutionMinutes, SLAOverrides: in.SLAOverrides})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toContractDTO(c))
}

func (h *MSPHandler) UpdateContract(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in contractInput
	if !mspDecode(w, r, &in) {
		return
	}
	c, err := h.svc.UpdateContract(r.Context(), mspTenant(r), &domain.Contract{ID: id, ClientID: in.ClientID, Name: in.Name, Status: in.Status, StartsAt: in.StartsAt, EndsAt: in.EndsAt, Coverage: in.Coverage, Currency: in.Currency, MonthlyFeeCents: in.MonthlyFeeCents, HourlyRateCents: in.HourlyRateCents, IncludedHoursMonthly: in.IncludedHoursMonthly, SLAResponseMinutes: in.SLAResponseMinutes, SLAResolutionMinutes: in.SLAResolutionMinutes, SLAOverrides: in.SLAOverrides})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toContractDTO(c))
}

func (h *MSPHandler) ListMSPTickets(w http.ResponseWriter, r *http.Request) {
	callerID, role := mspCaller(r)
	tickets, err := h.svc.ListMSPTickets(r.Context(), mspTenant(r), callerID, role, mspOptID(r.URL.Query().Get("client_id")))
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]mspTicketDTO, len(tickets))
	for i := range tickets {
		dtos[i] = toMSPTicketDTO(tickets[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateMSPTicket(w http.ResponseWriter, r *http.Request) {
	var in mspTicketInput
	if !mspDecode(w, r, &in) {
		return
	}
	callerID, role := mspCaller(r)
	priority := in.Priority
	if priority == "" {
		priority = "medium"
	}
	id, err := h.svc.CreateMSPTicket(r.Context(), mspTenant(r), callerID, role, appMSP.MSPTicketInput{ClientID: in.ClientID, ContractID: in.ContractID, Title: in.Title, Description: in.Description, Category: in.Category, Priority: priority})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"ticket_id": id})
}

func (h *MSPHandler) GetTicketContext(w http.ResponseWriter, r *http.Request) {
	ticketID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	callerID, role := mspCaller(r)
	tc, err := h.svc.GetTicketContext(r.Context(), mspTenant(r), callerID, role, ticketID)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, ticketContextDTO{TicketID: tc.TicketID, ClientID: tc.ClientID, ClientName: tc.ClientName, ContractID: tc.ContractID, SLAResponseMinutes: tc.SLAResponseMinutes, SLAResolutionMinutes: tc.SLAResolutionMinutes, CreatedAt: tc.CreatedAt})
}

// ── Handlers: technicians and planning ──────────────────────────────

func (h *MSPHandler) ListTechnicians(w http.ResponseWriter, r *http.Request) {
	techs, err := h.svc.ListTechnicians(r.Context(), mspTenant(r))
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]technicianDTO, len(techs))
	for i := range techs {
		dtos[i] = toTechnicianDTO(&techs[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateTechnician(w http.ResponseWriter, r *http.Request) {
	var in technicianInput
	if !mspDecode(w, r, &in) {
		return
	}
	active := true
	if in.Active != nil {
		active = *in.Active
	}
	t, err := h.svc.CreateTechnician(r.Context(), mspTenant(r), &domain.Technician{UserID: in.UserID, Skills: in.Skills, HourlyRateCents: in.HourlyRateCents, Active: active})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toTechnicianDTO(t))
}

func (h *MSPHandler) UpdateTechnician(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in technicianInput
	if !mspDecode(w, r, &in) {
		return
	}
	active := true
	if in.Active != nil {
		active = *in.Active
	}
	t, err := h.svc.UpdateTechnician(r.Context(), mspTenant(r), &domain.Technician{ID: id, UserID: in.UserID, Skills: in.Skills, HourlyRateCents: in.HourlyRateCents, Active: active})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toTechnicianDTO(t))
}

func (h *MSPHandler) ListAvailability(w http.ResponseWriter, r *http.Request) {
	techID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	av, err := h.svc.ListAvailability(r.Context(), mspTenant(r), techID)
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]availabilityDTO, len(av))
	for i := range av {
		dtos[i] = toAvailabilityDTO(&av[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateAvailability(w http.ResponseWriter, r *http.Request) {
	techID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in availabilityInput
	if !mspDecode(w, r, &in) {
		return
	}
	a, err := h.svc.CreateAvailability(r.Context(), mspTenant(r), &domain.Availability{TechnicianID: techID, DayOfWeek: in.DayOfWeek, StartTime: in.StartTime, EndTime: in.EndTime})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toAvailabilityDTO(a))
}

func (h *MSPHandler) ListInterventions(w http.ResponseWriter, r *http.Request) {
	callerID, role := mspCaller(r)
	vs, err := h.svc.ListInterventions(r.Context(), mspTenant(r), callerID, role, mspOptID(r.URL.Query().Get("client_id")), mspOptID(r.URL.Query().Get("technician_id")))
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]interventionDTO, len(vs))
	for i := range vs {
		dtos[i] = toInterventionDTO(&vs[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateIntervention(w http.ResponseWriter, r *http.Request) {
	var in interventionInput
	if !mspDecode(w, r, &in) {
		return
	}
	v, err := h.svc.CreateIntervention(r.Context(), mspTenant(r), &domain.Intervention{ClientID: in.ClientID, TechnicianID: in.TechnicianID, TicketID: in.TicketID, Title: in.Title, StartsAt: in.StartsAt, EndsAt: in.EndsAt})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toInterventionDTO(v))
}

func (h *MSPHandler) UpdateInterventionStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in struct {
		Status string `json:"status" validate:"required"`
	}
	if !mspDecode(w, r, &in) {
		return
	}
	v, err := h.svc.UpdateInterventionStatus(r.Context(), mspTenant(r), id, in.Status)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toInterventionDTO(v))
}

// ── Handlers: time and expenses ─────────────────────────────────────

func (h *MSPHandler) ListTimeEntries(w http.ResponseWriter, r *http.Request) {
	callerID, role := mspCaller(r)
	from, ok1 := mspOptDate(r.URL.Query().Get("from"))
	to, ok2 := mspOptDate(r.URL.Query().Get("to"))
	if !ok1 || !ok2 {
		response.Error(w, http.StatusBadRequest, "invalid date filter")
		return
	}
	entries, err := h.svc.ListTimeEntries(r.Context(), mspTenant(r), callerID, role, mspOptID(r.URL.Query().Get("client_id")), from, to)
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]timeEntryDTO, len(entries))
	for i := range entries {
		dtos[i] = toTimeEntryDTO(&entries[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateTimeEntry(w http.ResponseWriter, r *http.Request) {
	var in timeEntryInput
	if !mspDecode(w, r, &in) {
		return
	}
	occurred := time.Now()
	if in.OccurredAt != "" {
		t, ok := mspDate(in.OccurredAt)
		if !ok {
			response.Error(w, http.StatusBadRequest, "invalid occurred_at")
			return
		}
		occurred = t
	}
	billable := true
	if in.Billable != nil {
		billable = *in.Billable
	}
	e, err := h.svc.CreateTimeEntry(r.Context(), mspTenant(r), &domain.TimeEntry{ClientID: in.ClientID, TicketID: in.TicketID, InterventionID: in.InterventionID, TechnicianID: in.TechnicianID, Minutes: in.Minutes, Description: in.Description, Billable: billable, OccurredAt: occurred})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toTimeEntryDTO(e))
}

func (h *MSPHandler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	callerID, role := mspCaller(r)
	from, ok1 := mspOptDate(r.URL.Query().Get("from"))
	to, ok2 := mspOptDate(r.URL.Query().Get("to"))
	if !ok1 || !ok2 {
		response.Error(w, http.StatusBadRequest, "invalid date filter")
		return
	}
	expenses, err := h.svc.ListExpenses(r.Context(), mspTenant(r), callerID, role, mspOptID(r.URL.Query().Get("client_id")), from, to)
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]expenseDTO, len(expenses))
	for i := range expenses {
		dtos[i] = toExpenseDTO(&expenses[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	var in expenseInput
	if !mspDecode(w, r, &in) {
		return
	}
	incurred := time.Now()
	if in.IncurredAt != "" {
		t, ok := mspDate(in.IncurredAt)
		if !ok {
			response.Error(w, http.StatusBadRequest, "invalid incurred_at")
			return
		}
		incurred = t
	}
	category := in.Category
	if category == "" {
		category = "other"
	}
	e, err := h.svc.CreateExpense(r.Context(), mspTenant(r), &domain.Expense{ClientID: in.ClientID, TicketID: in.TicketID, ProjectID: in.ProjectID, Description: in.Description, Category: category, AmountCents: in.AmountCents, IncurredAt: incurred})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toExpenseDTO(e))
}

// ── Handlers: projects ──────────────────────────────────────────────

func (h *MSPHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	callerID, role := mspCaller(r)
	projects, err := h.svc.ListProjects(r.Context(), mspTenant(r), callerID, role, mspOptID(r.URL.Query().Get("client_id")))
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]projectDTO, len(projects))
	for i := range projects {
		dtos[i] = toProjectDTO(&projects[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var in projectInput
	if !mspDecode(w, r, &in) {
		return
	}
	p, err := h.svc.CreateProject(r.Context(), mspTenant(r), &domain.Project{ClientID: in.ClientID, Name: in.Name, BudgetCents: in.BudgetCents, StartsAt: in.StartsAt, EndsAt: in.EndsAt})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toProjectDTO(p))
}

func (h *MSPHandler) UpdateProjectStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in struct {
		Status string `json:"status" validate:"required"`
	}
	if !mspDecode(w, r, &in) {
		return
	}
	p, err := h.svc.UpdateProjectStatus(r.Context(), mspTenant(r), id, in.Status)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toProjectDTO(p))
}

func (h *MSPHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	projectID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	viewerID, viewerRole := mspCaller(r)
	tasks, err := h.svc.ListTasks(r.Context(), mspTenant(r), viewerID, viewerRole, projectID)
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]taskDTO, len(tasks))
	for i := range tasks {
		dtos[i] = toTaskDTO(&tasks[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	projectID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in taskInput
	if !mspDecode(w, r, &in) {
		return
	}
	t, err := h.svc.CreateTask(r.Context(), mspTenant(r), &domain.ProjectTask{ProjectID: projectID, Title: in.Title, TechnicianID: in.TechnicianID, EstimatedMinutes: in.EstimatedMinutes, Position: in.Position})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toTaskDTO(t))
}

func (h *MSPHandler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	projectID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	taskID, ok := mspID(r, "taskId")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid task id")
		return
	}
	var in struct {
		Status string `json:"status" validate:"required"`
	}
	if !mspDecode(w, r, &in) {
		return
	}
	t, err := h.svc.UpdateTaskStatus(r.Context(), mspTenant(r), projectID, taskID, in.Status)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toTaskDTO(t))
}

// ── Handlers: billing and profitability ─────────────────────────────

func (h *MSPHandler) ListBillingDrafts(w http.ResponseWriter, r *http.Request) {
	callerID, role := mspCaller(r)
	drafts, err := h.svc.ListBillingDrafts(r.Context(), mspTenant(r), callerID, role, mspOptID(r.URL.Query().Get("client_id")))
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]billingDraftDTO, len(drafts))
	for i := range drafts {
		dtos[i] = toBillingDraftDTO(&drafts[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) GenerateBillingDraft(w http.ResponseWriter, r *http.Request) {
	var in billingDraftInput
	if !mspDecode(w, r, &in) {
		return
	}
	from, ok1 := mspDate(in.PeriodStart)
	to, ok2 := mspDate(in.PeriodEnd)
	if !ok1 || !ok2 {
		response.Error(w, http.StatusBadRequest, "invalid period")
		return
	}
	d, err := h.svc.GenerateDraft(r.Context(), mspTenant(r), in.ClientID, from, to, in.Notes)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toBillingDraftDTO(d))
}

func (h *MSPHandler) UpdateBillingDraftStatus(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in struct {
		Status string `json:"status" validate:"required"`
	}
	if !mspDecode(w, r, &in) {
		return
	}
	d, err := h.svc.UpdateBillingDraftStatus(r.Context(), mspTenant(r), id, in.Status)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toBillingDraftDTO(d))
}

func (h *MSPHandler) Profitability(w http.ResponseWriter, r *http.Request) {
	callerID, role := mspCaller(r)
	p, err := h.svc.Profitability(r.Context(), mspTenant(r), callerID, role, mspOptID(r.URL.Query().Get("client_id")), mspOptID(r.URL.Query().Get("project_id")))
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, p)
}

func (h *MSPHandler) Report(w http.ResponseWriter, r *http.Request) {
	rep, err := h.svc.Report(r.Context(), mspTenant(r))
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, rep)
}

// ── Handlers: inventory ─────────────────────────────────────────────

func (h *MSPHandler) ListClientAssets(w http.ResponseWriter, r *http.Request) {
	clientID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	callerID, role := mspCaller(r)
	assets, err := h.svc.ListClientAssets(r.Context(), mspTenant(r), callerID, role, clientID)
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]clientAssetDTO, len(assets))
	for i := range assets {
		dtos[i] = clientAssetDTO{AssetID: assets[i].AssetID, ClientID: assets[i].ClientID, Name: assets[i].Name, CITypeName: assets[i].CITypeName, Status: assets[i].Status, Serial: assets[i].Serial}
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) LinkClientAsset(w http.ResponseWriter, r *http.Request) {
	clientID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in struct {
		AssetID uuid.UUID `json:"asset_id" validate:"required"`
	}
	if !mspDecode(w, r, &in) {
		return
	}
	if err := h.svc.LinkAssetToClient(r.Context(), mspTenant(r), clientID, in.AssetID); err != nil {
		mspErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MSPHandler) UnlinkClientAsset(w http.ResponseWriter, r *http.Request) {
	clientID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	assetID, ok := mspID(r, "assetId")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid asset id")
		return
	}
	if err := h.svc.UnlinkAssetFromClient(r.Context(), mspTenant(r), clientID, assetID); err != nil {
		mspErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MSPHandler) ListDiscoveredAssets(w http.ResponseWriter, r *http.Request) {
	clientID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	callerID, role := mspCaller(r)
	assets, err := h.svc.ListDiscoveredAssets(r.Context(), mspTenant(r), callerID, role, clientID)
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]discoveredDTO, len(assets))
	for i := range assets {
		dtos[i] = toDiscoveredDTO(&assets[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateDiscoveredAsset(w http.ResponseWriter, r *http.Request) {
	clientID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in discoveredInput
	if !mspDecode(w, r, &in) {
		return
	}
	kind := in.Kind
	if kind == "" {
		kind = "host"
	}
	d, err := h.svc.CreateDiscoveredAsset(r.Context(), mspTenant(r), clientID, &domain.DiscoveredAsset{ExternalID: in.ExternalID, Name: in.Name, Kind: kind, Payload: in.Payload})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toDiscoveredDTO(d))
}

func (h *MSPHandler) ImportDiscoveredAsset(w http.ResponseWriter, r *http.Request) {
	clientID, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	discoveredID, ok := mspID(r, "discoveredId")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid discovery id")
		return
	}
	var in struct {
		CITypeID  uuid.UUID `json:"ci_type_id" validate:"required"`
		AssetName string    `json:"asset_name" validate:"required,max=255"`
	}
	if !mspDecode(w, r, &in) {
		return
	}
	assetID, err := h.svc.ImportDiscoveredAsset(r.Context(), mspTenant(r), clientID, discoveredID, in.CITypeID, in.AssetName)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{"asset_id": assetID})
}

// ── Handlers: monitoring and remote sessions ───────────────────────

func (h *MSPHandler) ListMonitoringSources(w http.ResponseWriter, r *http.Request) {
	sources, err := h.svc.ListMonitoringSources(r.Context(), mspTenant(r))
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]monitoringSourceDTO, len(sources))
	for i := range sources {
		dtos[i] = toMonitoringSourceDTO(&sources[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) CreateMonitoringSource(w http.ResponseWriter, r *http.Request) {
	var in monitoringSourceInput
	if !mspDecode(w, r, &in) {
		return
	}
	s, err := h.svc.CreateMonitoringSource(r.Context(), mspTenant(r), &domain.MonitoringSource{Name: in.Name, BaseURL: in.BaseURL}, in.WebhookToken)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toMonitoringSourceDTO(s))
}

func (h *MSPHandler) RotateSourceToken(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in struct {
		WebhookToken string `json:"webhook_token" validate:"required"`
	}
	if !mspDecode(w, r, &in) {
		return
	}
	if err := h.svc.RotateSourceToken(r.Context(), mspTenant(r), id, in.WebhookToken); err != nil {
		mspErr(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *MSPHandler) ListMonitoringAlerts(w http.ResponseWriter, r *http.Request) {
	callerID, role := mspCaller(r)
	alerts, err := h.svc.ListMonitoringAlerts(r.Context(), mspTenant(r), callerID, role, mspOptID(r.URL.Query().Get("client_id")))
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]alertDTO, len(alerts))
	for i := range alerts {
		dtos[i] = toAlertDTO(&alerts[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) ResolveMonitoringAlert(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	a, err := h.svc.ResolveMonitoringAlert(r.Context(), mspTenant(r), id)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toAlertDTO(a))
}

// Webhook ingests one Zabbix delivery. It runs outside the session
// group: authentication is the per-source token in
// X-MSP-Webhook-Token plus the tenant from X-Tenant-ID (resolved by
// the Tenant middleware). A wrong token or tenant yields the same
// "not found" as a missing source, and the token itself is never
// logged: only its hash is stored and compared.
func (h *MSPHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-MSP-Webhook-Token")
	if token == "" {
		response.Error(w, http.StatusUnauthorized, "missing webhook token")
		return
	}
	tenantID := mspTenant(r)
	if tenantID == uuid.Nil {
		response.Error(w, http.StatusBadRequest, "missing tenant")
		return
	}
	var in alertWebhookInput
	if !mspDecode(w, r, &in) {
		return
	}
	res, err := h.svc.IngestAlert(r.Context(), tenantID, token, appMSP.AlertInput{ClientID: in.ClientID, ExternalResource: in.ExternalResource, Severity: in.Severity, Status: in.Status, Title: in.Title, Message: in.Message})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{"alert": toAlertDTO(res.Alert), "ticket_created": res.TicketCreated})
}

func (h *MSPHandler) ListRemoteSessions(w http.ResponseWriter, r *http.Request) {
	callerID, role := mspCaller(r)
	sessions, err := h.svc.ListRemoteSessions(r.Context(), mspTenant(r), callerID, role, mspOptID(r.URL.Query().Get("client_id")))
	if err != nil {
		mspErr(w, err)
		return
	}
	dtos := make([]remoteSessionDTO, len(sessions))
	for i := range sessions {
		dtos[i] = toRemoteSessionDTO(&sessions[i])
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *MSPHandler) RequestRemoteSession(w http.ResponseWriter, r *http.Request) {
	var in remoteSessionInput
	if !mspDecode(w, r, &in) {
		return
	}
	callerID, callerRole := mspCaller(r)
	s, err := h.svc.RequestRemoteSession(r.Context(), mspTenant(r), callerID, callerRole, &domain.RemoteSession{ClientID: in.ClientID, TicketID: in.TicketID, TechnicianID: in.TechnicianID, ExternalSessionID: in.ExternalSessionID, AuditNote: in.AuditNote})
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, toRemoteSessionDTO(s))
}

func (h *MSPHandler) AuthorizeRemoteSession(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var in struct {
		Authorize bool `json:"authorize"`
	}
	if !mspDecode(w, r, &in) {
		return
	}
	callerID, role := mspCaller(r)
	s, err := h.svc.AuthorizeRemoteSession(r.Context(), mspTenant(r), id, callerID, role, in.Authorize)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toRemoteSessionDTO(s))
}

func (h *MSPHandler) StartRemoteSession(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	callerID, role := mspCaller(r)
	s, err := h.svc.StartRemoteSession(r.Context(), mspTenant(r), id, callerID, role)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toRemoteSessionDTO(s))
}

func (h *MSPHandler) EndRemoteSession(w http.ResponseWriter, r *http.Request) {
	id, ok := mspID(r, "id")
	if !ok {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	s, err := h.svc.EndRemoteSession(r.Context(), mspTenant(r), id)
	if err != nil {
		mspErr(w, err)
		return
	}
	response.JSON(w, http.StatusOK, toRemoteSessionDTO(s))
}
