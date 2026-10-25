package ticket

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	appAuto "github.com/trickreport/backend/internal/application/automation"
	domainSLA "github.com/trickreport/backend/internal/domain/sla"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
)

// AutomationEvaluator is the port for evaluating automation rules on events.
// Implemented by the automation engine.
type AutomationEvaluator interface {
	Evaluate(ctx context.Context, tenantID uuid.UUID, event appAuto.Event) error
}

// TxManager is the port for running operations inside a database transaction.
// The transaction is propagated to repositories via the context.
//
// TODO(wire): wire TxManager into wire.go's NewServices and pass it to the
// ticket service via SetTxManager.
type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// SLAPolicyFetcher is the port for fetching SLA policies by priority.
//
// TODO(wire): wire SLAPolicyFetcher into wire.go's NewServices and pass it
// to the ticket service via SetSLAPolicyFetcher.
type SLAPolicyFetcher interface {
	GetByPriority(ctx context.Context, tenantID uuid.UUID, priority string) (*domainSLA.Policy, error)
}

// Filter holds query parameters for listing tickets.
type Filter struct {
	Status     string
	Category   string
	Priority   string
	AssignedTo string
	Search     string
	Limit      int
	Offset     int
}

// IdempotentCreator reserves a logical create and persists its result atomically.
type IdempotentCreator interface {
	CreateIdempotent(ctx context.Context, t *domainTicket.Ticket, key uuid.UUID, digest []byte) (bool, error)
}

// Repository is the port for ticket persistence.
type Repository interface {
	List(ctx context.Context, tenantID uuid.UUID, filter Filter, role string, userID uuid.UUID) ([]domainTicket.Ticket, error)
	GetByID(ctx context.Context, id, tenantID uuid.UUID) (*domainTicket.Ticket, error)
	Create(ctx context.Context, t *domainTicket.Ticket) error
	UpdateStatus(ctx context.Context, id, tenantID uuid.UUID, version int, status domainTicket.Status, userID uuid.UUID) (*domainTicket.Ticket, error)
	Assign(ctx context.Context, id, tenantID, assignedTo, userID uuid.UUID) error
	SetSLADeadline(ctx context.Context, id, tenantID uuid.UUID, deadline time.Time) error
	GetCreator(ctx context.Context, id, tenantID uuid.UUID) (uuid.UUID, error)
}

// CommentRepository is the port for comment persistence.
type CommentRepository interface {
	List(ctx context.Context, ticketID uuid.UUID, role string) ([]domainTicket.Comment, error)
	Create(ctx context.Context, c *domainTicket.Comment) error
}

// HistoryRepository is the port for history persistence.
type HistoryRepository interface {
	List(ctx context.Context, ticketID uuid.UUID) ([]domainTicket.HistoryEntry, error)
}

// EventBroadcaster is the port for real-time event broadcasting.
type EventBroadcaster interface {
	BroadcastEvent(tenantID uuid.UUID, eventType string, data any)
}

// EmailNotifier is the port for email notifications.
type EmailNotifier interface {
	Send(to, subject, body string) error
}

// UserEmailFetcher is the port for looking up a user's email address.
type UserEmailFetcher interface {
	GetEmail(ctx context.Context, tenantID, userID uuid.UUID) (string, error)
}

// EmailRenderer is the port for rendering an HTML email body from one of the
// infrastructure layer's named templates ("ticket_created", "ticket_updated").
// Kept narrow and application-owned so this package never imports the
// infrastructure/email package directly (hexagonal layering).
type EmailRenderer interface {
	RenderTicketEmail(templateName, ticketID, title, description, priority, status string) (string, error)
}

// UserService is the application service for ticket operations.
type UserService struct {
	repo          Repository
	comments      CommentRepository
	history       HistoryRepository
	hub           EventBroadcaster
	email         EmailNotifier
	emailFetcher  UserEmailFetcher
	emailRenderer EmailRenderer
	engine        AutomationEvaluator
	tx            TxManager
	slaFetcher    SLAPolicyFetcher
	notifSvc      *NotificationService
}

// NewService creates a new ticket application service.
func NewService(repo Repository, comments CommentRepository, history HistoryRepository, hub EventBroadcaster, email EmailNotifier) *UserService {
	return &UserService{
		repo:     repo,
		comments: comments,
		history:  history,
		hub:      hub,
		email:    email,
	}
}

