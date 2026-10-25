// Package msp holds the pure domain model for Stage 4 (MSP: multi-client
// operation). It has no database or HTTP dependency so its rules are
// unit-testable without PostgreSQL. Every repository query in the
// infrastructure layer scopes by tenant_id AND client_id where a client
// exists; these rules are the second, in-process check.
//
// Money is always integer cents. Time is always whole minutes.
package msp

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrNotFound is returned when a record does not exist or does not
	// belong to the caller's tenant (or, for a client contact, to one of
	// the caller's clients). Callers outside the scope see "not found",
	// never a hint that the record exists elsewhere.
	ErrNotFound = errors.New("msp record not found")
	// ErrValidation is returned when input fails domain validation.
	ErrValidation = errors.New("msp validation error")
	// ErrForbidden is returned when the caller may not perform the action.
	ErrForbidden = errors.New("msp forbidden")
	// ErrInvalidTransition is returned for a disallowed status change.
	ErrInvalidTransition = errors.New("msp invalid status transition")
	// ErrConflict is returned for duplicates and overlapping bookings.
	ErrConflict = errors.New("msp conflict")
	// ErrNoActiveContract is returned when a ticket needs a contract and
	// the client has none active in its validity window.
	ErrNoActiveContract = errors.New("msp no active contract")
)

// ── Clients ─────────────────────────────────────────────────────────

const (
	ClientActive    = "active"
	ClientInactive  = "inactive"
	ClientSuspended = "suspended"
)

func ValidClientStatus(s string) bool {
	switch s {
	case ClientActive, ClientInactive, ClientSuspended:
		return true
	}
	return false
}

type Client struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	Name         string
	Domain       string
	City         string
	Status       string
	ContactName  string
	ContactEmail string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func ValidateClient(c *Client) error {
	if strings.TrimSpace(c.Name) == "" || len(c.Name) > 255 {
		return ErrValidation
	}
	if len(c.Domain) > 255 || len(c.City) > 255 || len(c.ContactName) > 255 || len(c.ContactEmail) > 255 {
		return ErrValidation
	}
	if c.Status != "" && !ValidClientStatus(c.Status) {
		return ErrValidation
	}
	return nil
}

type ClientContact struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	ClientID  uuid.UUID
	UserID    *uuid.UUID
	Name      string
	Email     string
	Role      string // contact | client_admin
	CreatedAt time.Time
}

func ValidateClientContact(c *ClientContact) error {
	if strings.TrimSpace(c.Name) == "" || len(c.Name) > 255 {
		return ErrValidation
	}
	if strings.TrimSpace(c.Email) == "" || len(c.Email) > 255 || !strings.Contains(c.Email, "@") {
		return ErrValidation
	}
	if c.Role != "" && c.Role != "contact" && c.Role != "client_admin" {
		return ErrValidation
	}
	return nil
}

// CanAccessClient is the isolation rule for client contacts: provider
// staff (admin/agent) see every client of their tenant; a client
// contact (end_user) sees only the clients they are a contact of. The
// service passes the contact's client list; an empty list means no
// client is visible, never "all clients".
func CanAccessClient(viewerRole string, viewerClientIDs []uuid.UUID, clientID uuid.UUID) bool {
	if viewerRole == "admin" || viewerRole == "agent" {
		return true
	}
	for _, id := range viewerClientIDs {
		if id == clientID {
			return true
		}
	}
	return false
}

// ── Contracts and SLA ───────────────────────────────────────────────

const (
	ContractDraft     = "draft"
	ContractActive    = "active"
	ContractSuspended = "suspended"
	ContractExpired   = "expired"
	ContractCancelled = "cancelled"
)

func ValidContractStatus(s string) bool {
	switch s {
	case ContractDraft, ContractActive, ContractSuspended, ContractExpired, ContractCancelled:
		return true
	}
	return false
}

// SLAOverride holds per-priority response/resolution minutes. A zero
// value means "no override for this priority": the contract base
// minutes apply.
type SLAOverride struct {
	ResponseMinutes   int `json:"response_minutes"`
	ResolutionMinutes int `json:"resolution_minutes"`
}

