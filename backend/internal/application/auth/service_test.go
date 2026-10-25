package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/domain/user"
)

// --- Mocks ---

type mockUserRepo struct {
	userByEmail   *user.User
	errByEmail    error
	userByID      *user.User
	errByID       error
	gotEmail      string
	gotID         uuid.UUID
	gotTenantID   uuid.UUID
	usersByTenant map[uuid.UUID]*user.User
	byEmailCalls  int
	byIDCalls     int
}

func (m *mockUserRepo) GetByEmail(_ context.Context, tenantID uuid.UUID, email string) (*user.User, error) {
	m.gotTenantID = tenantID
	m.gotEmail = email
	m.byEmailCalls++
	// usersByTenant lets a test model the same email living in several
	// tenants; when it is nil the mock keeps its single-user behaviour.
	if m.usersByTenant != nil {
		u, ok := m.usersByTenant[tenantID]
		if !ok {
			return nil, m.errByEmail
		}
		return u, nil
	}
	return m.userByEmail, m.errByEmail
}

func (m *mockUserRepo) GetByIDNoTenant(_ context.Context, id uuid.UUID) (*user.User, error) {
	m.gotID = id
	m.byIDCalls++
	return m.userByID, m.errByID
}

type mockHasher struct {
	compareErr error
	gotHash    string
	gotPwd     string
}

func (m *mockHasher) Compare(hash, password string) error {
	m.gotHash = hash
	m.gotPwd = password
	return m.compareErr
}

type mockLDAP struct {
	err     error
	gotUser string
	gotPwd  string
}

func (m *mockLDAP) Authenticate(username, password string) (string, error) {
	m.gotUser = username
	m.gotPwd = password
	return username, m.err
}

type mockTokenGen struct {
	token            string
	err              error
	gotUserID        uuid.UUID
	gotTenant        uuid.UUID
	gotRole          string
	validateCl       *Claims
	validateEr       error
	refreshToken     string
	refreshErr       error
	refreshGotUserID uuid.UUID
	refreshGotTenant uuid.UUID
	refreshGotRole   string
}

func (m *mockTokenGen) Generate(userID, tenantID uuid.UUID, role string) (string, error) {
	m.gotUserID = userID
	m.gotTenant = tenantID
	m.gotRole = role
	return m.token, m.err
}

func (m *mockTokenGen) GenerateRefresh(userID, tenantID uuid.UUID, role string) (string, error) {
	m.refreshGotUserID = userID
	m.refreshGotTenant = tenantID
	m.refreshGotRole = role
	if m.refreshToken != "" {
		return m.refreshToken, m.refreshErr
	}
	return "refresh-" + m.token, m.refreshErr
}

func (m *mockTokenGen) Validate(_ string) (*Claims, error) {
	return m.validateCl, m.validateEr
}

// --- Helpers ---

func newTestUser() *user.User {
	return &user.User{
		ID:           uuid.New(),
		TenantID:     uuid.New(),
		Name:         "Jane Doe",
		Email:        "jane@example.com",
		Role:         user.RoleAgent,
		PasswordHash: "hashed-secret",
		Active:       true,
	}
}

// --- Login tests ---

