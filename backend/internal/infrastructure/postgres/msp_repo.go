package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	appMSP "github.com/trickreport/backend/internal/application/msp"
	domain "github.com/trickreport/backend/internal/domain/msp"
	domainSLA "github.com/trickreport/backend/internal/domain/sla"
)

// MSPRepo implements the MSP repository boundary over PostgreSQL.
// Every query filters by tenant_id, and every client-owned query
// additionally by client_id: a record from another organization — or
// from another client of the same provider — is indistinguishable
// from a missing record.
type MSPRepo struct {
	db *pgxpool.Pool
}

// NewMSPRepo creates an MSPRepo.
func NewMSPRepo(db *pgxpool.Pool) *MSPRepo { return &MSPRepo{db: db} }

// Compile-time assertion that MSPRepo implements the application boundary.
var _ appMSP.Repository = (*MSPRepo)(nil)

// DB exposes the pool for wiring that needs it (mirrors other repos).
func (r *MSPRepo) DB() *pgxpool.Pool { return r.db }

func mspNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func mspMapErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	if isUniqueViolation(err) {
		return domain.ErrConflict
	}
	if isFKViolation(err) {
		return domain.ErrConflict
	}
	return err
}

// clientClause appends a client restriction to a query that already
// uses $1..$n. A nil list means "staff, no restriction"; an empty
// non-nil list means "contact with no clients": match nothing, never
// everything.
func clientClause(col string, nextArg int, ids []uuid.UUID) (string, []any) {
	if ids == nil {
		return "", nil
	}
	if len(ids) == 0 {
		return " AND 1=0", nil
	}
	return fmt.Sprintf(" AND %s = ANY($%d)", col, nextArg), []any{ids}
}

func parseDayTime(s string) (pgtype.Time, error) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return pgtype.Time{}, err
	}
	return pgtype.Time{Microseconds: int64(t.Hour())*3600*1000000 + int64(t.Minute())*60*1000000, Valid: true}, nil
}

func formatDayTime(t pgtype.Time) string {
	if !t.Valid {
		return ""
	}
	total := t.Microseconds / 1000000
	return fmt.Sprintf("%02d:%02d", total/3600, (total%3600)/60)
}

func marshalSLAOverrides(m map[string]domain.SLAOverride) []byte {
	if m == nil {
		return []byte("{}")
	}
	raw, err := json.Marshal(m)
	if err != nil {
		return []byte("{}")
	}
	return raw
}

func unmarshalSLAOverrides(raw []byte) map[string]domain.SLAOverride {
	out := map[string]domain.SLAOverride{}
	if len(raw) == 0 {
		return out
	}
	_ = json.Unmarshal(raw, &out)
	return out
}

// ── Clients and contacts ────────────────────────────────────────────

func (r *MSPRepo) ContactClientIDs(ctx context.Context, tenantID, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `SELECT client_id FROM msp_client_contacts WHERE tenant_id = $1 AND user_id = $2`, tenantID, userID)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ContactClientIDs: %w", err)
	}
	defer rows.Close()
	ids := []uuid.UUID{}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func scanClientRow(row pgx.Row) (*domain.Client, error) {
	var c domain.Client
	err := row.Scan(&c.ID, &c.TenantID, &c.Name, &c.Domain, &c.City, &c.Status, &c.ContactName, &c.ContactEmail, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *MSPRepo) ListClients(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID) ([]domain.Client, error) {
	query := `SELECT id, tenant_id, name, domain, city, status, contact_name, contact_email, created_at, updated_at FROM msp_clients WHERE tenant_id = $1`
	clause, args := clientClause("id", 2, clientIDs)
	query += clause + ` ORDER BY name`
	rows, err := r.db.Query(ctx, query, append([]any{tenantID}, args...)...)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListClients: %w", err)
	}
	defer rows.Close()
	var out []domain.Client
	for rows.Next() {
		c, err := scanClientRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) GetClient(ctx context.Context, id, tenantID uuid.UUID) (*domain.Client, error) {
	c, err := scanClientRow(r.db.QueryRow(ctx, `SELECT id, tenant_id, name, domain, city, status, contact_name, contact_email, created_at, updated_at FROM msp_clients WHERE id = $1 AND tenant_id = $2`, id, tenantID))
	return c, mspNotFound(err)
}

func (r *MSPRepo) CreateClient(ctx context.Context, c *domain.Client) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_clients (tenant_id, name, domain, city, status, contact_name, contact_email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		c.TenantID, c.Name, c.Domain, c.City, c.Status, c.ContactName, c.ContactEmail).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt))
}

