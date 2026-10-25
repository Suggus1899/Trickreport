package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	appITSM "github.com/trickreport/backend/internal/application/itsm"
	domain "github.com/trickreport/backend/internal/domain/itsm"
	domainSLA "github.com/trickreport/backend/internal/domain/sla"
)

// ITSMRepo implements the ITSM repository boundary over PostgreSQL.
// Every query filters by tenant_id: a record from another organization
// is indistinguishable from a missing record.
type ITSMRepo struct {
	db *pgxpool.Pool
}

// NewITSMRepo creates an ITSMRepo.
func NewITSMRepo(db *pgxpool.Pool) *ITSMRepo { return &ITSMRepo{db: db} }

// Compile-time assertion that ITSMRepo implements the application boundary.
var _ appITSM.Repository = (*ITSMRepo)(nil)

// DB exposes the pool for wiring that needs it (mirrors other repos).
func (r *ITSMRepo) DB() *pgxpool.Pool { return r.db }

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func isFKViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}

func toPgUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

func toPgTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func fromPgTimestamptz(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func nonNil[T any](rows []T) []T {
	if rows == nil {
		return []T{}
	}
	return rows
}

func mapAttrs(raw []byte) map[string]any {
	if len(raw) == 0 {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return map[string]any{}
	}
	if out == nil {
		return map[string]any{}
	}
	return out
}

func marshalAttrs(attrs map[string]any) []byte {
	if attrs == nil {
		return []byte("{}")
	}
	raw, err := json.Marshal(attrs)
	if err != nil {
		return []byte("{}")
	}
	return raw
}

// ── Catalog ─────────────────────────────────────────────────────────

func scanCategoryRow(row pgx.Row) (*domain.CatalogCategory, error) {
	var c domain.CatalogCategory
	err := row.Scan(&c.ID, &c.TenantID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ITSMRepo) ListCategories(ctx context.Context, tenantID uuid.UUID) ([]domain.CatalogCategory, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, name, description, created_at, updated_at FROM catalog_categories WHERE tenant_id = $1 ORDER BY name`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.ListCategories: %w", err)
	}
	defer rows.Close()
	var out []domain.CatalogCategory
	for rows.Next() {
		c, err := scanCategoryRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return nonNil(out), rows.Err()
}

func (r *ITSMRepo) CreateCategory(ctx context.Context, c *domain.CatalogCategory) error {
	err := r.db.QueryRow(ctx, `INSERT INTO catalog_categories (tenant_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
		c.TenantID, c.Name, c.Description).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("itsm_repo.CreateCategory: %w", err)
	}
	return nil
}

func (r *ITSMRepo) UpdateCategory(ctx context.Context, c *domain.CatalogCategory) error {
	tag, err := r.db.Exec(ctx, `UPDATE catalog_categories SET name = $1, description = $2, updated_at = NOW() WHERE id = $3 AND tenant_id = $4`,
		c.Name, c.Description, c.ID, c.TenantID)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("itsm_repo.UpdateCategory: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ITSMRepo) DeleteCategory(ctx context.Context, id, tenantID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM catalog_categories WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.DeleteCategory: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

const serviceSelect = `SELECT s.id, s.tenant_id, s.category_id, COALESCE(cc.name, ''), s.name, s.description, s.requires_approval, s.active, s.created_at, s.updated_at FROM catalog_services s LEFT JOIN catalog_categories cc ON cc.id = s.category_id`

func scanServiceRow(row pgx.Row) (*domain.CatalogService, error) {
	var s domain.CatalogService
	var categoryID pgtype.UUID
	err := row.Scan(&s.ID, &s.TenantID, &categoryID, &s.CategoryName, &s.Name, &s.Description, &s.RequiresApproval, &s.Active, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	s.CategoryID = pgToUUIDPtr(categoryID)
	return &s, nil
}

func (r *ITSMRepo) ListServices(ctx context.Context, tenantID uuid.UUID, categoryID *uuid.UUID, activeOnly bool) ([]domain.CatalogService, error) {
	q := serviceSelect + ` WHERE s.tenant_id = $1`
	args := []any{tenantID}
	if categoryID != nil {
		q += ` AND s.category_id = $2`
		args = append(args, *categoryID)
	}
	if activeOnly {
		q += ` AND s.active = TRUE`
	}
	q += ` ORDER BY s.name`
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.ListServices: %w", err)
	}
	defer rows.Close()
	var out []domain.CatalogService
	for rows.Next() {
		s, err := scanServiceRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *s)
	}
	return nonNil(out), rows.Err()
}

func (r *ITSMRepo) GetService(ctx context.Context, id, tenantID uuid.UUID) (*domain.CatalogService, error) {
	s, err := scanServiceRow(r.db.QueryRow(ctx, serviceSelect+` WHERE s.id = $1 AND s.tenant_id = $2`, id, tenantID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("itsm_repo.GetService: %w", err)
	}
	return s, nil
}

func (r *ITSMRepo) CreateService(ctx context.Context, s *domain.CatalogService) error {
	err := r.db.QueryRow(ctx, `INSERT INTO catalog_services (tenant_id, category_id, name, description, requires_approval, active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`,
		s.TenantID, toPgUUID(s.CategoryID), s.Name, s.Description, s.RequiresApproval, s.Active).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		if isFKViolation(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("itsm_repo.CreateService: %w", err)
	}
	return nil
}

func (r *ITSMRepo) UpdateService(ctx context.Context, s *domain.CatalogService) error {
	tag, err := r.db.Exec(ctx, `UPDATE catalog_services SET category_id = $1, name = $2, description = $3, requires_approval = $4, active = $5, updated_at = NOW() WHERE id = $6 AND tenant_id = $7`,
		toPgUUID(s.CategoryID), s.Name, s.Description, s.RequiresApproval, s.Active, s.ID, s.TenantID)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		if isFKViolation(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("itsm_repo.UpdateService: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ITSMRepo) DeleteService(ctx context.Context, id, tenantID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM catalog_services WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.DeleteService: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ── Tickets / users (minimal surface for ITSM links) ────────────────

func (r *ITSMRepo) CreateTicketFromService(ctx context.Context, tenantID, requesterID uuid.UUID, title, description, category, priority string, dl domainSLA.Deadlines) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx, `INSERT INTO tickets (tenant_id, title, description, status, priority, category, created_by, sla_deadline, response_deadline, escalation_deadline) VALUES ($1, $2, $3, 'open', $4, $5, $6, $7, $8, $9) RETURNING id`,
		tenantID, title, description, priority, category, requesterID, dl.Resolution, dl.Response, dl.Escalation).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("itsm_repo.CreateTicketFromService: %w", err)
	}
	return id, nil
}

func (r *ITSMRepo) TicketExists(ctx context.Context, id, tenantID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM tickets WHERE id = $1 AND tenant_id = $2)`, id, tenantID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("itsm_repo.TicketExists: %w", err)
	}
	return exists, nil
}

func (r *ITSMRepo) UserExists(ctx context.Context, id, tenantID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND tenant_id = $2)`, id, tenantID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("itsm_repo.UserExists: %w", err)
	}
	return exists, nil
}

// ── Approvals ───────────────────────────────────────────────────────

func scanApprovalRow(row pgx.Row) (*domain.Approval, error) {
	var a domain.Approval
	var approverID, decidedBy pgtype.UUID
	var decidedAt pgtype.Timestamptz
	err := row.Scan(&a.ID, &a.TenantID, &a.TargetType, &a.TargetID, &a.RequesterID, &approverID, &a.Status, &a.Comment, &decidedBy, &decidedAt, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	a.ApproverID = pgToUUIDPtr(approverID)
	a.DecidedBy = pgToUUIDPtr(decidedBy)
	a.DecidedAt = fromPgTimestamptz(decidedAt)
	return &a, nil
}

const approvalSelect = `SELECT id, tenant_id, target_type, target_id, requester_id, approver_id, status, comment, decided_by, decided_at, created_at, updated_at FROM approvals`

func (r *ITSMRepo) ListApprovals(ctx context.Context, tenantID, viewerID uuid.UUID, viewerRole string) ([]domain.Approval, error) {
	q := approvalSelect + ` WHERE tenant_id = $1`
	args := []any{tenantID}
	if viewerRole != "admin" && viewerRole != "agent" {
		q += ` AND (requester_id = $2 OR approver_id = $2)`
		args = append(args, viewerID)
	}
	q += ` ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.ListApprovals: %w", err)
	}
	defer rows.Close()
	var out []domain.Approval
	for rows.Next() {
		a, err := scanApprovalRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return nonNil(out), rows.Err()
}

func (r *ITSMRepo) GetApproval(ctx context.Context, id, tenantID uuid.UUID) (*domain.Approval, error) {
	a, err := scanApprovalRow(r.db.QueryRow(ctx, approvalSelect+` WHERE id = $1 AND tenant_id = $2`, id, tenantID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("itsm_repo.GetApproval: %w", err)
	}
	return a, nil
}

func (r *ITSMRepo) CreateApproval(ctx context.Context, a *domain.Approval) error {
	err := r.db.QueryRow(ctx, `INSERT INTO approvals (tenant_id, target_type, target_id, requester_id, approver_id, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`,
		a.TenantID, a.TargetType, a.TargetID, a.RequesterID, toPgUUID(a.ApproverID), a.Status).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("itsm_repo.CreateApproval: %w", err)
	}
	return nil
}

// DecideApproval updates the approval and appends a history event in a
// single transaction, and only while the row is still pending, so two
// concurrent decisions cannot both succeed.
func (r *ITSMRepo) DecideApproval(ctx context.Context, a *domain.Approval, actorID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("itsm_repo.DecideApproval: begin: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	tag, err := tx.Exec(ctx, `UPDATE approvals SET status = $1, comment = $2, decided_by = $3, decided_at = $4, updated_at = NOW() WHERE id = $5 AND tenant_id = $6 AND status = 'pending'`,
		a.Status, a.Comment, actorID, a.DecidedAt, a.ID, a.TenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.DecideApproval: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrInvalidTransition
	}
	if _, err := tx.Exec(ctx, `INSERT INTO approval_history (approval_id, actor_id, from_status, to_status, comment) VALUES ($1, $2, 'pending', $3, $4)`,
		a.ID, actorID, a.Status, a.Comment); err != nil {
		return fmt.Errorf("itsm_repo.DecideApproval: history: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("itsm_repo.DecideApproval: commit: %w", err)
	}
	return nil
}

func (r *ITSMRepo) ApprovalHistory(ctx context.Context, id, tenantID uuid.UUID) ([]domain.ApprovalEvent, error) {
	rows, err := r.db.Query(ctx, `SELECT h.id, h.approval_id, h.actor_id, h.from_status, h.to_status, h.comment, h.created_at FROM approval_history h JOIN approvals a ON a.id = h.approval_id WHERE h.approval_id = $1 AND a.tenant_id = $2 ORDER BY h.created_at`,
		id, tenantID)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.ApprovalHistory: %w", err)
	}
	defer rows.Close()
	var out []domain.ApprovalEvent
	for rows.Next() {
		var ev domain.ApprovalEvent
		if err := rows.Scan(&ev.ID, &ev.ApprovalID, &ev.ActorID, &ev.FromStatus, &ev.ToStatus, &ev.Comment, &ev.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, ev)
	}
	return nonNil(out), rows.Err()
}

func (r *ITSMRepo) HasPendingApproval(ctx context.Context, tenantID uuid.UUID, targetType string, targetID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM approvals WHERE tenant_id = $1 AND target_type = $2 AND target_id = $3 AND status = 'pending')`,
		tenantID, targetType, targetID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("itsm_repo.HasPendingApproval: %w", err)
	}
	return exists, nil
}

func (r *ITSMRepo) HasApprovedApproval(ctx context.Context, tenantID uuid.UUID, targetType string, targetID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM approvals WHERE tenant_id = $1 AND target_type = $2 AND target_id = $3 AND status = 'approved')`,
		tenantID, targetType, targetID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("itsm_repo.HasApprovedApproval: %w", err)
	}
	return exists, nil
}

// ── Assets / CMDB ───────────────────────────────────────────────────

func (r *ITSMRepo) ListCITypes(ctx context.Context, tenantID uuid.UUID) ([]domain.CIType, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, name, description, created_at, updated_at FROM ci_types WHERE tenant_id = $1 ORDER BY name`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.ListCITypes: %w", err)
	}
	defer rows.Close()
	var out []domain.CIType
	for rows.Next() {
		var ct domain.CIType
		if err := rows.Scan(&ct.ID, &ct.TenantID, &ct.Name, &ct.Description, &ct.CreatedAt, &ct.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, ct)
	}
	return nonNil(out), rows.Err()
}