type Contract struct {
	ID                   uuid.UUID
	TenantID             uuid.UUID
	ClientID             uuid.UUID
	Name                 string
	Status               string
	StartsAt             time.Time
	EndsAt               *time.Time
	Coverage             string
	Currency             string
	MonthlyFeeCents      int64
	HourlyRateCents      int64
	IncludedHoursMonthly float64
	SLAResponseMinutes   int
	SLAResolutionMinutes int
	SLAOverrides         map[string]SLAOverride // keyed by priority
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func ValidateContract(c *Contract) error {
	if strings.TrimSpace(c.Name) == "" || len(c.Name) > 255 {
		return ErrValidation
	}
	if len(c.Coverage) > 10000 {
		return ErrValidation
	}
	if c.Status != "" && !ValidContractStatus(c.Status) {
		return ErrValidation
	}
	if c.MonthlyFeeCents < 0 || c.HourlyRateCents < 0 || c.IncludedHoursMonthly < 0 {
		return ErrValidation
	}
	if c.SLAResponseMinutes < 0 || c.SLAResolutionMinutes < 0 {
		return ErrValidation
	}
	if c.EndsAt != nil && c.EndsAt.Before(c.StartsAt) {
		return ErrValidation
	}
	for priority, o := range c.SLAOverrides {
		switch priority {
		case "low", "medium", "high", "critical":
		default:
			return ErrValidation
		}
		if o.ResponseMinutes < 0 || o.ResolutionMinutes < 0 {
			return ErrValidation
		}
	}
	return nil
}

// IsActiveAt reports whether the contract covers the given moment:
// status active and inside its validity window.
func (c *Contract) IsActiveAt(at time.Time) bool {
	if c.Status != ContractActive {
		return false
	}
	if at.Before(c.StartsAt) {
		return false
	}
	if c.EndsAt != nil && at.After(*c.EndsAt) {
		return false
	}
	return true
}

// ResolveSLA returns the SLA minutes a ticket of the given priority
// gets under this contract: the per-priority override when set, else
// the contract base minutes.
func (c *Contract) ResolveSLA(priority string) (responseMinutes, resolutionMinutes int) {
	if o, ok := c.SLAOverrides[priority]; ok {
		responseMinutes, resolutionMinutes = c.SLAResponseMinutes, c.SLAResolutionMinutes
		if o.ResponseMinutes > 0 {
			responseMinutes = o.ResponseMinutes
		}
		if o.ResolutionMinutes > 0 {
			resolutionMinutes = o.ResolutionMinutes
		}
		return responseMinutes, resolutionMinutes
	}
	return c.SLAResponseMinutes, c.SLAResolutionMinutes
}

// ── Multi-client tickets ────────────────────────────────────────────

type TicketContext struct {
	TicketID             uuid.UUID
	TenantID             uuid.UUID
	ClientID             uuid.UUID
	ClientName           string // derived, for display
	ContractID           *uuid.UUID
	SLAResponseMinutes   int
	SLAResolutionMinutes int
	CreatedAt            time.Time
}

// ── Technicians and planning ────────────────────────────────────────

type Technician struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	UserID          uuid.UUID
	UserName        string // derived, for display
	Skills          []string
	HourlyRateCents int64
	Active          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func ValidateTechnician(t *Technician) error {
	if t.UserID == uuid.Nil {
		return ErrValidation
	}
	if t.HourlyRateCents < 0 {
		return ErrValidation
	}
	for _, s := range t.Skills {
		if len(s) > 100 {
			return ErrValidation
		}
	}
	return nil
}

type Availability struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	TechnicianID uuid.UUID
	DayOfWeek    int    // 0 = Sunday … 6 = Saturday
	StartTime    string // "15:04"
	EndTime      string // "15:04"
}

func ValidateAvailability(a *Availability) error {
	if a.DayOfWeek < 0 || a.DayOfWeek > 6 {
		return ErrValidation
	}
	start, err1 := time.Parse("15:04", a.StartTime)
	end, err2 := time.Parse("15:04", a.EndTime)
	if err1 != nil || err2 != nil || !end.After(start) {
		return ErrValidation
	}
	return nil
}

