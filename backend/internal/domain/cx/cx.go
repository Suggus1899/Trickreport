// Package cx holds the pure domain model for Stage 5 (customer
// service: portal, email, live chat and WhatsApp; inbox, assignment,
// history, knowledge reuse and satisfaction). It has no database or
// HTTP dependency so its rules are unit-testable without PostgreSQL.
//
// Isolation rule, mirroring MSP: provider staff (admin/agent) see every
// customer of their tenant; an external customer (end_user) sees only
// the customers they are a contact of, and any other customer — or any
// conversation of another customer — is "not found", never a partial
// record. The server, never the UI, enforces assignment: taking a
// conversation already assigned to another agent fails with
// ErrConflict.
//
// No credentials are modelled here: email accounts store hosts and
// addresses only, WhatsApp accounts store the business/phone ids and
// only the SHA-256 hash of the webhook verify token. Real SMTP/IMAP
// passwords and the WhatsApp access token live in operator
// configuration (see docs/cx-integraciones.md), never in this model.
package cx

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrNotFound is returned when a record does not exist or does not
	// belong to the caller's tenant (or, for an external customer, to
	// one of the caller's customers).
	ErrNotFound = errors.New("cx record not found")
	// ErrValidation is returned when input fails domain validation.
	ErrValidation = errors.New("cx validation error")
	// ErrForbidden is returned when the caller may not perform the action.
	ErrForbidden = errors.New("cx forbidden")
	// ErrInvalidTransition is returned for a disallowed status change.
	ErrInvalidTransition = errors.New("cx invalid status transition")
	// ErrConflict is returned for duplicates and for taking a
	// conversation already assigned to another agent.
	ErrConflict = errors.New("cx conflict")
)

// ── Channels, statuses, priorities ──────────────────────────────────

const (
	ChannelEmail    = "email"
	ChannelPortal   = "portal"
	ChannelChat     = "chat"
	ChannelWhatsApp = "whatsapp"
)

func ValidChannel(c string) bool {
	switch c {
	case ChannelEmail, ChannelPortal, ChannelChat, ChannelWhatsApp:
		return true
	}
	return false
}

const (
	StatusOpen            = "open"
	StatusInProgress      = "in_progress"
	StatusWaitingCustomer = "waiting_customer"
	StatusResolved        = "resolved"
	StatusClosed          = "closed"
)

var ConversationTransitions = map[string][]string{
	StatusOpen:            {StatusInProgress, StatusWaitingCustomer, StatusResolved, StatusClosed},
	StatusInProgress:      {StatusWaitingCustomer, StatusResolved, StatusClosed},
	StatusWaitingCustomer: {StatusInProgress, StatusResolved, StatusClosed},
	StatusResolved:        {StatusInProgress, StatusClosed},
	StatusClosed:          {StatusInProgress},
}

func ValidConversationStatus(s string) bool {
	_, ok := ConversationTransitions[s]
	return ok
}

const (
	PriorityLow      = "low"
	PriorityMedium   = "medium"
	PriorityHigh     = "high"
	PriorityCritical = "critical"
)

func ValidPriority(p string) bool {
	switch p {
	case PriorityLow, PriorityMedium, PriorityHigh, PriorityCritical:
		return true
	}
	return false
}

// NextPriority raises one level; a critical conversation stays critical.
func NextPriority(p string) string {
	switch p {
	case PriorityLow:
		return PriorityMedium
	case PriorityMedium:
		return PriorityHigh
	case PriorityHigh:
		return PriorityCritical
	default:
		return PriorityCritical
	}
}

// ── Customers and contacts ──────────────────────────────────────────

const (
	CustomerActive   = "active"
	CustomerInactive = "inactive"
)

type Customer struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	Name             string
	Email            string
	Phone            string
	Company          string
	Status           string
	ConsentGiven     bool
	ErasureRequested bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func ValidateCustomer(c *Customer) error {
	if strings.TrimSpace(c.Name) == "" || len(c.Name) > 255 {
		return ErrValidation
	}
	if len(c.Email) > 255 || len(c.Phone) > 100 || len(c.Company) > 255 {
		return ErrValidation
	}
	if c.Email != "" && !strings.Contains(c.Email, "@") {
		return ErrValidation
	}
	if c.Status != "" && c.Status != CustomerActive && c.Status != CustomerInactive {
		return ErrValidation
	}
	return nil
}

type CustomerContact struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	CustomerID uuid.UUID
	UserID     *uuid.UUID
	Name       string
	Email      string
	Role       string // contact | client_admin
	CreatedAt  time.Time
}

