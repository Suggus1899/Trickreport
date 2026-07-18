package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
	"github.com/trickreport/backend/internal/interfaces/http/response"
)

// rateLimitConfig holds a token bucket limiter and the last time it was used.
type rateLimitConfig struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter is an in-memory rate limiter keyed by client identifier.
type RateLimiter struct {
	rate     rate.Limit
	burst    int
	mu       sync.Mutex
	clients  map[string]*rateLimitConfig
	maxAge   time.Duration
}

// NewRateLimiter creates a new per-client rate limiter.
// rate is the number of requests per second sustained, burst is the maximum
// burst size, and maxAge is how long to keep idle client entries.
func NewRateLimiter(r rate.Limit, burst int, maxAge time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:    r,
		burst:   burst,
		clients: make(map[string]*rateLimitConfig),
		maxAge:  maxAge,
	}
}

// cleanup periodically removes stale limiters. It should be started once.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.maxAge / 2)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, cfg := range rl.clients {
			if now.Sub(cfg.lastSeen) > rl.maxAge {
				delete(rl.clients, key)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cfg, ok := rl.clients[key]
	if !ok {
		cfg = &rateLimitConfig{limiter: rate.NewLimiter(rl.rate, rl.burst)}
		rl.clients[key] = cfg
	}
	cfg.lastSeen = time.Now()
	return cfg.limiter
}

// clientIP returns the client IP from the request, preferring X-Forwarded-For.
func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		if host, _, err := net.SplitHostPort(ip); err == nil {
			return host
		}
		return ip
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	if host != "" {
		return host
	}
	return r.RemoteAddr
}

// LimitByIP returns middleware that rate limits requests per IP address.
func (rl *RateLimiter) LimitByIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := clientIP(r)
		if !rl.getLimiter(key).Allow() {
			response.Error(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LimitByUser returns middleware that rate limits requests per authenticated user.
// It falls back to IP when no user is in context.
func (rl *RateLimiter) LimitByUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := clientIP(r)
		if claims, ok := ClaimsFromContext(r.Context()); ok && claims.UserID != uuid.Nil {
			key = claims.UserID.String()
		}
		if !rl.getLimiter(key).Allow() {
			response.Error(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Start initializes the background cleanup goroutine.
func (rl *RateLimiter) Start() {
	go rl.cleanup()
}