func (r *MSPRepo) UpdateClient(ctx context.Context, c *domain.Client) error {
	tag, err := r.db.Exec(ctx, `UPDATE msp_clients SET name = $1, domain = $2, city = $3, status = $4, contact_name = $5, contact_email = $6, updated_at = NOW() WHERE id = $7 AND tenant_id = $8`,
		c.Name, c.Domain, c.City, c.Status, c.ContactName, c.ContactEmail, c.ID, c.TenantID)
	if err != nil {
		return mspMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MSPRepo) DeleteClient(ctx context.Context, id, tenantID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM msp_clients WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return mspMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MSPRepo) ListContacts(ctx context.Context, tenantID, clientID uuid.UUID) ([]domain.ClientContact, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, client_id, user_id, name, email, role, created_at FROM msp_client_contacts WHERE tenant_id = $1 AND client_id = $2 ORDER BY name`, tenantID, clientID)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListContacts: %w", err)
	}
	defer rows.Close()
	var out []domain.ClientContact
	for rows.Next() {
		var c domain.ClientContact
		var userID pgtype.UUID
		if err := rows.Scan(&c.ID, &c.TenantID, &c.ClientID, &userID, &c.Name, &c.Email, &c.Role, &c.CreatedAt); err != nil {
			return nil, err
		}
		if userID.Valid {
			id := uuid.UUID(userID.Bytes)
			c.UserID = &id
		}
		out = append(out, c)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) CreateContact(ctx context.Context, c *domain.ClientContact) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_client_contacts (tenant_id, client_id, user_id, name, email, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		c.TenantID, c.ClientID, toPgUUID(c.UserID), c.Name, c.Email, c.Role).Scan(&c.ID, &c.CreatedAt))
}

func (r *MSPRepo) DeleteContact(ctx context.Context, id, tenantID, clientID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM msp_client_contacts WHERE id = $1 AND tenant_id = $2 AND client_id = $3`, id, tenantID, clientID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ── Contracts ───────────────────────────────────────────────────────

func scanContractRow(row pgx.Row) (*domain.Contract, error) {
	var c domain.Contract
	var endsAt pgtype.Timestamptz
	var overrides []byte
	err := row.Scan(&c.ID, &c.TenantID, &c.ClientID, &c.Name, &c.Status, &c.StartsAt, &endsAt, &c.Coverage, &c.Currency, &c.MonthlyFeeCents, &c.HourlyRateCents, &c.IncludedHoursMonthly, &c.SLAResponseMinutes, &c.SLAResolutionMinutes, &overrides, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	c.EndsAt = fromPgTimestamptz(endsAt)
	c.SLAOverrides = unmarshalSLAOverrides(overrides)
	return &c, nil
}

const contractColumns = `id, tenant_id, client_id, name, status, starts_at, ends_at, coverage, currency, monthly_fee_cents, hourly_rate_cents, included_hours_monthly, sla_response_minutes, sla_resolution_minutes, sla_overrides, created_at, updated_at`

func (r *MSPRepo) ListContracts(ctx context.Context, tenantID uuid.UUID, clientID *uuid.UUID) ([]domain.Contract, error) {
	query := `SELECT ` + contractColumns + ` FROM msp_contracts WHERE tenant_id = $1`
	args := []any{tenantID}
	if clientID != nil {
		query += ` AND client_id = $2`
		args = append(args, *clientID)
	}
	query += ` ORDER BY starts_at DESC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListContracts: %w", err)
	}
	defer rows.Close()
	var out []domain.Contract
	for rows.Next() {
		c, err := scanContractRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) GetContract(ctx context.Context, id, tenantID uuid.UUID) (*domain.Contract, error) {
	c, err := scanContractRow(r.db.QueryRow(ctx, `SELECT `+contractColumns+` FROM msp_contracts WHERE id = $1 AND tenant_id = $2`, id, tenantID))
	return c, mspNotFound(err)
}

func (r *MSPRepo) CreateContract(ctx context.Context, c *domain.Contract) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_contracts (tenant_id, client_id, name, status, starts_at, ends_at, coverage, currency, monthly_fee_cents, hourly_rate_cents, included_hours_monthly, sla_response_minutes, sla_resolution_minutes, sla_overrides) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, created_at, updated_at`,
		c.TenantID, c.ClientID, c.Name, c.Status, c.StartsAt, toPgTimestamptz(c.EndsAt), c.Coverage, c.Currency, c.MonthlyFeeCents, c.HourlyRateCents, c.IncludedHoursMonthly, c.SLAResponseMinutes, c.SLAResolutionMinutes, marshalSLAOverrides(c.SLAOverrides)).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt))
}

func (r *MSPRepo) UpdateContract(ctx context.Context, c *domain.Contract) error {
	tag, err := r.db.Exec(ctx, `UPDATE msp_contracts SET name = $1, status = $2, starts_at = $3, ends_at = $4, coverage = $5, currency = $6, monthly_fee_cents = $7, hourly_rate_cents = $8, included_hours_monthly = $9, sla_response_minutes = $10, sla_resolution_minutes = $11, sla_overrides = $12, updated_at = NOW() WHERE id = $13 AND tenant_id = $14`,
		c.Name, c.Status, c.StartsAt, toPgTimestamptz(c.EndsAt), c.Coverage, c.Currency, c.MonthlyFeeCents, c.HourlyRateCents, c.IncludedHoursMonthly, c.SLAResponseMinutes, c.SLAResolutionMinutes, marshalSLAOverrides(c.SLAOverrides), c.ID, c.TenantID)
	if err != nil {
		return mspMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ActiveContractForClient returns the most recently started active
// contract covering the moment, or ErrNotFound when there is none.
func (r *MSPRepo) ActiveContractForClient(ctx context.Context, tenantID, clientID uuid.UUID, at time.Time) (*domain.Contract, error) {
	c, err := scanContractRow(r.db.QueryRow(ctx, `SELECT `+contractColumns+` FROM msp_contracts WHERE tenant_id = $1 AND client_id = $2 AND status = 'active' AND starts_at <= $3 AND (ends_at IS NULL OR ends_at >= $3) ORDER BY starts_at DESC LIMIT 1`, tenantID, clientID, at))
	return c, mspNotFound(err)
}

// ── Tickets ─────────────────────────────────────────────────────────

func (r *MSPRepo) CreateTicketWithContext(ctx context.Context, tenantID, creatorID uuid.UUID, title, description, category, priority string, tc *domain.TicketContext, dl domainSLA.Deadlines) (uuid.UUID, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var ticketID uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO tickets (tenant_id, title, description, status, priority, category, created_by, sla_deadline, response_deadline, escalation_deadline) VALUES ($1, $2, $3, 'open', $4, $5, $6, $7, $8, $9) RETURNING id`,
		tenantID, title, description, priority, category, creatorID, dl.Resolution, dl.Response, dl.Escalation).Scan(&ticketID); err != nil {
		return uuid.Nil, mspMapErr(err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO msp_ticket_context (ticket_id, tenant_id, client_id, contract_id, sla_response_minutes, sla_resolution_minutes) VALUES ($1, $2, $3, $4, $5, $6)`,
		ticketID, tenantID, tc.ClientID, toPgUUID(tc.ContractID), tc.SLAResponseMinutes, tc.SLAResolutionMinutes); err != nil {
		return uuid.Nil, mspMapErr(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}
	return ticketID, nil
}

func (r *MSPRepo) AlertTicketCreator(ctx context.Context, tenantID uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx, `SELECT id FROM users WHERE tenant_id = $1 AND role IN ('admin', 'agent') AND active = TRUE ORDER BY created_at ASC LIMIT 1`, tenantID).Scan(&id)
	if err != nil {
		return uuid.Nil, mspNotFound(fmt.Errorf("msp_repo.AlertTicketCreator: no staff user to attribute the alert ticket: %w", err))
	}
	return id, nil
}

func (r *MSPRepo) TicketContext(ctx context.Context, tenantID, ticketID uuid.UUID) (*domain.TicketContext, error) {
	var tc domain.TicketContext
	var contractID pgtype.UUID
	err := r.db.QueryRow(ctx, `SELECT tc.ticket_id, tc.tenant_id, tc.client_id, tc.contract_id, tc.sla_response_minutes, tc.sla_resolution_minutes, tc.created_at, cl.name FROM msp_ticket_context tc JOIN msp_clients cl ON cl.id = tc.client_id WHERE tc.ticket_id = $1 AND tc.tenant_id = $2`, ticketID, tenantID).
		Scan(&tc.TicketID, &tc.TenantID, &tc.ClientID, &contractID, &tc.SLAResponseMinutes, &tc.SLAResolutionMinutes, &tc.CreatedAt, &tc.ClientName)
	if err != nil {
		return nil, mspNotFound(err)
	}
	if contractID.Valid {
		id := uuid.UUID(contractID.Bytes)
		tc.ContractID = &id
	}
	return &tc, nil
}

func (r *MSPRepo) ListMSPTickets(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID) ([]appMSP.MSPTicket, error) {
	query := `SELECT t.id, t.title, t.status, t.priority, t.category, t.created_at, tc.client_id, cl.name, tc.contract_id, tc.sla_response_minutes, tc.sla_resolution_minutes FROM tickets t JOIN msp_ticket_context tc ON tc.ticket_id = t.id JOIN msp_clients cl ON cl.id = tc.client_id WHERE t.tenant_id = $1`
	args := []any{tenantID}
	next := 2
	if clause, cargs := clientClause("tc.client_id", next, clientIDs); clause != "" {
		query += clause
		args = append(args, cargs...)
		next++
	}
	if clientFilter != nil {
		query += fmt.Sprintf(` AND tc.client_id = $%d`, next)
		args = append(args, *clientFilter)
	}
	query += ` ORDER BY t.created_at DESC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListMSPTickets: %w", err)
	}
	defer rows.Close()
	var out []appMSP.MSPTicket
	for rows.Next() {
		var tk appMSP.MSPTicket
		var contractID pgtype.UUID
		if err := rows.Scan(&tk.ID, &tk.Title, &tk.Status, &tk.Priority, &tk.Category, &tk.CreatedAt, &tk.ClientID, &tk.ClientName, &contractID, &tk.SLAResponseMinutes, &tk.SLAResolutionMinutes); err != nil {
			return nil, err
		}
		if contractID.Valid {
			id := uuid.UUID(contractID.Bytes)
			tk.ContractID = &id
		}
		out = append(out, tk)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) UserExists(ctx context.Context, id, tenantID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND tenant_id = $2)`, id, tenantID).Scan(&exists)
	return exists, err
}