var InterventionTransitions = map[string][]string{
	"scheduled":   {"in_progress", "cancelled"},
	"in_progress": {"done", "cancelled"},
	"done":        {},
	"cancelled":   {},
}

type Intervention struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	ClientID     uuid.UUID
	TechnicianID uuid.UUID
	TicketID     *uuid.UUID
	Title        string
	StartsAt     time.Time
	EndsAt       time.Time
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (v *Intervention) CanMoveTo(next string) bool {
	for _, allowed := range InterventionTransitions[v.Status] {
		if allowed == next {
			return true
		}
	}
	return false
}

func ValidateIntervention(v *Intervention) error {
	if strings.TrimSpace(v.Title) == "" || len(v.Title) > 255 {
		return ErrValidation
	}
	if !v.EndsAt.After(v.StartsAt) {
		return ErrValidation
	}
	if _, ok := InterventionTransitions[v.Status]; v.Status != "" && !ok {
		return ErrValidation
	}
	return nil
}

// IntervalsOverlap reports whether two half-open [start, end) windows
// overlap. Used to reject double-bookings of the same technician.
func IntervalsOverlap(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}

// ── Time, expenses and costs ────────────────────────────────────────

// TimeEntryCost values minutes at an hourly rate in cents, rounding
// half up to the nearest cent: 30 minutes at 6000 cents/hour = 3000.
func TimeEntryCost(minutes int, hourlyRateCents int64) int64 {
	if minutes <= 0 || hourlyRateCents <= 0 {
		return 0
	}
	return (int64(minutes)*hourlyRateCents + 30) / 60
}

type TimeEntry struct {
	ID                      uuid.UUID
	TenantID                uuid.UUID
	ClientID                uuid.UUID
	TicketID                *uuid.UUID
	InterventionID          *uuid.UUID
	TechnicianID            uuid.UUID
	Minutes                 int
	Description             string
	Billable                bool
	Billed                  bool
	HourlyRateCents         int64
	ContractHourlyRateCents *int64
	CostCents               int64
	OccurredAt              time.Time
	CreatedAt               time.Time
}

func ValidateTimeEntry(e *TimeEntry) error {
	if e.Minutes <= 0 || e.Minutes > 1440 {
		return ErrValidation
	}
	if len(e.Description) > 10000 {
		return ErrValidation
	}
	if e.HourlyRateCents < 0 {
		return ErrValidation
	}
	return nil
}

type Expense struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	ClientID    uuid.UUID
	TicketID    *uuid.UUID
	ProjectID   *uuid.UUID
	Description string
	Category    string
	AmountCents int64
	Billed      bool
	IncurredAt  time.Time
	CreatedAt   time.Time
}

func ValidateExpense(e *Expense) error {
	if e.AmountCents <= 0 {
		return ErrValidation
	}
	if len(e.Description) > 10000 || len(e.Category) > 100 {
		return ErrValidation
	}
	return nil
}

// ── Projects ────────────────────────────────────────────────────────

var ProjectTransitions = map[string][]string{
	"planned":   {"active", "cancelled"},
	"active":    {"on_hold", "done", "cancelled"},
	"on_hold":   {"active", "cancelled"},
	"done":      {},
	"cancelled": {},
}

var TaskTransitions = map[string][]string{
	"todo":  {"doing"},
	"doing": {"todo", "done"},
	"done":  {"doing"},
}

type Project struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	ClientID    uuid.UUID
	Name        string
	Status      string
	BudgetCents int64
	StartsAt    *time.Time
	EndsAt      *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (p *Project) CanMoveTo(next string) bool {
	for _, allowed := range ProjectTransitions[p.Status] {
		if allowed == next {
			return true
		}
	}
	return false
}

func ValidateProject(p *Project) error {
	if strings.TrimSpace(p.Name) == "" || len(p.Name) > 255 {
		return ErrValidation
	}
	if p.BudgetCents < 0 {
		return ErrValidation
	}
	if _, ok := ProjectTransitions[p.Status]; p.Status != "" && !ok {
		return ErrValidation
	}
	if p.StartsAt != nil && p.EndsAt != nil && p.EndsAt.Before(*p.StartsAt) {
		return ErrValidation
	}
	return nil
}

