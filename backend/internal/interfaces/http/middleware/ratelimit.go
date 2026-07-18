package middleware

import (
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
	"github.com/trickreport/backend/internal/interfaces/http/response"
)

// rateLimitConfig holds a token bucket limiter and the last time it was used.
type rateLimitConfig struct {
	limiter  *rate.Limiter
	lastSeen atomic.Int64 // unix nanoseconds
}

// Store is the abstraction for rate limiter state. The default in-memory
// implementation works for single-instance deployments. Implement a Redis
// store for distributed rate limiting.
type Store interface {
	GetOrCreate(key string, rate rate.Limit, burst int) *rate.Limiter
	Touch(key string)
	Cleanup(maxAge time.Duration)
}

// memoryStore is the default in-memory rate limiter store.
type memoryStore struct {
	mu      sync.RWMutex
	clients map[string]*rateLimitConfig
}

func newMemoryStore() *memoryStore {
	return &memoryStore{clients: make(map[string]*rateLimitConfig)}
}

func (s *memoryStore) GetOrCreate(key string, r rate.Limit, burst int) *rate.Limiter {
	// Fast path: read lock
	s.mu.RLock()
	if cfg, ok := s.clients[key]; ok {
		cfg.lastSeen.Store(time.Now().UnixNano())
		s.mu.RUnlock()
		return cfg.limiter
	}
	s.mu.RUnlock()

	// Slow path: write lock
	s.mu.Lock()
	defer s.mu.Unlock()
	// Double-check after acquiring write lock
	if cfg, ok := s.clients[key]; ok {
		cfg.lastSeen.Store(time.Now().UnixNano())
		return cfg.limiter
	}
	cfg := &rateLimitConfig{limiter: rate.NewLimiter(r, burst)}
	cfg.lastSeen.Store(time.Now().UnixNano())
	s.clients[key] = cfg
	return cfg.limiter
}

func (s *memoryStore) Touch(key string) {
	s.mu.RLock()
	if cfg, ok := s.clients[key]; ok {
		cfg.lastSeen.Store(time.Now().UnixNano())
	}
	s.mu.RUnlock()
}

func (s *memoryStore) Cleanup(maxAge time.Duration) {
	s.mu.Lock()
	now := time.Now()
	for key, cfg := range s.clients {
		last := time.Unix(0, cfg.lastSeen.Load())
		if now.Sub(last) > maxAge {
			delete(s.clients, key)
		}
	}
	s.mu.Unlock()
}

// RateLimiter is a per-client rate limiter. It uses an in-memory store by
// default and can be configured with a Redis store for distributed deployments.
type RateLimiter struct {
	rate  rate.Limit
	burst int
	store Store
	maxAge time.Duration
}

// NewRateLimiter creates a new per-client rate limiter with an in-memory store.
func NewRateLimiter(r rate.Limit, burst int, maxAge time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:   r,
		burst:  burst,
		store:  newMemoryStore(),
		maxAge: maxAge,
	}
}

// NewRateLimiterWithStore creates a rate limiter with a custom store (e.g. Redis).
func NewRateLimiterWithStore(r rate.Limit, burst int, maxAge time.Duration, store Store) *RateLimiter {
	return &RateLimiter{
		rate:   r,
		burst:  burst,
		store:  store,
		maxAge: maxAge,
	}
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.maxAge / 2)
	defer ticker.Stop()
	for range ticker.C {
		rl.store.Cleanup(rl.maxAge)
	}
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
		if !rl.store.GetOrCreate(key, rl.rate, rl.burst).Allow() {
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
		if !rl.store.GetOrCreate(key, rl.rate, rl.burst).Allow() {
			response.Error(w, http.StatusTooManyRequests, "rate limit exceeded")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// LimitByAccount returns middleware that rate limits requests by a given
// account identifier (e.g. email). This is intended to be used with a
// pre-extracted key, such as in the LoginRateLimiter wrapper.
func (rl *RateLimiter) LimitByAccount(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The account key is expected to be stored in context by the
		// LoginRateLimiter wrapper. Fall back to IP if not present.
		key := clientIP(r)
		if accountKey, ok := r.Context().Value(accountKeyCtxKey).(string); ok && accountKey != "" {
			key = accountKey
		}
		if !rl.store.GetOrCreate(key, rl.rate, rl.burst).Allow() {
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

// accountKeyCtxKey is the context key for the per-account rate limit key.
type accountCtxKey string

const accountKeyCtxKey accountCtxKey = "account_rate_limit_key"
