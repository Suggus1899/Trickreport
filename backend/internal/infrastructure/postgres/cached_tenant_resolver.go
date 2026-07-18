package postgres

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// cacheEntry holds a resolved tenant ID and its expiration time.
type cacheEntry struct {
	tenantID uuid.UUID
	expiresAt time.Time
}

// CachedTenantResolver wraps a TenantResolver with an in-memory cache
// (sync.Map with TTL) to reduce database lookups on every request.
//
// TODO(wire): wire CachedTenantResolver into wire.go's NewRepos / server.go
// so the middleware uses the cached resolver instead of the raw one.
type CachedTenantResolver struct {
	inner *TenantResolver
	ttl   time.Duration
	cache sync.Map
}

// NewCachedTenantResolver creates a new CachedTenantResolver with the given
// TTL. A TTL of 0 or less defaults to 5 minutes.
func NewCachedTenantResolver(inner *TenantResolver, ttl time.Duration) *CachedTenantResolver {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &CachedTenantResolver{
		inner: inner,
		ttl:   ttl,
	}
}

// Resolve returns the tenant UUID for the given slug or UUID, checking the
// cache first and falling back to the database. Results are cached for the
// configured TTL.
func (r *CachedTenantResolver) Resolve(ctx context.Context, slugOrID string) (uuid.UUID, error) {
	if v, ok := r.cache.Load(slugOrID); ok {
		entry := v.(cacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.tenantID, nil
		}
		// Expired — remove and fall through to DB.
		r.cache.Delete(slugOrID)
	}

	tenantID, err := r.inner.Resolve(ctx, slugOrID)
	if err != nil {
		return uuid.Nil, err
	}

	r.cache.Store(slugOrID, cacheEntry{
		tenantID:  tenantID,
		expiresAt: time.Now().Add(r.ttl),
	})
	return tenantID, nil
}

// Compile-time assertion that CachedTenantResolver implements the middleware resolver.
var _ interface {
	Resolve(context.Context, string) (uuid.UUID, error)
} = (*CachedTenantResolver)(nil)
