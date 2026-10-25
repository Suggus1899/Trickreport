package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	appCX "github.com/trickreport/backend/internal/application/cx"
	domain "github.com/trickreport/backend/internal/domain/cx"
)

// CXRepo implements the CX repository boundary over PostgreSQL.
// Every query filters by tenant_id, and every customer-owned query
// additionally by customer_id: a record from another organization —
// or from another external customer of the same company — is
// indistinguishable from a missing record. Internal notes are
// filtered by the service (includeInternal), never by the UI alone.
type CXRepo struct {
	db *pgxpool.Pool
}

// NewCXRepo creates a CXRepo.
func NewCXRepo(db *pgxpool.Pool) *CXRepo { return &CXRepo{db: db} }

// DB exposes the pool for wiring that needs it (mirrors other repos).
func (r *CXRepo) DB() *pgxpool.Pool { return r.db }

// Compile-time assertion that CXRepo implements the application boundary.
var _ appCX.Repository = (*CXRepo)(nil)

func cxNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	return err
}

func cxMapErr(err error) error {
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

// ── Customers and contacts ──────────────────────────────────────────

func (r *CXRepo) ContactCustomerIDs(ctx context.Context, tenantID, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `SELECT customer_id FROM cx_customer_contacts WHERE tenant_id = $1 AND user_id = $2`, tenantID, userID)
	if err != nil {
		return nil, fmt.Errorf("cx_repo.ContactCustomerIDs: %w", err)
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

func scanCustomerRow(row pgx.Row) (*domain.Customer, error) {
	var c domain.Customer
	err := row.Scan(&c.ID, &c.TenantID, &c.Name, &c.Email, &c.Phone, &c.Company, &c.Status, &c.ConsentGiven, &c.ErasureRequested, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CXRepo) ListCustomers(ctx context.Context, tenantID uuid.UUID, customerIDs []uuid.UUID) ([]domain.Customer, error) {
	query := `SELECT id, tenant_id, name, email, phone, company, status, consent_given, erasure_requested, created_at, updated_at FROM cx_customers WHERE tenant_id = $1`
	clause, args := clientClause("id", 2, customerIDs)
	query += clause + ` ORDER BY name`
	rows, err := r.db.Query(ctx, query, append([]any{tenantID}, args...)...)
	if err != nil {
		return nil, fmt.Errorf("cx_repo.ListCustomers: %w", err)
	}
	defer rows.Close()
	var out []domain.Customer
	for rows.Next() {
		c, err := scanCustomerRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return nonNil(out), rows.Err()
}

func (r *CXRepo) GetCustomer(ctx context.Context, id, tenantID uuid.UUID) (*domain.Customer, error) {
	c, err := scanCustomerRow(r.db.QueryRow(ctx, `SELECT id, tenant_id, name, email, phone, company, status, consent_given, erasure_requested, created_at, updated_at FROM cx_customers WHERE id = $1 AND tenant_id = $2`, id, tenantID))
	return c, cxNotFound(err)
}

func (r *CXRepo) GetCustomerByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*domain.Customer, error) {
	c, err := scanCustomerRow(r.db.QueryRow(ctx, `SELECT id, tenant_id, name, email, phone, company, status, consent_given, erasure_requested, created_at, updated_at FROM cx_customers WHERE tenant_id = $1 AND LOWER(email) = $2 ORDER BY created_at LIMIT 1`, tenantID, email))
	return c, cxNotFound(err)
}

func (r *CXRepo) GetCustomerByPhone(ctx context.Context, tenantID uuid.UUID, phone string) (*domain.Customer, error) {
	c, err := scanCustomerRow(r.db.QueryRow(ctx, `SELECT id, tenant_id, name, email, phone, company, status, consent_given, erasure_requested, created_at, updated_at FROM cx_customers WHERE tenant_id = $1 AND phone = $2 ORDER BY created_at LIMIT 1`, tenantID, phone))
	return c, cxNotFound(err)
}

func (r *CXRepo) CreateCustomer(ctx context.Context, c *domain.Customer) error {
	return cxMapErr(r.db.QueryRow(ctx, `INSERT INTO cx_customers (tenant_id, name, email, phone, company, status, consent_given, erasure_requested) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`,
		c.TenantID, c.Name, c.Email, c.Phone, c.Company, c.Status, c.ConsentGiven, c.ErasureRequested).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt))
}

func (r *CXRepo) UpdateCustomer(ctx context.Context, c *domain.Customer) error {
	tag, err := r.db.Exec(ctx, `UPDATE cx_customers SET name = $1, email = $2, phone = $3, company = $4, status = $5, consent_given = $6, erasure_requested = $7, updated_at = NOW() WHERE id = $8 AND tenant_id = $9`,
		c.Name, c.Email, c.Phone, c.Company, c.Status, c.ConsentGiven, c.ErasureRequested, c.ID, c.TenantID)
	if err != nil {
		return cxMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CXRepo) DeleteCustomer(ctx context.Context, id, tenantID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM cx_customers WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return cxMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanCustomerContactRow(row pgx.Row) (*domain.CustomerContact, error) {
	var c domain.CustomerContact
	err := row.Scan(&c.ID, &c.TenantID, &c.CustomerID, &c.UserID, &c.Name, &c.Email, &c.Role, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CXRepo) ListContacts(ctx context.Context, tenantID, customerID uuid.UUID) ([]domain.CustomerContact, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, customer_id, user_id, name, email, role, created_at FROM cx_customer_contacts WHERE tenant_id = $1 AND customer_id = $2 ORDER BY name`, tenantID, customerID)
	if err != nil {
		return nil, fmt.Errorf("cx_repo.ListContacts: %w", err)
	}
	defer rows.Close()
	var out []domain.CustomerContact
	for rows.Next() {
		c, err := scanCustomerContactRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return nonNil(out), rows.Err()
}

func (r *CXRepo) CreateContact(ctx context.Context, c *domain.CustomerContact) error {
	return cxMapErr(r.db.QueryRow(ctx, `INSERT INTO cx_customer_contacts (tenant_id, customer_id, user_id, name, email, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		c.TenantID, c.CustomerID, c.UserID, c.Name, c.Email, c.Role).Scan(&c.ID, &c.CreatedAt))
}

func (r *CXRepo) DeleteContact(ctx context.Context, id, tenantID, customerID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM cx_customer_contacts WHERE id = $1 AND tenant_id = $2 AND customer_id = $3`, id, tenantID, customerID)
	if err != nil {
		return cxMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CXRepo) UserExists(ctx context.Context, id, tenantID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND tenant_id = $2)`, id, tenantID).Scan(&exists)
	return exists, err
}

// ── Conversations and messages ──────────────────────────────────────

const cxConversationColumns = `c.id, c.tenant_id, c.customer_id, COALESCE(cu.name, ''), c.channel, c.subject, c.status, c.priority, c.assigned_to, COALESCE(u.name, ''), c.created_by, c.external_thread_id, c.sla_response_minutes, c.sla_resolution_minutes, c.response_due_at, c.resolution_due_at, c.first_response_at, c.escalated, c.escalated_at, c.created_at, c.updated_at`

func scanConversationRow(row pgx.Row) (*domain.Conversation, error) {
	var c domain.Conversation
	err := row.Scan(&c.ID, &c.TenantID, &c.CustomerID, &c.CustomerName, &c.Channel, &c.Subject, &c.Status, &c.Priority, &c.AssignedTo, &c.AssigneeName, &c.CreatedBy, &c.ExternalThreadID, &c.SLAResponseMinutes, &c.SLAResolutionMinutes, &c.ResponseDueAt, &c.ResolutionDueAt, &c.FirstResponseAt, &c.Escalated, &c.EscalatedAt, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func cxConversationFrom() string {
	return ` FROM cx_conversations c LEFT JOIN cx_customers cu ON cu.id = c.customer_id LEFT JOIN users u ON u.id = c.assigned_to`
}

func (r *CXRepo) ListConversations(ctx context.Context, tenantID uuid.UUID, customerIDs []uuid.UUID, customerFilter *uuid.UUID, channel, status string, assigneeFilter *uuid.UUID) ([]domain.Conversation, error) {
	query := `SELECT ` + cxConversationColumns + cxConversationFrom() + ` WHERE c.tenant_id = $1`
	args := []any{tenantID}
	next := 2
	clause, cargs := clientClause("c.customer_id", next, customerIDs)
	query += clause
	args = append(args, cargs...)
	next += len(cargs)
	if customerFilter != nil {
		query += fmt.Sprintf(` AND c.customer_id = $%d`, next)
		args = append(args, *customerFilter)
		next++
	}
	if channel != "" {
		query += fmt.Sprintf(` AND c.channel = $%d`, next)
		args = append(args, channel)
		next++
	}
	if status != "" {
		query += fmt.Sprintf(` AND c.status = $%d`, next)
		args = append(args, status)
		next++
	}
	if assigneeFilter != nil {
		query += fmt.Sprintf(` AND c.assigned_to = $%d`, next)
		args = append(args, *assigneeFilter)
		next++
	}
	query += ` ORDER BY c.created_at DESC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cx_repo.ListConversations: %w", err)
	}
	defer rows.Close()
	var out []domain.Conversation
	for rows.Next() {
		c, err := scanConversationRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return nonNil(out), rows.Err()
}

func (r *CXRepo) GetConversation(ctx context.Context, id, tenantID uuid.UUID) (*domain.Conversation, error) {
	c, err := scanConversationRow(r.db.QueryRow(ctx, `SELECT `+cxConversationColumns+cxConversationFrom()+` WHERE c.id = $1 AND c.tenant_id = $2`, id, tenantID))
	return c, cxNotFound(err)
}

func (r *CXRepo) OpenConversationForThread(ctx context.Context, tenantID, customerID uuid.UUID, threadID string) (*domain.Conversation, error) {
	c, err := scanConversationRow(r.db.QueryRow(ctx, `SELECT `+cxConversationColumns+cxConversationFrom()+` WHERE c.tenant_id = $1 AND c.customer_id = $2 AND c.external_thread_id = $3 AND c.status <> 'closed' ORDER BY c.created_at DESC LIMIT 1`, tenantID, customerID, threadID))
	return c, cxNotFound(err)
}

func (r *CXRepo) CreateConversation(ctx context.Context, c *domain.Conversation) error {
	return cxMapErr(r.db.QueryRow(ctx, `INSERT INTO cx_conversations (tenant_id, customer_id, channel, subject, status, priority, assigned_to, created_by, external_thread_id, sla_response_minutes, sla_resolution_minutes, response_due_at, resolution_due_at, first_response_at, escalated, escalated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id, created_at, updated_at`,
		c.TenantID, c.CustomerID, c.Channel, c.Subject, c.Status, c.Priority, c.AssignedTo, c.CreatedBy, c.ExternalThreadID, c.SLAResponseMinutes, c.SLAResolutionMinutes, c.ResponseDueAt, c.ResolutionDueAt, c.FirstResponseAt, c.Escalated, c.EscalatedAt).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt))
}

func (r *CXRepo) UpdateConversation(ctx context.Context, c *domain.Conversation) error {
	tag, err := r.db.Exec(ctx, `UPDATE cx_conversations SET subject = $1, status = $2, priority = $3, assigned_to = $4, external_thread_id = $5, sla_response_minutes = $6, sla_resolution_minutes = $7, response_due_at = $8, resolution_due_at = $9, first_response_at = $10, escalated = $11, escalated_at = $12, updated_at = NOW() WHERE id = $13 AND tenant_id = $14`,
		c.Subject, c.Status, c.Priority, c.AssignedTo, c.ExternalThreadID, c.SLAResponseMinutes, c.SLAResolutionMinutes, c.ResponseDueAt, c.ResolutionDueAt, c.FirstResponseAt, c.Escalated, c.EscalatedAt, c.ID, c.TenantID)
	if err != nil {
		return cxMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *CXRepo) AssignConversation(ctx context.Context, id, tenantID uuid.UUID, assignee *uuid.UUID, status string) error {
	tag, err := r.db.Exec(ctx, `UPDATE cx_conversations SET assigned_to = $1, status = $2, updated_at = NOW() WHERE id = $3 AND tenant_id = $4`, assignee, status, id, tenantID)
	if err != nil {
		return cxMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanMessageRow(row pgx.Row) (*domain.Message, error) {
	var m domain.Message
	err := row.Scan(&m.ID, &m.TenantID, &m.ConversationID, &m.SenderType, &m.SenderID, &m.SenderName, &m.Body, &m.IsInternal, &m.Channel, &m.ExternalMessageID, &m.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *CXRepo) ListMessages(ctx context.Context, tenantID, conversationID uuid.UUID, includeInternal bool) ([]domain.Message, error) {
	query := `SELECT m.id, m.tenant_id, m.conversation_id, m.sender_type, m.sender_id, COALESCE(u.name, ''), m.body, m.is_internal, m.channel, m.external_message_id, m.created_at FROM cx_messages m LEFT JOIN users u ON u.id = m.sender_id WHERE m.tenant_id = $1 AND m.conversation_id = $2`
	if !includeInternal {
		query += ` AND m.is_internal = FALSE`
	}
	query += ` ORDER BY m.created_at ASC`
	rows, err := r.db.Query(ctx, query, tenantID, conversationID)
	if err != nil {
		return nil, fmt.Errorf("cx_repo.ListMessages: %w", err)
	}
	defer rows.Close()
	var out []domain.Message
	for rows.Next() {
		m, err := scanMessageRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *m)
	}
	return nonNil(out), rows.Err()
}

func (r *CXRepo) CreateMessage(ctx context.Context, m *domain.Message) error {
	return cxMapErr(r.db.QueryRow(ctx, `INSERT INTO cx_messages (tenant_id, conversation_id, sender_type, sender_id, body, is_internal, channel, external_message_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at`,
		m.TenantID, m.ConversationID, m.SenderType, m.SenderID, m.Body, m.IsInternal, m.Channel, m.ExternalMessageID).Scan(&m.ID, &m.CreatedAt))
}

func (r *CXRepo) MessageExistsByExternalID(ctx context.Context, tenantID uuid.UUID, externalID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM cx_messages WHERE tenant_id = $1 AND external_message_id = $2)`, tenantID, externalID).Scan(&exists)
	return exists, err
}

// ── Channel configuration ───────────────────────────────────────────

func scanEmailAccountRow(row pgx.Row) (*domain.EmailAccount, error) {
	var a domain.EmailAccount
	err := row.Scan(&a.ID, &a.TenantID, &a.Name, &a.SMTPHost, &a.SMTPPort, &a.SMTPUsername, &a.SMTPFrom, &a.IMAPHost, &a.IMAPPort, &a.IMAPUsername, &a.Active, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *CXRepo) ListEmailAccounts(ctx context.Context, tenantID uuid.UUID) ([]domain.EmailAccount, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, name, smtp_host, smtp_port, smtp_username, smtp_from, imap_host, imap_port, imap_username, active, created_at, updated_at FROM cx_email_accounts WHERE tenant_id = $1 ORDER BY name`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("cx_repo.ListEmailAccounts: %w", err)
	}
	defer rows.Close()
	var out []domain.EmailAccount
	for rows.Next() {
		a, err := scanEmailAccountRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return nonNil(out), rows.Err()
}

func (r *CXRepo) CreateEmailAccount(ctx context.Context, a *domain.EmailAccount) error {
	return cxMapErr(r.db.QueryRow(ctx, `INSERT INTO cx_email_accounts (tenant_id, name, smtp_host, smtp_port, smtp_username, smtp_from, imap_host, imap_port, imap_username, active) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at, updated_at`,
		a.TenantID, a.Name, a.SMTPHost, a.SMTPPort, a.SMTPUsername, a.SMTPFrom, a.IMAPHost, a.IMAPPort, a.IMAPUsername, a.Active).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt))
}

func (r *CXRepo) UpdateEmailAccount(ctx context.Context, a *domain.EmailAccount) error {
	tag, err := r.db.Exec(ctx, `UPDATE cx_email_accounts SET name = $1, smtp_host = $2, smtp_port = $3, smtp_username = $4, smtp_from = $5, imap_host = $6, imap_port = $7, imap_username = $8, active = $9, updated_at = NOW() WHERE id = $10 AND tenant_id = $11`,
		a.Name, a.SMTPHost, a.SMTPPort, a.SMTPUsername, a.SMTPFrom, a.IMAPHost, a.IMAPPort, a.IMAPUsername, a.Active, a.ID, a.TenantID)
	if err != nil {
		return cxMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanWhatsAppAccountRow(row pgx.Row) (*domain.WhatsAppAccount, error) {
	var a domain.WhatsAppAccount
	err := row.Scan(&a.ID, &a.TenantID, &a.Name, &a.BusinessID, &a.PhoneNumberID, &a.DisplayName, &a.WebhookVerifyTokenHash, &a.Active, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *CXRepo) GetWhatsAppAccount(ctx context.Context, id, tenantID uuid.UUID) (*domain.WhatsAppAccount, error) {
	a, err := scanWhatsAppAccountRow(r.db.QueryRow(ctx, `SELECT id, tenant_id, name, business_id, phone_number_id, display_name, webhook_verify_token_hash, active, created_at, updated_at FROM cx_whatsapp_accounts WHERE id = $1 AND tenant_id = $2`, id, tenantID))
	return a, cxNotFound(err)
}

func (r *CXRepo) ListWhatsAppAccounts(ctx context.Context, tenantID uuid.UUID) ([]domain.WhatsAppAccount, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, name, business_id, phone_number_id, display_name, webhook_verify_token_hash, active, created_at, updated_at FROM cx_whatsapp_accounts WHERE tenant_id = $1 ORDER BY name`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("cx_repo.ListWhatsAppAccounts: %w", err)
	}
	defer rows.Close()
	var out []domain.WhatsAppAccount
	for rows.Next() {
		a, err := scanWhatsAppAccountRow(rows)
		if err != nil {
			return nil, err
		}
		// Never expose the verify-token hash through a list call.
		a.WebhookVerifyTokenHash = ""
		out = append(out, *a)
	}
	return nonNil(out), rows.Err()
}

func (r *CXRepo) FindWhatsAppAccountByTokenHash(ctx context.Context, tenantID uuid.UUID, tokenHash string) (*domain.WhatsAppAccount, error) {
	a, err := scanWhatsAppAccountRow(r.db.QueryRow(ctx, `SELECT id, tenant_id, name, business_id, phone_number_id, display_name, webhook_verify_token_hash, active, created_at, updated_at FROM cx_whatsapp_accounts WHERE tenant_id = $1 AND webhook_verify_token_hash = $2 AND active = TRUE`, tenantID, tokenHash))
	return a, cxNotFound(err)
}

func (r *CXRepo) CreateWhatsAppAccount(ctx context.Context, a *domain.WhatsAppAccount) error {
	return cxMapErr(r.db.QueryRow(ctx, `INSERT INTO cx_whatsapp_accounts (tenant_id, name, business_id, phone_number_id, display_name, webhook_verify_token_hash, active) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`,
		a.TenantID, a.Name, a.BusinessID, a.PhoneNumberID, a.DisplayName, a.WebhookVerifyTokenHash, a.Active).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt))
}

func (r *CXRepo) UpdateWhatsAppAccount(ctx context.Context, a *domain.WhatsAppAccount) error {
	tag, err := r.db.Exec(ctx, `UPDATE cx_whatsapp_accounts SET name = $1, business_id = $2, phone_number_id = $3, display_name = $4, active = $5, updated_at = NOW() WHERE id = $6 AND tenant_id = $7`,
		a.Name, a.BusinessID, a.PhoneNumberID, a.DisplayName, a.Active, a.ID, a.TenantID)
	if err != nil {
		return cxMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ── SLA policies, assignment rules, auto-replies ────────────────────

func scanSLAPolicyRow(row pgx.Row) (*domain.SLAPolicy, error) {
	var p domain.SLAPolicy
	err := row.Scan(&p.ID, &p.TenantID, &p.Channel, &p.Priority, &p.ResponseMinutes, &p.ResolutionMinutes, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *CXRepo) ListSLAPolicies(ctx context.Context, tenantID uuid.UUID) ([]domain.SLAPolicy, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, channel, priority, response_minutes, resolution_minutes, created_at, updated_at FROM cx_sla_policies WHERE tenant_id = $1 ORDER BY channel, priority`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("cx_repo.ListSLAPolicies: %w", err)
	}
	defer rows.Close()
	var out []domain.SLAPolicy
	for rows.Next() {
		p, err := scanSLAPolicyRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return nonNil(out), rows.Err()
}

func (r *CXRepo) UpsertSLAPolicy(ctx context.Context, p *domain.SLAPolicy) error {
	return cxMapErr(r.db.QueryRow(ctx, `INSERT INTO cx_sla_policies (tenant_id, channel, priority, response_minutes, resolution_minutes) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (tenant_id, channel, priority) DO UPDATE SET response_minutes = EXCLUDED.response_minutes, resolution_minutes = EXCLUDED.resolution_minutes, updated_at = NOW() RETURNING id, created_at, updated_at`,
		p.TenantID, p.Channel, p.Priority, p.ResponseMinutes, p.ResolutionMinutes).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt))
}

func scanAssignmentRuleRow(row pgx.Row) (*domain.AssignmentRule, error) {
	var r2 domain.AssignmentRule
	err := row.Scan(&r2.ID, &r2.TenantID, &r2.Channel, &r2.Priority, &r2.AssigneeUserID, &r2.Position, &r2.Active, &r2.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &r2, nil
}

func (r *CXRepo) ListAssignmentRules(ctx context.Context, tenantID uuid.UUID) ([]domain.AssignmentRule, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, channel, priority, assignee_user_id, position, active, created_at FROM cx_assignment_rules WHERE tenant_id = $1 ORDER BY position, created_at`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("cx_repo.ListAssignmentRules: %w", err)
	}
	defer rows.Close()
	var out []domain.AssignmentRule
	for rows.Next() {
		rule, err := scanAssignmentRuleRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *rule)
	}
	return nonNil(out), rows.Err()
}

func (r *CXRepo) CreateAssignmentRule(ctx context.Context, rule *domain.AssignmentRule) error {
	return cxMapErr(r.db.QueryRow(ctx, `INSERT INTO cx_assignment_rules (tenant_id, channel, priority, assignee_user_id, position, active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		rule.TenantID, rule.Channel, rule.Priority, rule.AssigneeUserID, rule.Position, rule.Active).Scan(&rule.ID, &rule.CreatedAt))
}

func (r *CXRepo) DeleteAssignmentRule(ctx context.Context, id, tenantID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM cx_assignment_rules WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return cxMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanAutoReplyRow(row pgx.Row) (*domain.AutoReply, error) {
	var a domain.AutoReply
	err := row.Scan(&a.ID, &a.TenantID, &a.TriggerKind, &a.Keyword, &a.Channel, &a.ResponseText, &a.Active, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *CXRepo) ListAutoReplies(ctx context.Context, tenantID uuid.UUID) ([]domain.AutoReply, error) {
	rows, err := r.db.Query(ctx, `SELECT id, tenant_id, trigger_kind, keyword, channel, response_text, active, created_at FROM cx_auto_replies WHERE tenant_id = $1 ORDER BY created_at`, tenantID)
	if err != nil {
		return nil, fmt.Errorf("cx_repo.ListAutoReplies: %w", err)
	}
	defer rows.Close()
	var out []domain.AutoReply
	for rows.Next() {
		a, err := scanAutoReplyRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *a)
	}
	return nonNil(out), rows.Err()
}

func (r *CXRepo) CreateAutoReply(ctx context.Context, a *domain.AutoReply) error {
	return cxMapErr(r.db.QueryRow(ctx, `INSERT INTO cx_auto_replies (tenant_id, trigger_kind, keyword, channel, response_text, active) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		a.TenantID, a.TriggerKind, a.Keyword, a.Channel, a.ResponseText, a.Active).Scan(&a.ID, &a.CreatedAt))
}

func (r *CXRepo) DeleteAutoReply(ctx context.Context, id, tenantID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM cx_auto_replies WHERE id = $1 AND tenant_id = $2`, id, tenantID)
	if err != nil {
		return cxMapErr(err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ── Surveys and reports ─────────────────────────────────────────────

func scanSurveyRow(row pgx.Row) (*domain.Survey, error) {
	var s domain.Survey
	err := row.Scan(&s.ID, &s.TenantID, &s.ConversationID, &s.CustomerID, &s.Rating, &s.Comment, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *CXRepo) CreateSurvey(ctx context.Context, s *domain.Survey) error {
	return cxMapErr(r.db.QueryRow(ctx, `INSERT INTO cx_surveys (tenant_id, conversation_id, customer_id, rating, comment) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`,
		s.TenantID, s.ConversationID, s.CustomerID, s.Rating, s.Comment).Scan(&s.ID, &s.CreatedAt))
}

func (r *CXRepo) ListSurveys(ctx context.Context, tenantID uuid.UUID, customerIDs []uuid.UUID, customerFilter *uuid.UUID) ([]domain.Survey, error) {
	query := `SELECT id, tenant_id, conversation_id, customer_id, rating, comment, created_at FROM cx_surveys WHERE tenant_id = $1`
	args := []any{tenantID}
	next := 2
	clause, cargs := clientClause("customer_id", next, customerIDs)
	query += clause
	args = append(args, cargs...)
	next += len(cargs)
	if customerFilter != nil {
		query += fmt.Sprintf(` AND customer_id = $%d`, next)
		args = append(args, *customerFilter)
		next++
	}
	query += ` ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("cx_repo.ListSurveys: %w", err)
	}
	defer rows.Close()
	var out []domain.Survey
	for rows.Next() {
		s, err := scanSurveyRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *s)
	}
	return nonNil(out), rows.Err()
}

func (r *CXRepo) Report(ctx context.Context, tenantID uuid.UUID, customerIDs []uuid.UUID) (*domain.CXReport, error) {
	report := &domain.CXReport{ConversationsByChannel: map[string]int{}}
	// Customer scope clauses shared by the queries below. Customers
	// are scoped by id; conversations and surveys by customer_id.
	custClause, cargs := clientClause("id", 2, customerIDs)
	custArgs := append([]any{tenantID}, cargs...)
	clause, cargs := clientClause("customer_id", 2, customerIDs)
	convArgs := append([]any{tenantID}, cargs...)
	var err error
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM cx_customers WHERE tenant_id = $1`+custClause, custArgs...).Scan(&report.CustomersTotal)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM cx_conversations WHERE tenant_id = $1`+clause+` AND status IN ('open', 'in_progress', 'waiting_customer')`, convArgs...).Scan(&report.OpenConversations)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM cx_conversations WHERE tenant_id = $1`+clause+` AND status IN ('resolved', 'closed')`, convArgs...).Scan(&report.ResolvedConversations)
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FROM cx_conversations WHERE tenant_id = $1`+clause+` AND status NOT IN ('resolved', 'closed') AND ((resolution_due_at IS NOT NULL AND resolution_due_at < NOW()) OR (first_response_at IS NULL AND response_due_at IS NOT NULL AND response_due_at < NOW()))`, convArgs...).Scan(&report.SLABreached)
	if err != nil {
		return nil, err
	}
	surveyClause, sargs := clientClause("customer_id", 2, customerIDs)
	surveyArgs := append([]any{tenantID}, sargs...)
	err = r.db.QueryRow(ctx, `SELECT COUNT(*), COALESCE(AVG(rating), 0) FROM cx_surveys WHERE tenant_id = $1`+surveyClause, surveyArgs...).Scan(&report.SurveyCount, &report.AvgSatisfaction)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, `SELECT channel, COUNT(*) FROM cx_conversations WHERE tenant_id = $1`+clause+` GROUP BY channel`, convArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var channel string
		var count int
		if err := rows.Scan(&channel, &count); err != nil {
			return nil, err
		}
		report.ConversationsByChannel[channel] = count
	}
	return report, rows.Err()
}
