package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	appAuth "github.com/trickreport/backend/internal/application/auth"
	"github.com/trickreport/backend/internal/application/user"
	domainuser "github.com/trickreport/backend/internal/domain/user"
)

// UserRepo implements both user.Repository and auth.UserRepository.
//
// NOTE: The auth.UserRepository interface declares GetByIDNoTenant(ctx, id)
// with a single id argument, while user.Repository declares GetByID(ctx, id,
// tenantID) with two arguments. Go does not permit two methods with the same
// name on a single struct, so the tenant-less lookup is named GetByIDNoTenant.
// GetByEmail is shared and matches both interfaces.
type UserRepo struct {
	db *pgxpool.Pool
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

// Compile-time assertion that UserRepo implements user.Repository.
var _ user.Repository = (*UserRepo)(nil)

// Compile-time assertion that UserRepo implements auth.UserRepository.
var _ appAuth.UserRepository = (*UserRepo)(nil)

// List returns all users for a tenant, ordered by created_at desc.
func (r *UserRepo) List(ctx context.Context, tenantID uuid.UUID) ([]domainuser.User, error) {
	const q = `SELECT id, tenant_id, name, email, role, avatar_url, active, created_at FROM users WHERE tenant_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, q, tenantID)
	if err != nil {
		return nil, fmt.Errorf("user_repo.List: query: %w", err)
	}
	defer rows.Close()

	var users []domainuser.User
	for rows.Next() {
		var u domainuser.User
		var role string
		var id, tid pgtype.UUID
		if err := rows.Scan(&id, &tid, &u.Name, &u.Email, &role, &u.AvatarURL, &u.Active, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("user_repo.List: scan: %w", err)
		}
		u.ID = pgToUUID(id)
		u.TenantID = pgToUUID(tid)
		u.Role = domainuser.Role(role)
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("user_repo.List: rows: %w", err)
	}
	return users, nil
}

// GetByID returns a user by id within a tenant.
func (r *UserRepo) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*domainuser.User, error) {
	const q = `SELECT id, tenant_id, name, email, role, COALESCE(password, ''), COALESCE(ldap_dn, ''), avatar_url, active, created_at, updated_at, COALESCE(failed_login_attempts, 0), COALESCE(locked_until, timestamp 'epoch'), COALESCE(mfa_secret, ''), COALESCE(mfa_enabled, false) FROM users WHERE id = $1 AND tenant_id = $2`

	return r.scanUser(ctx, q, id, tenantID)
}

// GetByIDNoTenant returns an active user by id without tenant scoping.
func (r *UserRepo) GetByIDNoTenant(ctx context.Context, id uuid.UUID) (*domainuser.User, error) {
	const q = `SELECT id, tenant_id, name, email, role, COALESCE(password, ''), COALESCE(ldap_dn, ''), avatar_url, active, created_at, updated_at, COALESCE(failed_login_attempts, 0), COALESCE(locked_until, timestamp 'epoch'), COALESCE(mfa_secret, ''), COALESCE(mfa_enabled, false) FROM users WHERE id = $1 AND active = TRUE`

	return r.scanUser(ctx, q, id)
}

// GetByEmail returns an active user by email.
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	const q = `SELECT id, tenant_id, name, email, role, COALESCE(password, ''), COALESCE(ldap_dn, ''), avatar_url, active, created_at, updated_at, COALESCE(failed_login_attempts, 0), COALESCE(locked_until, timestamp 'epoch'), COALESCE(mfa_secret, ''), COALESCE(mfa_enabled, false) FROM users WHERE email = $1 AND active = TRUE LIMIT 1`

	return r.scanUser(ctx, q, email)
}

// Create inserts a new user. If passwordHash is empty, nil is passed.
func (r *UserRepo) Create(ctx context.Context, u *domainuser.User, passwordHash string) error {
	const q = `INSERT INTO users (tenant_id, name, email, role, password) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	var pw any
	if passwordHash == "" {
		pw = nil
	} else {
		pw = passwordHash
	}

	var id pgtype.UUID
	if err := r.db.QueryRow(ctx, q, u.TenantID, u.Name, u.Email, string(u.Role), pw).
		Scan(&id, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return fmt.Errorf("user_repo.Create: %w", err)
	}
	u.ID = pgToUUID(id)
	return nil
}

// Update updates the mutable fields of a user.
func (r *UserRepo) Update(ctx context.Context, id, tenantID uuid.UUID, fields user.UpdateFields) (*domainuser.User, error) {
	const q = `UPDATE users SET name = COALESCE($1, name), role = COALESCE($2, role::text)::user_role, avatar_url = COALESCE($3, avatar_url), active = COALESCE($4, active), updated_at = NOW() WHERE id = $5 AND tenant_id = $6 RETURNING id, tenant_id, name, email, role, avatar_url, active, created_at, updated_at`

	var roleStr *string
	if fields.Role != nil {
		s := string(*fields.Role)
		roleStr = &s
	}

	var u domainuser.User
	var role string
	var idCol, tid pgtype.UUID
	err := r.db.QueryRow(ctx, q, fields.Name, roleStr, fields.AvatarURL, fields.Active, id, tenantID).
		Scan(&idCol, &tid, &u.Name, &u.Email, &role, &u.AvatarURL, &u.Active, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainuser.ErrNotFound
		}
		return nil, fmt.Errorf("user_repo.Update: %w", err)
	}
	u.ID = pgToUUID(idCol)
	u.TenantID = pgToUUID(tid)
	u.Role = domainuser.Role(role)
	return &u, nil
}