// SetEngine wires the automation engine. It is optional — when set, ticket
// events are evaluated against automation rules.
func (s *UserService) SetEngine(engine AutomationEvaluator) {
	s.engine = engine
}

// SetTxManager wires the transaction manager. When set, Create and
// ChangeStatus wrap their persistence calls in a single transaction.
func (s *UserService) SetTxManager(tx TxManager) {
	s.tx = tx
}

// SetSLAPolicyFetcher wires the SLA policy fetcher. When set, Create
// calculates and persists the SLA deadline based on the tenant's policy.
func (s *UserService) SetSLAPolicyFetcher(f SLAPolicyFetcher) {
	s.slaFetcher = f
}

// SetNotificationService wires the notification service. When set, ticket
// events create persistent notifications for relevant users.
func (s *UserService) SetNotificationService(n *NotificationService) {
	s.notifSvc = n
}

// SetUserEmailFetcher wires the user email lookup. Required (alongside
// SetEmailRenderer) for ticket events to also notify by email, not just
// in-app.
func (s *UserService) SetUserEmailFetcher(f UserEmailFetcher) {
	s.emailFetcher = f
}

// SetEmailRenderer wires the HTML email template renderer.
func (s *UserService) SetEmailRenderer(r EmailRenderer) {
	s.emailRenderer = r
}

// evaluateAutomation fires an automation event asynchronously.
func (s *UserService) evaluateAutomation(tenantID uuid.UUID, event appAuto.Event) {
	if s.engine == nil {
		return
	}
	go func() {
		ctx := context.Background()
		if err := s.engine.Evaluate(ctx, tenantID, event); err != nil {
			// Errors are logged by the engine; nothing to do here.
			_ = err
		}
	}()
}

// List returns tickets for the tenant, filtered by role and query params.
func (s *UserService) List(ctx context.Context, tenantID uuid.UUID, filter Filter, role string, userID uuid.UUID) ([]domainTicket.Ticket, error) {
	return s.repo.List(ctx, tenantID, filter, role, userID)
}

// Get returns a single ticket, checking role-based visibility.
func (s *UserService) Get(ctx context.Context, id, tenantID uuid.UUID, role string, userID uuid.UUID) (*domainTicket.Ticket, error) {
	t, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return nil, err
	}
	if !t.CanBeViewedBy(userID, role) {
		return nil, domainTicket.ErrForbidden
	}
	return t, nil
}

// CreateInput holds the data for creating a new ticket.
type CreateInput struct {
	TenantID       uuid.UUID
	Title          string
	Description    string
	Priority       string
	Category       string
	CreatedBy      uuid.UUID
	IdempotencyKey uuid.UUID
}