func (r *MSPRepo) AssetExists(ctx context.Context, id, tenantID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM assets WHERE id = $1 AND tenant_id = $2)`, id, tenantID).Scan(&exists)
	return exists, err
}

// ── Technicians and planning ────────────────────────────────────────

func scanTechnicianRow(row pgx.Row) (*domain.Technician, error) {
	var t domain.Technician
	var skills []string
	err := row.Scan(&t.ID, &t.TenantID, &t.UserID, &t.UserName, &skills, &t.HourlyRateCents, &t.Active, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	t.Skills = skills
	if t.Skills == nil {
		t.Skills = []string{}
	}
	return &t, nil
}

func (r *MSPRepo) ListTechnicians(ctx context.Context, tenantID uuid.UUID) ([]domain.Technician, error) {
	rows, err := r.db.Query(ctx, `SELECT te.id, te.tenant_id, te.user_id, u.name, te.skills, te.hourly_rate_cents, te.active, te.created_at, te.updated_at FROM msp_technicians te JOIN users u ON u.id = te.user_id WHERE te.tenant_id = $1 ORDER BY u.name`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListTechnicians: %w", err)
	}
	defer rows.Close()
	var out []domain.Technician
	for rows.Next() {
		t, err := scanTechnicianRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *t)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) GetTechnician(ctx context.Context, id, tenantID uuid.UUID) (*domain.Technician, error) {
	t, err := scanTechnicianRow(r.db.QueryRow(ctx, `SELECT te.id, te.tenant_id, te.user_id, u.name, te.skills, te.hourly_rate_cents, te.active, te.created_at, te.updated_at FROM msp_technicians te JOIN users u ON u.id = te.user_id WHERE te.id = $1 AND te.tenant_id = $2`, id, tenantID))
	return t, mspNotFound(err)
}

func (r *MSPRepo) CreateTechnician(ctx context.Context, t *domain.Technician) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_technicians (tenant_id, user_id, skills, hourly_rate_cents, active) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`,
		t.TenantID, t.UserID, t.Skills, t.HourlyRateCents, t.Active).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt))
}

func (r *MSPRepo) UpdateTechnician(ctx context.Context, t *domain.Technician) error {
	tag, err := r.db.Exec(ctx, `UPDATE msp_technicians SET skills = $1, hourly_rate_cents = $2, active = $3, updated_at = NOW() WHERE id = $4 AND tenant_id = $5`,
		t.Skills, t.HourlyRateCents, t.Active, t.ID, t.TenantID)
	if err != nil {
		return mspMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MSPRepo) ListAvailability(ctx context.Context, tenantID, technicianID uuid.UUID) ([]domain.Availability, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, technician_id, day_of_week, start_time, end_time FROM msp_availability WHERE tenant_id = $1 AND technician_id = $2 ORDER BY day_of_week, start_time`, tenantID, technicianID)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListAvailability: %w", err)
	}
	defer rows.Close()
	var out []domain.Availability
	for rows.Next() {
		var a domain.Availability
		var start, end pgtype.Time
		if err := rows.Scan(&a.ID, &a.TenantID, &a.TechnicianID, &a.DayOfWeek, &start, &end); err != nil {
			return nil, err
		}
		a.StartTime = formatDayTime(start)
		a.EndTime = formatDayTime(end)
		out = append(out, a)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) CreateAvailability(ctx context.Context, a *domain.Availability) error {
	start, err := parseDayTime(a.StartTime)
	if err != nil {
		return domain.ErrValidation
	}
	end, err := parseDayTime(a.EndTime)
	if err != nil {
		return domain.ErrValidation
	}
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_availability (tenant_id, technician_id, day_of_week, start_time, end_time) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		a.TenantID, a.TechnicianID, a.DayOfWeek, start, end).Scan(&a.ID))
}

func scanInterventionRow(row pgx.Row) (*domain.Intervention, error) {
	var v domain.Intervention
	var ticketID pgtype.UUID
	err := row.Scan(&v.ID, &v.TenantID, &v.ClientID, &v.TechnicianID, &ticketID, &v.Title, &v.StartsAt, &v.EndsAt, &v.Status, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if ticketID.Valid {
		id := uuid.UUID(ticketID.Bytes)
		v.TicketID = &id
	}
	return &v, nil
}

const interventionColumns = `id, tenant_id, client_id, technician_id, ticket_id, title, starts_at, ends_at, status, created_at, updated_at`

func (r *MSPRepo) ListInterventions(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter, technicianFilter *uuid.UUID) ([]domain.Intervention, error) {
	query := `SELECT ` + interventionColumns + ` FROM msp_interventions WHERE tenant_id = $1`
	args := []any{tenantID}
	next := 2
	if clause, cargs := clientClause("client_id", next, clientIDs); clause != "" {
		query += clause
		args = append(args, cargs...)
		next++
	}
	if clientFilter != nil {
		query += fmt.Sprintf(` AND client_id = $%d`, next)
		args = append(args, *clientFilter)
		next++
	}
	if technicianFilter != nil {
		query += fmt.Sprintf(` AND technician_id = $%d`, next)
		args = append(args, *technicianFilter)
	}
	query += ` ORDER BY starts_at`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListInterventions: %w", err)
	}
	defer rows.Close()
	var out []domain.Intervention
	for rows.Next() {
		v, err := scanInterventionRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *v)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) CreateIntervention(ctx context.Context, v *domain.Intervention) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_interventions (tenant_id, client_id, technician_id, ticket_id, title, starts_at, ends_at, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`,
		v.TenantID, v.ClientID, v.TechnicianID, toPgUUID(v.TicketID), v.Title, v.StartsAt, v.EndsAt, v.Status).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt))
}

