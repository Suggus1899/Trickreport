package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	appAuto "github.com/trickreport/backend/internal/application/automation"
)

// AutomationExecutor implements appAuto.TicketExecutor using direct SQL.
// It is used by the automation engine to execute actions on tickets.
type AutomationExecutor struct {
	db           *pgxpool.Pool
	systemUserID uuid.UUID
}

// NewAutomationExecutor creates a new AutomationExecutor.
func NewAutomationExecutor(db *pgxpool.Pool) *AutomationExecutor {
	return &AutomationExecutor{
		db:           db,
		systemUserID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
	}
}

// WithSystemUserID sets the user ID used for system-generated comments and history.
func (e *AutomationExecutor) WithSystemUserID(id uuid.UUID) *AutomationExecutor {
	e.systemUserID = id
	return e
}

// Compile-time assertion that AutomationExecutor implements appAuto.TicketExecutor.
var _ appAuto.TicketExecutor = (*AutomationExecutor)(nil)

// GetTicket returns a snapshot of the ticket for condition evaluation.
func (e *AutomationExecutor) GetTicket(ctx context.Context, tenantID, ticketID uuid.UUID) (*appAuto.TicketSnapshot, error) {
	const q = `SELECT id, status, priority, category, title FROM tickets WHERE id = $1 AND tenant_id = $2`

	var s appAuto.TicketSnapshot
	var id pgtype.UUID
	if err := e.db.QueryRow(ctx, q, ticketID, tenantID).Scan(&id, &s.Status, &s.Priority, &s.Category, &s.Title); err != nil {
		return nil, fmt.Errorf("automation_executor.GetTicket: %w", err)
	}
	s.ID = pgToUUID(id)
	return &s, nil
}

// SetPriority updates the ticket priority and records a history entry.
func (e *AutomationExecutor) SetPriority(ctx context.Context, tenantID, ticketID uuid.UUID, priority string) error {
	tx, err := e.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("automation_executor.SetPriority: begin: %w", err)
	}
	defer tx.Rollback(ctx)

	var oldPriority string
	if err := tx.QueryRow(ctx, `SELECT priority FROM tickets WHERE id = $1 AND tenant_id = $2`, ticketID, tenantID).Scan(&oldPriority); err != nil {
		return fmt.Errorf("automation_executor.SetPriority: select: %w", err)
	}

	if _, err := tx.Exec(ctx, `UPDATE tickets SET priority = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3`, priority, ticketID, tenantID); err != nil {
		return fmt.Errorf("automation_executor.SetPriority: update: %w", err)
	}

	if _, err := tx.Exec(ctx, `INSERT INTO ticket_history (ticket_id, user_id, field, old_value, new_value) VALUES ($1, $2, 'priority', $3, $4)`, ticketID, e.systemUserID, oldPriority, priority); err != nil {
		return fmt.Errorf("automation_executor.SetPriority: history: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("automation_executor.SetPriority: commit: %w", err)
	}
	return nil
}

// SetStatus updates the ticket status and records a history entry.
func (e *AutomationExecutor) SetStatus(ctx context.Context, tenantID, ticketID uuid.UUID, status string) error {
	tx, err := e.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("automation_executor.SetStatus: begin: %w", err)
	}
	defer tx.Rollback(ctx)

	var oldStatus string
	if err := tx.QueryRow(ctx, `SELECT status FROM tickets WHERE id = $1 AND tenant_id = $2`, ticketID, tenantID).Scan(&oldStatus); err != nil {
		return fmt.Errorf("automation_executor.SetStatus: select: %w", err)
	}

	if _, err := tx.Exec(ctx, `UPDATE tickets SET status = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3`, status, ticketID, tenantID); err != nil {
		return fmt.Errorf("automation_executor.SetStatus: update: %w", err)
	}

	if _, err := tx.Exec(ctx, `INSERT INTO ticket_history (ticket_id, user_id, field, old_value, new_value) VALUES ($1, $2, 'status', $3, $4)`, ticketID, e.systemUserID, oldStatus, status); err != nil {
		return fmt.Errorf("automation_executor.SetStatus: history: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("automation_executor.SetStatus: commit: %w", err)
	}
	return nil
}

// AssignTo assigns the ticket to a user and records a history entry.
func (e *AutomationExecutor) AssignTo(ctx context.Context, tenantID, ticketID, userID uuid.UUID) error {
	tx, err := e.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("automation_executor.AssignTo: begin: %w", err)
	}
	defer tx.Rollback(ctx)

	var currentAssignee pgtype.UUID
	if err := tx.QueryRow(ctx, `SELECT assigned_to FROM tickets WHERE id = $1 AND tenant_id = $2`, ticketID, tenantID).Scan(&currentAssignee); err != nil {
		return fmt.Errorf("automation_executor.AssignTo: select: %w", err)
	}

	if _, err := tx.Exec(ctx, `UPDATE tickets SET assigned_to = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3`, userID, ticketID, tenantID); err != nil {
		return fmt.Errorf("automation_executor.AssignTo: update: %w", err)
	}

	oldVal := ""
	if currentAssignee.Valid {
		oldVal = pgToUUID(currentAssignee).String()
	}
	if _, err := tx.Exec(ctx, `INSERT INTO ticket_history (ticket_id, user_id, field, old_value, new_value) VALUES ($1, $2, 'assigned_to', $3, $4)`, ticketID, e.systemUserID, oldVal, userID.String()); err != nil {
		return fmt.Errorf("automation_executor.AssignTo: history: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("automation_executor.AssignTo: commit: %w", err)
	}
	return nil
}

// AddComment adds an internal comment to the ticket.
func (e *AutomationExecutor) AddComment(ctx context.Context, tenantID, ticketID uuid.UUID, content string, isInternal bool) error {
	const q = `INSERT INTO ticket_comments (ticket_id, user_id, content, is_internal) VALUES ($1, $2, $3, $4)`
	if _, err := e.db.Exec(ctx, q, ticketID, e.systemUserID, content, isInternal); err != nil {
		return fmt.Errorf("automation_executor.AddComment: %w", err)
	}
	return nil
}

// AddTag adds a tag to the ticket. The tickets table does not currently have a
// tags column, so this is a no-op that logs a warning.
func (e *AutomationExecutor) AddTag(ctx context.Context, tenantID, ticketID uuid.UUID, tag string) error {
	log.Warn().
		Str("ticket_id", ticketID.String()).
		Str("tag", tag).
		Msg("automation_executor.AddTag: tickets table has no tags column, skipping")
	return nil
}
