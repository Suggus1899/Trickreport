package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	domainuser "github.com/trickreport/backend/internal/domain/user"
)

// AuthUserRepo implements auth.UserRepository for the auth service.
// Unlike UserRepo, it does not scope by tenant_id because user IDs are
// globally unique UUIDs and the auth flow identifies users by email/ID alone.
type AuthUserRepo struct {
	db *pgxpool.Pool
}

func NewAuthUserRepo(db *pgxpool.Pool) *AuthUserRepo {
	return &AuthUserRepo{db: db}
}

func (r *AuthUserRepo) GetByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	const q = `
		SELECT id, tenant_id, name, email, role, COALESCE(password, ''), COALESCE(ldap_dn, ''),
		       avatar_url, active, created_at, updated_at
		FROM users
		WHERE email = $1 AND active = TRUE
		LIMIT 1`

	var u domainuser.User
	err := r.db.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.TenantID, &u.Name, &u.Email, &u.Role, &u.PasswordHash, &u.LDAPDN,
		&u.AvatarURL, &u.Active, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appAuth.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &u, nil
}

func (r *AuthUserRepo) GetByID(ctx context.Context, id string) (*domainuser.User, error) {
	const q = `
		SELECT id, tenant_id, name, email, role, COALESCE(password, ''), COALESCE(ldap_dn, ''),
		       avatar_url, active, created_at, updated_at
		FROM users
		WHERE id = $1 AND active = TRUE`

	var u domainuser.User
	err := r.db.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.TenantID, &u.Name, &u.Email, &u.Role, &u.PasswordHash, &u.LDAPDN,
		&u.AvatarURL, &u.Active, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainuser.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &u, nil
}

// Compile-time assertion that AuthUserRepo implements auth.UserRepository.
var _ appAuth.UserRepository = (*AuthUserRepo)(nil)