func (r *MSPRepo) UpdateInterventionStatus(ctx context.Context, v *domain.Intervention) error {
	tag, err := r.db.Exec(ctx, `UPDATE msp_interventions SET status = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3`, v.Status, v.ID, v.TenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MSPRepo) TechnicianOverlappingIntervention(ctx context.Context, tenantID, technicianID uuid.UUID, startsAt, endsAt time.Time) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM msp_interventions WHERE tenant_id = $1 AND technician_id = $2 AND status <> 'cancelled' AND starts_at < $4 AND $3 < ends_at)`, tenantID, technicianID, startsAt, endsAt).Scan(&exists)
	return exists, err
}

// ── Time and expenses ───────────────────────────────────────────────

func scanTimeEntryRow(row pgx.Row) (*domain.TimeEntry, error) {
	var e domain.TimeEntry
	var ticketID, interventionID, contractRate pgtype.UUID
	var contractRateVal pgtype.Int8
	_ = contractRate
	err := row.Scan(&e.ID, &e.TenantID, &e.ClientID, &ticketID, &interventionID, &e.TechnicianID, &e.Minutes, &e.Description, &e.Billable, &e.Billed, &e.HourlyRateCents, &contractRateVal, &e.CostCents, &e.OccurredAt, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	if ticketID.Valid {
		id := uuid.UUID(ticketID.Bytes)
		e.TicketID = &id
	}
	if interventionID.Valid {
		id := uuid.UUID(interventionID.Bytes)
		e.InterventionID = &id
	}
	if contractRateVal.Valid {
		v := contractRateVal.Int64
		e.ContractHourlyRateCents = &v
	}
	return &e, nil
}

const timeEntryColumns = `id, tenant_id, client_id, ticket_id, intervention_id, technician_id, minutes, description, billable, billed, hourly_rate_cents, contract_hourly_rate_cents, cost_cents, occurred_at, created_at`

func (r *MSPRepo) CreateTimeEntry(ctx context.Context, e *domain.TimeEntry) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_time_entries (tenant_id, client_id, ticket_id, intervention_id, technician_id, minutes, description, billable, hourly_rate_cents, contract_hourly_rate_cents, cost_cents, occurred_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, created_at`,
		e.TenantID, e.ClientID, toPgUUID(e.TicketID), toPgUUID(e.InterventionID), e.TechnicianID, e.Minutes, e.Description, e.Billable, e.HourlyRateCents, toPgInt8(e.ContractHourlyRateCents), e.CostCents, e.OccurredAt).Scan(&e.ID, &e.CreatedAt))
}

func toPgInt8(v *int64) pgtype.Int8 {
	if v == nil {
		return pgtype.Int8{Valid: false}
	}
	return pgtype.Int8{Int64: *v, Valid: true}
}

func dateFilter(query string, args []any, next int, col string, from, to *time.Time) (string, []any) {
	if from != nil {
		query += fmt.Sprintf(` AND %s >= $%d::date`, col, next)
		args = append(args, *from)
		next++
	}
	if to != nil {
		query += fmt.Sprintf(` AND %s <= $%d::date`, col, next)
		args = append(args, *to)
	}
	return query, args
}

func (r *MSPRepo) ListTimeEntries(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID, from, to *time.Time) ([]domain.TimeEntry, error) {
	query := `SELECT ` + timeEntryColumns + ` FROM msp_time_entries WHERE tenant_id = $1`
	args := []any{tenantID}
	next := 2
	if clause, cargs := clientClause("client_id", next, clientIDs); clause != "" {
		query += clause
		args = append(args, cargs...)
		next++
	}
	if clientFilter != nil {
		query += fmt.Sprintf(` AND client_id = $%d`, next)
		args = append(args, *clientFilter)
		next++
	}
	query, args = dateFilter(query, args, next, "occurred_at", from, to)
	query += ` ORDER BY occurred_at DESC, created_at DESC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListTimeEntries: %w", err)
	}
	defer rows.Close()
	var out []domain.TimeEntry
	for rows.Next() {
		e, err := scanTimeEntryRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *e)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) UnbilledTimeEntries(ctx context.Context, tenantID, clientID uuid.UUID, from, to time.Time) ([]domain.TimeEntry, error) {
	rows, err := r.db.Query(ctx, `SELECT `+timeEntryColumns+` FROM msp_time_entries WHERE tenant_id = $1 AND client_id = $2 AND billed = FALSE AND occurred_at >= $3::date AND occurred_at <= $4::date ORDER BY occurred_at`, tenantID, clientID, from, to)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.UnbilledTimeEntries: %w", err)
	}
	defer rows.Close()
	var out []domain.TimeEntry
	for rows.Next() {
		e, err := scanTimeEntryRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *e)
	}
	return nonNil(out), rows.Err()
}

func scanExpenseRow(row pgx.Row) (*domain.Expense, error) {
	var e domain.Expense
	var ticketID, projectID pgtype.UUID
	err := row.Scan(&e.ID, &e.TenantID, &e.ClientID, &ticketID, &projectID, &e.Description, &e.Category, &e.AmountCents, &e.Billed, &e.IncurredAt, &e.CreatedAt)
	if err != nil {
		return nil, err
	}
	if ticketID.Valid {
		id := uuid.UUID(ticketID.Bytes)
		e.TicketID = &id
	}
	if projectID.Valid {
		id := uuid.UUID(projectID.Bytes)
		e.ProjectID = &id
	}
	return &e, nil
}

const expenseColumns = `id, tenant_id, client_id, ticket_id, project_id, description, category, amount_cents, billed, incurred_at, created_at`

func (r *MSPRepo) CreateExpense(ctx context.Context, e *domain.Expense) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_expenses (tenant_id, client_id, ticket_id, project_id, description, category, amount_cents, incurred_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at`,
		e.TenantID, e.ClientID, toPgUUID(e.TicketID), toPgUUID(e.ProjectID), e.Description, e.Category, e.AmountCents, e.IncurredAt).Scan(&e.ID, &e.CreatedAt))
}

