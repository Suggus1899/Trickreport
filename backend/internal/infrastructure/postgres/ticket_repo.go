package postgres

import (
	"bytes"
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

// conn returns the executor to use for the given context: an active
// transaction if one is present, otherwise the connection pool.
func (r *TicketRepo) conn(ctx context.Context) DBTX {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}
	return r.db
}

func (r *TicketRepo) List(ctx context.Context, tenantID uuid.UUID, filter appTicket.Filter, role string, userID uuid.UUID) ([]domainTicket.Ticket, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	q := psql.Select(
		"t.id", "t.tenant_id", "t.title", "t.description", "t.status", "t.priority", "t.category",
		"t.created_by", "t.assigned_to", "t.sla_deadline", "t.sla_breached", "t.version", "t.created_at", "t.updated_at",
		"c.name AS creator_name", "a.name AS assignee_name", "t.response_deadline", "t.escalation_deadline", "t.first_responded_at", "t.resolved_at", "t.response_sla_breached",
	).
		From("tickets t").
		Join("users c ON t.created_by = c.id").
		LeftJoin("users a ON t.assigned_to = a.id").
		Where(squirrel.Eq{"t.tenant_id": tenantID}).
		OrderBy("t.created_at DESC", "t.id DESC")

	if role == "end_user" {
		q = q.Where(squirrel.Eq{"t.created_by": userID})
	}
	if filter.Status != "" {
		q = q.Where(squirrel.Eq{"t.status": filter.Status})
	}
	if filter.Priority != "" {
		q = q.Where(squirrel.Eq{"t.priority": filter.Priority})
	}
	if filter.Category != "" {
		q = q.Where(squirrel.Eq{"t.category": filter.Category})
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
			&createdBy, &assignedTo, &t.SLADeadline, &t.SLABreached, &t.Version, &t.CreatedAt, &t.UpdatedAt,
			&t.CreatorName, &t.AssigneeName,
			&t.ResponseDeadline, &t.EscalationDeadline, &t.FirstRespondedAt, &t.ResolvedAt, &t.ResponseSLABreached,
		); err != nil {
			return nil, fmt.Errorf("failed to scan ticket: %w", err)
		}
		t.ID = pgToUUID(id)
		t.TenantID = pgToUUID(tenantIDCol)
		t.CreatedBy = pgToUUID(createdBy)
		t.AssignedTo = pgToUUIDPtr(assignedTo)
		tickets = append(tickets, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list tickets: %w", err)
	}
	return tickets, nil
}

