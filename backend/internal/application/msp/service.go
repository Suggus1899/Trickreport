// Package msp implements the Stage 4 MSP application logic: clients and
// their contacts, contracts and SLA resolution, technicians and
// planning, time and expenses, projects, billing drafts, profitability,
// client inventory, monitoring ingestion and remote sessions.
//
// Authorization is enforced here, in Go, not in the frontend: route
// middleware only narrows write endpoints to staff, while the
// per-client isolation rule (domain.CanAccessClient) and the remote
// session authorization rule are checked per record in this service.
// A client contact (end_user) that is not linked to a client sees that
// client as "not found" in every method below.
package msp

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	domain "github.com/trickreport/backend/internal/domain/msp"
	domainSLA "github.com/trickreport/backend/internal/domain/sla"
)

// Repository is the persistence boundary for MSP aggregates. Every
// method scopes by tenantID, and every client-owned method additionally
// by clientID; implementations must return ErrNotFound for records
// outside the tenant (or outside the caller's clients).
type Repository interface {
	// Clients and contacts
	ContactClientIDs(ctx context.Context, tenantID, userID uuid.UUID) ([]uuid.UUID, error)
	ListClients(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID) ([]domain.Client, error)
	GetClient(ctx context.Context, id, tenantID uuid.UUID) (*domain.Client, error)
	CreateClient(ctx context.Context, c *domain.Client) error
	UpdateClient(ctx context.Context, c *domain.Client) error
	DeleteClient(ctx context.Context, id, tenantID uuid.UUID) error
	ListContacts(ctx context.Context, tenantID, clientID uuid.UUID) ([]domain.ClientContact, error)
	CreateContact(ctx context.Context, c *domain.ClientContact) error
	DeleteContact(ctx context.Context, id, tenantID, clientID uuid.UUID) error
	// Contracts
	ListContracts(ctx context.Context, tenantID uuid.UUID, clientID *uuid.UUID) ([]domain.Contract, error)
	GetContract(ctx context.Context, id, tenantID uuid.UUID) (*domain.Contract, error)
	CreateContract(ctx context.Context, c *domain.Contract) error
	UpdateContract(ctx context.Context, c *domain.Contract) error
	ActiveContractForClient(ctx context.Context, tenantID, clientID uuid.UUID, at time.Time) (*domain.Contract, error)
	// Tickets
	CreateTicketWithContext(ctx context.Context, tenantID, creatorID uuid.UUID, title, description, category, priority string, tc *domain.TicketContext, dl domainSLA.Deadlines) (uuid.UUID, error)
	AlertTicketCreator(ctx context.Context, tenantID uuid.UUID) (uuid.UUID, error)
	TicketContext(ctx context.Context, tenantID, ticketID uuid.UUID) (*domain.TicketContext, error)
	ListMSPTickets(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID) ([]MSPTicket, error)
	// Technicians and planning
	UserExists(ctx context.Context, id, tenantID uuid.UUID) (bool, error)
	ListTechnicians(ctx context.Context, tenantID uuid.UUID) ([]domain.Technician, error)
	GetTechnician(ctx context.Context, id, tenantID uuid.UUID) (*domain.Technician, error)
	CreateTechnician(ctx context.Context, t *domain.Technician) error
	UpdateTechnician(ctx context.Context, t *domain.Technician) error
	ListAvailability(ctx context.Context, tenantID, technicianID uuid.UUID) ([]domain.Availability, error)
	CreateAvailability(ctx context.Context, a *domain.Availability) error
	ListInterventions(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter, technicianFilter *uuid.UUID) ([]domain.Intervention, error)
	CreateIntervention(ctx context.Context, v *domain.Intervention) error
	UpdateInterventionStatus(ctx context.Context, v *domain.Intervention) error
	TechnicianOverlappingIntervention(ctx context.Context, tenantID, technicianID uuid.UUID, startsAt, endsAt time.Time) (bool, error)
	// Time and expenses
	CreateTimeEntry(ctx context.Context, e *domain.TimeEntry) error
	ListTimeEntries(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID, from, to *time.Time) ([]domain.TimeEntry, error)
	UnbilledTimeEntries(ctx context.Context, tenantID, clientID uuid.UUID, from, to time.Time) ([]domain.TimeEntry, error)
	UnbilledExpenses(ctx context.Context, tenantID, clientID uuid.UUID, from, to time.Time) ([]domain.Expense, error)
	CreateExpense(ctx context.Context, e *domain.Expense) error
	ListExpenses(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID, from, to *time.Time) ([]domain.Expense, error)
	// Projects
	ListProjects(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID) ([]domain.Project, error)
	GetProject(ctx context.Context, id, tenantID uuid.UUID) (*domain.Project, error)
	CreateProject(ctx context.Context, p *domain.Project) error
	UpdateProject(ctx context.Context, p *domain.Project) error
	ListTasks(ctx context.Context, tenantID, projectID uuid.UUID) ([]domain.ProjectTask, error)
	CreateTask(ctx context.Context, t *domain.ProjectTask) error
	UpdateTask(ctx context.Context, t *domain.ProjectTask) error
	// Billing and profitability
	CreateBillingDraft(ctx context.Context, d *domain.BillingDraft, items []domain.BillingItem, billedEntryIDs, billedExpenseIDs []uuid.UUID) error
	ListBillingDrafts(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID) ([]domain.BillingDraft, error)
	UpdateBillingDraftStatus(ctx context.Context, d *domain.BillingDraft) error
	Profitability(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter, projectFilter *uuid.UUID) (*domain.Profitability, error)
	// Inventory
	LinkAssetToClient(ctx context.Context, tenantID, clientID, assetID uuid.UUID) error
	UnlinkAssetFromClient(ctx context.Context, tenantID, clientID, assetID uuid.UUID) error
	ListClientAssets(ctx context.Context, tenantID, clientID uuid.UUID) ([]domain.ClientAssetView, error)
	AssetExists(ctx context.Context, id, tenantID uuid.UUID) (bool, error)
	CreateDiscoveredAsset(ctx context.Context, d *domain.DiscoveredAsset) error
	ListDiscoveredAssets(ctx context.Context, tenantID, clientID uuid.UUID) ([]domain.DiscoveredAsset, error)
	ImportDiscoveredAsset(ctx context.Context, tenantID, clientID, discoveredID, ciTypeID uuid.UUID, assetName string) (uuid.UUID, error)
	// Monitoring
	CreateMonitoringSource(ctx context.Context, s *domain.MonitoringSource) error
	ListMonitoringSources(ctx context.Context, tenantID uuid.UUID) ([]domain.MonitoringSource, error)
	FindSourceByTokenHash(ctx context.Context, tenantID uuid.UUID, tokenHash string) (*domain.MonitoringSource, error)
	RotateSourceToken(ctx context.Context, id, tenantID uuid.UUID, tokenHash string) error
	OpenAlertForResource(ctx context.Context, tenantID, sourceID uuid.UUID, resource string) (*domain.MonitoringAlert, error)
	CreateMonitoringAlert(ctx context.Context, a *domain.MonitoringAlert) error
	UpdateMonitoringAlert(ctx context.Context, a *domain.MonitoringAlert) error
	ListMonitoringAlerts(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID) ([]domain.MonitoringAlert, error)
	GetMonitoringAlert(ctx context.Context, id, tenantID uuid.UUID) (*domain.MonitoringAlert, error)
	// Remote sessions
	CreateRemoteSession(ctx context.Context, s *domain.RemoteSession) error
	GetRemoteSession(ctx context.Context, id, tenantID uuid.UUID) (*domain.RemoteSession, error)
	UpdateRemoteSession(ctx context.Context, s *domain.RemoteSession) error
	ListRemoteSessions(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID) ([]domain.RemoteSession, error)
	// Reports
	Report(ctx context.Context, tenantID uuid.UUID) (*domain.MSPReport, error)
}