func (r *MSPRepo) ListExpenses(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID, from, to *time.Time) ([]domain.Expense, error) {
	query := `SELECT ` + expenseColumns + ` FROM msp_expenses WHERE tenant_id = $1`
	args := []any{tenantID}
	next := 2
	if clause, cargs := clientClause("client_id", next, clientIDs); clause != "" {
		query += clause
		args = append(args, cargs...)
		next++
	}
	if clientFilter != nil {
		query += fmt.Sprintf(` AND client_id = $%d`, next)
		args = append(args, *clientFilter)
		next++
	}
	query, args = dateFilter(query, args, next, "incurred_at", from, to)
	query += ` ORDER BY incurred_at DESC, created_at DESC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListExpenses: %w", err)
	}
	defer rows.Close()
	var out []domain.Expense
	for rows.Next() {
		e, err := scanExpenseRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *e)
	}
	return nonNil(out), rows.Err()
}

// ── Projects ────────────────────────────────────────────────────────

func scanProjectRow(row pgx.Row) (*domain.Project, error) {
	var p domain.Project
	var startsAt, endsAt pgtype.Timestamptz
	err := row.Scan(&p.ID, &p.TenantID, &p.ClientID, &p.Name, &p.Status, &p.BudgetCents, &startsAt, &endsAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	p.StartsAt = fromPgTimestamptz(startsAt)
	p.EndsAt = fromPgTimestamptz(endsAt)
	return &p, nil
}

const projectColumns = `id, tenant_id, client_id, name, status, budget_cents, starts_at, ends_at, created_at, updated_at`

func (r *MSPRepo) ListProjects(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID) ([]domain.Project, error) {
	query := `SELECT ` + projectColumns + ` FROM msp_projects WHERE tenant_id = $1`
	args := []any{tenantID}
	next := 2
	if clause, cargs := clientClause("client_id", next, clientIDs); clause != "" {
		query += clause
		args = append(args, cargs...)
		next++
	}
	if clientFilter != nil {
		query += fmt.Sprintf(` AND client_id = $%d`, next)
		args = append(args, *clientFilter)
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListProjects: %w", err)
	}
	defer rows.Close()
	var out []domain.Project
	for rows.Next() {
		p, err := scanProjectRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) GetProject(ctx context.Context, id, tenantID uuid.UUID) (*domain.Project, error) {
	p, err := scanProjectRow(r.db.QueryRow(ctx, `SELECT `+projectColumns+` FROM msp_projects WHERE id = $1 AND tenant_id = $2`, id, tenantID))
	return p, mspNotFound(err)
}

func (r *MSPRepo) CreateProject(ctx context.Context, p *domain.Project) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_projects (tenant_id, client_id, name, status, budget_cents, starts_at, ends_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		p.TenantID, p.ClientID, p.Name, p.Status, p.BudgetCents, toPgTimestamptz(p.StartsAt), toPgTimestamptz(p.EndsAt)).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt))
}

func (r *MSPRepo) UpdateProject(ctx context.Context, p *domain.Project) error {
	tag, err := r.db.Exec(ctx, `UPDATE msp_projects SET name = $1, status = $2, budget_cents = $3, starts_at = $4, ends_at = $5, updated_at = NOW() WHERE id = $6 AND tenant_id = $7`,
		p.Name, p.Status, p.BudgetCents, toPgTimestamptz(p.StartsAt), toPgTimestamptz(p.EndsAt), p.ID, p.TenantID)
	if err != nil {
		return mspMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanTaskRow(row pgx.Row) (*domain.ProjectTask, error) {
	var t domain.ProjectTask
	var techID pgtype.UUID
	err := row.Scan(&t.ID, &t.TenantID, &t.ProjectID, &t.Title, &t.Status, &techID, &t.EstimatedMinutes, &t.Position, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if techID.Valid {
		id := uuid.UUID(techID.Bytes)
		t.TechnicianID = &id
	}
	return &t, nil
}

func (r *MSPRepo) ListTasks(ctx context.Context, tenantID, projectID uuid.UUID) ([]domain.ProjectTask, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, project_id, title, status, technician_id, estimated_minutes, position, created_at, updated_at FROM msp_project_tasks WHERE tenant_id = $1 AND project_id = $2 ORDER BY position, created_at`, tenantID, projectID)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListTasks: %w", err)
	}
	defer rows.Close()
	var out []domain.ProjectTask
	for rows.Next() {
		t, err := scanTaskRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *t)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) CreateTask(ctx context.Context, t *domain.ProjectTask) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_project_tasks (tenant_id, project_id, title, status, technician_id, estimated_minutes, position) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		t.TenantID, t.ProjectID, t.Title, t.Status, toPgUUID(t.TechnicianID), t.EstimatedMinutes, t.Position).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt))
}

