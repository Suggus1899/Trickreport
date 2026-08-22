package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/trickreport/backend/internal/interfaces/http/response"
)

const (
	// csrfCookieName is the name of the CSRF cookie.
	csrfCookieName = "trickreport_csrf"
	// csrfHeaderName is the name of the header that must match the cookie.
	csrfHeaderName = "X-CSRF-Token"
	// csrfTokenLength is the number of random bytes used for the token.
	csrfTokenLength = 32
)

// generateCSRFToken generates a cryptographically random hex token.
func generateCSRFToken() string {
	b := make([]byte, csrfTokenLength)
	if _, err := rand.Read(b); err != nil {
		// Fallback — should not happen in practice
		return "csrf-fallback-token"
	}
	return hex.EncodeToString(b)
}

// isAuthPath returns true for paths that should be exempt from CSRF protection
// (login, refresh, password reset — these don't have a session yet).
func isAuthPath(path string) bool {
	return strings.HasPrefix(path, "/api/v1/auth/") &&
		(strings.HasSuffix(path, "/login") ||
			strings.HasSuffix(path, "/refresh") ||
			strings.HasSuffix(path, "/password-reset") ||
			strings.HasSuffix(path, "/password-reset/confirm") ||
			strings.HasSuffix(path, "/mfa/login"))
}

// isWebSocket returns true for WebSocket upgrade requests.
func isWebSocket(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket")
}

// CSRF returns middleware that generates a CSRF token on safe methods (GET,
// HEAD, OPTIONS) and validates it on state-changing methods (POST, PUT, PATCH,
// DELETE). Auth endpoints and WebSocket connections are exempt.
func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip CSRF for auth endpoints and WebSocket
		if isAuthPath(r.URL.Path) || isWebSocket(r) {
			next.ServeHTTP(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			// Generate or refresh the CSRF token cookie on safe methods
			token := generateCSRFToken()
			http.SetCookie(w, &http.Cookie{
				Name:     csrfCookieName,
				Value:    token,
				Path:     "/",
				HttpOnly: false, // Must be readable by JavaScript
				SameSite: http.SameSiteStrictMode,
			})
			// Also expose in header for convenience
			w.Header().Set("X-CSRF-Token", token)
			next.ServeHTTP(w, r)

		default:
			// Validate the CSRF token on state-changing methods
			cookie, err := r.Cookie(csrfCookieName)
			if err != nil || cookie.Value == "" {
				response.Error(w, http.StatusForbidden, "missing csrf token")
				return
			}
			headerToken := r.Header.Get(csrfHeaderName)
			if headerToken == "" || headerToken != cookie.Value {
				response.Error(w, http.StatusForbidden, "invalid csrf token")
				return
			}
			next.ServeHTTP(w, r)
		}
	})
}
