package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"io"
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
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
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

// isMonitoringWebhook returns true for the Zabbix webhook path. Zabbix
// cannot hold a browser session or a CSRF cookie, so this one path is
// exempt; it authenticates instead with the per-source token header
// (see the MSP webhook handler), which is the only credential it
// accepts.
func isMonitoringWebhook(path string) bool {
	return path == "/api/v1/msp/monitoring/webhook"
}

// isWhatsAppWebhook returns true for the WhatsApp webhook path. Meta
// cannot hold a browser session or a CSRF cookie, so this one path is
// exempt; it authenticates instead with the verification token on
// GET and the HMAC signature on POST (see the CX webhook handler).
func isWhatsAppWebhook(path string) bool {
	return path == "/api/v1/cx/whatsapp/webhook"
}

// CSRF returns middleware that generates a CSRF token on safe methods (GET,
// HEAD, OPTIONS) and validates it on state-changing methods (POST, PUT, PATCH,
// DELETE). Auth endpoints and WebSocket connections are exempt. secure should
// mirror the auth cookie's COOKIE_SECURE setting (cfg.SecureCookie).
func CSRF(secure bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip CSRF for auth endpoints, WebSocket and the
			// token-authenticated monitoring webhook and WhatsApp webhook
			if isAuthPath(r.URL.Path) || isWebSocket(r) || isMonitoringWebhook(r.URL.Path) || isWhatsAppWebhook(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			switch r.Method {
			case http.MethodGet, http.MethodHead, http.MethodOptions:
				token := ""
				if cookie, err := r.Cookie(csrfCookieName); err == nil {
					if decoded, err := hex.DecodeString(cookie.Value); err == nil && len(decoded) == csrfTokenLength {
						token = cookie.Value
					}
				}
				if token == "" {
					token = generateCSRFToken()
					if token == "" {
						response.Error(w, http.StatusServiceUnavailable, "CSRF initialization unavailable")
						return
					}
					http.SetCookie(w, &http.Cookie{
						Name:     csrfCookieName,
						Value:    token,
						Path:     "/",
						HttpOnly: false, // Must be readable by JavaScript
						Secure:   secure,
						SameSite: http.SameSiteStrictMode,
					})
				}
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
}
