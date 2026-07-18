package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/application/auth"
)

// --- Mocks ---

type mockTenantResolver struct {
	id  uuid.UUID
	err error
	got string
}

func (m *mockTenantResolver) Resolve(ctx context.Context, slugOrID string) (uuid.UUID, error) {
	m.got = slugOrID
	return m.id, m.err
}

type mockTokenGen struct {
	claims *auth.Claims
	err    error
	got    string
}

func (m *mockTokenGen) Generate(userID, tenantID uuid.UUID, role string) (string, error) {
	return "token", nil
}

func (m *mockTokenGen) GenerateRefresh(userID, tenantID uuid.UUID, role string) (string, error) {
	return "refresh-token", nil
}

func (m *mockTokenGen) Validate(token string) (*auth.Claims, error) {
	m.got = token
	return m.claims, m.err
}

// nextHandler returns a handler that sets a flag and writes 200.
func nextHandler(called *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
		w.WriteHeader(http.StatusOK)
	})
}

// --- Context helper tests ---

func TestContextWithClaims_And_ClaimsFromContext(t *testing.T) {
	claims := &auth.Claims{UserID: uuid.New(), TenantID: uuid.New(), Role: "agent"}
	ctx := ContextWithClaims(context.Background(), claims)

	got, ok := ClaimsFromContext(ctx)
	if !ok {
		t.Fatal("expected claims in context")
	}
	if got.UserID != claims.UserID {
		t.Errorf("userID = %v", got.UserID)
	}
}

func TestClaimsFromContext_Missing(t *testing.T) {
	_, ok := ClaimsFromContext(context.Background())
	if ok {
		t.Error("expected ok=false when no claims")
	}
}

func TestTenantFromContext_Missing(t *testing.T) {
	if got := TenantFromContext(context.Background()); got != uuid.Nil {
		t.Errorf("expected uuid.Nil, got %v", got)
	}
}

func TestTenantFromContext_Set(t *testing.T) {
	id := uuid.New()
	ctx := context.WithValue(context.Background(), tenantKey, id)
	if got := TenantFromContext(ctx); got != id {
		t.Errorf("got %v, want %v", got, id)
	}
}

// --- Tenant middleware tests ---

func TestTenant_ValidHeader(t *testing.T) {
	tenantID := uuid.New()
	resolver := &mockTenantResolver{id: tenantID}
	called := false
	h := Tenant(resolver)(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", "my-tenant")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called {
		t.Error("expected next handler to be called")
	}
	if resolver.got != "my-tenant" {
		t.Errorf("resolver got %q", resolver.got)
	}
	if got := TenantFromContext(req.Context()); got != uuid.Nil {
		// req.Context() is the original; the modified one is inside rec
	}
	// Verify tenant is in the context that reaches next handler
}

func TestTenant_ValidHeader_TenantInContext(t *testing.T) {
	tenantID := uuid.New()
	resolver := &mockTenantResolver{id: tenantID}
	var ctxTenant uuid.UUID
	h := Tenant(resolver)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxTenant = TenantFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", "my-tenant")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if ctxTenant != tenantID {
		t.Errorf("context tenant = %v, want %v", ctxTenant, tenantID)
	}
}

func TestTenant_MissingHeader(t *testing.T) {
	resolver := &mockTenantResolver{}
	called := false
	h := Tenant(resolver)(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "localhost"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTenant_SubdomainFallback(t *testing.T) {
	tenantID := uuid.New()
	resolver := &mockTenantResolver{id: tenantID}
	called := false
	h := Tenant(resolver)(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "acme.example.com"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called {
		t.Error("expected next handler to be called via subdomain")
	}
	if resolver.got != "acme" {
		t.Errorf("resolver got %q, want acme", resolver.got)
	}
}

func TestTenant_SubdomainWwwIgnored(t *testing.T) {
	resolver := &mockTenantResolver{}
	called := false
	h := Tenant(resolver)(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "www.example.com"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if called {
		t.Error("expected next handler NOT to be called for www subdomain")
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTenant_SubdomainAppIgnored(t *testing.T) {
	resolver := &mockTenantResolver{}
	called := false
	h := Tenant(resolver)(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Host = "app.example.com"
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if called {
		t.Error("expected next handler NOT to be called for app subdomain")
	}
}

func TestTenant_InvalidUUID(t *testing.T) {
	resolver := &mockTenantResolver{err: errors.New("not found")}
	called := false
	h := Tenant(resolver)(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Tenant-ID", "bogus")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// --- Authenticate middleware tests ---

func TestAuthenticate_ValidToken(t *testing.T) {
	claims := &auth.Claims{UserID: uuid.New(), TenantID: uuid.New(), Role: "agent"}
	tok := &mockTokenGen{claims: claims}
	called := false
	h := Authenticate(tok)(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called {
		t.Error("expected next handler to be called")
	}
	if tok.got != "valid-token" {
		t.Errorf("token gen got %q", tok.got)
	}
}

func TestAuthenticate_MissingToken(t *testing.T) {
	tok := &mockTokenGen{}
	called := false
	h := Authenticate(tok)(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	tok := &mockTokenGen{err: errors.New("bad token")}
	called := false
	h := Authenticate(tok)(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthenticate_ExpiredToken(t *testing.T) {
	tok := &mockTokenGen{err: errors.New("token is expired")}
	called := false
	h := Authenticate(tok)(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer expired")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthenticate_CookieToken(t *testing.T) {
	claims := &auth.Claims{UserID: uuid.New(), Role: "agent"}
	tok := &mockTokenGen{claims: claims}
	called := false
	h := Authenticate(tok)(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "trickreport_token", Value: "cookie-token"})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called {
		t.Error("expected next handler to be called via cookie")
	}
	if tok.got != "cookie-token" {
		t.Errorf("token gen got %q", tok.got)
	}
}

// --- RequireRole middleware tests ---

func TestRequireRole_Allowed(t *testing.T) {
	called := false
	h := RequireRole("admin", "agent")(nextHandler(&called))

	claims := &auth.Claims{Role: "agent"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(ContextWithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if !called {
		t.Error("expected next handler to be called")
	}
}

func TestRequireRole_Forbidden(t *testing.T) {
	called := false
	h := RequireRole("admin")(nextHandler(&called))

	claims := &auth.Claims{Role: "end_user"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(ContextWithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}
	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestRequireRole_NoClaims(t *testing.T) {
	called := false
	h := RequireRole("admin")(nextHandler(&called))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}
	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// --- extractToken tests ---

func TestExtractToken_BearerHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer abc123")
	if got := extractToken(req); got != "abc123" {
		t.Errorf("extractToken = %q, want abc123", got)
	}
}

func TestExtractToken_Cookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "trickreport_token", Value: "cookie-val"})
	if got := extractToken(req); got != "cookie-val" {
		t.Errorf("extractToken = %q, want cookie-val", got)
	}
}

func TestExtractToken_None(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if got := extractToken(req); got != "" {
		t.Errorf("extractToken = %q, want empty", got)
	}
}

func TestExtractToken_NonBearerHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic abc")
	if got := extractToken(req); got != "" {
		t.Errorf("extractToken = %q, want empty for non-Bearer", got)
	}
}
