package ticket

import (
	"context"
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
	Priority   string
	AssignedTo string
	Search     string
	Limit      int
	Offset     int
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

// UserService is the application service for ticket operations.
type UserService struct {
	repo        Repository
	comments    CommentRepository
	history     HistoryRepository
	hub         EventBroadcaster
	email       EmailNotifier
	engine      AutomationEvaluator
	tx          TxManager
	slaFetcher  SLAPolicyFetcher
	notifSvc    *NotificationService
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
	TenantID    uuid.UUID
	Title       string
	Description string
	Priority    string
	Category    string
	CreatedBy   uuid.UUID
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
	var slaDeadline *time.Time
	if s.slaFetcher != nil {
		policy, err := s.slaFetcher.GetByPriority(ctx, input.TenantID, string(priority))
		if err == nil && policy != nil && policy.ResolutionTimeMinutes > 0 {
			d := t.CreatedAt.Add(time.Duration(policy.ResolutionTimeMinutes) * time.Minute)
			slaDeadline = &d
		}
	}

	// Persist the ticket (and SLA deadline) inside a transaction.
	persist := func(ctx context.Context) error {
		if err := s.repo.Create(ctx, t); err != nil {
			return err
		}
		if slaDeadline != nil {
			if err := s.repo.SetSLADeadline(ctx, t.ID, input.TenantID, *slaDeadline); err != nil {
				return err
			}
		}
		return nil
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

	if slaDeadline != nil {
		t.SLADeadline = slaDeadline
	}

	if s.hub != nil {
		s.hub.BroadcastEvent(input.TenantID, "TICKET_CREATED", t)
	}

	// Create a persistent notification for the ticket creator.
	if s.notifSvc != nil {
		title := "Ticket created: " + t.Title
		body := "Your ticket has been created successfully."
		s.createNotificationSafe(context.Background(), input.TenantID, input.CreatedBy, title, body, "ticket_created", &t.ID, "ticket")
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
func (s *UserService) ChangeStatus(ctx context.Context, id, tenantID, userID uuid.UUID, role, newStatus string) error {
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

	if err := t.ChangeStatus(status); err != nil {
		return err
	}

	// Persist the status update (and history insertion) inside a transaction.
	persist := func(ctx context.Context) error {
		_, err := s.repo.UpdateStatus(ctx, id, tenantID, version, status, userID)
		return err
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

	// Notify the ticket creator about the status change (if they're not the one changing it).
	if s.notifSvc != nil {
		if creatorID, err := s.repo.GetCreator(ctx, id, tenantID); err == nil && creatorID != userID {
			title := "Ticket updated: status changed to " + string(status)
			body := "The status of your ticket has been updated."
			s.createNotificationSafe(context.Background(), tenantID, creatorID, title, body, "ticket_updated", &id, "ticket")
		}
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

	if err := s.comments.Create(ctx, c); err != nil {
		return nil, err
	}

	if s.hub != nil {
		s.hub.BroadcastEvent(input.TenantID, "NEW_COMMENT", c)
	}

	// Notify the ticket creator about the new comment (if they're not the commenter).
	if s.notifSvc != nil {
		if creatorID, err := s.repo.GetCreator(ctx, input.TicketID, input.TenantID); err == nil && creatorID != input.UserID {
			title := "New comment on your ticket"
			body := "Someone commented on your ticket."
			s.createNotificationSafe(context.Background(), input.TenantID, creatorID, title, body, "new_comment", &input.TicketID, "ticket")
		}
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
func (s *UserService) createNotificationSafe(ctx context.Context, tenantID, userID uuid.UUID, title, body, notifType string, refID *uuid.UUID, refType string) {
	if s.notifSvc == nil {
		return
	}
	n, err := s.notifSvc.CreateForUser(ctx, tenantID, userID, title, body, notifType, refID, refType)
	if err != nil {
		return
	}
	// Broadcast the notification in real-time via WebSocket.
	if s.hub != nil {
		s.hub.BroadcastEvent(tenantID, "NOTIFICATION", n)
	}
}
