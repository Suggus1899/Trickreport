package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	appAuto "github.com/trickreport/backend/internal/application/automation"
	domainAuto "github.com/trickreport/backend/internal/domain/automation"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
	"strings"
)

// AutomationExecutor uses the same transactional ticket mutations as HTTP callers.
type AutomationExecutor struct {
	db           *pgxpool.Pool
	systemUserID uuid.UUID
}

func NewAutomationExecutor(db *pgxpool.Pool) *AutomationExecutor {
	return &AutomationExecutor{db: db, systemUserID: uuid.MustParse("00000000-0000-0000-0000-000000000002")}
}
func (e *AutomationExecutor) WithSystemUserID(id uuid.UUID) *AutomationExecutor {
	e.systemUserID = id
	return e
}

var _ appAuto.TicketExecutor = (*AutomationExecutor)(nil)

func (e *AutomationExecutor) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	return withTx(ctx, e.db, fn)
}
func (e *AutomationExecutor) GetTicket(ctx context.Context, tenantID, ticketID uuid.UUID) (*appAuto.TicketSnapshot, error) {
	t, err := NewTicketRepo(e.db).GetByID(ctx, ticketID, tenantID)
	if err != nil {
		return nil, err
	}
	return &appAuto.TicketSnapshot{ID: t.ID, Status: string(t.Status), Priority: string(t.Priority), Category: t.Category, Title: t.Title}, nil
}
func (e *AutomationExecutor) SetPriority(ctx context.Context, tenantID, ticketID uuid.UUID, priority string) error {
	return NewTicketRepo(e.db).UpdatePriority(ctx, ticketID, tenantID, e.systemUserID, priority)
}
func (e *AutomationExecutor) SetStatus(ctx context.Context, tenantID, ticketID uuid.UUID, value string) error {
	status, err := domainTicket.ParseStatus(value)
	if err != nil {
		return domainTicket.ErrValidation
	}
	repo := NewTicketRepo(e.db)
	return withTx(ctx, e.db, func(ctx context.Context) error {
		var version int
		if err := repo.conn(ctx).QueryRow(ctx, `SELECT version FROM tickets WHERE id=$1 AND tenant_id=$2 FOR UPDATE`, ticketID, tenantID).Scan(&version); err != nil {
			return err
		}
		_, err := repo.UpdateStatus(ctx, ticketID, tenantID, version, status, e.systemUserID)
		return err
	})
}
func (e *AutomationExecutor) AssignTo(ctx context.Context, tenantID, ticketID, userID uuid.UUID) error {
	return NewTicketRepo(e.db).Assign(ctx, ticketID, tenantID, userID, e.systemUserID)
}
func (e *AutomationExecutor) AddComment(ctx context.Context, tenantID, ticketID uuid.UUID, content string, isInternal bool) error {
	if strings.TrimSpace(content) == "" {
		return domainTicket.ErrValidation
	}
	return withTx(ctx, e.db, func(ctx context.Context) error {
		if _, err := NewTicketRepo(e.db).GetByID(ctx, ticketID, tenantID); err != nil {
			return err
		}
		return NewCommentRepo(e.db).Create(ctx, &domainTicket.Comment{TicketID: ticketID, UserID: e.systemUserID, Content: content, IsInternal: isInternal})
	})
}
func (e *AutomationExecutor) AddTag(_ context.Context, _, _ uuid.UUID, _ string) error {
	return fmt.Errorf("%w: add_tag is not supported", domainAuto.ErrValidation)
}
