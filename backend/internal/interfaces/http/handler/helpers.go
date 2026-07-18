package handler

import (
	"net/http"

	"github.com/trickreport/backend/internal/application/auth"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
)

// getClaimsFromContext extracts JWT claims from the request context.
func getClaimsFromContext(r *http.Request) (*auth.Claims, bool) {
	return middleware.ClaimsFromContext(r.Context())
}
