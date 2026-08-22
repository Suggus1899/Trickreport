package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	"github.com/trickreport/backend/internal/config"
	"github.com/trickreport/backend/internal/email"
	infraAuth "github.com/trickreport/backend/internal/infrastructure/auth"
	infraRealtime "github.com/trickreport/backend/internal/infrastructure/realtime"
	"github.com/trickreport/backend/internal/realtime"
)

func TestNewRepos(t *testing.T) {
	repos := NewRepos(nil)
	if repos == nil {
		t.Fatal("expected non-nil repos")
	}
	if repos.Ticket == nil || repos.User == nil || repos.Article == nil {
		t.Error("expected non-nil repos in struct")
	}
}

func TestNewServices(t *testing.T) {
	repos := NewRepos(nil)
	hasher := infraAuth.NewBcryptHasher()
	tokenGen := infraAuth.NewJWTGenerator("test-secret", 8)

	hub := realtime.NewHub()
	hubAdapter := infraRealtime.NewHubAdapter(hub)

	sender := email.NewConsoleSender()
	authCfg := appAuth.Config{
		JWTSecret:    "test-secret",
		JWTExpHours:  8,
		SecureCookie: false,
	}

	services := NewServices(repos, hasher, tokenGen, hubAdapter, sender, nil, authCfg)
	if services == nil {
		t.Fatal("expected non-nil services")
	}
	if services.Auth == nil || services.Ticket == nil || services.User == nil {
		t.Error("expected non-nil services in struct")
	}
	if services.Engine == nil {
		t.Error("expected non-nil engine")
	}
	if services.Attachment == nil {
		t.Error("expected non-nil attachment service")
	}
}

func TestNewHandlers(t *testing.T) {
	repos := NewRepos(nil)
	hasher := infraAuth.NewBcryptHasher()
	tokenGen := infraAuth.NewJWTGenerator("test-secret", 8)

	hub := realtime.NewHub()
	hubAdapter := infraRealtime.NewHubAdapter(hub)

	sender := email.NewConsoleSender()
	authCfg := appAuth.Config{
		JWTSecret:    "test-secret",
		JWTExpHours:  8,
		SecureCookie: false,
	}

	services := NewServices(repos, hasher, tokenGen, hubAdapter, sender, nil, authCfg)
	handlers := NewHandlers(services)
	if handlers == nil {
		t.Fatal("expected non-nil handlers")
	}
	if handlers.Auth == nil || handlers.Ticket == nil || handlers.User == nil {
		t.Error("expected non-nil handlers in struct")
	}
	if handlers.Article == nil || handlers.SLA == nil || handlers.Automation == nil {
		t.Error("expected non-nil handlers in struct")
	}
	if handlers.Analytics == nil || handlers.Attachment == nil {
		t.Error("expected non-nil handlers in struct")
	}
}

func TestNew_Server(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)
	if srv == nil {
		t.Fatal("expected non-nil server")
	}
	if srv.router == nil {
		t.Error("expected non-nil router")
	}
	if srv.cfg != cfg {
		t.Error("expected config to be set")
	}

	// Cancel context to stop the worker goroutine
	cancel()
	time.Sleep(50 * time.Millisecond)
}

func TestServer_Router_SecurityHeaders(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	// Make a request to any route — middleware will add security headers.
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	// Security headers should be set regardless of the route existing or not.
	if rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("expected X-Content-Type-Options header")
	}
	if rec.Header().Get("X-Frame-Options") != "DENY" {
		t.Error("expected X-Frame-Options header")
	}
}