// Deactivate sets active=false for a user.
func (r *UserRepo) Deactivate(ctx context.Context, id, tenantID uuid.UUID) error {
	const q = `UPDATE users SET active = FALSE, updated_at = NOW() WHERE id = $1 AND tenant_id = $2`

	ct, err := r.db.Exec(ctx, q, id, tenantID)
	if err != nil {
		return fmt.Errorf("user_repo.Deactivate: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return domainuser.ErrNotFound
	}
	return nil
}

// scanUser is a helper that scans a single user row with the full column set.
func (r *UserRepo) scanUser(ctx context.Context, query string, args ...any) (*domainuser.User, error) {
	var u domainuser.User
	var role string
	var id, tid pgtype.UUID
	var lockedUntil pgtype.Timestamptz
	err := r.db.QueryRow(ctx, query, args...).
		Scan(&id, &tid, &u.Name, &u.Email, &role, &u.PasswordHash, &u.LDAPDN, &u.AvatarURL, &u.Active, &u.CreatedAt, &u.UpdatedAt, &u.FailedLoginAttempts, &lockedUntil, &u.MFASecret, &u.MFAEnabled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainuser.ErrNotFound
		}
		return nil, fmt.Errorf("user_repo: scan: %w", err)
	}
	u.ID = pgToUUID(id)
	u.TenantID = pgToUUID(tid)
	u.Role = domainuser.Role(role)
	if lockedUntil.Valid && !lockedUntil.Time.IsZero() {
		u.LockedUntil = lockedUntil.Time
	}
	return &u, nil
}

// --- Account Lockout ---

// IncrementFailedAttempts increments the failed login attempt counter for a user.
func (r *UserRepo) IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	const q = `UPDATE users SET failed_login_attempts = COALESCE(failed_login_attempts, 0) + 1, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("user_repo.IncrementFailedAttempts: %w", err)
	}
	return nil
}

// ResetFailedAttempts resets the failed login attempt counter to zero.
func (r *UserRepo) ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	const q = `UPDATE users SET failed_login_attempts = 0, locked_until = NULL, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("user_repo.ResetFailedAttempts: %w", err)
	}
	return nil
}

// LockAccount sets the locked_until timestamp for a user.
func (r *UserRepo) LockAccount(ctx context.Context, userID uuid.UUID, until time.Time) error {
	const q = `UPDATE users SET locked_until = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, userID, until)
	if err != nil {
		return fmt.Errorf("user_repo.LockAccount: %w", err)
	}
	return nil
}

// --- Password Update ---

// UpdatePassword sets a new password hash for a user.
func (r *UserRepo) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	const q = `UPDATE users SET password = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("user_repo.UpdatePassword: %w", err)
	}
	return nil
}

// --- MFA ---

// GetMFASecret returns the MFA secret and enabled status for a user.
func (r *UserRepo) GetMFASecret(ctx context.Context, userID uuid.UUID) (string, bool, error) {
	const q = `SELECT COALESCE(mfa_secret, ''), COALESCE(mfa_enabled, false) FROM users WHERE id = $1`
	var secret string
	var enabled bool
	err := r.db.QueryRow(ctx, q, userID).Scan(&secret, &enabled)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, domainuser.ErrNotFound
		}
		return "", false, fmt.Errorf("user_repo.GetMFASecret: %w", err)
	}
	return secret, enabled, nil
}

// SetMFASecret stores the TOTP secret for a user (MFA not yet enabled).
func (r *UserRepo) SetMFASecret(ctx context.Context, userID uuid.UUID, secret string) error {
	const q = `UPDATE users SET mfa_secret = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, userID, secret)
	if err != nil {
		return fmt.Errorf("user_repo.SetMFASecret: %w", err)
	}
	return nil
}

// EnableMFA marks MFA as enabled for a user.
func (r *UserRepo) EnableMFA(ctx context.Context, userID uuid.UUID) error {
	const q = `UPDATE users SET mfa_enabled = TRUE, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("user_repo.EnableMFA: %w", err)
	}
	return nil
}

// DisableMFA clears the MFA secret and disables MFA for a user.
func (r *UserRepo) DisableMFA(ctx context.Context, userID uuid.UUID) error {
	const q = `UPDATE users SET mfa_enabled = FALSE, mfa_secret = NULL, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("user_repo.DisableMFA: %w", err)
	}
	return nil
}

// Compile-time assertions that UserRepo implements the new auth interfaces.
var (
	_ appAuth.AccountLockoutRepository = (*UserRepo)(nil)
	_ appAuth.PasswordUpdater          = (*UserRepo)(nil)
	_ appAuth.MFARepository            = (*UserRepo)(nil)
)
