package ticket

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
)

// NotificationRepository is the port for notification persistence.
type NotificationRepository interface {
	Create(ctx context.Context, n *domainTicket.Notification) error
	ListByUser(ctx context.Context, tenantID, userID uuid.UUID, limit int) ([]*domainTicket.Notification, error)
	MarkRead(ctx context.Context, tenantID, userID, notifID uuid.UUID) error
	MarkAllRead(ctx context.Context, tenantID, userID uuid.UUID) error
	UnreadCount(ctx context.Context, tenantID, userID uuid.UUID) (int, error)
}

// NotificationService handles notification operations.
type NotificationService struct {
	repo NotificationRepository
}

// NewNotificationService creates a new NotificationService.
func NewNotificationService(repo NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// List returns notifications for a user.
func (s *NotificationService) List(ctx context.Context, tenantID, userID uuid.UUID) ([]*domainTicket.Notification, error) {
	return s.repo.ListByUser(ctx, tenantID, userID, 50)
}

// MarkRead marks a notification as read.
func (s *NotificationService) MarkRead(ctx context.Context, tenantID, userID, notifID uuid.UUID) error {
	return s.repo.MarkRead(ctx, tenantID, userID, notifID)
}

// MarkAllRead marks all notifications as read for a user.
func (s *NotificationService) MarkAllRead(ctx context.Context, tenantID, userID uuid.UUID) error {
	return s.repo.MarkAllRead(ctx, tenantID, userID)
}

// UnreadCount returns the unread count for a user.
func (s *NotificationService) UnreadCount(ctx context.Context, tenantID, userID uuid.UUID) (int, error) {
	return s.repo.UnreadCount(ctx, tenantID, userID)
}

// CreateForUser creates a notification for a specific user and returns it.
func (s *NotificationService) CreateForUser(ctx context.Context, tenantID, userID uuid.UUID, title, body, notifType string, refID *uuid.UUID, refType string) (*domainTicket.Notification, error) {
	n := &domainTicket.Notification{
		TenantID: tenantID,
		UserID:   userID,
		Title:    title,
		Body:     body,
		Type:     notifType,
		RefID:    refID,
		RefType:  refType,
	}
	if err := s.repo.Create(ctx, n); err != nil {
		return nil, fmt.Errorf("create notification: %w", err)
	}
	return n, nil
}