func (r *MSPRepo) UpdateTask(ctx context.Context, t *domain.ProjectTask) error {
	tag, err := r.db.Exec(ctx, `UPDATE msp_project_tasks SET title = $1, status = $2, technician_id = $3, estimated_minutes = $4, position = $5, updated_at = NOW() WHERE id = $6 AND tenant_id = $7`,
		t.Title, t.Status, toPgUUID(t.TechnicianID), t.EstimatedMinutes, t.Position, t.ID, t.TenantID)
	if err != nil {
		return mspMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ── Billing and profitability ───────────────────────────────────────

func scanBillingDraftRow(row pgx.Row) (*domain.BillingDraft, error) {
	var d domain.BillingDraft
	var contractID pgtype.UUID
	err := row.Scan(&d.ID, &d.TenantID, &d.ClientID, &contractID, &d.PeriodStart, &d.PeriodEnd, &d.Status, &d.TotalCents, &d.Notes, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if contractID.Valid {
		id := uuid.UUID(contractID.Bytes)
		d.ContractID = &id
	}
	return &d, nil
}

const billingDraftColumns = `id, tenant_id, client_id, contract_id, period_start, period_end, status, total_cents, notes, created_at, updated_at`

// CreateBillingDraft inserts the draft, its items and the billed marks
// in one transaction: a draft never bills the same minutes or expenses
// twice, and a failed draft never leaves entries marked as billed.
func (r *MSPRepo) CreateBillingDraft(ctx context.Context, d *domain.BillingDraft, items []domain.BillingItem, billedEntryIDs, billedExpenseIDs []uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if err := tx.QueryRow(ctx, `INSERT INTO msp_billing_drafts (tenant_id, client_id, contract_id, period_start, period_end, status, total_cents, notes) VALUES ($1, $2, $3, $4, $5, 'draft', $6, $7) RETURNING id, created_at, updated_at`,
		d.TenantID, d.ClientID, toPgUUID(d.ContractID), d.PeriodStart, d.PeriodEnd, d.TotalCents, d.Notes).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return mspMapErr(err)
	}
	for i := range items {
		items[i].DraftID = d.ID
		if err := tx.QueryRow(ctx, `INSERT INTO msp_billing_items (tenant_id, draft_id, kind, description, amount_cents) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
			d.TenantID, d.ID, items[i].Kind, items[i].Description, items[i].AmountCents).Scan(&items[i].ID, &items[i].CreatedAt); err != nil {
			return mspMapErr(err)
		}
	}
	if len(billedEntryIDs) > 0 {
		if _, err := tx.Exec(ctx, `UPDATE msp_time_entries SET billed = TRUE WHERE tenant_id = $1 AND id = ANY($2)`, d.TenantID, billedEntryIDs); err != nil {
			return err
		}
	}
	if len(billedExpenseIDs) > 0 {
		if _, err := tx.Exec(ctx, `UPDATE msp_expenses SET billed = TRUE WHERE tenant_id = $1 AND id = ANY($2)`, d.TenantID, billedExpenseIDs); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}

func (r *MSPRepo) UnbilledExpenses(ctx context.Context, tenantID, clientID uuid.UUID, from, to time.Time) ([]domain.Expense, error) {
	rows, err := r.db.Query(ctx, `SELECT `+expenseColumns+` FROM msp_expenses WHERE tenant_id = $1 AND client_id = $2 AND billed = FALSE AND incurred_at >= $3::date AND incurred_at <= $4::date ORDER BY incurred_at`, tenantID, clientID, from, to)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.UnbilledExpenses: %w", err)
	}
	defer rows.Close()
	var out []domain.Expense
	for rows.Next() {
		e, err := scanExpenseRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *e)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) ListBillingDrafts(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID) ([]domain.BillingDraft, error) {
	query := `SELECT ` + billingDraftColumns + ` FROM msp_billing_drafts WHERE tenant_id = $1`
	args := []any{tenantID}
	next := 2
	if clause, cargs := clientClause("client_id", next, clientIDs); clause != "" {
		query += clause
		args = append(args, cargs...)
		next++
	}
	if clientFilter != nil {
		query += fmt.Sprintf(` AND client_id = $%d`, next)
		args = append(args, *clientFilter)
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListBillingDrafts: %w", err)
	}
	defer rows.Close()
	var out []domain.BillingDraft
	for rows.Next() {
		d, err := scanBillingDraftRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *d)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) UpdateBillingDraftStatus(ctx context.Context, d *domain.BillingDraft) error {
	tag, err := r.db.Exec(ctx, `UPDATE msp_billing_drafts SET status = $1, updated_at = NOW() WHERE id = $2 AND tenant_id = $3`, d.Status, d.ID, d.TenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MSPRepo) Profitability(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter, projectFilter *uuid.UUID) (*domain.Profitability, error) {
	p := &domain.Profitability{}
	if clientFilter != nil {
		p.ClientID = *clientFilter
	}
	// Time: billable minutes valued at the contract rate when present,
	// else at the technician rate; cost always at the technician rate.
	timeQuery := `SELECT COALESCE(SUM(minutes), 0), COALESCE(SUM(minutes * COALESCE(contract_hourly_rate_cents, hourly_rate_cents) / 60), 0), COALESCE(SUM(cost_cents), 0) FROM msp_time_entries WHERE tenant_id = $1 AND billable = TRUE`
	timeArgs := []any{tenantID}
	next := 2
	if clause, cargs := clientClause("client_id", next, clientIDs); clause != "" {
		timeQuery += clause
		timeArgs = append(timeArgs, cargs...)
		next++
	}
	if clientFilter != nil {
		timeQuery += fmt.Sprintf(` AND client_id = $%d`, next)
		timeArgs = append(timeArgs, *clientFilter)
	}
	var minutes int
	if err := r.db.QueryRow(ctx, timeQuery, timeArgs...).Scan(&minutes, &p.TimeRevenue, &p.TimeCost); err != nil {
		return nil, fmt.Errorf("msp_repo.Profitability time: %w", err)
	}
	p.Minutes = minutes
	// Expenses: optionally restricted to one project.
	expQuery := `SELECT COALESCE(SUM(amount_cents), 0) FROM msp_expenses WHERE tenant_id = $1`
	expArgs := []any{tenantID}
	next = 2
	if clause, cargs := clientClause("client_id", next, clientIDs); clause != "" {
		expQuery += clause
		expArgs = append(expArgs, cargs...)
		next++
	}
	if clientFilter != nil {
		expQuery += fmt.Sprintf(` AND client_id = $%d`, next)
		expArgs = append(expArgs, *clientFilter)
		next++
	}
	if projectFilter != nil {
		expQuery += fmt.Sprintf(` AND project_id = $%d`, next)
		expArgs = append(expArgs, *projectFilter)
	}
	if err := r.db.QueryRow(ctx, expQuery, expArgs...).Scan(&p.Expenses); err != nil {
		return nil, fmt.Errorf("msp_repo.Profitability expenses: %w", err)
	}
	// Billed drafts: issued drafts count as revenue already billed.
	draftQuery := `SELECT COALESCE(SUM(total_cents), 0) FROM msp_billing_drafts WHERE tenant_id = $1 AND status = 'issued'`
	draftArgs := []any{tenantID}
	next = 2
	if clause, cargs := clientClause("client_id", next, clientIDs); clause != "" {
		draftQuery += clause
		draftArgs = append(draftArgs, cargs...)
		next++
	}
	if clientFilter != nil {
		draftQuery += fmt.Sprintf(` AND client_id = $%d`, next)
		draftArgs = append(draftArgs, *clientFilter)
	}
	if err := r.db.QueryRow(ctx, draftQuery, draftArgs...).Scan(&p.BilledDrafts); err != nil {
		return nil, fmt.Errorf("msp_repo.Profitability drafts: %w", err)
	}
	return p, nil
}

// ── Inventory ───────────────────────────────────────────────────────

func (r *MSPRepo) LinkAssetToClient(ctx context.Context, tenantID, clientID, assetID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `INSERT INTO msp_client_assets (tenant_id, client_id, asset_id) VALUES ($1, $2, $3) ON CONFLICT (client_id, asset_id) DO NOTHING`, tenantID, clientID, assetID)
	return mspMapErr(err)
}

func (r *MSPRepo) UnlinkAssetFromClient(ctx context.Context, tenantID, clientID, assetID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM msp_client_assets WHERE tenant_id = $1 AND client_id = $2 AND asset_id = $3`, tenantID, clientID, assetID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MSPRepo) ListClientAssets(ctx context.Context, tenantID, clientID uuid.UUID) ([]domain.ClientAssetView, error) {
	rows, err := r.db.Query(ctx, `SELECT a.id, ca.client_id, a.name, ct.name, a.status, a.serial_number FROM msp_client_assets ca JOIN assets a ON a.id = ca.asset_id JOIN ci_types ct ON ct.id = a.ci_type_id WHERE ca.tenant_id = $1 AND ca.client_id = $2 ORDER BY a.name`, tenantID, clientID)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListClientAssets: %w", err)
	}
	defer rows.Close()
	var out []domain.ClientAssetView
	for rows.Next() {
		var v domain.ClientAssetView
		if err := rows.Scan(&v.AssetID, &v.ClientID, &v.Name, &v.CITypeName, &v.Status, &v.Serial); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) CreateDiscoveredAsset(ctx context.Context, d *domain.DiscoveredAsset) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_discovered_assets (tenant_id, client_id, external_id, name, kind, payload) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, discovered_at`,
		d.TenantID, d.ClientID, d.ExternalID, d.Name, d.Kind, marshalAttrs(d.Payload)).Scan(&d.ID, &d.DiscoveredAt))
}

func scanDiscoveredRow(row pgx.Row) (*domain.DiscoveredAsset, error) {
	var d domain.DiscoveredAsset
	var payload []byte
	var importedID pgtype.UUID
	err := row.Scan(&d.ID, &d.TenantID, &d.ClientID, &d.ExternalID, &d.Name, &d.Kind, &payload, &d.Status, &importedID, &d.DiscoveredAt)
	if err != nil {
		return nil, err
	}
	d.Payload = mapAttrs(payload)
	if importedID.Valid {
		id := uuid.UUID(importedID.Bytes)
		d.ImportedAssetID = &id
	}
	return &d, nil
}

func (r *MSPRepo) ListDiscoveredAssets(ctx context.Context, tenantID, clientID uuid.UUID) ([]domain.DiscoveredAsset, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, client_id, external_id, name, kind, payload, status, imported_asset_id, discovered_at FROM msp_discovered_assets WHERE tenant_id = $1 AND client_id = $2 ORDER BY discovered_at DESC`, tenantID, clientID)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListDiscoveredAssets: %w", err)
	}
	defer rows.Close()
	var out []domain.DiscoveredAsset
	for rows.Next() {
		d, err := scanDiscoveredRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *d)
	}
	return nonNil(out), rows.Err()
}

