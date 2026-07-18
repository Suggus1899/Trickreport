package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/trickreport/backend/internal/application/auth"
	"github.com/trickreport/backend/internal/interfaces/http/response"
)

type contextKey string

const claimsKey contextKey = "claims"
const tenantKey contextKey = "tenant_id"

// TenantResolver translates a tenant slug or UUID into the canonical tenant UUID.
type TenantResolver interface {
	Resolve(ctx context.Context, slugOrID string) (string, error)
}

// ContextWithClaims stores JWT claims in the request context.
func ContextWithClaims(ctx context.Context, claims *auth.Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

// ClaimsFromContext retrieves JWT claims from the request context.
func ClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	c, ok := ctx.Value(claimsKey).(*auth.Claims)
	return c, ok
}

// TenantFromContext retrieves the tenant identifier from the request context.
// Returns empty string if no tenant is set — callers must check.
func TenantFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(tenantKey).(string); ok {
		return v
	}
	return ""
}

// Authenticate validates the JWT from the Authorization header or cookie.
func Authenticate(tokenGen auth.TokenGenerator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				response.Error(w, http.StatusUnauthorized, "authentication required")
				return
			}

			claims, err := tokenGen.Validate(token)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			ctx := ContextWithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole restricts access to specific roles.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok {
				response.Error(w, http.StatusUnauthorized, "authentication required")
				return
			}

			if !allowed[claims.Role] {
				response.Error(w, http.StatusForbidden, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Tenant resolves the current tenant from the X-Tenant-ID header or subdomain
// and validates/translates it to the canonical tenant UUID.
// Returns 400 if no tenant can be resolved — no silent fallback.
func Tenant(resolver TenantResolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slugOrID := r.Header.Get("X-Tenant-ID")

			if slugOrID == "" {
				host := r.Host
				parts := strings.SplitN(host, ".", 2)
				if len(parts) > 1 && parts[0] != "www" && parts[0] != "app" {
					slugOrID = parts[0]
				}
			}

			if slugOrID == "" {
				response.Error(w, http.StatusBadRequest, "missing tenant: provide X-Tenant-ID header or use a subdomain")
				return
			}

			tenantID, err := resolver.Resolve(r.Context(), slugOrID)
			if err != nil {
				response.Error(w, http.StatusBadRequest, "invalid tenant")
				return
			}

			ctx := context.WithValue(r.Context(), tenantKey, tenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	if c, err := r.Cookie("trickreport_token"); err == nil {
		return c.Value
	}
	return ""
}