type ProjectTask struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	ProjectID        uuid.UUID
	Title            string
	Status           string
	TechnicianID     *uuid.UUID
	EstimatedMinutes int
	Position         int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (t *ProjectTask) CanMoveTo(next string) bool {
	for _, allowed := range TaskTransitions[t.Status] {
		if allowed == next {
			return true
		}
	}
	return false
}

func ValidateProjectTask(t *ProjectTask) error {
	if strings.TrimSpace(t.Title) == "" || len(t.Title) > 255 {
		return ErrValidation
	}
	if t.EstimatedMinutes < 0 {
		return ErrValidation
	}
	if _, ok := TaskTransitions[t.Status]; t.Status != "" && !ok {
		return ErrValidation
	}
	return nil
}

// ProjectProgress is the share of finished tasks, as a whole percent.
// A project without tasks is 0% done, never undefined.
func ProjectProgress(tasks []ProjectTask) int {
	if len(tasks) == 0 {
		return 0
	}
	done := 0
	for _, t := range tasks {
		if t.Status == "done" {
			done++
		}
	}
	return done * 100 / len(tasks)
}

// ── Profitability and billing drafts ────────────────────────────────
//
// A draft is a snapshot, never an automatic charge: generating it only
// records what would be billed. Included billable time entries are
// marked billed so a later draft cannot bill the same minutes twice.

type BillingItem struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	DraftID     uuid.UUID
	Kind        string // time | expense | fee
	Description string
	AmountCents int64
	CreatedAt   time.Time
}

func ValidBillingKind(kind string) bool {
	switch kind {
	case "time", "expense", "fee":
		return true
	}
	return false
}

type BillingDraft struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	ClientID    uuid.UUID
	ContractID  *uuid.UUID
	PeriodStart time.Time
	PeriodEnd   time.Time
	Status      string // draft | issued | void
	TotalCents  int64
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func DraftTotal(items []BillingItem) int64 {
	var total int64
	for _, it := range items {
		total += it.AmountCents
	}
	return total
}

// Profitability summarizes one client (or one project, when the repo
// filters by it): revenue values billable time at the contract rate
// when the entry carries one, else at the technician rate; cost values
// the same minutes at the technician rate.
type Profitability struct {
	ClientID     uuid.UUID `json:"client_id"`
	Minutes      int       `json:"minutes"`
	TimeRevenue  int64     `json:"time_revenue_cents"`
	TimeCost     int64     `json:"time_cost_cents"`
	Expenses     int64     `json:"expenses_cents"`
	BilledDrafts int64     `json:"billed_drafts_cents"`
	Profit       int64     `json:"profit_cents"`
}

func (p *Profitability) ComputeProfit() {
	p.Profit = p.TimeRevenue + p.BilledDrafts - p.TimeCost - p.Expenses
}

// ── Inventory and discovery ─────────────────────────────────────────

type ClientAsset struct {
	TenantID  uuid.UUID
	ClientID  uuid.UUID
	AssetID   uuid.UUID
	CreatedAt time.Time
}

// ClientAssetView is an Etapa 3 asset as seen from one MSP client.
type ClientAssetView struct {
	AssetID    uuid.UUID
	ClientID   uuid.UUID
	Name       string
	CITypeName string
	Status     string
	Serial     string
}

type DiscoveredAsset struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	ClientID        uuid.UUID
	ExternalID      string
	Name            string
	Kind            string
	Payload         map[string]any
	Status          string // pending | imported | ignored
	ImportedAssetID *uuid.UUID
	DiscoveredAt    time.Time
}

func ValidateDiscoveredAsset(d *DiscoveredAsset) error {
	if strings.TrimSpace(d.ExternalID) == "" || len(d.ExternalID) > 255 {
		return ErrValidation
	}
	if strings.TrimSpace(d.Name) == "" || len(d.Name) > 255 {
		return ErrValidation
	}
	if len(d.Kind) > 100 {
		return ErrValidation
	}
	return nil
}

// ── Monitoring (Zabbix) ─────────────────────────────────────────────