// ImportDiscoveredAsset creates the Etapa 3 asset, links it to the
// client and marks the discovery imported, in one transaction. Only a
// still-pending discovery can be imported.
func (r *MSPRepo) ImportDiscoveredAsset(ctx context.Context, tenantID, clientID, discoveredID, ciTypeID uuid.UUID, assetName string) (uuid.UUID, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var assetID uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO assets (tenant_id, ci_type_id, name, status) SELECT $1, $2, $3, 'active' WHERE EXISTS (SELECT 1 FROM ci_types WHERE id = $2 AND tenant_id = $1) RETURNING id`, tenantID, ciTypeID, assetName).Scan(&assetID); err != nil {
		return uuid.Nil, mspMapErr(err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO msp_client_assets (tenant_id, client_id, asset_id) VALUES ($1, $2, $3) ON CONFLICT (client_id, asset_id) DO NOTHING`, tenantID, clientID, assetID); err != nil {
		return uuid.Nil, err
	}
	tag, err := tx.Exec(ctx, `UPDATE msp_discovered_assets SET status = 'imported', imported_asset_id = $1 WHERE id = $2 AND tenant_id = $3 AND client_id = $4 AND status = 'pending'`, assetID, discoveredID, tenantID, clientID)
	if err != nil {
		return uuid.Nil, err
	}
	if tag.RowsAffected() == 0 {
		return uuid.Nil, domain.ErrNotFound
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}
	return assetID, nil
}

// ── Monitoring ──────────────────────────────────────────────────────

