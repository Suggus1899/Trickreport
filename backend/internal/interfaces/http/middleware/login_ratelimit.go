package middleware

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"

	"github.com/trickreport/backend/internal/interfaces/http/response"
)

// TODO: Wire in server.go after the auth middleware — wrap the login handler:
//   r.With(loginLimiter.LoginLimit).Post("/login", authHandler.Login)

// loginEmailReq is the minimal request body shape needed to extract the email.
type loginEmailReq struct {
	Email string `json:"email"`
}

// LoginRateLimiter wraps the login handler with per-account (email+IP) rate
// limiting. It reads the email from the request body, builds a composite key
// of email+IP, and checks the rate limiter before delegating to the wrapped
// handler.
type LoginRateLimiter struct {
	limiter *RateLimiter
}

// NewLoginRateLimiter creates a new LoginRateLimiter with the given rate and
// burst. A typical configuration is rate=5 attempts/minute, burst=10.
func NewLoginRateLimiter(r rate.Limit, burst int, maxAge int) *LoginRateLimiter {
	// maxAge is in minutes for the cleanup interval
	return &LoginRateLimiter{
		limiter: NewRateLimiter(r, burst, durationMinutes(maxAge)),
	}
}

// durationMinutes converts minutes to a time.Duration.
func durationMinutes(minutes int) time.Duration {
	return time.Duration(minutes) * time.Minute
}

// LoginLimit wraps the given handler with per-account rate limiting. It
// peeks at the request body to extract the email, then restores the body
// for the downstream handler.
func (l *LoginRateLimiter) LoginLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the body to extract the email
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request body")
			return
		}
		// Restore the body for the downstream handler
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		r.ContentLength = int64(len(bodyBytes))

		// Extract email from the body
		email := extractEmailFromBody(bodyBytes)
		ip := clientIP(r)

		// Build a composite key: email+IP
		var key string
		if email != "" {
			key = strings.ToLower(email) + "|" + ip
		} else {
			key = ip
		}

		if !l.limiter.store.GetOrCreate(key, l.limiter.rate, l.limiter.burst).Allow() {
			response.Error(w, http.StatusTooManyRequests, "too many login attempts")
			return
		}

		// Store the key in context for any downstream LimitByAccount usage
		ctx := context.WithValue(r.Context(), accountKeyCtxKey, key)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Start initializes the background cleanup goroutine for the underlying limiter.
func (l *LoginRateLimiter) Start() {
	l.limiter.Start()
}

// extractEmailFromBody attempts to parse the email field from a JSON body.
// Returns an empty string if parsing fails.
func extractEmailFromBody(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var req loginEmailReq
	if err := json.Unmarshal(body, &req); err != nil {
		return ""
	}
	return req.Email
}