// MSPTicket is a ticket row joined with its MSP context, as listed by
// the multi-client desk.
type MSPTicket struct {
	ID                   uuid.UUID
	Title                string
	Status               string
	Priority             string
	Category             string
	CreatedAt            time.Time
	ClientID             uuid.UUID
	ClientName           string
	ContractID           *uuid.UUID
	SLAResponseMinutes   int
	SLAResolutionMinutes int
}

// Service is the MSP application service.
type Service struct {
	repo Repository
}

// NewService creates a Service from a Repository implementation.
func NewService(repo Repository) *Service { return &Service{repo: repo} }

// visibleClients resolves the client scope for a viewer: staff see all
// (nil = no restriction); a client contact sees only their clients.
func (s *Service) visibleClients(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string) ([]uuid.UUID, error) {
	if viewerRole == "admin" || viewerRole == "agent" {
		return nil, nil
	}
	ids, err := s.repo.ContactClientIDs(ctx, tenantID, viewerID)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// checkClient enforces the isolation rule for one client. Staff pass;
// contacts must hold the client in their list, else NotFound.
func (s *Service) checkClient(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientID uuid.UUID) error {
	if viewerRole == "admin" || viewerRole == "agent" {
		return nil
	}
	ids, err := s.repo.ContactClientIDs(ctx, tenantID, viewerID)
	if err != nil {
		return err
	}
	if !domain.CanAccessClient(viewerRole, ids, clientID) {
		return domain.ErrNotFound
	}
	return nil
}

// ── Clients and contacts ────────────────────────────────────────────

func (s *Service) ListClients(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string) ([]domain.Client, error) {
	ids, err := s.visibleClients(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	return s.repo.ListClients(ctx, tenantID, ids)
}

func (s *Service) GetClient(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, id uuid.UUID) (*domain.Client, error) {
	if err := s.checkClient(ctx, tenantID, viewerID, viewerRole, id); err != nil {
		return nil, err
	}
	return s.repo.GetClient(ctx, id, tenantID)
}

func (s *Service) CreateClient(ctx context.Context, tenantID uuid.UUID, c *domain.Client) (*domain.Client, error) {
	c.TenantID = tenantID
	c.Name = strings.TrimSpace(c.Name)
	if c.Status == "" {
		c.Status = domain.ClientActive
	}
	if err := domain.ValidateClient(c); err != nil {
		return nil, err
	}
	if err := s.repo.CreateClient(ctx, c); err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

func (s *Service) UpdateClient(ctx context.Context, tenantID uuid.UUID, c *domain.Client) (*domain.Client, error) {
	c.TenantID = tenantID
	c.Name = strings.TrimSpace(c.Name)
	if err := domain.ValidateClient(c); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateClient(ctx, c); err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

func (s *Service) DeleteClient(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteClient(ctx, id, tenantID)
}

func (s *Service) ListContacts(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientID uuid.UUID) ([]domain.ClientContact, error) {
	if err := s.checkClient(ctx, tenantID, viewerID, viewerRole, clientID); err != nil {
		return nil, err
	}
	return s.repo.ListContacts(ctx, tenantID, clientID)
}

func (s *Service) CreateContact(ctx context.Context, tenantID, clientID uuid.UUID, c *domain.ClientContact) (*domain.ClientContact, error) {
	c.TenantID = tenantID
	c.ClientID = clientID
	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)
	if c.Role == "" {
		c.Role = "contact"
	}
	if err := domain.ValidateClientContact(c); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetClient(ctx, clientID, tenantID); err != nil {
		return nil, err
	}
	if c.UserID != nil {
		ok, err := s.repo.UserExists(ctx, *c.UserID, tenantID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, domain.ErrNotFound
		}
	}
	if err := s.repo.CreateContact(ctx, c); err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

func (s *Service) DeleteContact(ctx context.Context, tenantID, clientID, id uuid.UUID) error {
	return s.repo.DeleteContact(ctx, id, tenantID, clientID)
}

// ── Contracts ───────────────────────────────────────────────────────

func (s *Service) ListContracts(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientID *uuid.UUID) ([]domain.Contract, error) {
	if clientID != nil {
		if err := s.checkClient(ctx, tenantID, viewerID, viewerRole, *clientID); err != nil {
			return nil, err
		}
	} else if viewerRole != "admin" && viewerRole != "agent" {
		// Contacts without a client filter only see their own clients'
		// contracts: filter in the repo would need the list, so resolve
		// per client here and merge.
		ids, err := s.visibleClients(ctx, tenantID, viewerID, viewerRole)
		if err != nil {
			return nil, err
		}
		var out []domain.Contract
		for _, id := range ids {
			contracts, err := s.repo.ListContracts(ctx, tenantID, &id)
			if err != nil {
				return nil, err
			}
			out = append(out, contracts...)
		}
		return out, nil
	}
	return s.repo.ListContracts(ctx, tenantID, clientID)
}

func (s *Service) CreateContract(ctx context.Context, tenantID uuid.UUID, c *domain.Contract) (*domain.Contract, error) {
	c.TenantID = tenantID
	c.Name = strings.TrimSpace(c.Name)
	if c.Status == "" {
		c.Status = domain.ContractDraft
	}
	if c.Currency == "" {
		c.Currency = "EUR"
	}
	if err := domain.ValidateContract(c); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetClient(ctx, c.ClientID, tenantID); err != nil {
		return nil, err
	}
	if err := s.repo.CreateContract(ctx, c); err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

func (s *Service) UpdateContract(ctx context.Context, tenantID uuid.UUID, c *domain.Contract) (*domain.Contract, error) {
	c.TenantID = tenantID
	c.Name = strings.TrimSpace(c.Name)
	if err := domain.ValidateContract(c); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateContract(ctx, c); err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

// ResolveContractForTicket picks the contract a new ticket bills
// against: the explicitly given one (verified to belong to the client),
// else the client's active contract in its validity window, else none
// (the ticket is still created, but without SLA minutes).
func (s *Service) ResolveContractForTicket(ctx context.Context, tenantID, clientID uuid.UUID, contractID *uuid.UUID, priority string, at time.Time) (*domain.Contract, int, int, error) {
	if contractID != nil {
		c, err := s.repo.GetContract(ctx, *contractID, tenantID)
		if err != nil {
			return nil, 0, 0, err
		}
		if c.ClientID != clientID {
			return nil, 0, 0, domain.ErrNotFound
		}
		resp, res := c.ResolveSLA(priority)
		return c, resp, res, nil
	}
	c, err := s.repo.ActiveContractForClient(ctx, tenantID, clientID, at)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, 0, 0, nil
		}
		return nil, 0, 0, err
	}
	resp, res := c.ResolveSLA(priority)
	return c, resp, res, nil
}

// ── Multi-client tickets ────────────────────────────────────────────

type MSPTicketInput struct {
	ClientID    uuid.UUID
	ContractID  *uuid.UUID
	Title       string
	Description string
	Category    string
	Priority    string
}

func (s *Service) CreateMSPTicket(ctx context.Context, tenantID, creatorID uuid.UUID, creatorRole string, in MSPTicketInput) (uuid.UUID, error) {
	if strings.TrimSpace(in.Title) == "" || len(in.Title) > 255 || len(in.Description) > 10000 {
		return uuid.Nil, domain.ErrValidation
	}
	switch in.Priority {
	case "low", "medium", "high", "critical":
	default:
		return uuid.Nil, domain.ErrValidation
	}
	if err := s.checkClient(ctx, tenantID, creatorID, creatorRole, in.ClientID); err != nil {
		return uuid.Nil, err
	}
	if _, err := s.repo.GetClient(ctx, in.ClientID, tenantID); err != nil {
		return uuid.Nil, err
	}
	category := in.Category
	if category == "" {
		category = "general"
	}
	if len(category) > 100 {
		category = category[:100]
	}
	now := time.Now()
	contract, resp, res, err := s.ResolveContractForTicket(ctx, tenantID, in.ClientID, in.ContractID, in.Priority, now)
	if err != nil {
		return uuid.Nil, err
	}
	tc := &domain.TicketContext{TenantID: tenantID, ClientID: in.ClientID, SLAResponseMinutes: resp, SLAResolutionMinutes: res}
	if contract != nil {
		tc.ContractID = &contract.ID
	}
	dl := domainSLA.NewDeadlines(now, resp, res, 0)
	return s.repo.CreateTicketWithContext(ctx, tenantID, creatorID, strings.TrimSpace(in.Title), in.Description, category, in.Priority, tc, dl)
}

func (s *Service) ListMSPTickets(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientFilter *uuid.UUID) ([]MSPTicket, error) {
	ids, err := s.visibleClients(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	if clientFilter != nil && !domain.CanAccessClient(viewerRole, ids, *clientFilter) {
		return []MSPTicket{}, nil
	}
	return s.repo.ListMSPTickets(ctx, tenantID, ids, clientFilter)
}

func (s *Service) GetTicketContext(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, ticketID uuid.UUID) (*domain.TicketContext, error) {
	tc, err := s.repo.TicketContext(ctx, tenantID, ticketID)
	if err != nil {
		return nil, err
	}
	if err := s.checkClient(ctx, tenantID, viewerID, viewerRole, tc.ClientID); err != nil {
		return nil, err
	}
	return tc, nil
}

// ── Technicians and planning ────────────────────────────────────────

func (s *Service) ListTechnicians(ctx context.Context, tenantID uuid.UUID) ([]domain.Technician, error) {
	return s.repo.ListTechnicians(ctx, tenantID)
}

func (s *Service) CreateTechnician(ctx context.Context, tenantID uuid.UUID, t *domain.Technician) (*domain.Technician, error) {
	t.TenantID = tenantID
	if err := domain.ValidateTechnician(t); err != nil {
		return nil, err
	}
	ok, err := s.repo.UserExists(ctx, t.UserID, tenantID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrNotFound
	}
	if err := s.repo.CreateTechnician(ctx, t); err != nil {
		return nil, mapRepoErr(err)
	}
	return t, nil
}

func (s *Service) UpdateTechnician(ctx context.Context, tenantID uuid.UUID, t *domain.Technician) (*domain.Technician, error) {
	t.TenantID = tenantID
	if err := domain.ValidateTechnician(t); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateTechnician(ctx, t); err != nil {
		return nil, mapRepoErr(err)
	}
	return t, nil
}

func (s *Service) ListAvailability(ctx context.Context, tenantID, technicianID uuid.UUID) ([]domain.Availability, error) {
	if _, err := s.repo.GetTechnician(ctx, technicianID, tenantID); err != nil {
		return nil, err
	}
	return s.repo.ListAvailability(ctx, tenantID, technicianID)
}

func (s *Service) CreateAvailability(ctx context.Context, tenantID uuid.UUID, a *domain.Availability) (*domain.Availability, error) {
	a.TenantID = tenantID
	if err := domain.ValidateAvailability(a); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetTechnician(ctx, a.TechnicianID, tenantID); err != nil {
		return nil, err
	}
	if err := s.repo.CreateAvailability(ctx, a); err != nil {
		return nil, mapRepoErr(err)
	}
	return a, nil
}

func (s *Service) ListInterventions(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientFilter, technicianFilter *uuid.UUID) ([]domain.Intervention, error) {
	ids, err := s.visibleClients(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	if clientFilter != nil && !domain.CanAccessClient(viewerRole, ids, *clientFilter) {
		return []domain.Intervention{}, nil
	}
	return s.repo.ListInterventions(ctx, tenantID, ids, clientFilter, technicianFilter)
}

func (s *Service) CreateIntervention(ctx context.Context, tenantID uuid.UUID, v *domain.Intervention) (*domain.Intervention, error) {
	v.TenantID = tenantID
	if v.Status == "" {
		v.Status = "scheduled"
	}
	if err := domain.ValidateIntervention(v); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetClient(ctx, v.ClientID, tenantID); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetTechnician(ctx, v.TechnicianID, tenantID); err != nil {
		return nil, err
	}
	if v.TicketID != nil {
		tc, err := s.repo.TicketContext(ctx, tenantID, *v.TicketID)
		if err != nil {
			return nil, err
		}
		if tc.ClientID != v.ClientID {
			return nil, fmt.Errorf("%w: ticket belongs to another client", domain.ErrValidation)
		}
	}
	overlap, err := s.repo.TechnicianOverlappingIntervention(ctx, tenantID, v.TechnicianID, v.StartsAt, v.EndsAt)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, domain.ErrConflict
	}
	if err := s.repo.CreateIntervention(ctx, v); err != nil {
		return nil, mapRepoErr(err)
	}
	return v, nil
}

func (s *Service) UpdateInterventionStatus(ctx context.Context, tenantID, id uuid.UUID, next string) (*domain.Intervention, error) {
	// The repo loads the intervention for the tenant; the status machine
	// itself is the domain rule.
	list, err := s.repo.ListInterventions(ctx, tenantID, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	var found *domain.Intervention
	for i := range list {
		if list[i].ID == id {
			found = &list[i]
			break
		}
	}
	if found == nil {
		return nil, domain.ErrNotFound
	}
	if !found.CanMoveTo(next) {
		return nil, domain.ErrInvalidTransition
	}
	found.Status = next
	if err := s.repo.UpdateInterventionStatus(ctx, found); err != nil {
		return nil, err
	}
	return found, nil
}

// ── Time and expenses ───────────────────────────────────────────────

func (s *Service) CreateTimeEntry(ctx context.Context, tenantID uuid.UUID, e *domain.TimeEntry) (*domain.TimeEntry, error) {
	e.TenantID = tenantID
	tech, err := s.repo.GetTechnician(ctx, e.TechnicianID, tenantID)
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.GetClient(ctx, e.ClientID, tenantID); err != nil {
		return nil, err
	}
	if e.TicketID != nil {
		tc, err := s.repo.TicketContext(ctx, tenantID, *e.TicketID)
		if err != nil {
			return nil, err
		}
		if tc.ClientID != e.ClientID {
			return nil, fmt.Errorf("%w: ticket belongs to another client", domain.ErrValidation)
		}
	}
	e.HourlyRateCents = tech.HourlyRateCents
	e.CostCents = domain.TimeEntryCost(e.Minutes, tech.HourlyRateCents)
	if err := domain.ValidateTimeEntry(e); err != nil {
		return nil, err
	}
	if err := s.repo.CreateTimeEntry(ctx, e); err != nil {
		return nil, mapRepoErr(err)
	}
	return e, nil
}

func (s *Service) ListTimeEntries(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientFilter *uuid.UUID, from, to *time.Time) ([]domain.TimeEntry, error) {
	ids, err := s.visibleClients(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	if clientFilter != nil && !domain.CanAccessClient(viewerRole, ids, *clientFilter) {
		return []domain.TimeEntry{}, nil
	}
	return s.repo.ListTimeEntries(ctx, tenantID, ids, clientFilter, from, to)
}

func (s *Service) CreateExpense(ctx context.Context, tenantID uuid.UUID, e *domain.Expense) (*domain.Expense, error) {
	e.TenantID = tenantID
	if err := domain.ValidateExpense(e); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetClient(ctx, e.ClientID, tenantID); err != nil {
		return nil, err
	}
	if e.TicketID != nil {
		tc, err := s.repo.TicketContext(ctx, tenantID, *e.TicketID)
		if err != nil {
			return nil, err
		}
		if tc.ClientID != e.ClientID {
			return nil, fmt.Errorf("%w: ticket belongs to another client", domain.ErrValidation)
		}
	}
	if e.ProjectID != nil {
		p, err := s.repo.GetProject(ctx, *e.ProjectID, tenantID)
		if err != nil {
			return nil, err
		}
		if p.ClientID != e.ClientID {
			return nil, fmt.Errorf("%w: project belongs to another client", domain.ErrValidation)
		}
	}
	if err := s.repo.CreateExpense(ctx, e); err != nil {
		return nil, mapRepoErr(err)
	}
	return e, nil
}

func (s *Service) ListExpenses(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientFilter *uuid.UUID, from, to *time.Time) ([]domain.Expense, error) {
	ids, err := s.visibleClients(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	if clientFilter != nil && !domain.CanAccessClient(viewerRole, ids, *clientFilter) {
		return []domain.Expense{}, nil
	}
	return s.repo.ListExpenses(ctx, tenantID, ids, clientFilter, from, to)
}

// ── Projects ────────────────────────────────────────────────────────

func (s *Service) ListProjects(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientFilter *uuid.UUID) ([]domain.Project, error) {
	ids, err := s.visibleClients(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	if clientFilter != nil && !domain.CanAccessClient(viewerRole, ids, *clientFilter) {
		return []domain.Project{}, nil
	}
	return s.repo.ListProjects(ctx, tenantID, ids, clientFilter)
}

func (s *Service) CreateProject(ctx context.Context, tenantID uuid.UUID, p *domain.Project) (*domain.Project, error) {
	p.TenantID = tenantID
	p.Name = strings.TrimSpace(p.Name)
	if p.Status == "" {
		p.Status = "planned"
	}
	if err := domain.ValidateProject(p); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetClient(ctx, p.ClientID, tenantID); err != nil {
		return nil, err
	}
	if err := s.repo.CreateProject(ctx, p); err != nil {
		return nil, mapRepoErr(err)
	}
	return p, nil
}

func (s *Service) UpdateProjectStatus(ctx context.Context, tenantID, id uuid.UUID, next string) (*domain.Project, error) {
	p, err := s.repo.GetProject(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	if !p.CanMoveTo(next) {
		return nil, domain.ErrInvalidTransition
	}
	p.Status = next
	if err := s.repo.UpdateProject(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ListTasks(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, projectID uuid.UUID) ([]domain.ProjectTask, error) {
	p, err := s.repo.GetProject(ctx, projectID, tenantID)
	if err != nil {
		return nil, err
	}
	if err := s.checkClient(ctx, tenantID, viewerID, viewerRole, p.ClientID); err != nil {
		return nil, err
	}
	return s.repo.ListTasks(ctx, tenantID, projectID)
}

func (s *Service) CreateTask(ctx context.Context, tenantID uuid.UUID, t *domain.ProjectTask) (*domain.ProjectTask, error) {
	t.TenantID = tenantID
	t.Title = strings.TrimSpace(t.Title)
	if t.Status == "" {
		t.Status = "todo"
	}
	if err := domain.ValidateProjectTask(t); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetProject(ctx, t.ProjectID, tenantID); err != nil {
		return nil, err
	}
	if t.TechnicianID != nil {
		if _, err := s.repo.GetTechnician(ctx, *t.TechnicianID, tenantID); err != nil {
			return nil, err
		}
	}
	if err := s.repo.CreateTask(ctx, t); err != nil {
		return nil, mapRepoErr(err)
	}
	return t, nil
}

func (s *Service) UpdateTaskStatus(ctx context.Context, tenantID, projectID, taskID uuid.UUID, next string) (*domain.ProjectTask, error) {
	tasks, err := s.repo.ListTasks(ctx, tenantID, projectID)
	if err != nil {
		return nil, err
	}
	var found *domain.ProjectTask
	for i := range tasks {
		if tasks[i].ID == taskID {
			found = &tasks[i]
			break
		}
	}
	if found == nil {
		return nil, domain.ErrNotFound
	}
	if !found.CanMoveTo(next) {
		return nil, domain.ErrInvalidTransition
	}
	found.Status = next
	if err := s.repo.UpdateTask(ctx, found); err != nil {
		return nil, err
	}
	return found, nil
}

// ── Billing and profitability ───────────────────────────────────────
//
// GenerateDraft snapshots one period for one client: unbilled billable
// time valued at the contract rate (or the technician rate when the
// entry carries no contract rate), unbilled expenses of the period,
// and the monthly fee of the active contract, if any. It never charges
// anything; the draft stays in "draft" until staff issue it. One draft
// per client and period: the unique constraint rejects a second draft
// for the same period instead of billing it twice.

func (s *Service) GenerateDraft(ctx context.Context, tenantID, clientID uuid.UUID, from, to time.Time, notes string) (*domain.BillingDraft, error) {
	if to.Before(from) {
		return nil, domain.ErrValidation
	}
	if _, err := s.repo.GetClient(ctx, clientID, tenantID); err != nil {
		return nil, err
	}
	entries, err := s.repo.UnbilledTimeEntries(ctx, tenantID, clientID, from, to)
	if err != nil {
		return nil, err
	}
	expenses, err := s.repo.UnbilledExpenses(ctx, tenantID, clientID, from, to)
	if err != nil {
		return nil, err
	}
	contract, err := s.repo.ActiveContractForClient(ctx, tenantID, clientID, time.Now())
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	var items []domain.BillingItem
	var entryIDs []uuid.UUID
	for _, e := range entries {
		if !e.Billable {
			continue
		}
		rate := e.HourlyRateCents
		if e.ContractHourlyRateCents != nil {
			rate = *e.ContractHourlyRateCents
		}
		items = append(items, domain.BillingItem{
			TenantID: tenantID, Kind: "time",
			Description: fmt.Sprintf("%d min — %s", e.Minutes, e.Description),
			AmountCents: domain.TimeEntryCost(e.Minutes, rate),
		})
		entryIDs = append(entryIDs, e.ID)
	}
	var expenseIDs []uuid.UUID
	for _, e := range expenses {
		items = append(items, domain.BillingItem{
			TenantID: tenantID, Kind: "expense",
			Description: e.Description, AmountCents: e.AmountCents,
		})
		expenseIDs = append(expenseIDs, e.ID)
	}
	draft := &domain.BillingDraft{
		TenantID: tenantID, ClientID: clientID, PeriodStart: from, PeriodEnd: to,
		Status: "draft", TotalCents: domain.DraftTotal(items), Notes: notes,
	}
	if contract != nil {
		draft.ContractID = &contract.ID
		if contract.MonthlyFeeCents > 0 {
			items = append(items, domain.BillingItem{
				TenantID: tenantID, Kind: "fee",
				Description: fmt.Sprintf("Cuota mensual — %s", contract.Name),
				AmountCents: contract.MonthlyFeeCents,
			})
			draft.TotalCents = domain.DraftTotal(items)
		}
	}
	if err := s.repo.CreateBillingDraft(ctx, draft, items, entryIDs, expenseIDs); err != nil {
		return nil, mapRepoErr(err)
	}
	return draft, nil
}

func (s *Service) ListBillingDrafts(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientFilter *uuid.UUID) ([]domain.BillingDraft, error) {
	ids, err := s.visibleClients(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	if clientFilter != nil && !domain.CanAccessClient(viewerRole, ids, *clientFilter) {
		return []domain.BillingDraft{}, nil
	}
	return s.repo.ListBillingDrafts(ctx, tenantID, ids, clientFilter)
}

func (s *Service) UpdateBillingDraftStatus(ctx context.Context, tenantID, id uuid.UUID, next string) (*domain.BillingDraft, error) {
	drafts, err := s.repo.ListBillingDrafts(ctx, tenantID, nil, nil)
	if err != nil {
		return nil, err
	}
	var found *domain.BillingDraft
	for i := range drafts {
		if drafts[i].ID == id {
			found = &drafts[i]
			break
		}
	}
	if found == nil {
		return nil, domain.ErrNotFound
	}
	switch next {
	case "issued", "void":
		// A void draft must never be re-issued.
		if found.Status == "void" {
			return nil, domain.ErrInvalidTransition
		}
	default:
		return nil, domain.ErrValidation
	}
	found.Status = next
	if err := s.repo.UpdateBillingDraftStatus(ctx, found); err != nil {
		return nil, err
	}
	return found, nil
}

func (s *Service) Profitability(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientFilter, projectFilter *uuid.UUID) (*domain.Profitability, error) {
	ids, err := s.visibleClients(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	if clientFilter != nil && !domain.CanAccessClient(viewerRole, ids, *clientFilter) {
		return nil, domain.ErrNotFound
	}
	p, err := s.repo.Profitability(ctx, tenantID, ids, clientFilter, projectFilter)
	if err != nil {
		return nil, err
	}
	p.ComputeProfit()
	return p, nil
}

// ── Inventory ───────────────────────────────────────────────────────

func (s *Service) LinkAssetToClient(ctx context.Context, tenantID, clientID, assetID uuid.UUID) error {
	if _, err := s.repo.GetClient(ctx, clientID, tenantID); err != nil {
		return err
	}
	ok, err := s.repo.AssetExists(ctx, assetID, tenantID)
	if err != nil {
		return err
	}
	if !ok {
		return domain.ErrNotFound
	}
	return mapRepoErr(s.repo.LinkAssetToClient(ctx, tenantID, clientID, assetID))
}

func (s *Service) UnlinkAssetFromClient(ctx context.Context, tenantID, clientID, assetID uuid.UUID) error {
	return s.repo.UnlinkAssetFromClient(ctx, tenantID, clientID, assetID)
}

func (s *Service) ListClientAssets(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientID uuid.UUID) ([]domain.ClientAssetView, error) {
	if err := s.checkClient(ctx, tenantID, viewerID, viewerRole, clientID); err != nil {
		return nil, err
	}
	return s.repo.ListClientAssets(ctx, tenantID, clientID)
}

func (s *Service) CreateDiscoveredAsset(ctx context.Context, tenantID, clientID uuid.UUID, d *domain.DiscoveredAsset) (*domain.DiscoveredAsset, error) {
	d.TenantID = tenantID
	d.ClientID = clientID
	d.Name = strings.TrimSpace(d.Name)
	if d.Kind == "" {
		d.Kind = "host"
	}
	if d.Status == "" {
		d.Status = "pending"
	}
	if err := domain.ValidateDiscoveredAsset(d); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetClient(ctx, clientID, tenantID); err != nil {
		return nil, err
	}
	if err := s.repo.CreateDiscoveredAsset(ctx, d); err != nil {
		return nil, mapRepoErr(err)
	}
	return d, nil
}

func (s *Service) ListDiscoveredAssets(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientID uuid.UUID) ([]domain.DiscoveredAsset, error) {
	if err := s.checkClient(ctx, tenantID, viewerID, viewerRole, clientID); err != nil {
		return nil, err
	}
	return s.repo.ListDiscoveredAssets(ctx, tenantID, clientID)
}

func (s *Service) ImportDiscoveredAsset(ctx context.Context, tenantID, clientID, discoveredID, ciTypeID uuid.UUID, assetName string) (uuid.UUID, error) {
	if strings.TrimSpace(assetName) == "" {
		return uuid.Nil, domain.ErrValidation
	}
	if _, err := s.repo.GetClient(ctx, clientID, tenantID); err != nil {
		return uuid.Nil, err
	}
	return s.repo.ImportDiscoveredAsset(ctx, tenantID, clientID, discoveredID, ciTypeID, strings.TrimSpace(assetName))
}

// ── Monitoring ──────────────────────────────────────────────────────

func HashWebhookToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *Service) CreateMonitoringSource(ctx context.Context, tenantID uuid.UUID, src *domain.MonitoringSource, webhookToken string) (*domain.MonitoringSource, error) {
	src.TenantID = tenantID
	src.Name = strings.TrimSpace(src.Name)
	if src.Kind == "" {
		src.Kind = "zabbix"
	}
	if err := domain.ValidateMonitoringSource(src); err != nil {
		return nil, err
	}
	if webhookToken != "" {
		src.TokenHash = HashWebhookToken(webhookToken)
		src.Configured = true
	}
	if err := s.repo.CreateMonitoringSource(ctx, src); err != nil {
		return nil, mapRepoErr(err)
	}
	src.TokenHash = ""
	return src, nil
}

func (s *Service) ListMonitoringSources(ctx context.Context, tenantID uuid.UUID) ([]domain.MonitoringSource, error) {
	return s.repo.ListMonitoringSources(ctx, tenantID)
}

func (s *Service) RotateSourceToken(ctx context.Context, tenantID, id uuid.UUID, webhookToken string) error {
	if strings.TrimSpace(webhookToken) == "" {
		return domain.ErrValidation
	}
	return s.repo.RotateSourceToken(ctx, id, tenantID, HashWebhookToken(webhookToken))
}

// AlertIngest is the result of one Zabbix webhook delivery.
type AlertIngest struct {
	Alert         *domain.MonitoringAlert
	TicketCreated bool
}

// IngestAlert is idempotent: a repeated delivery for a resource with an
// open alert updates that alert instead of creating a second ticket.
// The webhook authenticates with the per-source token (hashed), never
// with a user session, so Zabbix can deliver without a browser login.
func (s *Service) IngestAlert(ctx context.Context, tenantID uuid.UUID, token string, in AlertInput) (*AlertIngest, error) {
	src, err := s.repo.FindSourceByTokenHash(ctx, tenantID, HashWebhookToken(token))
	if err != nil {
		return nil, err
	}
	if _, err := s.repo.GetClient(ctx, in.ClientID, tenantID); err != nil {
		return nil, err
	}
	alert := &domain.MonitoringAlert{
		TenantID: tenantID, SourceID: src.ID, ClientID: in.ClientID,
		ExternalResource: strings.TrimSpace(in.ExternalResource),
		Severity:         in.Severity, Status: in.Status,
		Title: in.Title, Message: in.Message,
	}
	if alert.Severity == "" {
		alert.Severity = "warning"
	}
	if alert.Status == "" {
		alert.Status = "open"
	}
	if err := domain.ValidateMonitoringAlert(alert); err != nil {
		return nil, err
	}
	existing, err := s.repo.OpenAlertForResource(ctx, tenantID, src.ID, alert.ExternalResource)
	if err != nil {
		return nil, err
	}
	openExists := existing != nil
	if openExists && alert.Status == "open" {
		existing.Severity = alert.Severity
		existing.Title = alert.Title
		existing.Message = alert.Message
		if err := s.repo.UpdateMonitoringAlert(ctx, existing); err != nil {
			return nil, err
		}
		return &AlertIngest{Alert: existing, TicketCreated: false}, nil
	}
	if openExists && alert.Status == "resolved" {
		existing.Status = "resolved"
		existing.Severity = alert.Severity
		existing.Message = alert.Message
		if err := s.repo.UpdateMonitoringAlert(ctx, existing); err != nil {
			return nil, err
		}
		return &AlertIngest{Alert: existing, TicketCreated: false}, nil
	}
	// No open alert: a fresh "open" alert creates (or reuses) its ticket.
	// Tickets require a user as creator, and a webhook carries none, so
	// the ticket is attributed to the tenant's first admin (the audit
	// trail still shows the alert itself as the origin). Without any
	// staff user the ingestion fails instead of inventing an author.
	if domain.ShouldCreateTicketForAlert(false, alert.Status) {
		now := time.Now()
		contract, resp, res, err := s.ResolveContractForTicket(ctx, tenantID, alert.ClientID, nil, "high", now)
		if err != nil {
			return nil, err
		}
		creatorID, err := s.repo.AlertTicketCreator(ctx, tenantID)
		if err != nil {
			return nil, err
		}
		title := alert.Title
		if title == "" {
			title = fmt.Sprintf("Alerta de monitorización: %s", alert.ExternalResource)
		}
		tc := &domain.TicketContext{TenantID: tenantID, ClientID: alert.ClientID, SLAResponseMinutes: resp, SLAResolutionMinutes: res}
		if contract != nil {
			tc.ContractID = &contract.ID
		}
		dl := domainSLA.NewDeadlines(now, resp, res, 0)
		ticketID, err := s.repo.CreateTicketWithContext(ctx, tenantID, creatorID, title, alert.Message, "monitoring", "high", tc, dl)
		if err != nil {
			return nil, err
		}
		alert.TicketID = &ticketID
	}
	if err := s.repo.CreateMonitoringAlert(ctx, alert); err != nil {
		return nil, mapRepoErr(err)
	}
	return &AlertIngest{Alert: alert, TicketCreated: alert.TicketID != nil}, nil
}

// AlertInput is one webhook delivery from Zabbix.
type AlertInput struct {
	ClientID         uuid.UUID
	ExternalResource string
	Severity         string
	Status           string
	Title            string
	Message          string
}

func (s *Service) ListMonitoringAlerts(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientFilter *uuid.UUID) ([]domain.MonitoringAlert, error) {
	ids, err := s.visibleClients(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	if clientFilter != nil && !domain.CanAccessClient(viewerRole, ids, *clientFilter) {
		return []domain.MonitoringAlert{}, nil
	}
	return s.repo.ListMonitoringAlerts(ctx, tenantID, ids, clientFilter)
}

func (s *Service) ResolveMonitoringAlert(ctx context.Context, tenantID, id uuid.UUID) (*domain.MonitoringAlert, error) {
	a, err := s.repo.GetMonitoringAlert(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	a.Status = "resolved"
	if err := s.repo.UpdateMonitoringAlert(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

// ── Remote sessions ─────────────────────────────────────────────────

func (s *Service) RequestRemoteSession(ctx context.Context, tenantID, callerID uuid.UUID, callerRole string, sess *domain.RemoteSession) (*domain.RemoteSession, error) {
	sess.TenantID = tenantID
	if sess.Provider == "" {
		sess.Provider = "meshcentral"
	}
	if sess.Status == "" {
		sess.Status = domain.SessionRequested
	}
	if err := domain.ValidateRemoteSession(sess); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetClient(ctx, sess.ClientID, tenantID); err != nil {
		return nil, err
	}
	if err := s.checkClient(ctx, tenantID, callerID, callerRole, sess.ClientID); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetTechnician(ctx, sess.TechnicianID, tenantID); err != nil {
		return nil, err
	}
	if sess.TicketID != nil {
		tc, err := s.repo.TicketContext(ctx, tenantID, *sess.TicketID)
		if err != nil {
			return nil, err
		}
		if tc.ClientID != sess.ClientID {
			return nil, fmt.Errorf("%w: ticket belongs to another client", domain.ErrValidation)
		}
	}
	if err := s.repo.CreateRemoteSession(ctx, sess); err != nil {
		return nil, mapRepoErr(err)
	}
	return sess, nil
}

func (s *Service) AuthorizeRemoteSession(ctx context.Context, tenantID, id, callerID uuid.UUID, callerRole string, authorize bool) (*domain.RemoteSession, error) {
	if !domain.CanAuthorizeSession(callerRole) {
		return nil, domain.ErrForbidden
	}
	sess, err := s.repo.GetRemoteSession(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	next := domain.SessionAuthorized
	if !authorize {
		next = domain.SessionDenied
	}
	if !sess.CanMoveTo(next) {
		return nil, domain.ErrInvalidTransition
	}
	sess.Status = next
	sess.AuthorizedBy = &callerID
	if err := s.repo.UpdateRemoteSession(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Service) StartRemoteSession(ctx context.Context, tenantID, id, callerID uuid.UUID, callerRole string) (*domain.RemoteSession, error) {
	sess, err := s.repo.GetRemoteSession(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	if !domain.CanStartSession(sess, callerID, callerRole) {
		return nil, domain.ErrForbidden
	}
	sess.Status = domain.SessionActive
	now := time.Now()
	sess.StartedAt = &now
	if err := s.repo.UpdateRemoteSession(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Service) EndRemoteSession(ctx context.Context, tenantID, id uuid.UUID) (*domain.RemoteSession, error) {
	sess, err := s.repo.GetRemoteSession(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	if !sess.CanMoveTo(domain.SessionEnded) {
		return nil, domain.ErrInvalidTransition
	}
	sess.Status = domain.SessionEnded
	now := time.Now()
	sess.EndedAt = &now
	if err := s.repo.UpdateRemoteSession(ctx, sess); err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Service) ListRemoteSessions(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, clientFilter *uuid.UUID) ([]domain.RemoteSession, error) {
	ids, err := s.visibleClients(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	if clientFilter != nil && !domain.CanAccessClient(viewerRole, ids, *clientFilter) {
		return []domain.RemoteSession{}, nil
	}
	return s.repo.ListRemoteSessions(ctx, tenantID, ids, clientFilter)
}

// ── Reports ─────────────────────────────────────────────────────────

func (s *Service) Report(ctx context.Context, tenantID uuid.UUID) (*domain.MSPReport, error) {
	return s.repo.Report(ctx, tenantID)
}

func mapRepoErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrConflict) || errors.Is(err, domain.ErrValidation) {
		return err
	}
	return err
}