type MonitoringSource struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	Name       string
	Kind       string // zabbix
	BaseURL    string
	TokenHash  string // SHA-256 hex of the webhook token; never the token
	Configured bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func ValidateMonitoringSource(s *MonitoringSource) error {
	if strings.TrimSpace(s.Name) == "" || len(s.Name) > 255 {
		return ErrValidation
	}
	if s.Kind != "" && s.Kind != "zabbix" {
		return ErrValidation
	}
	if len(s.BaseURL) > 500 {
		return ErrValidation
	}
	return nil
}

func ValidAlertSeverity(s string) bool {
	switch s {
	case "info", "warning", "average", "high", "disaster":
		return true
	}
	return false
}

type MonitoringAlert struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	SourceID         uuid.UUID
	ClientID         uuid.UUID
	ExternalResource string
	Severity         string
	Status           string // open | resolved
	Title            string
	Message          string
	TicketID         *uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func ValidateMonitoringAlert(a *MonitoringAlert) error {
	if strings.TrimSpace(a.ExternalResource) == "" || len(a.ExternalResource) > 255 {
		return ErrValidation
	}
	if !ValidAlertSeverity(a.Severity) {
		return ErrValidation
	}
	if len(a.Title) > 255 || len(a.Message) > 10000 {
		return ErrValidation
	}
	return nil
}

// ShouldCreateTicketForAlert is the de-duplication rule: a new ticket
// is created only when no open alert for the same resource exists.
// The partial unique index in migration 000015 enforces the same rule
// in the database, so a race still cannot duplicate the ticket.
func ShouldCreateTicketForAlert(openAlertExists bool, status string) bool {
	return !openAlertExists && status == "open"
}

// ── Remote access (MeshCentral) ─────────────────────────────────────

const (
	SessionRequested  = "requested"
	SessionAuthorized = "authorized"
	SessionActive     = "active"
	SessionEnded      = "ended"
	SessionDenied     = "denied"
)

var RemoteSessionTransitions = map[string][]string{
	SessionRequested:  {SessionAuthorized, SessionDenied},
	SessionAuthorized: {SessionActive, SessionDenied},
	SessionActive:     {SessionEnded},
	SessionEnded:      {},
	SessionDenied:     {},
}

type RemoteSession struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	ClientID          uuid.UUID
	TicketID          *uuid.UUID
	TechnicianID      uuid.UUID
	TechnicianUserID  uuid.UUID // derived, for the authorization check
	Provider          string    // meshcentral
	ExternalSessionID string
	Status            string
	AuthorizedBy      *uuid.UUID
	AuditNote         string
	StartedAt         *time.Time
	EndedAt           *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (s *RemoteSession) CanMoveTo(next string) bool {
	for _, allowed := range RemoteSessionTransitions[s.Status] {
		if allowed == next {
			return true
		}
	}
	return false
}

// CanAuthorizeSession: only provider staff authorize a remote session.
// The technician themself never self-authorizes.
func CanAuthorizeSession(callerRole string) bool {
	return callerRole == "admin" || callerRole == "agent"
}

// CanStartSession: the owning technician (after authorization) or
// staff. Starting an unauthorized session is always denied.
func CanStartSession(s *RemoteSession, callerID uuid.UUID, callerRole string) bool {
	if s.Status != SessionAuthorized {
		return false
	}
	if callerRole == "admin" || callerRole == "agent" {
		return true
	}
	return s.TechnicianUserID == callerID
}

func ValidateRemoteSession(s *RemoteSession) error {
	if s.Provider != "" && s.Provider != "meshcentral" {
		return ErrValidation
	}
	if len(s.ExternalSessionID) > 255 || len(s.AuditNote) > 10000 {
		return ErrValidation
	}
	if _, ok := RemoteSessionTransitions[s.Status]; s.Status != "" && !ok {
		return ErrValidation
	}
	return nil
}

// ── Reports ─────────────────────────────────────────────────────────

type MSPReport struct {
	ClientsTotal    int `json:"clients_total"`
	ActiveClients   int `json:"active_clients"`
	ActiveContracts int `json:"active_contracts"`
	OpenAlerts      int `json:"open_alerts"`
	ActiveSessions  int `json:"active_sessions"`
	OpenMSPTickets  int `json:"open_msp_tickets"`
}
