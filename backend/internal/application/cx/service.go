// Package cx implements the Stage 5 customer-service application
// logic: external customers and contacts, the inbox across portal,
// email, live chat and WhatsApp, assignment, SLA, escalations,
// channel configuration, assignment rules, auto-replies, surveys and
// reports.
//
// Authorization is enforced here, in Go, not in the frontend: route
// middleware only narrows write endpoints to staff, while the
// per-customer isolation rule (domain.CanAccessCustomer) and the
// assignment race rule (domain.CanTakeConversation) are checked per
// record in this service. An external customer that is not linked to
// a customer sees that customer — and its conversations — as "not
// found" in every method below. Internal notes (is_internal) are
// filtered out for external customers in ListMessages, so a note
// written by staff can never leak through a list call.
package cx

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	domain "github.com/trickreport/backend/internal/domain/cx"
)

// Repository is the persistence boundary for CX aggregates. Every
// method scopes by tenantID, and every customer-owned method
// additionally by customerID; implementations must return
// ErrNotFound for records outside the tenant (or outside the
// caller's customers).
type Repository interface {
	// Customers and contacts
	ContactCustomerIDs(ctx context.Context, tenantID, userID uuid.UUID) ([]uuid.UUID, error)
	ListCustomers(ctx context.Context, tenantID uuid.UUID, customerIDs []uuid.UUID) ([]domain.Customer, error)
	GetCustomer(ctx context.Context, id, tenantID uuid.UUID) (*domain.Customer, error)
	GetCustomerByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*domain.Customer, error)
	GetCustomerByPhone(ctx context.Context, tenantID uuid.UUID, phone string) (*domain.Customer, error)
	CreateCustomer(ctx context.Context, c *domain.Customer) error
	UpdateCustomer(ctx context.Context, c *domain.Customer) error
	DeleteCustomer(ctx context.Context, id, tenantID uuid.UUID) error
	ListContacts(ctx context.Context, tenantID, customerID uuid.UUID) ([]domain.CustomerContact, error)
	CreateContact(ctx context.Context, c *domain.CustomerContact) error
	DeleteContact(ctx context.Context, id, tenantID, customerID uuid.UUID) error
	UserExists(ctx context.Context, id, tenantID uuid.UUID) (bool, error)
	// Conversations and messages
	ListConversations(ctx context.Context, tenantID uuid.UUID, customerIDs []uuid.UUID, customerFilter *uuid.UUID, channel, status string, assigneeFilter *uuid.UUID) ([]domain.Conversation, error)
	GetConversation(ctx context.Context, id, tenantID uuid.UUID) (*domain.Conversation, error)
	OpenConversationForThread(ctx context.Context, tenantID, customerID uuid.UUID, threadID string) (*domain.Conversation, error)
	CreateConversation(ctx context.Context, c *domain.Conversation) error
	UpdateConversation(ctx context.Context, c *domain.Conversation) error
	AssignConversation(ctx context.Context, id, tenantID uuid.UUID, assignee *uuid.UUID, status string) error
	ListMessages(ctx context.Context, tenantID, conversationID uuid.UUID, includeInternal bool) ([]domain.Message, error)
	CreateMessage(ctx context.Context, m *domain.Message) error
	MessageExistsByExternalID(ctx context.Context, tenantID uuid.UUID, externalID string) (bool, error)
	// Channel configuration
	ListEmailAccounts(ctx context.Context, tenantID uuid.UUID) ([]domain.EmailAccount, error)
	CreateEmailAccount(ctx context.Context, a *domain.EmailAccount) error
	UpdateEmailAccount(ctx context.Context, a *domain.EmailAccount) error
	GetWhatsAppAccount(ctx context.Context, id, tenantID uuid.UUID) (*domain.WhatsAppAccount, error)
	ListWhatsAppAccounts(ctx context.Context, tenantID uuid.UUID) ([]domain.WhatsAppAccount, error)
	FindWhatsAppAccountByTokenHash(ctx context.Context, tenantID uuid.UUID, tokenHash string) (*domain.WhatsAppAccount, error)
	CreateWhatsAppAccount(ctx context.Context, a *domain.WhatsAppAccount) error
	UpdateWhatsAppAccount(ctx context.Context, a *domain.WhatsAppAccount) error
	// SLA, rules, auto-replies
	ListSLAPolicies(ctx context.Context, tenantID uuid.UUID) ([]domain.SLAPolicy, error)
	UpsertSLAPolicy(ctx context.Context, p *domain.SLAPolicy) error
	ListAssignmentRules(ctx context.Context, tenantID uuid.UUID) ([]domain.AssignmentRule, error)
	CreateAssignmentRule(ctx context.Context, r *domain.AssignmentRule) error
	DeleteAssignmentRule(ctx context.Context, id, tenantID uuid.UUID) error
	ListAutoReplies(ctx context.Context, tenantID uuid.UUID) ([]domain.AutoReply, error)
	CreateAutoReply(ctx context.Context, a *domain.AutoReply) error
	DeleteAutoReply(ctx context.Context, id, tenantID uuid.UUID) error
	// Surveys and reports
	CreateSurvey(ctx context.Context, s *domain.Survey) error
	ListSurveys(ctx context.Context, tenantID uuid.UUID, customerIDs []uuid.UUID, customerFilter *uuid.UUID) ([]domain.Survey, error)
	Report(ctx context.Context, tenantID uuid.UUID, customerIDs []uuid.UUID) (*domain.CXReport, error)
}