func (r *MSPRepo) CreateMonitoringSource(ctx context.Context, s *domain.MonitoringSource) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_monitoring_sources (tenant_id, name, kind, base_url, webhook_token_hash, configured) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`,
		s.TenantID, s.Name, s.Kind, s.BaseURL, s.TokenHash, s.Configured).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt))
}

func (r *MSPRepo) ListMonitoringSources(ctx context.Context, tenantID uuid.UUID) ([]domain.MonitoringSource, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, name, kind, base_url, configured, created_at, updated_at FROM msp_monitoring_sources WHERE tenant_id = $1 ORDER BY name`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListMonitoringSources: %w", err)
	}
	defer rows.Close()
	var out []domain.MonitoringSource
	for rows.Next() {
		var s domain.MonitoringSource
		if err := rows.Scan(&s.ID, &s.TenantID, &s.Name, &s.Kind, &s.BaseURL, &s.Configured, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) FindSourceByTokenHash(ctx context.Context, tenantID uuid.UUID, tokenHash string) (*domain.MonitoringSource, error) {
	var s domain.MonitoringSource
	err := r.db.QueryRow(ctx, `SELECT id, tenant_id, name, kind, base_url, configured, created_at, updated_at FROM msp_monitoring_sources WHERE tenant_id = $1 AND webhook_token_hash = $2 AND webhook_token_hash <> ''`, tenantID, tokenHash).
		Scan(&s.ID, &s.TenantID, &s.Name, &s.Kind, &s.BaseURL, &s.Configured, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		// A wrong token is indistinguishable from a missing source.
		return nil, domain.ErrNotFound
	}
	return &s, nil
}

func (r *MSPRepo) RotateSourceToken(ctx context.Context, id, tenantID uuid.UUID, tokenHash string) error {
	tag, err := r.db.Exec(ctx, `UPDATE msp_monitoring_sources SET webhook_token_hash = $1, configured = TRUE, updated_at = NOW() WHERE id = $2 AND tenant_id = $3`, tokenHash, id, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanAlertRow(row pgx.Row) (*domain.MonitoringAlert, error) {
	var a domain.MonitoringAlert
	var ticketID pgtype.UUID
	err := row.Scan(&a.ID, &a.TenantID, &a.SourceID, &a.ClientID, &a.ExternalResource, &a.Severity, &a.Status, &a.Title, &a.Message, &ticketID, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if ticketID.Valid {
		id := uuid.UUID(ticketID.Bytes)
		a.TicketID = &id
	}
	return &a, nil
}

const alertColumns = `id, tenant_id, source_id, client_id, external_resource, severity, status, title, message, ticket_id, created_at, updated_at`

func (r *MSPRepo) OpenAlertForResource(ctx context.Context, tenantID, sourceID uuid.UUID, resource string) (*domain.MonitoringAlert, error) {
	a, err := scanAlertRow(r.db.QueryRow(ctx, `SELECT `+alertColumns+` FROM msp_monitoring_alerts WHERE tenant_id = $1 AND source_id = $2 AND external_resource = $3 AND status = 'open' LIMIT 1`, tenantID, sourceID, resource))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return a, nil
}

func (r *MSPRepo) CreateMonitoringAlert(ctx context.Context, a *domain.MonitoringAlert) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_monitoring_alerts (tenant_id, source_id, client_id, external_resource, severity, status, title, message, ticket_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`,
		a.TenantID, a.SourceID, a.ClientID, a.ExternalResource, a.Severity, a.Status, a.Title, a.Message, toPgUUID(a.TicketID)).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt))
}

func (r *MSPRepo) UpdateMonitoringAlert(ctx context.Context, a *domain.MonitoringAlert) error {
	tag, err := r.db.Exec(ctx, `UPDATE msp_monitoring_alerts SET severity = $1, status = $2, title = $3, message = $4, ticket_id = $5, updated_at = NOW() WHERE id = $6 AND tenant_id = $7`,
		a.Severity, a.Status, a.Title, a.Message, toPgUUID(a.TicketID), a.ID, a.TenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MSPRepo) ListMonitoringAlerts(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID) ([]domain.MonitoringAlert, error) {
	query := `SELECT ` + alertColumns + ` FROM msp_monitoring_alerts WHERE tenant_id = $1`
	args := []any{tenantID}
	next := 2
	if clause, cargs := clientClause("client_id", next, clientIDs); clause != "" {
		query += clause
		args = append(args, cargs...)
		next++
	}
	if clientFilter != nil {
		query += fmt.Sprintf(` AND client_id = $%d`, next)
		args = append(args, *clientFilter)
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListMonitoringAlerts: %w", err)
	}
	defer rows.Close()
	var out []domain.MonitoringAlert
	for rows.Next() {
		a, err := scanAlertRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return nonNil(out), rows.Err()
}

func (r *MSPRepo) GetMonitoringAlert(ctx context.Context, id, tenantID uuid.UUID) (*domain.MonitoringAlert, error) {
	a, err := scanAlertRow(r.db.QueryRow(ctx, `SELECT `+alertColumns+` FROM msp_monitoring_alerts WHERE id = $1 AND tenant_id = $2`, id, tenantID))
	return a, mspNotFound(err)
}

// ── Remote sessions ─────────────────────────────────────────────────

func scanRemoteSessionRow(row pgx.Row) (*domain.RemoteSession, error) {
	var s domain.RemoteSession
	var ticketID, authorizedBy pgtype.UUID
	var startedAt, endedAt pgtype.Timestamptz
	err := row.Scan(&s.ID, &s.TenantID, &s.ClientID, &ticketID, &s.TechnicianID, &s.TechnicianUserID, &s.Provider, &s.ExternalSessionID, &s.Status, &authorizedBy, &s.AuditNote, &startedAt, &endedAt, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if ticketID.Valid {
		id := uuid.UUID(ticketID.Bytes)
		s.TicketID = &id
	}
	if authorizedBy.Valid {
		id := uuid.UUID(authorizedBy.Bytes)
		s.AuthorizedBy = &id
	}
	s.StartedAt = fromPgTimestamptz(startedAt)
	s.EndedAt = fromPgTimestamptz(endedAt)
	return &s, nil
}

const remoteSessionColumns = `rs.id, rs.tenant_id, rs.client_id, rs.ticket_id, rs.technician_id, te.user_id, rs.provider, rs.external_session_id, rs.status, rs.authorized_by, rs.audit_note, rs.started_at, rs.ended_at, rs.created_at, rs.updated_at`

func (r *MSPRepo) CreateRemoteSession(ctx context.Context, s *domain.RemoteSession) error {
	return mspMapErr(r.db.QueryRow(ctx, `INSERT INTO msp_remote_sessions (tenant_id, client_id, ticket_id, technician_id, provider, external_session_id, status, audit_note) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`,
		s.TenantID, s.ClientID, toPgUUID(s.TicketID), s.TechnicianID, s.Provider, s.ExternalSessionID, s.Status, s.AuditNote).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt))
}

func (r *MSPRepo) GetRemoteSession(ctx context.Context, id, tenantID uuid.UUID) (*domain.RemoteSession, error) {
	s, err := scanRemoteSessionRow(r.db.QueryRow(ctx, `SELECT `+remoteSessionColumns+` FROM msp_remote_sessions rs JOIN msp_technicians te ON te.id = rs.technician_id WHERE rs.id = $1 AND rs.tenant_id = $2`, id, tenantID))
	return s, mspNotFound(err)
}

func (r *MSPRepo) UpdateRemoteSession(ctx context.Context, s *domain.RemoteSession) error {
	tag, err := r.db.Exec(ctx, `UPDATE msp_remote_sessions SET status = $1, authorized_by = $2, started_at = $3, ended_at = $4, updated_at = NOW() WHERE id = $5 AND tenant_id = $6`,
		s.Status, toPgUUID(s.AuthorizedBy), toPgTimestamptz(s.StartedAt), toPgTimestamptz(s.EndedAt), s.ID, s.TenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *MSPRepo) ListRemoteSessions(ctx context.Context, tenantID uuid.UUID, clientIDs []uuid.UUID, clientFilter *uuid.UUID) ([]domain.RemoteSession, error) {
	query := `SELECT ` + remoteSessionColumns + ` FROM msp_remote_sessions rs JOIN msp_technicians te ON te.id = rs.technician_id WHERE rs.tenant_id = $1`
	args := []any{tenantID}
	next := 2
	if clause, cargs := clientClause("rs.client_id", next, clientIDs); clause != "" {
		query += clause
		args = append(args, cargs...)
		next++
	}
	if clientFilter != nil {
		query += fmt.Sprintf(` AND rs.client_id = $%d`, next)
		args = append(args, *clientFilter)
	}
	query += ` ORDER BY rs.created_at DESC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("msp_repo.ListRemoteSessions: %w", err)
	}
	defer rows.Close()
	var out []domain.RemoteSession
	for rows.Next() {
		s, err := scanRemoteSessionRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *s)
	}
	return nonNil(out), rows.Err()
}

// ── Reports ─────────────────────────────────────────────────────────

func (r *MSPRepo) Report(ctx context.Context, tenantID uuid.UUID) (*domain.MSPReport, error) {
	rep := &domain.MSPReport{}
	queries := []struct {
		query string
		dst   *int
	}{
		{`SELECT COUNT(*) FROM msp_clients WHERE tenant_id = $1`, &rep.ClientsTotal},
		{`SELECT COUNT(*) FROM msp_clients WHERE tenant_id = $1 AND status = 'active'`, &rep.ActiveClients},
		{`SELECT COUNT(*) FROM msp_contracts WHERE tenant_id = $1 AND status = 'active'`, &rep.ActiveContracts},
		{`SELECT COUNT(*) FROM msp_monitoring_alerts WHERE tenant_id = $1 AND status = 'open'`, &rep.OpenAlerts},
		{`SELECT COUNT(*) FROM msp_remote_sessions WHERE tenant_id = $1 AND status IN ('authorized', 'active')`, &rep.ActiveSessions},
		{`SELECT COUNT(*) FROM tickets t JOIN msp_ticket_context tc ON tc.ticket_id = t.id WHERE t.tenant_id = $1 AND t.status IN ('open', 'in_progress', 'waiting_client')`, &rep.OpenMSPTickets},
	}
	for _, q := range queries {
		if err := r.db.QueryRow(ctx, q.query, tenantID).Scan(q.dst); err != nil {
			return nil, fmt.Errorf("msp_repo.Report: %w", err)
		}
	}
	return rep, nil
}
