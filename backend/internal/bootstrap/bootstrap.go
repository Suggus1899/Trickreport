// Package bootstrap ensures the default tenant and an initial admin user
// exist on first startup. The admin credentials come from environment
// variables (ADMIN_EMAIL / ADMIN_PASSWORD) — never from the repository.
package bootstrap

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var (
	defaultTenantID   = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	defaultTenantSlug = "default"
)

// EnsureAdmin creates the initial admin user from the given credentials if
// no users exist yet. It is safe to call on every startup — it is a no-op
// once any user exists.
func EnsureAdmin(ctx context.Context, pool *pgxpool.Pool, adminEmail, adminPassword string) error {
	if adminEmail == "" || adminPassword == "" {
		log.Warn().Msg("ADMIN_EMAIL/ADMIN_PASSWORD not set — skipping admin bootstrap")
		return nil
	}

	var userCount int
	err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&userCount)
	if err != nil {
		return fmt.Errorf("bootstrap: failed to count users: %w", err)
	}
	if userCount > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bootstrap: failed to hash admin password: %w", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO users (tenant_id, name, email, role, password)
		VALUES ($1, 'Super Admin', $2, 'admin', $3)`,
		defaultTenantID, adminEmail, string(hash),
	)
	if err != nil {
		return fmt.Errorf("bootstrap: failed to create admin user: %w", err)
	}

	log.Info().Str("email", adminEmail).Msg("bootstrap: initial admin user created")
	return nil
}
