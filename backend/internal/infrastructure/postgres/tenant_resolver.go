package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
)

// TenantResolver translates a tenant slug or UUID into the canonical tenant UUID.
type TenantResolver struct {
	db *pgxpool.Pool
}

// NewTenantResolver creates a new TenantResolver.
func NewTenantResolver(db *pgxpool.Pool) *TenantResolver {
	return &TenantResolver{db: db}
}

// Resolve returns the tenant UUID for the given slug or UUID.
// If the input is already a UUID present in the tenants table, it is returned as-is.
func (r *TenantResolver) Resolve(ctx context.Context, slugOrID string) (string, error) {
	var id string
	err := r.db.QueryRow(ctx, `SELECT id FROM tenants WHERE slug = $1 OR id::text = $1 LIMIT 1`, slugOrID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("tenant not found")
		}
		return "", fmt.Errorf("failed to resolve tenant: %w", err)
	}
	return id, nil
}

// Compile-time assertion that TenantResolver implements the middleware resolver.
var _ middleware.TenantResolver = (*TenantResolver)(nil)