func (r *TicketRepo) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*domainTicket.Ticket, error) {
	q := `
		SELECT t.id, t.tenant_id, t.title, t.description, t.status, t.priority, t.category,
		       t.created_by, t.assigned_to, t.sla_deadline, t.sla_breached, t.version, t.created_at, t.updated_at,
		       c.name AS creator_name, a.name AS assignee_name,
		       t.response_deadline, t.escalation_deadline, t.first_responded_at, t.resolved_at, t.response_sla_breached
		FROM tickets t
		JOIN users c ON t.created_by = c.id
		LEFT JOIN users a ON t.assigned_to = a.id
		WHERE t.id = $1 AND t.tenant_id = $2`

	var t domainTicket.Ticket
	var idCol, tenantIDCol, createdBy pgtype.UUID
	var assignedTo pgtype.UUID
	err := r.conn(ctx).QueryRow(ctx, q, id, tenantID).Scan(
		&idCol, &tenantIDCol, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Category,
		&createdBy, &assignedTo, &t.SLADeadline, &t.SLABreached, &t.Version, &t.CreatedAt, &t.UpdatedAt,
		&t.CreatorName, &t.AssigneeName,
		&t.ResponseDeadline, &t.EscalationDeadline, &t.FirstRespondedAt, &t.ResolvedAt, &t.ResponseSLABreached,
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
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	q := `
		INSERT INTO tickets (id, tenant_id, title, description, priority, category, created_by, sla_deadline, response_deadline, escalation_deadline)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, status, version, created_at, updated_at, sla_deadline, response_deadline, escalation_deadline`

	var id pgtype.UUID
	err := r.conn(ctx).QueryRow(ctx, q, t.ID, t.TenantID, t.Title, t.Description, t.Priority, t.Category, t.CreatedBy, t.SLADeadline, t.ResponseDeadline, t.EscalationDeadline).
		Scan(&id, &t.Status, &t.Version, &t.CreatedAt, &t.UpdatedAt, &t.SLADeadline, &t.ResponseDeadline, &t.EscalationDeadline)
	if err != nil {
		return fmt.Errorf("failed to create ticket: %w", err)
	}
	t.ID = pgToUUID(id)
	return nil
}

func (r *TicketRepo) CreateIdempotent(ctx context.Context, t *domainTicket.Ticket, key uuid.UUID, digest []byte) (bool, error) {
	created := false
	err := withTx(ctx, r.db, func(ctx context.Context) error {
		t.ID = uuid.New()
		conn := r.conn(ctx)
		tag, err := conn.Exec(ctx, `INSERT INTO ticket_idempotency(tenant_id,user_id,operation,key,payload_hash,ticket_id) VALUES($1,$2,'create_ticket',$3,$4,$5) ON CONFLICT DO NOTHING`, t.TenantID, t.CreatedBy, key, digest, t.ID)
		if err != nil {
			return fmt.Errorf("reserve ticket idempotency: %w", err)
		}
		if tag.RowsAffected() == 0 {
			var stored []byte
			var id uuid.UUID
			if err := conn.QueryRow(ctx, `SELECT payload_hash,ticket_id FROM ticket_idempotency WHERE tenant_id=$1 AND user_id=$2 AND operation='create_ticket' AND key=$3`, t.TenantID, t.CreatedBy, key).Scan(&stored, &id); err != nil {
				return err
			}
			if !bytes.Equal(stored, digest) {
				return domainTicket.ErrIdempotencyConflict
			}
			previous, err := r.GetByID(ctx, id, t.TenantID)
			if err != nil {
				return err
			}
			*t = *previous
			return nil
		}
		if err := r.Create(ctx, t); err != nil {
			return err
		}
		created = true
		return nil
	})
	return created, err
}

func (r *TicketRepo) UpdateStatus(ctx context.Context, id, tenantID uuid.UUID, version int, status domainTicket.Status, userID uuid.UUID) (*domainTicket.Ticket, error) {
	var updated *domainTicket.Ticket

	err := withTx(ctx, r.db, func(ctx context.Context) error {
		conn := r.conn(ctx)
		var oldStatus domainTicket.Status
		var currentVersion int
		if err := conn.QueryRow(ctx, `SELECT status,version FROM tickets WHERE id=$1 AND tenant_id=$2 FOR UPDATE`, id, tenantID).Scan(&oldStatus, &currentVersion); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domainTicket.ErrNotFound
			}
			return err
		}
		if currentVersion != version {
			return domainTicket.ErrConcurrentModification
		}
		if oldStatus == status {
			var err error
			updated, err = r.GetByID(ctx, id, tenantID)
			return err
		}
		if !status.IsValid() {
			return domainTicket.ErrValidation
		}
		if !oldStatus.CanTransitionTo(status) {
			return domainTicket.ErrInvalidTransition
		}

		tag, err := conn.Exec(ctx, `
			UPDATE tickets
			SET status = $1::ticket_status, version = version + 1, updated_at = NOW(),
			    resolved_at = CASE WHEN $1::ticket_status = 'resolved' THEN NOW() WHEN $1::ticket_status IN ('open','in_progress','waiting_client') THEN NULL ELSE resolved_at END
			WHERE id = $2 AND tenant_id = $3 AND version = $4`,
			status, id, tenantID, version,
		)
		if err != nil {
			return fmt.Errorf("failed to update ticket status: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return domainTicket.ErrConcurrentModification
		}

		_, err = conn.Exec(ctx, `
			INSERT INTO ticket_history (ticket_id, user_id, field, old_value, new_value)
			VALUES ($1, $2, 'status', $3, $4)`,
			id, userID, string(oldStatus), string(status),
		)
		if err != nil {
			return fmt.Errorf("failed to insert history: %w", err)
		}

		t, err := r.GetByID(ctx, id, tenantID)
		if err != nil {
			return err
		}
		updated = t
		return nil
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// SetSLADeadline updates the SLA deadline for a ticket.
func (r *TicketRepo) SetSLADeadline(ctx context.Context, id, tenantID uuid.UUID, deadline time.Time) error {
	_, err := r.conn(ctx).Exec(ctx, `
		UPDATE tickets SET sla_deadline = $1, updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3`,
		deadline, id, tenantID,
	)
	if err != nil {
		return fmt.Errorf("failed to set SLA deadline: %w", err)
	}
	return nil
}

func (r *TicketRepo) Assign(ctx context.Context, id, tenantID, assignedTo, userID uuid.UUID) error {
	return withTx(ctx, r.db, func(ctx context.Context) error {
		conn := r.conn(ctx)

		var currentAssignee pgtype.UUID
		err := conn.QueryRow(ctx, `SELECT assigned_to FROM tickets WHERE id = $1 AND tenant_id = $2 FOR UPDATE`, id, tenantID).Scan(&currentAssignee)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domainTicket.ErrNotFound
			}
			return fmt.Errorf("failed to fetch ticket: %w", err)
		}
		var assignee any
		newVal := ""
		if assignedTo != uuid.Nil {
			var eligible bool
			if err := conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1 AND tenant_id=$2 AND active=TRUE AND role IN ('agent','admin'))`, assignedTo, tenantID).Scan(&eligible); err != nil {
				return err
			}
			if !eligible {
				return fmt.Errorf("%w: assignee must be active tenant staff", domainTicket.ErrValidation)
			}
			assignee = assignedTo
			newVal = assignedTo.String()
		}
		oldVal := ""
		if currentAssignee.Valid {
			oldVal = pgToUUID(currentAssignee).String()
		}
		if oldVal == newVal {
			return nil
		}

		tag, err := conn.Exec(ctx, `
			UPDATE tickets SET assigned_to = $1, version=version+1, updated_at = NOW()
			WHERE id = $2 AND tenant_id=$3`,
			assignee, id, tenantID)
		if err != nil {
			return fmt.Errorf("failed to assign ticket: %w", err)
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("%w: assignee does not belong to this tenant", domainTicket.ErrValidation)
		}

		if oldVal != newVal {
			_, err = conn.Exec(ctx, `INSERT INTO ticket_history (ticket_id, user_id, field, old_value, new_value) VALUES ($1, $2, 'assigned_to', $3, $4)`, id, userID, oldVal, newVal)
			if err != nil {
				return fmt.Errorf("failed to insert history: %w", err)
			}
		}
		return nil
	})
}

func (r *TicketRepo) GetCreator(ctx context.Context, id, tenantID uuid.UUID) (uuid.UUID, error) {
	var createdBy pgtype.UUID
	err := r.conn(ctx).QueryRow(ctx, `SELECT created_by FROM tickets WHERE id = $1 AND tenant_id = $2`, id, tenantID).Scan(&createdBy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, domainTicket.ErrNotFound
		}
		return uuid.Nil, fmt.Errorf("failed to get ticket creator: %w", err)
	}
	return pgToUUID(createdBy), nil
}

// UpdatePriority preserves the ticket's contracted deadlines while recording the change.
func (r *TicketRepo) UpdatePriority(ctx context.Context, id, tenantID, userID uuid.UUID, value string) error {
	priority, err := domainTicket.ParsePriority(value)
	if err != nil {
		return domainTicket.ErrValidation
	}
	return withTx(ctx, r.db, func(ctx context.Context) error {
		conn := r.conn(ctx)
		var old string
		if err := conn.QueryRow(ctx, `SELECT priority FROM tickets WHERE id=$1 AND tenant_id=$2 FOR UPDATE`, id, tenantID).Scan(&old); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domainTicket.ErrNotFound
			}
			return err
		}
		if old == value {
			return nil
		}
		if _, err := conn.Exec(ctx, `UPDATE tickets SET priority=$1,version=version+1,updated_at=NOW() WHERE id=$2 AND tenant_id=$3`, priority, id, tenantID); err != nil {
			return err
		}
		_, err := conn.Exec(ctx, `INSERT INTO ticket_history(ticket_id,user_id,field,old_value,new_value)VALUES($1,$2,'priority',$3,$4)`, id, userID, old, value)
		return err
	})
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
	return withTx(ctx, r.db, func(ctx context.Context) error {
		conn := DBTX(r.db)
		if tx, ok := TxFromContext(ctx); ok {
			conn = tx
		}
		var role string
		if err := conn.QueryRow(ctx, `SELECT u.name,u.role FROM tickets t JOIN users u ON u.id=$2 AND u.tenant_id=t.tenant_id AND u.active=TRUE WHERE t.id=$1 FOR UPDATE OF t`, c.TicketID, c.UserID).Scan(&c.UserName, &role); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domainTicket.ErrNotFound
			}
			return err
		}
		q := `
		INSERT INTO ticket_comments (ticket_id, user_id, content, is_internal)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

		var id pgtype.UUID
		err := conn.QueryRow(ctx, q, c.TicketID, c.UserID, c.Content, c.IsInternal).Scan(&id, &c.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to create comment: %w", err)
		}
		c.ID = pgToUUID(id)

		if !c.IsInternal && (role == "admin" || role == "agent") {
			if _, err := conn.Exec(ctx, `UPDATE tickets SET first_responded_at=COALESCE(first_responded_at,$2),updated_at=NOW() WHERE id=$1`, c.TicketID, c.CreatedAt); err != nil {
				return err
			}
		}
		return nil
	})
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
