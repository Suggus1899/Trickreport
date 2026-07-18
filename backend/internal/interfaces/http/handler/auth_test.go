package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	"github.com/trickreport/backend/internal/domain/user"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
)

// --- Mocks for auth service dependencies ---

type authMockUserRepo struct {
	u      *user.User
	err    error
	gotID  uuid.UUID
}

func (m *authMockUserRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	return m.u, m.err
}

func (m *authMockUserRepo) GetByIDNoTenant(ctx context.Context, id uuid.UUID) (*user.User, error) {
	m.gotID = id
	return m.u, m.err
}

type authMockHasher struct{ err error }

func (m *authMockHasher) Compare(hash, password string) error { return m.err }

type authMockTokenGen struct {
	token string
	err   error
}

func (m *authMockTokenGen) Generate(userID, tenantID uuid.UUID, role string) (string, error) {
	return m.token, m.err
}
func (m *authMockTokenGen) Validate(token string) (*appAuth.Claims, error) {
	return nil, errors.New("invalid")
}

// --- Helpers ---

func newAuthHandler(repo appAuth.UserRepository, hasher appAuth.PasswordHasher, ldap appAuth.LDAPAuthenticator, tok appAuth.TokenGenerator) *AuthHandler {
	svc := appAuth.NewService(repo, hasher, ldap, tok, appAuth.Config{JWTExpHours: 24, SecureCookie: true})
	return NewAuthHandler(svc)
}

func newTestUser() *user.User {
	return &user.User{
		ID:       uuid.New(),
		TenantID: uuid.New(),
		Name:     "Jane Doe",
		Email:    "jane@example.com",
		Role:     user.RoleAgent,
	}
}

func decodeBody(t *testing.T, body string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.NewDecoder(strings.NewReader(body)).Decode(&m); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	return m
}

// --- Login tests ---

func TestAuthHandler_Login_Success(t *testing.T) {
	u := newTestUser()
	u.PasswordHash = "hashed"
	repo := &authMockUserRepo{u: u}
	hasher := &authMockHasher{}
	tok := &authMockTokenGen{token: "jwt-abc"}
	h := newAuthHandler(repo, hasher, nil, tok)

	body := `{"email":"jane@example.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	resp := decodeBody(t, rec.Body.String())
	if resp["token"] != "jwt-abc" {
		t.Errorf("token = %v", resp["token"])
	}
	userMap, ok := resp["user"].(map[string]any)
	if !ok {
		t.Fatal("expected user object")
	}
	if userMap["email"] != "jane@example.com" {
		t.Errorf("user email = %v", userMap["email"])
	}

	// Cookie should be set
	found := false
	for _, c := range rec.Result().Cookies() {
		if c.Name == "trickreport_token" && c.Value == "jwt-abc" {
			found = true
		}
	}
	if !found {
		t.Error("expected trickreport_token cookie to be set")
	}
}

func TestAuthHandler_Login_InvalidJSON(t *testing.T) {
	h := newAuthHandler(&authMockUserRepo{}, &authMockHasher{}, nil, &authMockTokenGen{})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader("{bad json"))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAuthHandler_Login_ValidationError(t *testing.T) {
	h := newAuthHandler(&authMockUserRepo{}, &authMockHasher{}, nil, &authMockTokenGen{})

	// password too short (min=6)
	body := `{"email":"jane@example.com","password":"x"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	repo := &authMockUserRepo{err: errors.New("not found")}
	h := newAuthHandler(repo, &authMockHasher{}, nil, &authMockTokenGen{})

	body := `{"email":"nope@example.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Login_PasswordNotConfigured(t *testing.T) {
	u := newTestUser()
	u.PasswordHash = ""
	repo := &authMockUserRepo{u: u}
	h := newAuthHandler(repo, &authMockHasher{}, nil, &authMockTokenGen{})

	body := `{"email":"jane@example.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Login_TokenGenError(t *testing.T) {
	u := newTestUser()
	u.PasswordHash = "hashed"
	repo := &authMockUserRepo{u: u}
	tok := &authMockTokenGen{err: errors.New("sign fail")}
	h := newAuthHandler(repo, &authMockHasher{}, nil, tok)

	body := `{"email":"jane@example.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.Login(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- Me tests ---

func TestAuthHandler_Me_Success(t *testing.T) {
	u := newTestUser()
	repo := &authMockUserRepo{u: u}
	h := newAuthHandler(repo, &authMockHasher{}, nil, &authMockTokenGen{})

	claims := &appAuth.Claims{UserID: u.ID, TenantID: u.TenantID, Role: "agent"}
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), claims))

	rec := httptest.NewRecorder()
	h.Me(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	resp := decodeBody(t, rec.Body.String())
	if resp["email"] != "jane@example.com" {
		t.Errorf("email = %v", resp["email"])
	}
	if repo.gotID != u.ID {
		t.Errorf("repo got id %v, want %v", repo.gotID, u.ID)
	}
}

func TestAuthHandler_Me_Unauthorized(t *testing.T) {
	h := newAuthHandler(&authMockUserRepo{}, &authMockHasher{}, nil, &authMockTokenGen{})

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()

	h.Me(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Me_NotFound(t *testing.T) {
	repo := &authMockUserRepo{err: errors.New("not found")}
	h := newAuthHandler(repo, &authMockHasher{}, nil, &authMockTokenGen{})

	claims := &appAuth.Claims{UserID: uuid.New(), Role: "agent"}
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), claims))
	rec := httptest.NewRecorder()

	h.Me(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// --- Logout tests ---

func TestAuthHandler_Logout_ClearsCookie(t *testing.T) {
	h := newAuthHandler(&authMockUserRepo{}, &authMockHasher{}, nil, &authMockTokenGen{})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	rec := httptest.NewRecorder()

	h.Logout(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	found := false
	for _, c := range rec.Result().Cookies() {
		if c.Name == "trickreport_token" && c.MaxAge < 0 {
			found = true
		}
	}
	if !found {
		t.Error("expected trickreport_token cookie to be cleared")
	}
}