// Create creates a new ticket and broadcasts a real-time event.
func (s *UserService) Create(ctx context.Context, input CreateInput) (*domainTicket.Ticket, error) {
	if input.Title == "" || input.Description == "" {
		return nil, domainTicket.ErrValidation
	}

	priority := domainTicket.PriorityLow
	if input.Priority != "" {
		p, err := domainTicket.ParsePriority(input.Priority)
		if err != nil {
			return nil, domainTicket.ErrValidation
		}
		priority = p
	}

	category := input.Category
	if category == "" {
		category = "general"
	}

	t := &domainTicket.Ticket{
		TenantID:    input.TenantID,
		Title:       input.Title,
		Description: input.Description,
		Status:      domainTicket.StatusOpen,
		Priority:    priority,
		Category:    category,
		CreatedBy:   input.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Fetch the SLA policy (read, outside the transaction) to calculate the
	// resolution deadline.
	if s.slaFetcher != nil {
		policy, err := s.slaFetcher.GetByPriority(ctx, input.TenantID, string(priority))
		if err != nil && !errors.Is(err, domainSLA.ErrNotFound) {
			return nil, err
		}
		if policy != nil {
			if err := policy.Validate(); err != nil {
				return nil, domainTicket.ErrValidation
			}
			d := domainSLA.NewDeadlines(t.CreatedAt, policy.ResponseTimeMinutes, policy.ResolutionTimeMinutes, policy.EscalationMinutes)
			t.ResponseDeadline, t.SLADeadline, t.EscalationDeadline = d.Response, d.Resolution, d.Escalation
		}
	}

	// Persist the ticket (and SLA deadline) inside a transaction.
	created := true
	persist := func(ctx context.Context) error {
		if input.IdempotencyKey != uuid.Nil {
			repo, ok := s.repo.(IdempotentCreator)
			if !ok {
				return errors.New("ticket repository does not support idempotency")
			}
			payload, _ := json.Marshal([]string{input.Title, input.Description, input.Priority, input.Category})
			digest := sha256.Sum256(payload)
			var err error
			created, err = repo.CreateIdempotent(ctx, t, input.IdempotencyKey, digest[:])
			if err != nil {
				return err
			}
		} else if err := s.repo.Create(ctx, t); err != nil {
			return err
		}
		if !created {
			return nil
		}
		if err := s.createNotificationSafe(ctx, input.TenantID, input.CreatedBy, "Ticket created: "+t.Title, "Your ticket has been created successfully.", "ticket_created", &t.ID, "ticket"); err != nil {
			return err
		}
		return s.notifyByEmailSafe(ctx, input.TenantID, input.CreatedBy, "ticket_created", t.ID.String(), t.Title, t.Description, string(t.Priority), string(t.Status))
	}

	if s.tx != nil {
		if err := s.tx.RunInTx(ctx, persist); err != nil {
			return nil, err
		}
	} else {
		if err := persist(ctx); err != nil {
			return nil, err
		}
	}

	if !created {
		return t, nil
	}

	if s.hub != nil {
		s.hub.BroadcastEvent(input.TenantID, "TICKET_CREATED", t)
	}

	// Fire automation rules for ticket creation.
	s.evaluateAutomation(input.TenantID, appAuto.Event{
		Type:     "ticket_created",
		TenantID: input.TenantID,
		TicketID: t.ID,
		UserID:   input.CreatedBy,
		NewValue: string(t.Priority),
	})

	return t, nil
}

// ChangeStatus transitions a ticket to a new status with authorization checks.
func (s *UserService) ChangeStatus(ctx context.Context, id, tenantID, userID uuid.UUID, role, newStatus string, expectedVersion ...int) error {
	status, err := domainTicket.ParseStatus(newStatus)
	if err != nil {
		return domainTicket.ErrValidation
	}

	t, err := s.repo.GetByID(ctx, id, tenantID)
	if err != nil {
		return err
	}

	if !domainTicket.CanStatusBeChangedBy(role, t, userID, status) {
		return domainTicket.ErrForbidden
	}

	oldStatus := string(t.Status)
	version := t.Version
	if len(expectedVersion) > 0 && expectedVersion[0] != version {
		return domainTicket.ErrConcurrentModification
	}
	if t.Status == status {
		return nil
	}

	if err := t.ChangeStatus(status); err != nil {
		return err
	}

	// Persist the status update (and history insertion) inside a transaction.
	persist := func(ctx context.Context) error {
		_, err := s.repo.UpdateStatus(ctx, id, tenantID, version, status, userID)
		if err != nil {
			return err
		}
		if t.CreatedBy == userID {
			return nil
		}
		if err := s.createNotificationSafe(ctx, tenantID, t.CreatedBy, "Ticket updated: status changed to "+string(status), "The status of your ticket has been updated.", "ticket_updated", &id, "ticket"); err != nil {
			return err
		}
		return s.notifyByEmailSafe(ctx, tenantID, t.CreatedBy, "ticket_updated", id.String(), t.Title, t.Description, string(t.Priority), string(t.Status))
	}

	if s.tx != nil {
		if err := s.tx.RunInTx(ctx, persist); err != nil {
			return err
		}
	} else {
		if err := persist(ctx); err != nil {
			return err
		}
	}

	if s.hub != nil {
		s.hub.BroadcastEvent(tenantID, "TICKET_UPDATED", map[string]any{
			"ticket_id": id,
			"status":    string(status),
		})
	}

	// Fire automation rules for status changes.
	s.evaluateAutomation(tenantID, appAuto.Event{
		Type:     "ticket_status_changed",
		TenantID: tenantID,
		TicketID: id,
		UserID:   userID,
		OldValue: oldStatus,
		NewValue: string(status),
	})

	return nil
}

// Assign assigns a ticket to a user.
func (s *UserService) Assign(ctx context.Context, id, tenantID, assignedTo, userID uuid.UUID) error {
	return s.repo.Assign(ctx, id, tenantID, assignedTo, userID)
}

// ListComments returns comments for a ticket, filtered by role visibility.
func (s *UserService) ListComments(ctx context.Context, ticketID, tenantID uuid.UUID, role string, userID uuid.UUID) ([]domainTicket.Comment, error) {
	// Verify access
	t, err := s.repo.GetByID(ctx, ticketID, tenantID)
	if err != nil {
		return nil, err
	}
	if !t.CanBeViewedBy(userID, role) {
		return nil, domainTicket.ErrForbidden
	}
	return s.comments.List(ctx, ticketID, role)
}

// AddCommentInput holds the data for creating a comment.
type AddCommentInput struct {
	TicketID   uuid.UUID
	TenantID   uuid.UUID
	UserID     uuid.UUID
	Role       string
	Content    string
	IsInternal bool
}

// AddComment adds a comment to a ticket and broadcasts a real-time event.
func (s *UserService) AddComment(ctx context.Context, input AddCommentInput) (*domainTicket.Comment, error) {
	if input.Content == "" {
		return nil, domainTicket.ErrValidation
	}

	// End users cannot create internal comments
	if input.Role == "end_user" {
		input.IsInternal = false
	}

	// Verify access
	t, err := s.repo.GetByID(ctx, input.TicketID, input.TenantID)
	if err != nil {
		return nil, err
	}
	if !t.CanBeViewedBy(input.UserID, input.Role) {
		return nil, domainTicket.ErrForbidden
	}

	c := &domainTicket.Comment{
		TicketID:   input.TicketID,
		UserID:     input.UserID,
		Content:    input.Content,
		IsInternal: input.IsInternal,
		CreatedAt:  time.Now(),
	}

	persist := func(ctx context.Context) error {
		if err := s.comments.Create(ctx, c); err != nil {
			return err
		}
		if input.IsInternal || t.CreatedBy == input.UserID {
			return nil
		}
		if err := s.createNotificationSafe(ctx, input.TenantID, t.CreatedBy, "New comment on your ticket", "Someone commented on your ticket.", "new_comment", &input.TicketID, "ticket"); err != nil {
			return err
		}
		return s.notifyByEmailSafe(ctx, input.TenantID, t.CreatedBy, "ticket_updated", input.TicketID.String(), t.Title, input.Content, string(t.Priority), string(t.Status))
	}
	if s.tx != nil {
		err = s.tx.RunInTx(ctx, persist)
	} else {
		err = persist(ctx)
	}
	if err != nil {
		return nil, err
	}

	if s.hub != nil {
		s.hub.BroadcastEvent(input.TenantID, "NEW_COMMENT", c)
	}

	return c, nil
}

// ListHistory returns the audit log for a ticket.
func (s *UserService) ListHistory(ctx context.Context, ticketID, tenantID uuid.UUID, role string, userID uuid.UUID) ([]domainTicket.HistoryEntry, error) {
	// Verify access
	t, err := s.repo.GetByID(ctx, ticketID, tenantID)
	if err != nil {
		return nil, err
	}
	if !t.CanBeViewedBy(userID, role) {
		return nil, domainTicket.ErrForbidden
	}
	return s.history.List(ctx, ticketID)
}

// createNotificationSafe creates a notification and broadcasts it via WebSocket.
// Errors are silently ignored to avoid failing the main operation.
func (s *UserService) createNotificationSafe(ctx context.Context, tenantID, userID uuid.UUID, title, body, notifType string, refID *uuid.UUID, refType string) error {
	if s.notifSvc == nil {
		return nil
	}
	_, err := s.notifSvc.CreateForUser(ctx, tenantID, userID, title, body, notifType, refID, refType)
	// Ticket events are emitted after commit and refresh the recipient's notifications.
	return err
}

// notifyByEmailSafe renders and sends a ticket-event email. Best-effort, same
// as createNotificationSafe: a missing dependency or a failed lookup/send
// never fails the calling operation, it just means no email goes out.
func (s *UserService) notifyByEmailSafe(ctx context.Context, tenantID, userID uuid.UUID, templateName, ticketID, title, description, priority, status string) error {
	if s.email == nil || s.emailFetcher == nil || s.emailRenderer == nil {
		return nil
	}
	to, err := s.emailFetcher.GetEmail(ctx, tenantID, userID)
	if err != nil {
		return err
	}
	if to == "" {
		return nil
	}
	body, err := s.emailRenderer.RenderTicketEmail(templateName, ticketID, title, description, priority, status)
	if err != nil {
		return err
	}
	if sender, ok := s.email.(interface {
		SendContext(context.Context, string, string, string) error
	}); ok {
		return sender.SendContext(ctx, to, "Trickreport: "+title, body)
	}
	return s.email.Send(to, "Trickreport: "+title, body)
}