func TestServer_Router_CORS_AllowedOrigin(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	req := httptest.NewRequest(http.MethodOptions, "/api/tickets", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Errorf("Access-Control-Allow-Origin = %q", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if rec.Code != http.StatusNoContent {
		t.Errorf("OPTIONS status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestServer_Router_CORS_RejectedOrigin(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	req := httptest.NewRequest(http.MethodOptions, "/api/tickets", nil)
	req.Header.Set("Origin", "https://evil.com")
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Error("expected no Access-Control-Allow-Origin for rejected origin")
	}
	if rec.Code != http.StatusForbidden {
		t.Errorf("OPTIONS status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestServer_Router_ProductionSecurityHeaders(t *testing.T) {
	cfg := &config.Config{
		Port:          "0",
		Env:           "production",
		JWTSecret:     "prod-secret-not-change",
		JWTExpHours:   8,
		CORSOrigins:   "http://localhost:3000",
		AdminPassword: "prod-admin-pass",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	if rec.Header().Get("Strict-Transport-Security") == "" {
		t.Error("expected Strict-Transport-Security header in production")
	}
	if rec.Header().Get("Content-Security-Policy") == "" {
		t.Error("expected Content-Security-Policy header in production")
	}
}

func TestServer_Router_RootRedirect(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusSeeOther)
	}
	if rec.Header().Get("Location") != "http://localhost:3000" {
		t.Errorf("Location = %q, want http://localhost:3000", rec.Header().Get("Location"))
	}
}

func TestServer_Router_RootJSON_NoOrigins(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestServer_Router_HealthCheck_NilPool(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	// Pool is nil, so Ping panics; chi's Recoverer catches it and returns 500.
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestServer_Router_ReadyCheck_NilPool(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/ready", nil)
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	// Pool is nil, so Ping panics; chi's Recoverer catches it and returns 500.
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestServer_Router_CSRF_MissingToken_Rejected(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	// A mutating request on a protected (non-auth) API route without a CSRF
	// cookie/header must be rejected before it ever reaches the auth middleware.
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tickets", nil)
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestServer_Router_CSRF_SafeMethod_IssuesCookie(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tickets", nil)
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	found := false
	for _, c := range rec.Result().Cookies() {
		if c.Name == "trickreport_csrf" && c.Value != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected trickreport_csrf cookie to be set on a safe-method API request")
	}
}

func TestServer_Router_CSRF_AuthLoginExempt(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	// POST /auth/login has no CSRF cookie/header attached; it must not be
	// rejected with 403 by the CSRF middleware (isAuthPath exemption).
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	if rec.Code == http.StatusForbidden {
		t.Errorf("login should be exempt from CSRF, got status %d", rec.Code)
	}
}

func TestServer_Start_GracefulShutdown(t *testing.T) {
	cfg := &config.Config{
		Port:        "0", // random free port
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	done := make(chan struct{})
	go func() {
		srv.Start()
		close(done)
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Cancel context to trigger graceful shutdown
	cancel()

	select {
	case <-done:
		// success — Start returned after graceful shutdown
	case <-time.After(5 * time.Second):
		t.Fatal("Start did not return after context cancellation")
	}
}

// TestServer_Router_NotificationsAreAuthenticated is a regression test: the
// /notifications routes were registered as a sibling of the authenticated
// group instead of inside it, so Authenticate/Tenant never ran and every
// request was rejected with 401 no matter how valid the caller's token was —
// the notification feature could not work at all.
//
// A request carrying a valid token must therefore get past the auth
// middleware. Repos are nil here, so it then fails inside the handler and
// Recoverer turns that into a 500; anything other than 401 proves the
// middleware chain ran. Rejecting a tokenless request is checked too, so a
// route that simply skips auth entirely cannot pass this test either.
func TestServer_Router_NotificationsAreAuthenticated(t *testing.T) {
	cfg := &config.Config{
		Port:        "0",
		Env:         "development",
		JWTSecret:   "test-secret",
		JWTExpHours: 8,
		CORSOrigins: "http://localhost:3000",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	srv := New(ctx, cfg, nil)

	token, err := infraAuth.NewJWTGenerator(cfg.JWTSecret, cfg.JWTExpHours).
		Generate(uuid.New(), uuid.New(), "agent")
	if err != nil {
		t.Fatalf("failed to mint test token: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)

	if rec.Code == http.StatusUnauthorized {
		t.Error("a valid token was rejected with 401 — /notifications is outside the authenticated group")
	}

	anon := httptest.NewRequest(http.MethodGet, "/api/v1/notifications/", nil)
	anonRec := httptest.NewRecorder()
	srv.router.ServeHTTP(anonRec, anon)

	if anonRec.Code != http.StatusUnauthorized {
		t.Errorf("tokenless request: status = %d, want %d", anonRec.Code, http.StatusUnauthorized)
	}
}
