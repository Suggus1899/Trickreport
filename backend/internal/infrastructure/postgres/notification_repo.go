package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
)

// NotificationRepo persists notifications in PostgreSQL.
type NotificationRepo struct {
	db *pgxpool.Pool
}

// NewNotificationRepo creates a new NotificationRepo.
func NewNotificationRepo(db *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{db: db}
}

func scanNotification(row pgx.Row) (*domainTicket.Notification, error) {
	n := &domainTicket.Notification{}
	var refID *uuid.UUID
	err := row.Scan(
		&n.ID, &n.TenantID, &n.UserID, &n.Title, &n.Body, &n.Type,
		&refID, &n.RefType, &n.Read, &n.CreatedAt,
	)
	if refID != nil {
		n.RefID = refID
	}
	return n, err
}

// Create inserts a new notification.
func (r *NotificationRepo) Create(ctx context.Context, n *domainTicket.Notification) error {
	query := `
		INSERT INTO notifications (tenant_id, user_id, title, body, type, ref_id, ref_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`
	return r.db.QueryRow(ctx, query,
		n.TenantID, n.UserID, n.Title, n.Body, n.Type, n.RefID, n.RefType,
	).Scan(&n.ID, &n.CreatedAt)
}

// ListByUser returns notifications for a user, newest first, limited to `limit`.
func (r *NotificationRepo) ListByUser(ctx context.Context, tenantID, userID uuid.UUID, limit int) ([]*domainTicket.Notification, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	query := `
		SELECT id, tenant_id, user_id, title, body, type, ref_id, ref_type, read, created_at
		FROM notifications
		WHERE tenant_id = $1 AND user_id = $2
		ORDER BY created_at DESC
		LIMIT $3`
	rows, err := r.db.Query(ctx, query, tenantID, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("notification list: %w", err)
	}
	defer rows.Close()

	var notifications []*domainTicket.Notification
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return nil, fmt.Errorf("notification scan: %w", err)
		}
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

// MarkRead marks a single notification as read.
func (r *NotificationRepo) MarkRead(ctx context.Context, tenantID, userID, notifID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE notifications SET read = TRUE
		WHERE id = $1 AND tenant_id = $2 AND user_id = $3`,
		notifID, tenantID, userID,
	)
	if err != nil {
		return fmt.Errorf("notification mark read: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("notification not found")
	}
	return nil
}

// MarkAllRead marks all unread notifications for a user as read.
func (r *NotificationRepo) MarkAllRead(ctx context.Context, tenantID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications SET read = TRUE
		WHERE tenant_id = $1 AND user_id = $2 AND read = FALSE`,
		tenantID, userID,
	)
	if err != nil {
		return fmt.Errorf("notification mark all read: %w", err)
	}
	return nil
}

// UnreadCount returns the number of unread notifications for a user.
func (r *NotificationRepo) UnreadCount(ctx context.Context, tenantID, userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM notifications
		WHERE tenant_id = $1 AND user_id = $2 AND read = FALSE`,
		tenantID, userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("notification unread count: %w", err)
	}
	return count, nil
}

// DeleteOld removes notifications older than the given duration.
func (r *NotificationRepo) DeleteOld(ctx context.Context, before time.Time) (int64, error) {
	tag, err := r.db.Exec(ctx, `DELETE FROM notifications WHERE created_at < $1`, before)
	if err != nil {
		return 0, fmt.Errorf("notification cleanup: %w", err)
	}
	return tag.RowsAffected(), nil
}