func ValidateCustomerContact(c *CustomerContact) error {
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

// CanAccessCustomer is the isolation rule for external customers:
// staff see every customer of the tenant; an end_user sees only the
// customers they are a contact of. An empty list means no customer is
// visible, never "all customers".
func CanAccessCustomer(viewerRole string, viewerCustomerIDs []uuid.UUID, customerID uuid.UUID) bool {
	if viewerRole == "admin" || viewerRole == "agent" {
		return true
	}
	for _, id := range viewerCustomerIDs {
		if id == customerID {
			return true
		}
	}
	return false
}

// ── Conversations ───────────────────────────────────────────────────

type Conversation struct {
	ID                    uuid.UUID
	TenantID              uuid.UUID
	CustomerID            uuid.UUID
	CustomerName          string // derived, for display
	Channel               string
	Subject               string
	Status                string
	Priority              string
	AssignedTo            *uuid.UUID
	AssigneeName          string // derived, for display
	CreatedBy             *uuid.UUID
	ExternalThreadID      string
	SLAResponseMinutes    int
	SLAResolutionMinutes  int
	ResponseDueAt         *time.Time
	ResolutionDueAt       *time.Time
	FirstResponseAt       *time.Time
	Escalated             bool
	EscalatedAt           *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (c *Conversation) CanMoveTo(next string) bool {
	for _, allowed := range ConversationTransitions[c.Status] {
		if allowed == next {
			return true
		}
	}
	return false
}

func ValidateConversation(c *Conversation) error {
	if strings.TrimSpace(c.Subject) == "" || len(c.Subject) > 255 {
		return ErrValidation
	}
	if !ValidChannel(c.Channel) {
		return ErrValidation
	}
	if c.Priority != "" && !ValidPriority(c.Priority) {
		return ErrValidation
	}
	if c.Status != "" && !ValidConversationStatus(c.Status) {
		return ErrValidation
	}
	if c.SLAResponseMinutes < 0 || c.SLAResolutionMinutes < 0 {
		return ErrValidation
	}
	if len(c.ExternalThreadID) > 255 {
		return ErrValidation
	}
	return nil
}

// CanTakeConversation is the server-side assignment race rule: a
// conversation with no assignee, or already assigned to the caller,
// can be taken; one assigned to another agent cannot, and the second
// attempt must fail in the service (ErrConflict), not just be hidden
// in the UI.
func CanTakeConversation(c *Conversation, callerID uuid.UUID) bool {
	if c.AssignedTo == nil {
		return true
	}
	return *c.AssignedTo == callerID
}

// IsSLABreached reports a breach at the given moment: the resolution
// deadline passed while the conversation is not resolved/closed, or
// the first-response deadline passed with no agent reply yet.
func (c *Conversation) IsSLABreached(now time.Time) bool {
	if c.Status == StatusResolved || c.Status == StatusClosed {
		return false
	}
	if c.ResolutionDueAt != nil && now.After(*c.ResolutionDueAt) {
		return true
	}
	if c.FirstResponseAt == nil && c.ResponseDueAt != nil && now.After(*c.ResponseDueAt) {
		return true
	}
	return false
}

// DefaultSLA gives the fallback response/resolution minutes per
// channel: live chat answers fastest, email slowest. A configured
// policy (SLAPolicy) always wins over these defaults.
func DefaultSLA(channel string) (responseMinutes, resolutionMinutes int) {
	switch channel {
	case ChannelChat:
		return 15, 240
	case ChannelWhatsApp:
		return 60, 720
	case ChannelPortal:
		return 120, 1440
	default: // email
		return 240, 1440
	}
}

// ApplySLA fills the conversation's SLA minutes and deadlines from a
// policy when one matches, else from the channel defaults. Critical
// priority halves both deadlines (floor 15 minutes for response).
func ApplySLA(c *Conversation, policy *SLAPolicy, now time.Time) {
	resp, res := DefaultSLA(c.Channel)
	if policy != nil {
		resp, res = policy.ResponseMinutes, policy.ResolutionMinutes
	}
	if c.Priority == PriorityCritical {
		resp = max(resp/2, 15)
		res = max(res/2, 15)
	}
	c.SLAResponseMinutes, c.SLAResolutionMinutes = resp, res
	if resp > 0 {
		due := now.Add(time.Duration(resp) * time.Minute)
		c.ResponseDueAt = &due
	}
	if res > 0 {
		due := now.Add(time.Duration(res) * time.Minute)
		c.ResolutionDueAt = &due
	}
}

// ── Messages ────────────────────────────────────────────────────────

const (
	SenderCustomer = "customer"
	SenderAgent    = "agent"
	SenderSystem   = "system"
)

type Message struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	ConversationID    uuid.UUID
	SenderType        string // customer | agent | system
	SenderID          *uuid.UUID
	SenderName        string // derived, for display
	Body              string
	IsInternal        bool // staff-only note; never shown to the customer
	Channel           string
	ExternalMessageID string
	CreatedAt         time.Time
}

func ValidateMessage(m *Message) error {
	if strings.TrimSpace(m.Body) == "" || len(m.Body) > 20000 {
		return ErrValidation
	}
	switch m.SenderType {
	case SenderCustomer, SenderAgent, SenderSystem:
	default:
		return ErrValidation
	}
	if len(m.ExternalMessageID) > 255 {
		return ErrValidation
	}
	return nil
}

// MessageReopens reports the catalog rule CX-037: a customer reply to
// a resolved conversation reopens it for the agent.
func MessageReopens(convStatus, senderType string) bool {
	return convStatus == StatusResolved && senderType == SenderCustomer
}

// ── Channel configuration (no credentials stored) ───────────────────

type EmailAccount struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	Name         string
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPFrom     string
	IMAPHost     string
	IMAPPort     int
	IMAPUsername string
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func ValidateEmailAccount(a *EmailAccount) error {
	if strings.TrimSpace(a.Name) == "" || len(a.Name) > 255 {
		return ErrValidation
	}
	if len(a.SMTPHost) > 500 || len(a.IMAPHost) > 500 || len(a.SMTPUsername) > 255 ||
		len(a.IMAPUsername) > 255 || len(a.SMTPFrom) > 255 {
		return ErrValidation
	}
	if a.SMTPPort < 0 || a.SMTPPort > 65535 || a.IMAPPort < 0 || a.IMAPPort > 65535 {
		return ErrValidation
	}
	return nil
}

type WhatsAppAccount struct {
	ID                   uuid.UUID
	TenantID             uuid.UUID
	Name                 string
	BusinessID           string
	PhoneNumberID        string
	DisplayName          string
	WebhookVerifyTokenHash string // SHA-256 hex; never the token itself
	Active               bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func ValidateWhatsAppAccount(a *WhatsAppAccount) error {
	if strings.TrimSpace(a.Name) == "" || len(a.Name) > 255 {
		return ErrValidation
	}
	if len(a.BusinessID) > 255 || len(a.PhoneNumberID) > 255 || len(a.DisplayName) > 255 {
		return ErrValidation
	}
	return nil
}

// ── SLA policies, assignment rules, auto-replies ────────────────────

type SLAPolicy struct {
	ID                uuid.UUID
	TenantID          uuid.UUID
	Channel           string
	Priority          string
	ResponseMinutes   int
	ResolutionMinutes int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func ValidateSLAPolicy(p *SLAPolicy) error {
	if !ValidChannel(p.Channel) || !ValidPriority(p.Priority) {
		return ErrValidation
	}
	if p.ResponseMinutes < 0 || p.ResolutionMinutes < 0 {
		return ErrValidation
	}
	return nil
}

type AssignmentRule struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	Channel        string // "" = any channel
	Priority       string // "" = any priority
	AssigneeUserID uuid.UUID
	Position       int
	Active         bool
	CreatedAt      time.Time
}

func (r *AssignmentRule) Matches(channel, priority string) bool {
	if !r.Active {
		return false
	}
	if r.Channel != "" && r.Channel != channel {
		return false
	}
	if r.Priority != "" && r.Priority != priority {
		return false
	}
	return true
}

type AutoReply struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	TriggerKind  string // keyword | channel
	Keyword      string
	Channel      string // "" = any channel (keyword triggers)
	ResponseText string
	Active       bool
	CreatedAt    time.Time
}

func ValidateAutoReply(a *AutoReply) error {
	if a.TriggerKind != "keyword" && a.TriggerKind != "channel" {
		return ErrValidation
	}
	if a.TriggerKind == "keyword" && strings.TrimSpace(a.Keyword) == "" {
		return ErrValidation
	}
	if a.TriggerKind == "channel" && !ValidChannel(a.Channel) {
		return ErrValidation
	}
	if strings.TrimSpace(a.ResponseText) == "" || len(a.ResponseText) > 5000 {
		return ErrValidation
	}
	return nil
}

// MatchAutoReply returns the first active auto-reply whose trigger
// fires for this channel/body: a keyword trigger fires when the body
// contains the keyword (case-insensitive); a channel trigger fires
// for any message on that channel.
func MatchAutoReply(replies []AutoReply, channel, body string) *AutoReply {
	lower := strings.ToLower(body)
	for i := range replies {
		r := &replies[i]
		if !r.Active {
			continue
		}
		switch r.TriggerKind {
		case "keyword":
			if strings.Contains(lower, strings.ToLower(r.Keyword)) {
				return r
			}
		case "channel":
			if r.Channel == channel {
				return r
			}
		}
	}
	return nil
}

// ── Surveys and reports ─────────────────────────────────────────────

type Survey struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	ConversationID uuid.UUID
	CustomerID     uuid.UUID
	Rating         int // 1..5
	Comment        string
	CreatedAt      time.Time
}

func ValidateSurvey(s *Survey) error {
	if s.Rating < 1 || s.Rating > 5 {
		return ErrValidation
	}
	if len(s.Comment) > 5000 {
		return ErrValidation
	}
	return nil
}

type CXReport struct {
	CustomersTotal         int             `json:"customers_total"`
	OpenConversations      int             `json:"open_conversations"`
	ResolvedConversations  int             `json:"resolved_conversations"`
	SLABreached            int             `json:"sla_breached"`
	SurveyCount            int             `json:"survey_count"`
	AvgSatisfaction        float64         `json:"avg_satisfaction"`
	ConversationsByChannel map[string]int `json:"conversations_by_channel"`
}