func (r *ITSMRepo) CreateCIType(ctx context.Context, ct *domain.CIType) error {
	err := r.db.QueryRow(ctx, `INSERT INTO ci_types (tenant_id, name, description) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`,
		ct.TenantID, ct.Name, ct.Description).Scan(&ct.ID, &ct.CreatedAt, &ct.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("itsm_repo.CreateCIType: %w", err)
	}
	return nil
}

func (r *ITSMRepo) DeleteCIType(ctx context.Context, id, tenantID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM ci_types WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		if isFKViolation(err) {
			// RESTRICT: assets still reference this type.
			return domain.ErrConflict
		}
		return fmt.Errorf("itsm_repo.DeleteCIType: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

const assetSelect = `SELECT a.id, a.tenant_id, a.ci_type_id, COALESCE(ct.name, ''), a.name, a.serial_number, a.status, a.attributes, a.owner_id, COALESCE(u.name, ''), a.created_at, a.updated_at FROM assets a LEFT JOIN ci_types ct ON ct.id = a.ci_type_id LEFT JOIN users u ON u.id = a.owner_id`

func scanAssetRow(row pgx.Row) (*domain.Asset, error) {
	var a domain.Asset
	var rawAttrs []byte
	var ownerID pgtype.UUID
	err := row.Scan(&a.ID, &a.TenantID, &a.CITypeID, &a.CITypeName, &a.Name, &a.SerialNumber, &a.Status, &rawAttrs, &ownerID, &a.OwnerName, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	a.Attributes = mapAttrs(rawAttrs)
	a.OwnerID = pgToUUIDPtr(ownerID)
	return &a, nil
}

func (r *ITSMRepo) ListAssets(ctx context.Context, tenantID uuid.UUID, ciTypeID *uuid.UUID, status string) ([]domain.Asset, error) {
	q := assetSelect + ` WHERE a.tenant_id = $1`
	args := []any{tenantID}
	idx := 2
	if ciTypeID != nil {
		q += fmt.Sprintf(` AND a.ci_type_id = $%d`, idx)
		args = append(args, *ciTypeID)
		idx++
	}
	if status != "" {
		q += fmt.Sprintf(` AND a.status = $%d`, idx)
		args = append(args, status)
	}
	q += ` ORDER BY a.name`
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.ListAssets: %w", err)
	}
	defer rows.Close()
	var out []domain.Asset
	for rows.Next() {
		a, err := scanAssetRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return nonNil(out), rows.Err()
}

func (r *ITSMRepo) GetAsset(ctx context.Context, id, tenantID uuid.UUID) (*domain.Asset, error) {
	a, err := scanAssetRow(r.db.QueryRow(ctx, assetSelect+` WHERE a.id = $1 AND a.tenant_id = $2`, id, tenantID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("itsm_repo.GetAsset: %w", err)
	}
	return a, nil
}

func (r *ITSMRepo) CreateAsset(ctx context.Context, a *domain.Asset) error {
	err := r.db.QueryRow(ctx, `INSERT INTO assets (tenant_id, ci_type_id, name, serial_number, status, attributes, owner_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		a.TenantID, a.CITypeID, a.Name, a.SerialNumber, a.Status, marshalAttrs(a.Attributes), toPgUUID(a.OwnerID)).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		if isFKViolation(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("itsm_repo.CreateAsset: %w", err)
	}
	return nil
}

func (r *ITSMRepo) UpdateAsset(ctx context.Context, a *domain.Asset) error {
	tag, err := r.db.Exec(ctx, `UPDATE assets SET ci_type_id = $1, name = $2, serial_number = $3, status = $4, attributes = $5, owner_id = $6, updated_at = NOW() WHERE id = $7 AND tenant_id = $8`,
		a.CITypeID, a.Name, a.SerialNumber, a.Status, marshalAttrs(a.Attributes), toPgUUID(a.OwnerID), a.ID, a.TenantID)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		if isFKViolation(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("itsm_repo.UpdateAsset: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ITSMRepo) DeleteAsset(ctx context.Context, id, tenantID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM assets WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.DeleteAsset: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ITSMRepo) ListRelationships(ctx context.Context, tenantID, assetID uuid.UUID) ([]domain.CIRelationship, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, source_ci_id, target_ci_id, relationship_type, created_at FROM ci_relationships WHERE tenant_id = $1 AND (source_ci_id = $2 OR target_ci_id = $2) ORDER BY created_at`,
		tenantID, assetID)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.ListRelationships: %w", err)
	}
	defer rows.Close()
	var out []domain.CIRelationship
	for rows.Next() {
		var rel domain.CIRelationship
		if err := rows.Scan(&rel.ID, &rel.TenantID, &rel.SourceCIID, &rel.TargetCIID, &rel.RelationshipType, &rel.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, rel)
	}
	return nonNil(out), rows.Err()
}

func (r *ITSMRepo) CreateRelationship(ctx context.Context, rel *domain.CIRelationship) error {
	// Both endpoints must exist in this tenant; otherwise the insert
	// would silently link foreign records.
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM assets WHERE tenant_id = $1 AND id IN ($2, $3)`, rel.TenantID, rel.SourceCIID, rel.TargetCIID).Scan(&count)
	if err != nil {
		return fmt.Errorf("itsm_repo.CreateRelationship: check: %w", err)
	}
	if count != 2 {
		return domain.ErrNotFound
	}
	err = r.db.QueryRow(ctx, `INSERT INTO ci_relationships (tenant_id, source_ci_id, target_ci_id, relationship_type) VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		rel.TenantID, rel.SourceCIID, rel.TargetCIID, rel.RelationshipType).Scan(&rel.ID, &rel.CreatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrConflict
		}
		return fmt.Errorf("itsm_repo.CreateRelationship: %w", err)
	}
	return nil
}

func (r *ITSMRepo) DeleteRelationship(ctx context.Context, id, tenantID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM ci_relationships WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.DeleteRelationship: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ITSMRepo) LinkTicketAsset(ctx context.Context, tenantID, ticketID, assetID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `INSERT INTO ticket_assets (ticket_id, asset_id, tenant_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, ticketID, assetID, tenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.LinkTicketAsset: %w", err)
	}
	return nil
}

func (r *ITSMRepo) UnlinkTicketAsset(ctx context.Context, tenantID, ticketID, assetID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM ticket_assets WHERE ticket_id = $1 AND asset_id = $2 AND tenant_id = $3`, ticketID, assetID, tenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.UnlinkTicketAsset: %w", err)
	}
	return nil
}

