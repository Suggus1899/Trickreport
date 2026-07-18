package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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
func (r *TenantResolver) Resolve(ctx context.Context, slugOrID string) (uuid.UUID, error) {
	var id pgtype.UUID
	err := r.db.QueryRow(ctx, `SELECT id FROM tenants WHERE slug = $1 OR id::text = $1 LIMIT 1`, slugOrID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("tenant not found")
		}
		return uuid.Nil, fmt.Errorf("failed to resolve tenant: %w", err)
	}
	return pgToUUID(id), nil
}

// Compile-time assertion that TenantResolver implements the middleware resolver.
var _ middleware.TenantResolver = (*TenantResolver)(nil)
