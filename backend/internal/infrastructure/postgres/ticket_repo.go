package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	appTicket "github.com/trickreport/backend/internal/application/ticket"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
)

// TicketRepo implements ticket.Repository using PostgreSQL.
type TicketRepo struct {
	db *pgxpool.Pool
}

func NewTicketRepo(db *pgxpool.Pool) *TicketRepo {
	return &TicketRepo{db: db}
}

func (r *TicketRepo) List(ctx context.Context, tenantID uuid.UUID, filter appTicket.Filter, role string, userID uuid.UUID) ([]domainTicket.Ticket, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	q := psql.Select(
		"t.id", "t.tenant_id", "t.title", "t.description", "t.status", "t.priority", "t.category",
		"t.created_by", "t.assigned_to", "t.sla_deadline", "t.sla_breached", "t.created_at", "t.updated_at",
		"c.name AS creator_name", "a.name AS assignee_name",
	).
		From("tickets t").
		Join("users c ON t.created_by = c.id").
		LeftJoin("users a ON t.assigned_to = a.id").
		Where(squirrel.Eq{"t.tenant_id": tenantID}).
		OrderBy("t.created_at DESC")

	if role == "end_user" {
		q = q.Where(squirrel.Eq{"t.created_by": userID})
	}
	if filter.Status != "" {
		q = q.Where(squirrel.Eq{"t.status": filter.Status})
	}
	if filter.Priority != "" {
		q = q.Where(squirrel.Eq{"t.priority": filter.Priority})
	}
	if filter.AssignedTo != "" {
		q = q.Where(squirrel.Eq{"t.assigned_to": filter.AssignedTo})
	}
	if filter.Search != "" {
		q = q.Where("t.search_vector @@ plainto_tsquery('english', ?)", filter.Search)
	}
	if filter.Limit > 0 {
		q = q.Limit(uint64(filter.Limit))
	}
	if filter.Offset > 0 {
		q = q.Offset(uint64(filter.Offset))
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets: %w", err)
	}
	defer rows.Close()

	tickets := make([]domainTicket.Ticket, 0)
	for rows.Next() {
		var t domainTicket.Ticket
		var id, tenantIDCol, createdBy pgtype.UUID
		var assignedTo pgtype.UUID
		if err := rows.Scan(
			&id, &tenantIDCol, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Category,
			&createdBy, &assignedTo, &t.SLADeadline, &t.SLABreached, &t.CreatedAt, &t.UpdatedAt,
			&t.CreatorName, &t.AssigneeName,
		); err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %w", err)
		}
		t.ID = pgToUUID(id)
		t.TenantID = pgToUUID(tenantIDCol)
		t.CreatedBy = pgToUUID(createdBy)
		t.AssignedTo = pgToUUIDPtr(assignedTo)
		tickets = append(tickets, t)
	}

	return tickets, nil
}

func (r *TicketRepo) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*domainTicket.Ticket, error) {
	q := `
		SELECT t.id, t.tenant_id, t.title, t.description, t.status, t.priority, t.category,
		       t.created_by, t.assigned_to, t.sla_deadline, t.sla_breached, t.created_at, t.updated_at,
		       c.name AS creator_name, a.name AS assignee_name
		FROM tickets t
		JOIN users c ON t.created_by = c.id
		LEFT JOIN users a ON t.assigned_to = a.id
		WHERE t.id = $1 AND t.tenant_id = $2`

	var t domainTicket.Ticket
	var idCol, tenantIDCol, createdBy pgtype.UUID
	var assignedTo pgtype.UUID
	err := r.db.QueryRow(ctx, q, id, tenantID).Scan(
		&idCol, &tenantIDCol, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Category,
		&createdBy, &assignedTo, &t.SLADeadline, &t.SLABreached, &t.CreatedAt, &t.UpdatedAt,
		&t.CreatorName, &t.AssigneeName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainTicket.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}
	t.ID = pgToUUID(idCol)
	t.TenantID = pgToUUID(tenantIDCol)
	t.CreatedBy = pgToUUID(createdBy)
	t.AssignedTo = pgToUUIDPtr(assignedTo)
	return &t, nil
}

func (r *TicketRepo) Create(ctx context.Context, t *domainTicket.Ticket) error {
	q := `
		INSERT INTO tickets (tenant_id, title, description, priority, category, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, status, created_at, updated_at`

	var id pgtype.UUID
	err := r.db.QueryRow(ctx, q, t.TenantID, t.Title, t.Description, t.Priority, t.Category, t.CreatedBy).
		Scan(&id, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}
	t.ID = pgToUUID(id)
	return nil
}

