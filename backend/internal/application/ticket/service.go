package ticket

import (
	"context"
	"time"

	"github.com/google/uuid"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
)

// Filter holds query parameters for listing tickets.
type Filter struct {
	Status     string
	Priority   string
	AssignedTo string
	Limit      int
	Offset     int
}

// Repository is the port for ticket persistence.
type Repository interface {
	List(ctx context.Context, tenantID uuid.UUID, filter Filter, role string, userID uuid.UUID) ([]domainTicket.Ticket, error)
	GetByID(ctx context.Context, id, tenantID uuid.UUID) (*domainTicket.Ticket, error)
	Create(ctx context.Context, t *domainTicket.Ticket) error
	UpdateStatus(ctx context.Context, id, tenantID uuid.UUID, status domainTicket.Status, userID uuid.UUID) (*domainTicket.Ticket, error)
	Assign(ctx context.Context, id, tenantID, assignedTo, userID uuid.UUID) error
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
	repo     Repository
	comments CommentRepository
	history  HistoryRepository
	hub      EventBroadcaster
	email    EmailNotifier
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

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	if s.hub != nil {
		s.hub.BroadcastEvent(input.TenantID, "TICKET_CREATED", t)
	}

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

	if err := t.ChangeStatus(status); err != nil {
		return err
	}

	if _, err := s.repo.UpdateStatus(ctx, id, tenantID, status, userID); err != nil {
		return err
	}

	if s.hub != nil {
		s.hub.BroadcastEvent(tenantID, "TICKET_UPDATED", map[string]any{
			"ticket_id": id,
			"status":    string(status),
		})
	}

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
