package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/trickreport/backend/internal/application/automation"
	"github.com/trickreport/backend/internal/application/sla"
	domainautomation "github.com/trickreport/backend/internal/domain/automation"
	domainsla "github.com/trickreport/backend/internal/domain/sla"
)

// ---------------------------------------------------------------------------
// SLARepo
// ---------------------------------------------------------------------------

// SLARepo implements sla.Repository.
type SLARepo struct {
	db *pgxpool.Pool
}

// NewSLARepo creates a new SLARepo.
func NewSLARepo(db *pgxpool.Pool) *SLARepo {
	return &SLARepo{db: db}
}

// Compile-time assertion that SLARepo implements sla.Repository.
var _ sla.Repository = (*SLARepo)(nil)

// List returns all SLA policies for a tenant, ordered by priority.
func (r *SLARepo) List(ctx context.Context, tenantID uuid.UUID) ([]domainsla.Policy, error) {
	const q = `SELECT id, tenant_id, priority, response_time_minutes, resolution_time_minutes, escalation_minutes, created_at, updated_at FROM sla_policies WHERE tenant_id = $1 ORDER BY priority`

	rows, err := r.db.Query(ctx, q, tenantID)
	if err != nil {
		return nil, fmt.Errorf("sla_repo.List: query: %w", err)
	}
	defer rows.Close()

	var policies []domainsla.Policy
	for rows.Next() {
		var p domainsla.Policy
		var id, tid pgtype.UUID
		if err := rows.Scan(&id, &tid, &p.Priority, &p.ResponseTimeMinutes, &p.ResolutionTimeMinutes, &p.EscalationMinutes, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("sla_repo.List: scan: %w", err)
		}
		p.ID = pgToUUID(id)
		p.TenantID = pgToUUID(tid)
		policies = append(policies, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sla_repo.List: rows: %w", err)
	}
	return policies, nil
}

// Upsert inserts or updates an SLA policy on (tenant_id, priority) conflict.
func (r *SLARepo) Upsert(ctx context.Context, p *domainsla.Policy) error {
	const q = `INSERT INTO sla_policies (tenant_id, priority, response_time_minutes, resolution_time_minutes, escalation_minutes) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (tenant_id, priority) DO UPDATE SET response_time_minutes = EXCLUDED.response_time_minutes, resolution_time_minutes = EXCLUDED.resolution_time_minutes, escalation_minutes = EXCLUDED.escalation_minutes, updated_at = NOW() RETURNING id, created_at, updated_at`

	var id pgtype.UUID
	if err := r.db.QueryRow(ctx, q, p.TenantID, p.Priority, p.ResponseTimeMinutes, p.ResolutionTimeMinutes, p.EscalationMinutes).
		Scan(&id, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return fmt.Errorf("sla_repo.Upsert: %w", err)
	}
	p.ID = pgToUUID(id)
	return nil
}

// ---------------------------------------------------------------------------
// AutomationRepo
// ---------------------------------------------------------------------------

// AutomationRepo implements automation.Repository.
type AutomationRepo struct {
	db *pgxpool.Pool
}

// NewAutomationRepo creates a new AutomationRepo.
func NewAutomationRepo(db *pgxpool.Pool) *AutomationRepo {
	return &AutomationRepo{db: db}
}

// DB returns the underlying connection pool (used by the automation executor).
func (r *AutomationRepo) DB() *pgxpool.Pool {
	return r.db
}

// Compile-time assertion that AutomationRepo implements automation.Repository.
var _ automation.Repository = (*AutomationRepo)(nil)

// List returns all automation rules for a tenant, ordered by created_at desc.
func (r *AutomationRepo) List(ctx context.Context, tenantID uuid.UUID) ([]domainautomation.Rule, error) {
	const q = `SELECT id, tenant_id, name, description, trigger_type, conditions, actions, is_active, created_at, updated_at FROM automation_rules WHERE tenant_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, q, tenantID)
	if err != nil {
		return nil, fmt.Errorf("automation_repo.List: query: %w", err)
	}
	defer rows.Close()

	var rules []domainautomation.Rule
	for rows.Next() {
		rl, err := scanRule(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, *rl)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("automation_repo.List: rows: %w", err)
	}
	return rules, nil
}

// Create inserts a new automation rule.
func (r *AutomationRepo) Create(ctx context.Context, rl *domainautomation.Rule) error {
	const q = `INSERT INTO automation_rules (tenant_id, name, description, trigger_type, conditions, actions, is_active) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`

	var id pgtype.UUID
	if err := r.db.QueryRow(ctx, q, rl.TenantID, rl.Name, rl.Description, rl.TriggerType, rl.Conditions, rl.Actions, rl.IsActive).
		Scan(&id, &rl.CreatedAt, &rl.UpdatedAt); err != nil {
		return fmt.Errorf("automation_repo.Create: %w", err)
	}
	rl.ID = pgToUUID(id)
	return nil
}

// Update updates an existing automation rule.
func (r *AutomationRepo) Update(ctx context.Context, rl *domainautomation.Rule) error {
	const q = `UPDATE automation_rules SET name = $1, description = $2, trigger_type = $3, conditions = $4, actions = $5, is_active = $6, updated_at = NOW() WHERE id = $7 AND tenant_id = $8 RETURNING id, created_at, updated_at`

	var id pgtype.UUID
	if err := r.db.QueryRow(ctx, q, rl.Name, rl.Description, rl.TriggerType, rl.Conditions, rl.Actions, rl.IsActive, rl.ID, rl.TenantID).
		Scan(&id, &rl.CreatedAt, &rl.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainautomation.ErrNotFound
		}
		return fmt.Errorf("automation_repo.Update: %w", err)
	}
	rl.ID = pgToUUID(id)
	return nil
}

// Delete removes an automation rule by id within a tenant.
func (r *AutomationRepo) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
	const q = `DELETE FROM automation_rules WHERE id = $1 AND tenant_id = $2`

	ct, err := r.db.Exec(ctx, q, id, tenantID)
	if err != nil {
		return fmt.Errorf("automation_repo.Delete: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return domainautomation.ErrNotFound
	}
	return nil
}

// scanRule scans an automation rule from a pgx.Row-like scanner.
func scanRule(row pgx.Row) (*domainautomation.Rule, error) {
	var rl domainautomation.Rule
	var id, tid pgtype.UUID
	if err := row.Scan(&id, &tid, &rl.Name, &rl.Description, &rl.TriggerType, &rl.Conditions, &rl.Actions, &rl.IsActive, &rl.CreatedAt, &rl.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainautomation.ErrNotFound
		}
		return nil, fmt.Errorf("automation_repo: scan: %w", err)
	}
	rl.ID = pgToUUID(id)
	rl.TenantID = pgToUUID(tid)
	return &rl, nil
}
