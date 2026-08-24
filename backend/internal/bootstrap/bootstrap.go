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

	// systemUserID is the fixed identity attributed to automation- and
	// worker-generated ticket history/comments (SLA breach, escalation,
	// automation rule actions). Both infrastructure/postgres.AutomationExecutor
	// and worker.Worker default to this exact UUID; EnsureSystemUser makes
	// sure a real row backs it so their INSERTs don't fail the
	// ticket_history.user_id / ticket_comments.user_id foreign keys.
	systemUserID = uuid.MustParse("00000000-0000-0000-0000-000000000002")
)

// EnsureSystemUser creates the fixed system user row that automation- and
// worker-generated history/comments are attributed to, if it doesn't already
// exist. Safe to call on every startup.
func EnsureSystemUser(ctx context.Context, pool *pgxpool.Pool) error {
	var exists bool
	if err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`, systemUserID).Scan(&exists); err != nil {
		return fmt.Errorf("bootstrap: failed to check system user: %w", err)
	}
	if exists {
		return nil
	}

	_, err := pool.Exec(ctx, `
		INSERT INTO users (id, tenant_id, name, email, role, password, active)
		VALUES ($1, $2, 'Trickreport Automation', 'automation@trickreport.internal', 'agent', NULL, FALSE)
		ON CONFLICT (id) DO NOTHING`,
		systemUserID, defaultTenantID,
	)
	if err != nil {
		return fmt.Errorf("bootstrap: failed to create system user: %w", err)
	}

	log.Info().Str("id", systemUserID.String()).Msg("bootstrap: system user created")
	return nil
}

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