func TestService_Login_Success(t *testing.T) {
	u := newTestUser()
	repo := &mockUserRepo{userByEmail: u}
	hasher := &mockHasher{}
	tok := &mockTokenGen{token: "jwt-token"}

	svc := NewService(repo, hasher, nil, tok, Config{JWTExpHours: 24, SecureCookie: true})

	res, err := svc.Login(context.Background(), LoginInput{Email: "jane@example.com", Password: "secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Token != "jwt-token" {
		t.Errorf("token = %q, want %q", res.Token, "jwt-token")
	}
	if res.User != u {
		t.Error("expected returned user to match repo user")
	}
	if !res.SecureCookie {
		t.Error("expected SecureCookie=true")
	}
	if res.ExpiresAt.IsZero() {
		t.Error("expected non-zero ExpiresAt")
	}
	if repo.gotEmail != "jane@example.com" {
		t.Errorf("repo got email %q", repo.gotEmail)
	}
	if hasher.gotHash != "hashed-secret" || hasher.gotPwd != "secret" {
		t.Errorf("hasher got hash=%q pwd=%q", hasher.gotHash, hasher.gotPwd)
	}
	if tok.gotUserID != u.ID || tok.gotTenant != u.TenantID || tok.gotRole != "agent" {
		t.Errorf("token gen got userID=%v tenant=%v role=%q", tok.gotUserID, tok.gotTenant, tok.gotRole)
	}
}

func TestService_Login_EmptyCredentials(t *testing.T) {
	svc := NewService(&mockUserRepo{}, &mockHasher{}, nil, &mockTokenGen{}, Config{})

	tests := []struct {
		name string
		in   LoginInput
	}{
		{"empty email", LoginInput{Password: "x"}},
		{"empty password", LoginInput{Email: "a@b.com"}},
		{"both empty", LoginInput{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := svc.Login(context.Background(), tt.in); !errors.Is(err, ErrInvalidCredentials) {
				t.Errorf("expected ErrInvalidCredentials, got %v", err)
			}
		})
	}
}

func TestService_Login_WrongEmail(t *testing.T) {
	repo := &mockUserRepo{errByEmail: errors.New("not found")}
	svc := NewService(repo, &mockHasher{}, nil, &mockTokenGen{}, Config{})

	_, err := svc.Login(context.Background(), LoginInput{Email: "nope@example.com", Password: "secret"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestService_Login_WrongPassword(t *testing.T) {
	u := newTestUser()
	repo := &mockUserRepo{userByEmail: u}
	hasher := &mockHasher{compareErr: errors.New("mismatch")}
	svc := NewService(repo, hasher, nil, &mockTokenGen{}, Config{})

	_, err := svc.Login(context.Background(), LoginInput{Email: "jane@example.com", Password: "wrong"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestService_Login_LDAPPath(t *testing.T) {
	u := newTestUser()
	repo := &mockUserRepo{userByEmail: u}
	ldap := &mockLDAP{}
	tok := &mockTokenGen{token: "ldap-jwt"}
	svc := NewService(repo, &mockHasher{}, ldap, tok, Config{JWTExpHours: 1})

	res, err := svc.Login(context.Background(), LoginInput{Email: "jane@example.com", Password: "secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Token != "ldap-jwt" {
		t.Errorf("token = %q", res.Token)
	}
	if ldap.gotUser != "jane@example.com" || ldap.gotPwd != "secret" {
		t.Errorf("ldap got user=%q pwd=%q", ldap.gotUser, ldap.gotPwd)
	}
}

func TestService_Login_LDAPFailure(t *testing.T) {
	u := newTestUser()
	repo := &mockUserRepo{userByEmail: u}
	ldap := &mockLDAP{err: errors.New("ldap down")}
	svc := NewService(repo, &mockHasher{}, ldap, &mockTokenGen{}, Config{})

	_, err := svc.Login(context.Background(), LoginInput{Email: "jane@example.com", Password: "secret"})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestService_Login_PasswordNotConfigured(t *testing.T) {
	u := newTestUser()
	u.PasswordHash = ""
	repo := &mockUserRepo{userByEmail: u}
	svc := NewService(repo, &mockHasher{}, nil, &mockTokenGen{}, Config{})

	_, err := svc.Login(context.Background(), LoginInput{Email: "jane@example.com", Password: "secret"})
	if !errors.Is(err, ErrPasswordNotConfigured) {
		t.Errorf("expected ErrPasswordNotConfigured, got %v", err)
	}
}

func TestService_Login_TokenGenerationFails(t *testing.T) {
	u := newTestUser()
	repo := &mockUserRepo{userByEmail: u}
	tok := &mockTokenGen{err: errors.New("signing failed")}
	svc := NewService(repo, &mockHasher{}, nil, tok, Config{})

	_, err := svc.Login(context.Background(), LoginInput{Email: "jane@example.com", Password: "secret"})
	if err == nil {
		t.Fatal("expected error from token generation")
	}
	if !errors.Is(err, tok.err) {
		t.Errorf("expected token gen error, got %v", err)
	}
}

// --- GetProfile tests ---

func TestService_GetProfile_Success(t *testing.T) {
	u := newTestUser()
	id := u.ID
	repo := &mockUserRepo{userByID: u}
	svc := NewService(repo, &mockHasher{}, nil, &mockTokenGen{}, Config{})

	got, err := svc.GetProfile(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != id {
		t.Errorf("got user id %v, want %v", got.ID, id)
	}
	if repo.gotID != id {
		t.Errorf("repo got id %v, want %v", repo.gotID, id)
	}
}

func TestService_GetProfile_NotFound(t *testing.T) {
	repo := &mockUserRepo{errByID: errors.New("not found")}
	svc := NewService(repo, &mockHasher{}, nil, &mockTokenGen{}, Config{})

	_, err := svc.GetProfile(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// --- Cross-tenant isolation ---

type mockUserCreator struct {
	created *user.User
	err     error
}

func (m *mockUserCreator) Create(_ context.Context, u *user.User, _ string) error {
	m.created = u
	return m.err
}

type mockFullHasher struct{ err error }

func (m *mockFullHasher) Hash(password string) (string, error) { return "hashed-" + password, m.err }
func (m *mockFullHasher) Compare(_, _ string) error            { return nil }

// TestService_Login_IsTenantScoped is the regression guard for the login
// ambiguity: the schema allows the same address in several tenants
// (users_email_tenant_unique is on (tenant_id, email)), so a login must return
// the account belonging to the tenant it was issued against — never whichever
// row the database happened to return first.
func TestService_Login_IsTenantScoped(t *testing.T) {
	const sharedEmail = "gustavo@example.com"

	tenantA, tenantB := uuid.New(), uuid.New()
	userA := &user.User{ID: uuid.New(), TenantID: tenantA, Email: sharedEmail, Role: user.RoleAdmin, PasswordHash: "hash-a", Active: true}
	userB := &user.User{ID: uuid.New(), TenantID: tenantB, Email: sharedEmail, Role: user.RoleEndUser, PasswordHash: "hash-b", Active: true}

	tests := []struct {
		name     string
		tenantID uuid.UUID
		want     *user.User
	}{
		{"tenant A gets its own account", tenantA, userA},
		{"tenant B gets its own account", tenantB, userB},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockUserRepo{usersByTenant: map[uuid.UUID]*user.User{
				tenantA: userA,
				tenantB: userB,
			}}
			tok := &mockTokenGen{token: "jwt"}
			svc := NewService(repo, &mockHasher{}, nil, tok, Config{JWTExpHours: 1})

			res, err := svc.Login(context.Background(), LoginInput{
				TenantID: tt.tenantID,
				Email:    sharedEmail,
				Password: "secret",
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.User.ID != tt.want.ID {
				t.Errorf("logged in as user %v, want %v — cross-tenant leak", res.User.ID, tt.want.ID)
			}
			if repo.gotTenantID != tt.tenantID {
				t.Errorf("repo queried tenant %v, want %v", repo.gotTenantID, tt.tenantID)
			}
			// The JWT must be minted for the tenant that was asked for.
			if tok.gotTenant != tt.tenantID {
				t.Errorf("token minted for tenant %v, want %v", tok.gotTenant, tt.tenantID)
			}
		})
	}
}

// TestService_Login_UnknownTenant covers the account existing, but not in the
// tenant being logged into.
func TestService_Login_UnknownTenant(t *testing.T) {
	known := uuid.New()
	repo := &mockUserRepo{
		usersByTenant: map[uuid.UUID]*user.User{
			known: {ID: uuid.New(), TenantID: known, Email: "a@b.com", PasswordHash: "h", Active: true},
		},
		errByEmail: errors.New("not found"),
	}
	svc := NewService(repo, &mockHasher{}, nil, &mockTokenGen{}, Config{})

	_, err := svc.Login(context.Background(), LoginInput{
		TenantID: uuid.New(), // a different tenant
		Email:    "a@b.com",
		Password: "secret",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

// TestService_Register_RequiresTenant guards the registration path: without a
// resolved tenant there is no organization to create the account in, and
// silently falling back to a default one is how accounts end up in the wrong
// place.
func TestService_Register_RequiresTenant(t *testing.T) {
	repo := &mockUserRepo{}
	svc := NewService(repo, &mockHasher{}, nil, &mockTokenGen{}, Config{})
	svc.SetUserCreator(&mockUserCreator{})
	svc.SetPasswordHasherFull(&mockFullHasher{})

	_, err := svc.Register(context.Background(), RegisterInput{
		Name:     "Gustavo",
		Email:    "gustavo@example.com",
		Password: "secret123",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials for a nil tenant, got %v", err)
	}
	if repo.byEmailCalls != 0 {
		t.Error("must reject before touching the repository")
	}
}