// Service is the CX application service.
type Service struct {
	repo Repository
}

// NewService creates a Service from a Repository implementation.
func NewService(repo Repository) *Service { return &Service{repo: repo} }

func mapRepoErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, domain.ErrNotFound) || errors.Is(err, domain.ErrConflict) || errors.Is(err, domain.ErrValidation) {
		return err
	}
	return err
}

// visibleCustomers resolves the customer scope for a viewer: staff see
// all (nil = no restriction); an external customer sees only their
// customers.
func (s *Service) visibleCustomers(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string) ([]uuid.UUID, error) {
	if viewerRole == "admin" || viewerRole == "agent" {
		return nil, nil
	}
	ids, err := s.repo.ContactCustomerIDs(ctx, tenantID, viewerID)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// checkCustomer enforces the isolation rule for one customer. Staff
// pass; contacts must hold the customer in their list, else NotFound.
func (s *Service) checkCustomer(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, customerID uuid.UUID) error {
	if viewerRole == "admin" || viewerRole == "agent" {
		return nil
	}
	ids, err := s.repo.ContactCustomerIDs(ctx, tenantID, viewerID)
	if err != nil {
		return err
	}
	if !domain.CanAccessCustomer(viewerRole, ids, customerID) {
		return domain.ErrNotFound
	}
	return nil
}

// checkConversation loads a conversation and enforces isolation on its
// customer, so a conversation of another customer is NotFound.
func (s *Service) checkConversation(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, id uuid.UUID) (*domain.Conversation, error) {
	c, err := s.repo.GetConversation(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	if err := s.checkCustomer(ctx, tenantID, viewerID, viewerRole, c.CustomerID); err != nil {
		return nil, err
	}
	return c, nil
}

// ── Customers and contacts ──────────────────────────────────────────

func (s *Service) ListCustomers(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string) ([]domain.Customer, error) {
	ids, err := s.visibleCustomers(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	return s.repo.ListCustomers(ctx, tenantID, ids)
}

func (s *Service) GetCustomer(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, id uuid.UUID) (*domain.Customer, error) {
	if err := s.checkCustomer(ctx, tenantID, viewerID, viewerRole, id); err != nil {
		return nil, err
	}
	return s.repo.GetCustomer(ctx, id, tenantID)
}

func (s *Service) CreateCustomer(ctx context.Context, tenantID uuid.UUID, c *domain.Customer) (*domain.Customer, error) {
	c.TenantID = tenantID
	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)
	if c.Status == "" {
		c.Status = domain.CustomerActive
	}
	if err := domain.ValidateCustomer(c); err != nil {
		return nil, err
	}
	if err := s.repo.CreateCustomer(ctx, c); err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

func (s *Service) UpdateCustomer(ctx context.Context, tenantID uuid.UUID, c *domain.Customer) (*domain.Customer, error) {
	c.TenantID = tenantID
	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)
	if err := domain.ValidateCustomer(c); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateCustomer(ctx, c); err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

func (s *Service) DeleteCustomer(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteCustomer(ctx, id, tenantID)
}

func (s *Service) ListContacts(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, customerID uuid.UUID) ([]domain.CustomerContact, error) {
	if err := s.checkCustomer(ctx, tenantID, viewerID, viewerRole, customerID); err != nil {
		return nil, err
	}
	return s.repo.ListContacts(ctx, tenantID, customerID)
}

func (s *Service) CreateContact(ctx context.Context, tenantID, customerID uuid.UUID, c *domain.CustomerContact) (*domain.CustomerContact, error) {
	c.TenantID = tenantID
	c.CustomerID = customerID
	c.Name = strings.TrimSpace(c.Name)
	c.Email = strings.TrimSpace(c.Email)
	if c.Role == "" {
		c.Role = "contact"
	}
	if err := domain.ValidateCustomerContact(c); err != nil {
		return nil, err
	}
	if c.UserID != nil {
		ok, err := s.repo.UserExists(ctx, *c.UserID, tenantID)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, domain.ErrValidation
		}
	}
	if err := s.repo.CreateContact(ctx, c); err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

func (s *Service) DeleteContact(ctx context.Context, tenantID, customerID, id uuid.UUID) error {
	return s.repo.DeleteContact(ctx, id, tenantID, customerID)
}

// ── Conversations ───────────────────────────────────────────────────

// CreateConversation resolves the SLA (policy wins, else channel
// defaults), applies the first matching auto-assignment rule, and
// stores the conversation. An external customer may only open a
// conversation for their own customer.
func (s *Service) CreateConversation(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, c *domain.Conversation) (*domain.Conversation, error) {
	if err := s.checkCustomer(ctx, tenantID, viewerID, viewerRole, c.CustomerID); err != nil {
		return nil, err
	}
	c.TenantID = tenantID
	c.Subject = strings.TrimSpace(c.Subject)
	if c.Status == "" {
		c.Status = domain.StatusOpen
	}
	if c.Priority == "" {
		c.Priority = domain.PriorityMedium
	}
	if err := domain.ValidateConversation(c); err != nil {
		return nil, err
	}
	policies, err := s.repo.ListSLAPolicies(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	var policy *domain.SLAPolicy
	for i := range policies {
		if policies[i].Channel == c.Channel && policies[i].Priority == c.Priority {
			policy = &policies[i]
			break
		}
	}
	domain.ApplySLA(c, policy, time.Now())
	rules, err := s.repo.ListAssignmentRules(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	for i := range rules {
		if rules[i].Matches(c.Channel, c.Priority) {
			assignee := rules[i].AssigneeUserID
			c.AssignedTo = &assignee
			c.Status = domain.StatusInProgress
			break
		}
	}
	if viewerRole == "admin" || viewerRole == "agent" || viewerRole == "end_user" {
		c.CreatedBy = &viewerID
	}
	if err := s.repo.CreateConversation(ctx, c); err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

func (s *Service) ListConversations(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, customerFilter *uuid.UUID, channel, status string, assigneeFilter *uuid.UUID) ([]domain.Conversation, error) {
	ids, err := s.visibleCustomers(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	// A contact filtering by a customer they do not hold gets an empty
	// list, never an error that reveals the customer exists.
	if customerFilter != nil && viewerRole != "admin" && viewerRole != "agent" {
		if !domain.CanAccessCustomer(viewerRole, ids, *customerFilter) {
			return []domain.Conversation{}, nil
		}
	}
	return s.repo.ListConversations(ctx, tenantID, ids, customerFilter, channel, status, assigneeFilter)
}

func (s *Service) GetConversation(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, id uuid.UUID) (*domain.Conversation, error) {
	return s.checkConversation(ctx, tenantID, viewerID, viewerRole, id)
}

// UpdateConversationStatus moves a conversation through its state
// machine. Staff only: the route guards this, and the service rejects
// any transition the domain does not allow.
func (s *Service) UpdateConversationStatus(ctx context.Context, tenantID uuid.UUID, c *domain.Conversation, next string) (*domain.Conversation, error) {
	current, err := s.repo.GetConversation(ctx, c.ID, tenantID)
	if err != nil {
		return nil, err
	}
	if !current.CanMoveTo(next) {
		return nil, domain.ErrInvalidTransition
	}
	current.Status = next
	if err := s.repo.UpdateConversation(ctx, current); err != nil {
		return nil, mapRepoErr(err)
	}
	return current, nil
}

// TakeConversation is the assignment race rule: taking a conversation
// already assigned to another agent fails with ErrConflict here, in
// the server, even if a client raced past the UI state.
func (s *Service) TakeConversation(ctx context.Context, tenantID, callerID uuid.UUID, id uuid.UUID) (*domain.Conversation, error) {
	c, err := s.repo.GetConversation(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	if !domain.CanTakeConversation(c, callerID) {
		return nil, domain.ErrConflict
	}
	status := c.Status
	if status == domain.StatusOpen {
		status = domain.StatusInProgress
	}
	if err := s.repo.AssignConversation(ctx, id, tenantID, &callerID, status); err != nil {
		return nil, mapRepoErr(err)
	}
	c.AssignedTo = &callerID
	c.Status = status
	return c, nil
}

// ReleaseConversation unassigns and returns the conversation to the
// open queue; only the current assignee or an admin may release it.
func (s *Service) ReleaseConversation(ctx context.Context, tenantID, callerID uuid.UUID, callerRole string, id uuid.UUID) (*domain.Conversation, error) {
	c, err := s.repo.GetConversation(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	if callerRole != "admin" && (c.AssignedTo == nil || *c.AssignedTo != callerID) {
		return nil, domain.ErrForbidden
	}
	if err := s.repo.AssignConversation(ctx, id, tenantID, nil, domain.StatusOpen); err != nil {
		return nil, mapRepoErr(err)
	}
	c.AssignedTo = nil
	c.Status = domain.StatusOpen
	return c, nil
}

// EscalateConversation raises the priority one level (critical stays
// critical) and marks the conversation escalated. A second escalation
// of an already-critical conversation is a no-op, not an error.
func (s *Service) EscalateConversation(ctx context.Context, tenantID uuid.UUID, id uuid.UUID) (*domain.Conversation, error) {
	c, err := s.repo.GetConversation(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	c.Priority = domain.NextPriority(c.Priority)
	c.Escalated = true
	c.EscalatedAt = &now
	if err := s.repo.UpdateConversation(ctx, c); err != nil {
		return nil, mapRepoErr(err)
	}
	return c, nil
}

// ── Messages ────────────────────────────────────────────────────────

// AddMessage appends a message and maintains the conversation state:
// an agent's first visible reply stamps FirstResponseAt; an agent
// reply to a waiting conversation moves it back in progress; a
// customer reply to a resolved conversation reopens it; internal
// notes are staff-only and change nothing else. A closed
// conversation accepts no new messages.
func (s *Service) AddMessage(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, convID uuid.UUID, m *domain.Message) (*domain.Message, error) {
	c, err := s.checkConversation(ctx, tenantID, viewerID, viewerRole, convID)
	if err != nil {
		return nil, err
	}
	if c.Status == domain.StatusClosed {
		return nil, domain.ErrInvalidTransition
	}
	m.TenantID = tenantID
	m.ConversationID = convID
	m.Body = strings.TrimSpace(m.Body)
	m.Channel = c.Channel
	if viewerRole == "end_user" {
		// External customers write only visible messages as themselves.
		m.SenderType = domain.SenderCustomer
		m.IsInternal = false
		m.SenderID = nil
	} else if m.SenderType == "" {
		m.SenderType = domain.SenderAgent
		m.SenderID = &viewerID
	}
	if m.IsInternal && viewerRole == "end_user" {
		return nil, domain.ErrForbidden
	}
	if err := domain.ValidateMessage(m); err != nil {
		return nil, err
	}
	if err := s.repo.CreateMessage(ctx, m); err != nil {
		return nil, mapRepoErr(err)
	}
	// Maintain conversation state from the message, unless it was an
	// internal note.
	if !m.IsInternal {
		now := time.Now()
		changed := false
		if m.SenderType == domain.SenderAgent && c.FirstResponseAt == nil {
			c.FirstResponseAt = &now
			changed = true
		}
		if m.SenderType == domain.SenderAgent && c.Status == domain.StatusWaitingCustomer {
			c.Status = domain.StatusInProgress
			changed = true
		}
		if domain.MessageReopens(c.Status, m.SenderType) {
			c.Status = domain.StatusInProgress
			changed = true
		}
		if changed {
			if err := s.repo.UpdateConversation(ctx, c); err != nil {
				return nil, err
			}
		}
		// Auto-reply: a customer message that matches a rule gets an
		// immediate system reply, so the customer is never left
		// without an answer on a known topic.
		if m.SenderType == domain.SenderCustomer {
			replies, err := s.repo.ListAutoReplies(ctx, tenantID)
			if err == nil {
				if reply := domain.MatchAutoReply(replies, c.Channel, m.Body); reply != nil {
					_ = s.repo.CreateMessage(ctx, &domain.Message{
						TenantID: tenantID, ConversationID: convID,
						SenderType: domain.SenderSystem, Body: reply.ResponseText,
						Channel: c.Channel,
					})
				}
			}
		}
	}
	return m, nil
}

func (s *Service) ListMessages(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, convID uuid.UUID) ([]domain.Message, error) {
	if _, err := s.checkConversation(ctx, tenantID, viewerID, viewerRole, convID); err != nil {
		return nil, err
	}
	includeInternal := viewerRole == "admin" || viewerRole == "agent"
	return s.repo.ListMessages(ctx, tenantID, convID, includeInternal)
}

// ── Channel ingestion (email / WhatsApp) ────────────────────────────

// IngestEmail appends an inbound email to the open conversation for
// its thread, or opens a new one; the sender becomes (or is matched
// to) a customer by email. A repeated external_message_id is a
// no-op: ingestion is idempotent.
func (s *Service) IngestEmail(ctx context.Context, tenantID uuid.UUID, accountID uuid.UUID, fromEmail, subject, body, threadID, externalMessageID string) (*domain.Conversation, error) {
	if fromEmail == "" || body == "" {
		return nil, domain.ErrValidation
	}
	if externalMessageID != "" {
		exists, err := s.repo.MessageExistsByExternalID(ctx, tenantID, externalMessageID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrConflict
		}
	}
	customer, err := s.repo.GetCustomerByEmail(ctx, tenantID, strings.ToLower(strings.TrimSpace(fromEmail)))
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	if customer == nil {
		name := fromEmail
		if at := strings.Index(fromEmail, "@"); at > 0 {
			name = fromEmail[:at]
		}
		customer = &domain.Customer{TenantID: tenantID, Name: name, Email: strings.ToLower(strings.TrimSpace(fromEmail)), Status: domain.CustomerActive}
		if err := s.repo.CreateCustomer(ctx, customer); err != nil {
			return nil, mapRepoErr(err)
		}
	}
	return s.ingest(ctx, tenantID, customer, domain.ChannelEmail, subject, body, threadID, externalMessageID)
}

// IngestWhatsApp is the WhatsApp counterpart: the sender is matched by
// phone number, the thread is the WhatsApp conversation id, and the
// webhook verify token is checked by the handler against the stored
// hash before this method is ever called.
func (s *Service) IngestWhatsApp(ctx context.Context, tenantID uuid.UUID, accountID uuid.UUID, fromPhone, body, threadID, externalMessageID string) (*domain.Conversation, error) {
	if fromPhone == "" || body == "" {
		return nil, domain.ErrValidation
	}
	if externalMessageID != "" {
		exists, err := s.repo.MessageExistsByExternalID(ctx, tenantID, externalMessageID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, domain.ErrConflict
		}
	}
	customer, err := s.repo.GetCustomerByPhone(ctx, tenantID, strings.TrimSpace(fromPhone))
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	if customer == nil {
		customer = &domain.Customer{TenantID: tenantID, Name: strings.TrimSpace(fromPhone), Phone: strings.TrimSpace(fromPhone), Status: domain.CustomerActive}
		if err := s.repo.CreateCustomer(ctx, customer); err != nil {
			return nil, mapRepoErr(err)
		}
	}
	subject := body
	if len(subject) > 80 {
		subject = subject[:80]
	}
	return s.ingest(ctx, tenantID, customer, domain.ChannelWhatsApp, subject, body, threadID, externalMessageID)
}

func (s *Service) ingest(ctx context.Context, tenantID uuid.UUID, customer *domain.Customer, channel, subject, body, threadID, externalMessageID string) (*domain.Conversation, error) {
	var conv *domain.Conversation
	var err error
	if threadID != "" {
		conv, err = s.repo.OpenConversationForThread(ctx, tenantID, customer.ID, threadID)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
	}
	if conv == nil {
		if strings.TrimSpace(subject) == "" {
			subject = "Conversación " + channel
		}
		conv = &domain.Conversation{
			TenantID: tenantID, CustomerID: customer.ID, Channel: channel,
			Subject: strings.TrimSpace(subject), Status: domain.StatusOpen, Priority: domain.PriorityMedium,
			ExternalThreadID: threadID,
		}
		policies, err := s.repo.ListSLAPolicies(ctx, tenantID)
		if err != nil {
			return nil, err
		}
		var policy *domain.SLAPolicy
		for i := range policies {
			if policies[i].Channel == channel && policies[i].Priority == conv.Priority {
				policy = &policies[i]
				break
			}
		}
		domain.ApplySLA(conv, policy, time.Now())
		if err := s.repo.CreateConversation(ctx, conv); err != nil {
			return nil, mapRepoErr(err)
		}
	} else if conv.Status == domain.StatusResolved {
		// A new inbound message on a resolved thread reopens it.
		conv.Status = domain.StatusInProgress
		if err := s.repo.UpdateConversation(ctx, conv); err != nil {
			return nil, err
		}
	}
	msg := &domain.Message{
		TenantID: tenantID, ConversationID: conv.ID, SenderType: domain.SenderCustomer,
		Body: strings.TrimSpace(body), Channel: channel, ExternalMessageID: externalMessageID,
	}
	if err := s.repo.CreateMessage(ctx, msg); err != nil {
		return nil, mapRepoErr(err)
	}
	return conv, nil
}

// ── Channel configuration ───────────────────────────────────────────

func (s *Service) ListEmailAccounts(ctx context.Context, tenantID uuid.UUID) ([]domain.EmailAccount, error) {
	return s.repo.ListEmailAccounts(ctx, tenantID)
}

func (s *Service) CreateEmailAccount(ctx context.Context, tenantID uuid.UUID, a *domain.EmailAccount) (*domain.EmailAccount, error) {
	a.TenantID = tenantID
	a.Name = strings.TrimSpace(a.Name)
	if a.SMTPPort == 0 {
		a.SMTPPort = 587
	}
	if a.IMAPPort == 0 {
		a.IMAPPort = 993
	}
	if err := domain.ValidateEmailAccount(a); err != nil {
		return nil, err
	}
	if err := s.repo.CreateEmailAccount(ctx, a); err != nil {
		return nil, mapRepoErr(err)
	}
	return a, nil
}

func (s *Service) UpdateEmailAccount(ctx context.Context, tenantID uuid.UUID, a *domain.EmailAccount) (*domain.EmailAccount, error) {
	a.TenantID = tenantID
	a.Name = strings.TrimSpace(a.Name)
	if err := domain.ValidateEmailAccount(a); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateEmailAccount(ctx, a); err != nil {
		return nil, mapRepoErr(err)
	}
	return a, nil
}

// HashWebhookToken hashes a webhook verify token the same way the
// repository stores it, so the clear token is never persisted.
func HashWebhookToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *Service) ListWhatsAppAccounts(ctx context.Context, tenantID uuid.UUID) ([]domain.WhatsAppAccount, error) {
	return s.repo.ListWhatsAppAccounts(ctx, tenantID)
}

func (s *Service) CreateWhatsAppAccount(ctx context.Context, tenantID uuid.UUID, a *domain.WhatsAppAccount, verifyToken string) (*domain.WhatsAppAccount, error) {
	a.TenantID = tenantID
	a.Name = strings.TrimSpace(a.Name)
	if verifyToken != "" {
		a.WebhookVerifyTokenHash = HashWebhookToken(verifyToken)
	}
	if err := domain.ValidateWhatsAppAccount(a); err != nil {
		return nil, err
	}
	if err := s.repo.CreateWhatsAppAccount(ctx, a); err != nil {
		return nil, mapRepoErr(err)
	}
	a.WebhookVerifyTokenHash = ""
	return a, nil
}

func (s *Service) UpdateWhatsAppAccount(ctx context.Context, tenantID uuid.UUID, a *domain.WhatsAppAccount) (*domain.WhatsAppAccount, error) {
	a.TenantID = tenantID
	a.Name = strings.TrimSpace(a.Name)
	if err := domain.ValidateWhatsAppAccount(a); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateWhatsAppAccount(ctx, a); err != nil {
		return nil, mapRepoErr(err)
	}
	return a, nil
}

// VerifyWhatsAppWebhook checks Meta's verification challenge: the
// mode must be "subscribe" and the token must hash to a stored
// account's hash. It returns the account on success.
func (s *Service) VerifyWhatsAppWebhook(ctx context.Context, tenantID uuid.UUID, mode, token string) (*domain.WhatsAppAccount, error) {
	if mode != "subscribe" || token == "" {
		return nil, domain.ErrForbidden
	}
	account, err := s.repo.FindWhatsAppAccountByTokenHash(ctx, tenantID, HashWebhookToken(token))
	if err != nil {
		return nil, domain.ErrForbidden
	}
	return account, nil
}

// ── SLA policies, assignment rules, auto-replies ────────────────────

func (s *Service) ListSLAPolicies(ctx context.Context, tenantID uuid.UUID) ([]domain.SLAPolicy, error) {
	return s.repo.ListSLAPolicies(ctx, tenantID)
}

func (s *Service) UpsertSLAPolicy(ctx context.Context, tenantID uuid.UUID, p *domain.SLAPolicy) (*domain.SLAPolicy, error) {
	p.TenantID = tenantID
	if err := domain.ValidateSLAPolicy(p); err != nil {
		return nil, err
	}
	if err := s.repo.UpsertSLAPolicy(ctx, p); err != nil {
		return nil, mapRepoErr(err)
	}
	return p, nil
}

func (s *Service) ListAssignmentRules(ctx context.Context, tenantID uuid.UUID) ([]domain.AssignmentRule, error) {
	return s.repo.ListAssignmentRules(ctx, tenantID)
}

func (s *Service) CreateAssignmentRule(ctx context.Context, tenantID uuid.UUID, r *domain.AssignmentRule) (*domain.AssignmentRule, error) {
	r.TenantID = tenantID
	if r.AssigneeUserID == uuid.Nil {
		return nil, domain.ErrValidation
	}
	ok, err := s.repo.UserExists(ctx, r.AssigneeUserID, tenantID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrValidation
	}
	if err := s.repo.CreateAssignmentRule(ctx, r); err != nil {
		return nil, mapRepoErr(err)
	}
	return r, nil
}

func (s *Service) DeleteAssignmentRule(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteAssignmentRule(ctx, id, tenantID)
}

func (s *Service) ListAutoReplies(ctx context.Context, tenantID uuid.UUID) ([]domain.AutoReply, error) {
	return s.repo.ListAutoReplies(ctx, tenantID)
}

func (s *Service) CreateAutoReply(ctx context.Context, tenantID uuid.UUID, a *domain.AutoReply) (*domain.AutoReply, error) {
	a.TenantID = tenantID
	a.ResponseText = strings.TrimSpace(a.ResponseText)
	if err := domain.ValidateAutoReply(a); err != nil {
		return nil, err
	}
	if err := s.repo.CreateAutoReply(ctx, a); err != nil {
		return nil, mapRepoErr(err)
	}
	return a, nil
}

func (s *Service) DeleteAutoReply(ctx context.Context, tenantID, id uuid.UUID) error {
	return s.repo.DeleteAutoReply(ctx, id, tenantID)
}

// ── Surveys and reports ─────────────────────────────────────────────

// CreateSurvey records a customer's rating for a conversation: only
// the owning customer, only once, and only after the conversation is
// resolved or closed.
func (s *Service) CreateSurvey(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, convID uuid.UUID, survey *domain.Survey) (*domain.Survey, error) {
	c, err := s.checkConversation(ctx, tenantID, viewerID, viewerRole, convID)
	if err != nil {
		return nil, err
	}
	if c.Status != domain.StatusResolved && c.Status != domain.StatusClosed {
		return nil, domain.ErrInvalidTransition
	}
	survey.TenantID = tenantID
	survey.ConversationID = convID
	survey.CustomerID = c.CustomerID
	if err := domain.ValidateSurvey(survey); err != nil {
		return nil, err
	}
	if err := s.repo.CreateSurvey(ctx, survey); err != nil {
		return nil, mapRepoErr(err)
	}
	return survey, nil
}

func (s *Service) ListSurveys(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string, customerFilter *uuid.UUID) ([]domain.Survey, error) {
	ids, err := s.visibleCustomers(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	return s.repo.ListSurveys(ctx, tenantID, ids, customerFilter)
}

func (s *Service) Report(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string) (*domain.CXReport, error) {
	ids, err := s.visibleCustomers(ctx, tenantID, viewerID, viewerRole)
	if err != nil {
		return nil, err
	}
	return s.repo.Report(ctx, tenantID, ids)
}