func (r *ITSMRepo) AssetsForTicket(ctx context.Context, tenantID, ticketID uuid.UUID) ([]domain.Asset, error) {
	rows, err := r.db.Query(ctx, assetSelect+` JOIN ticket_assets ta ON ta.asset_id = a.id WHERE ta.ticket_id = $1 AND ta.tenant_id = $2 AND a.tenant_id = $2 ORDER BY a.name`, ticketID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.AssetsForTicket: %w", err)
	}
	defer rows.Close()
	var out []domain.Asset
	for rows.Next() {
		a, err := scanAssetRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return nonNil(out), rows.Err()
}

// ── Problems ────────────────────────────────────────────────────────

func scanProblemRow(row pgx.Row) (*domain.Problem, error) {
	var p domain.Problem
	err := row.Scan(&p.ID, &p.TenantID, &p.Title, &p.Description, &p.Status, &p.RootCause, &p.Workaround, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

const problemSelect = `SELECT id, tenant_id, title, description, status, root_cause, workaround, created_by, created_at, updated_at FROM problems`

func (r *ITSMRepo) ListProblems(ctx context.Context, tenantID uuid.UUID, status string) ([]domain.Problem, error) {
	q := problemSelect + ` WHERE tenant_id = $1`
	args := []any{tenantID}
	if status != "" {
		q += ` AND status = $2`
		args = append(args, status)
	}
	q += ` ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.ListProblems: %w", err)
	}
	defer rows.Close()
	var out []domain.Problem
	for rows.Next() {
		p, err := scanProblemRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return nonNil(out), rows.Err()
}

func (r *ITSMRepo) GetProblem(ctx context.Context, id, tenantID uuid.UUID) (*domain.Problem, error) {
	p, err := scanProblemRow(r.db.QueryRow(ctx, problemSelect+` WHERE id = $1 AND tenant_id = $2`, id, tenantID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("itsm_repo.GetProblem: %w", err)
	}
	return p, nil
}

func (r *ITSMRepo) CreateProblem(ctx context.Context, p *domain.Problem) error {
	err := r.db.QueryRow(ctx, `INSERT INTO problems (tenant_id, title, description, status, root_cause, workaround, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		p.TenantID, p.Title, p.Description, p.Status, p.RootCause, p.Workaround, p.CreatedBy).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("itsm_repo.CreateProblem: %w", err)
	}
	return nil
}

func (r *ITSMRepo) UpdateProblem(ctx context.Context, p *domain.Problem) error {
	tag, err := r.db.Exec(ctx, `UPDATE problems SET title = $1, description = $2, status = $3, root_cause = $4, workaround = $5, updated_at = NOW() WHERE id = $6 AND tenant_id = $7`,
		p.Title, p.Description, p.Status, p.RootCause, p.Workaround, p.ID, p.TenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.UpdateProblem: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ITSMRepo) LinkIncident(ctx context.Context, tenantID, problemID, ticketID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `INSERT INTO problem_incidents (problem_id, ticket_id, tenant_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, problemID, ticketID, tenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.LinkIncident: %w", err)
	}
	return nil
}

func (r *ITSMRepo) UnlinkIncident(ctx context.Context, tenantID, problemID, ticketID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM problem_incidents WHERE problem_id = $1 AND ticket_id = $2 AND tenant_id = $3`, problemID, ticketID, tenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.UnlinkIncident: %w", err)
	}
	return nil
}

func (r *ITSMRepo) IncidentsForProblem(ctx context.Context, tenantID, problemID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `SELECT ticket_id FROM problem_incidents WHERE problem_id = $1 AND tenant_id = $2 ORDER BY created_at`, problemID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.IncidentsForProblem: %w", err)
	}
	defer rows.Close()
	var out []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return nonNil(out), rows.Err()
}

// ── Changes ─────────────────────────────────────────────────────────

func scanChangeRow(row pgx.Row) (*domain.Change, error) {
	var c domain.Change
	var problemID pgtype.UUID
	var plannedStart, plannedEnd, implementedAt, closedAt pgtype.Timestamptz
	err := row.Scan(&c.ID, &c.TenantID, &c.Title, &c.Description, &problemID, &c.Status, &c.Risk, &plannedStart, &plannedEnd, &implementedAt, &closedAt, &c.CreatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	c.ProblemID = pgToUUIDPtr(problemID)
	c.PlannedStart = fromPgTimestamptz(plannedStart)
	c.PlannedEnd = fromPgTimestamptz(plannedEnd)
	c.ImplementedAt = fromPgTimestamptz(implementedAt)
	c.ClosedAt = fromPgTimestamptz(closedAt)
	return &c, nil
}

const changeSelect = `SELECT id, tenant_id, title, description, problem_id, status, risk, planned_start, planned_end, implemented_at, closed_at, created_by, created_at, updated_at FROM changes`

func (r *ITSMRepo) ListChanges(ctx context.Context, tenantID uuid.UUID, status string) ([]domain.Change, error) {
	q := changeSelect + ` WHERE tenant_id = $1`
	args := []any{tenantID}
	if status != "" {
		q += ` AND status = $2`
		args = append(args, status)
	}
	q += ` ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.ListChanges: %w", err)
	}
	defer rows.Close()
	var out []domain.Change
	for rows.Next() {
		c, err := scanChangeRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return nonNil(out), rows.Err()
}

func (r *ITSMRepo) GetChange(ctx context.Context, id, tenantID uuid.UUID) (*domain.Change, error) {
	c, err := scanChangeRow(r.db.QueryRow(ctx, changeSelect+` WHERE id = $1 AND tenant_id = $2`, id, tenantID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("itsm_repo.GetChange: %w", err)
	}
	return c, nil
}

func (r *ITSMRepo) CreateChange(ctx context.Context, c *domain.Change) error {
	err := r.db.QueryRow(ctx, `INSERT INTO changes (tenant_id, title, description, problem_id, status, risk, planned_start, planned_end, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`,
		c.TenantID, c.Title, c.Description, toPgUUID(c.ProblemID), c.Status, c.Risk, toPgTimestamptz(c.PlannedStart), toPgTimestamptz(c.PlannedEnd), c.CreatedBy).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if isFKViolation(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("itsm_repo.CreateChange: %w", err)
	}
	return nil
}

func (r *ITSMRepo) UpdateChange(ctx context.Context, c *domain.Change) error {
	tag, err := r.db.Exec(ctx, `UPDATE changes SET title = $1, description = $2, problem_id = $3, status = $4, risk = $5, planned_start = $6, planned_end = $7, implemented_at = $8, closed_at = $9, updated_at = NOW() WHERE id = $10 AND tenant_id = $11`,
		c.Title, c.Description, toPgUUID(c.ProblemID), c.Status, c.Risk, toPgTimestamptz(c.PlannedStart), toPgTimestamptz(c.PlannedEnd), toPgTimestamptz(c.ImplementedAt), toPgTimestamptz(c.ClosedAt), c.ID, c.TenantID)
	if err != nil {
		if isFKViolation(err) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("itsm_repo.UpdateChange: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ITSMRepo) LinkChangeTicket(ctx context.Context, tenantID, changeID, ticketID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `INSERT INTO change_tickets (change_id, ticket_id, tenant_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, changeID, ticketID, tenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.LinkChangeTicket: %w", err)
	}
	return nil
}

func (r *ITSMRepo) UnlinkChangeTicket(ctx context.Context, tenantID, changeID, ticketID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM change_tickets WHERE change_id = $1 AND ticket_id = $2 AND tenant_id = $3`, changeID, ticketID, tenantID)
	if err != nil {
		return fmt.Errorf("itsm_repo.UnlinkChangeTicket: %w", err)
	}
	return nil
}

func (r *ITSMRepo) TicketsForChange(ctx context.Context, tenantID, changeID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `SELECT ticket_id FROM change_tickets WHERE change_id = $1 AND tenant_id = $2 ORDER BY created_at`, changeID, tenantID)
	if err != nil {
		return nil, fmt.Errorf("itsm_repo.TicketsForChange: %w", err)
	}
	defer rows.Close()
	var out []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return nonNil(out), rows.Err()
}

// ── Reports ─────────────────────────────────────────────────────────

func (r *ITSMRepo) Report(ctx context.Context, tenantID uuid.UUID) (*domain.Report, error) {
	rep := &domain.Report{
		ChangesByStatus:   []domain.StatusCount{},
		ProblemsByStatus:  []domain.StatusCount{},
		AssetsByType:      []domain.TypeCount{},
		ApprovalsByStatus: []domain.StatusCount{},
	}
	collect := func(q string, into any) error {
		rows, err := r.db.Query(ctx, q, tenantID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var label string
			var count int
			if err := rows.Scan(&label, &count); err != nil {
				return err
			}
			switch dst := into.(type) {
			case *[]domain.StatusCount:
				*dst = append(*dst, domain.StatusCount{Status: label, Count: count})
			case *[]domain.TypeCount:
				*dst = append(*dst, domain.TypeCount{Type: label, Count: count})
			}
		}
		return rows.Err()
	}
	if err := collect(`SELECT status, COUNT(*) FROM changes WHERE tenant_id = $1 GROUP BY status ORDER BY status`, &rep.ChangesByStatus); err != nil {
		return nil, fmt.Errorf("itsm_repo.Report: changes: %w", err)
	}
	if err := collect(`SELECT status, COUNT(*) FROM problems WHERE tenant_id = $1 GROUP BY status ORDER BY status`, &rep.ProblemsByStatus); err != nil {
		return nil, fmt.Errorf("itsm_repo.Report: problems: %w", err)
	}
	if err := collect(`SELECT COALESCE(ct.name, ''), COUNT(*) FROM assets a LEFT JOIN ci_types ct ON ct.id = a.ci_type_id WHERE a.tenant_id = $1 GROUP BY ct.name ORDER BY ct.name`, &rep.AssetsByType); err != nil {
		return nil, fmt.Errorf("itsm_repo.Report: assets: %w", err)
	}
	if err := collect(`SELECT status, COUNT(*) FROM approvals WHERE tenant_id = $1 GROUP BY status ORDER BY status`, &rep.ApprovalsByStatus); err != nil {
		return nil, fmt.Errorf("itsm_repo.Report: approvals: %w", err)
	}
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM approvals WHERE tenant_id = $1 AND status = 'pending'`, tenantID).Scan(&rep.PendingApprovals); err != nil {
		return nil, fmt.Errorf("itsm_repo.Report: pending: %w", err)
	}
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM catalog_services WHERE tenant_id = $1 AND active = TRUE`, tenantID).Scan(&rep.ActiveServices); err != nil {
		return nil, fmt.Errorf("itsm_repo.Report: services: %w", err)
	}
	return rep, nil
}
