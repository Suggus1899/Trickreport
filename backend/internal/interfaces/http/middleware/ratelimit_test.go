package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/application/auth"
	"golang.org/x/time/rate"
)

// --- memoryStore tests ---

func TestMemoryStore_GetOrCreate(t *testing.T) {
	s := newMemoryStore()
	l1 := s.GetOrCreate("key1", rate.Limit(10), 5)
	if l1 == nil {
		t.Fatal("expected non-nil limiter")
	}
	// Same key returns same limiter
	l2 := s.GetOrCreate("key1", rate.Limit(10), 5)
	if l1 != l2 {
		t.Error("expected same limiter for same key")
	}
	// Different key returns different limiter
	l3 := s.GetOrCreate("key2", rate.Limit(10), 5)
	if l1 == l3 {
		t.Error("expected different limiter for different key")
	}
}

func TestMemoryStore_Touch(t *testing.T) {
	s := newMemoryStore()
	s.GetOrCreate("key1", rate.Limit(10), 5)
	before := s.clients["key1"].lastSeen.Load()
	time.Sleep(1 * time.Millisecond)
	s.Touch("key1")
	after := s.clients["key1"].lastSeen.Load()
	if after <= before {
		t.Errorf("expected lastSeen to update: before=%d after=%d", before, after)
	}
}

func TestMemoryStore_TouchMissingKey(t *testing.T) {
	s := newMemoryStore()
	// Should not panic on missing key
	s.Touch("nonexistent")
}

func TestMemoryStore_Cleanup(t *testing.T) {
	s := newMemoryStore()
	s.GetOrCreate("old", rate.Limit(10), 5)
	// Force old lastSeen
	s.clients["old"].lastSeen.Store(time.Now().Add(-2 * time.Hour).UnixNano())

	s.GetOrCreate("new", rate.Limit(10), 5)

	s.Cleanup(1 * time.Hour)

	s.mu.RLock()
	_, oldExists := s.clients["old"]
	_, newExists := s.clients["new"]
	s.mu.RUnlock()

	if oldExists {
		t.Error("expected old key to be cleaned up")
	}
	if !newExists {
		t.Error("expected new key to remain")
	}
}

// --- clientIP tests ---

func TestClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.5:1234")
	if got := clientIP(req); got != "203.0.113.5" {
		t.Errorf("clientIP = %q, want 203.0.113.5", got)
	}
}

func TestClientIP_XForwardedForNoPort(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.5")
	if got := clientIP(req); got != "203.0.113.5" {
		t.Errorf("clientIP = %q, want 203.0.113.5", got)
	}
}

func TestClientIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:5678"
	if got := clientIP(req); got != "192.168.1.1" {
		t.Errorf("clientIP = %q, want 192.168.1.1", got)
	}
}

func TestClientIP_RemoteAddrNoPort(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1"
	if got := clientIP(req); got != "192.168.1.1" {
		t.Errorf("clientIP = %q, want 192.168.1.1", got)
	}
}

// --- LimitByIP tests ---

func TestLimitByIP_AllowsUnderLimit(t *testing.T) {
	rl := NewRateLimiter(rate.Inf, 100, time.Hour)
	called := false
	h := rl.LimitByIP(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called {
		t.Error("expected next handler to be called")
	}
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestLimitByIP_BlocksOverLimit(t *testing.T) {
	// rate=1 per second, burst=1
	rl := NewRateLimiter(rate.Limit(1), 1, time.Hour)
	called := false
	h := rl.LimitByIP(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.2:1234"

	// First request allowed (uses the burst token)
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req)
	if rec1.Code != http.StatusOK {
		t.Fatalf("first request: status = %d, want %d", rec1.Code, http.StatusOK)
	}

	// Second request immediately should be blocked
	called = false
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req)
	if called {
		t.Error("expected second request to be blocked")
	}
	if rec2.Code != http.StatusTooManyRequests {
		t.Errorf("second request: status = %d, want %d", rec2.Code, http.StatusTooManyRequests)
	}
}

func TestLimitByIP_DifferentIPsIndependent(t *testing.T) {
	rl := NewRateLimiter(rate.Limit(1), 1, time.Hour)
	called := false
	h := rl.LimitByIP(nextHandler(&called))

	// IP 1 uses its token
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.RemoteAddr = "10.0.0.10:1234"
	h.ServeHTTP(httptest.NewRecorder(), req1)

	// IP 2 should still be allowed
	called = false
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.RemoteAddr = "10.0.0.20:1234"
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req2)
	if !called {
		t.Error("expected different IP to be allowed")
	}
}

// --- LimitByUser tests ---