func (r *TicketRepo) UpdateStatus(ctx context.Context, id, tenantID uuid.UUID, status domainTicket.Status, userID uuid.UUID) (*domainTicket.Ticket, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE tickets SET status = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3`, status, id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to update ticket status: %w", err)
	}

	_, err = tx.Exec(ctx, `INSERT INTO ticket_history (ticket_id, user_id, field, old_value, new_value) VALUES ($1, $2, 'status', $3, $4)`, id, userID, "", string(status))
	if err != nil {
		return nil, fmt.Errorf("failed to insert history: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return r.GetByID(ctx, id, tenantID)
}

func (r *TicketRepo) Assign(ctx context.Context, id, tenantID, assignedTo, userID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var currentAssignee pgtype.UUID
	err = tx.QueryRow(ctx, `SELECT assigned_to FROM tickets WHERE id = $1 AND tenant_id = $2 FOR UPDATE`, id, tenantID).Scan(&currentAssignee)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainTicket.ErrNotFound
		}
		return fmt.Errorf("failed to fetch ticket: %w", err)
	}

	_, err = tx.Exec(ctx, `UPDATE tickets SET assigned_to = $1, updated_at = NOW() WHERE id = $2`, assignedTo, id)
	if err != nil {
		return fmt.Errorf("failed to assign ticket: %w", err)
	}

	oldVal := ""
	if currentAssignee.Valid {
		oldVal = pgToUUID(currentAssignee).String()
	}
	newVal := assignedTo.String()
	if oldVal != newVal {
		_, err = tx.Exec(ctx, `INSERT INTO ticket_history (ticket_id, user_id, field, old_value, new_value) VALUES ($1, $2, 'assigned_to', $3, $4)`, id, userID, oldVal, newVal)
		if err != nil {
			return fmt.Errorf("failed to insert history: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}
	return nil
}

func (r *TicketRepo) GetCreator(ctx context.Context, id, tenantID uuid.UUID) (uuid.UUID, error) {
	var createdBy pgtype.UUID
	err := r.db.QueryRow(ctx, `SELECT created_by FROM tickets WHERE id = $1 AND tenant_id = $2`, id, tenantID).Scan(&createdBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, domainTicket.ErrNotFound
		}
		return uuid.Nil, fmt.Errorf("failed to get ticket creator: %w", err)
	}
	return pgToUUID(createdBy), nil
}

// CommentRepo implements ticket.CommentRepository using PostgreSQL.
type CommentRepo struct {
	db *pgxpool.Pool
}

func NewCommentRepo(db *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{db: db}
}

func (r *CommentRepo) List(ctx context.Context, ticketID uuid.UUID, role string) ([]domainTicket.Comment, error) {
	q := `
		SELECT c.id, c.ticket_id, c.user_id, c.content, c.is_internal, c.created_at, u.name
		FROM ticket_comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.ticket_id = $1`
	args := []any{ticketID}

	if role == "end_user" {
		q += ` AND c.is_internal = FALSE`
	}
	q += ` ORDER BY c.created_at ASC`

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}
	defer rows.Close()

	comments := make([]domainTicket.Comment, 0)
	for rows.Next() {
		var c domainTicket.Comment
		var id, ticketIDCol, userID pgtype.UUID
		if err := rows.Scan(&id, &ticketIDCol, &userID, &c.Content, &c.IsInternal, &c.CreatedAt, &c.UserName); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		c.ID = pgToUUID(id)
		c.TicketID = pgToUUID(ticketIDCol)
		c.UserID = pgToUUID(userID)
		comments = append(comments, c)
	}
	return comments, nil
}

func (r *CommentRepo) Create(ctx context.Context, c *domainTicket.Comment) error {
	q := `
		INSERT INTO ticket_comments (ticket_id, user_id, content, is_internal)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	var id pgtype.UUID
	err := r.db.QueryRow(ctx, q, c.TicketID, c.UserID, c.Content, c.IsInternal).Scan(&id, &c.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	c.ID = pgToUUID(id)

	// Fetch user name for response
	_ = r.db.QueryRow(ctx, `SELECT name FROM users WHERE id = $1`, c.UserID).Scan(&c.UserName)
	return nil
}

// HistoryRepo implements ticket.HistoryRepository using PostgreSQL.
type HistoryRepo struct {
	db *pgxpool.Pool
}

func NewHistoryRepo(db *pgxpool.Pool) *HistoryRepo {
	return &HistoryRepo{db: db}
}

func (r *HistoryRepo) List(ctx context.Context, ticketID uuid.UUID) ([]domainTicket.HistoryEntry, error) {
	q := `
		SELECT h.id, h.ticket_id, h.user_id, h.field, h.old_value, h.new_value, h.created_at, u.name
		FROM ticket_history h
		JOIN users u ON h.user_id = u.id
		WHERE h.ticket_id = $1
		ORDER BY h.created_at DESC`

	rows, err := r.db.Query(ctx, q, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to list history: %w", err)
	}
	defer rows.Close()

	entries := make([]domainTicket.HistoryEntry, 0)
	for rows.Next() {
		var e domainTicket.HistoryEntry
		var id, ticketIDCol, userID pgtype.UUID
		if err := rows.Scan(&id, &ticketIDCol, &userID, &e.Field, &e.OldValue, &e.NewValue, &e.CreatedAt, &e.UserName); err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		e.ID = pgToUUID(id)
		e.TicketID = pgToUUID(ticketIDCol)
		e.UserID = pgToUUID(userID)
		entries = append(entries, e)
	}
	return entries, nil
}

// Ensure time is referenced (used in domain types)
var _ = time.Now