func TestLimitByUser_UsesUserIDWhenAuthenticated(t *testing.T) {
	rl := NewRateLimiter(rate.Limit(1), 1, time.Hour)
	called := false
	h := rl.LimitByUser(nextHandler(&called))

	claims := &auth.Claims{UserID: uuid.New(), Role: "agent"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(ContextWithClaims(req.Context(), claims))
	req.RemoteAddr = "10.0.0.30:1234"

	// First request allowed
	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req)
	if rec1.Code != http.StatusOK {
		t.Fatalf("first: status = %d, want %d", rec1.Code, http.StatusOK)
	}

	// Second request blocked (same user)
	called = false
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req)
	if called {
		t.Error("expected second request to be blocked")
	}
	if rec2.Code != http.StatusTooManyRequests {
		t.Errorf("second: status = %d, want %d", rec2.Code, http.StatusTooManyRequests)
	}
}

func TestLimitByUser_FallsBackToIP(t *testing.T) {
	rl := NewRateLimiter(rate.Limit(1), 1, time.Hour)
	called := false
	h := rl.LimitByUser(nextHandler(&called))

	// No claims -> falls back to IP
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.40:1234"

	rec1 := httptest.NewRecorder()
	h.ServeHTTP(rec1, req)
	if rec1.Code != http.StatusOK {
		t.Fatalf("first: status = %d, want %d", rec1.Code, http.StatusOK)
	}

	called = false
	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, req)
	if called {
		t.Error("expected second request to be blocked by IP")
	}
	if rec2.Code != http.StatusTooManyRequests {
		t.Errorf("second: status = %d, want %d", rec2.Code, http.StatusTooManyRequests)
	}
}

func TestLimitByUser_NilUserIDFallsBackToIP(t *testing.T) {
	rl := NewRateLimiter(rate.Inf, 100, time.Hour)
	called := false
	h := rl.LimitByUser(nextHandler(&called))

	// Claims with uuid.Nil userID -> falls back to IP
	claims := &auth.Claims{UserID: uuid.Nil, Role: "agent"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(ContextWithClaims(req.Context(), claims))
	req.RemoteAddr = "10.0.0.50:1234"

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if !called {
		t.Error("expected next handler to be called")
	}
}

// --- NewRateLimiterWithStore tests ---

type mockStore struct {
	mu       sync.Mutex
	keys     map[string]*rate.Limiter
	touched  []string
	cleanups int
}

func newMockStore() *mockStore {
	return &mockStore{keys: make(map[string]*rate.Limiter)}
}

func (m *mockStore) GetOrCreate(key string, r rate.Limit, burst int) *rate.Limiter {
	m.mu.Lock()
	defer m.mu.Unlock()
	if l, ok := m.keys[key]; ok {
		return l
	}
	l := rate.NewLimiter(r, burst)
	m.keys[key] = l
	return l
}

func (m *mockStore) Touch(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.touched = append(m.touched, key)
}

func (m *mockStore) Cleanup(maxAge time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cleanups++
}

func TestNewRateLimiterWithStore(t *testing.T) {
	store := newMockStore()
	rl := NewRateLimiterWithStore(rate.Limit(10), 5, time.Hour, store)
	if rl == nil {
		t.Fatal("expected non-nil rate limiter")
	}
	if rl.store != store {
		t.Error("expected custom store to be used")
	}
}

func TestRateLimiter_UsesCustomStore(t *testing.T) {
	store := newMockStore()
	rl := NewRateLimiterWithStore(rate.Inf, 100, time.Hour, store)
	called := false
	h := rl.LimitByIP(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "10.0.0.60:1234"
	h.ServeHTTP(httptest.NewRecorder(), req)

	if !called {
		t.Error("expected next handler to be called")
	}
	store.mu.Lock()
	if len(store.keys) != 1 {
		t.Errorf("expected 1 key in store, got %d", len(store.keys))
	}
	store.mu.Unlock()
}

// --- Start tests ---

func TestRateLimiter_Start(t *testing.T) {
	// Just verify Start doesn't block and launches goroutine.
	// Use a short maxAge so cleanup ticker fires quickly.
	rl := NewRateLimiter(rate.Limit(10), 5, 10*time.Millisecond)
	rl.Start()
	// Give the goroutine time to run a few ticks.
	time.Sleep(30 * time.Millisecond)
	// No assertion needed — just ensure no panic/block.
}

// --- Concurrency test ---

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	s := newMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			s.GetOrCreate("concurrent-key", rate.Limit(100), 10)
		}(i)
	}
	wg.Wait()
	s.mu.RLock()
	l := s.clients["concurrent-key"]
	s.mu.RUnlock()
	if l == nil {
		t.Error("expected limiter to exist after concurrent access")
	}
}

// Ensure context import is used (for ContextWithClaims in LimitByUser tests).
var _ = context.Background
