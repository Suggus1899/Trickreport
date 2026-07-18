package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	domainUser "github.com/trickreport/backend/internal/domain/user"
)

// --- Mocks ---

type mockRepo struct {
	users        []domainUser.User
	user         *domainUser.User
	err          error
	createErr    error
	updateErr    error
	deactErr     error
	gotEmail     string
	gotCreated   *domainUser.User
	gotHash      string
	gotFields    UpdateFields
	createCalls  int
	updateCalls  int
	deactCalls   int
}

func (m *mockRepo) List(ctx context.Context, tenantID uuid.UUID) ([]domainUser.User, error) {
	return m.users, m.err
}

func (m *mockRepo) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*domainUser.User, error) {
	return m.user, m.err
}

func (m *mockRepo) GetByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	m.gotEmail = email
	return m.user, m.err
}

func (m *mockRepo) Create(ctx context.Context, u *domainUser.User, passwordHash string) error {
	m.createCalls++
	m.gotCreated = u
	m.gotHash = passwordHash
	return m.createErr
}

func (m *mockRepo) Update(ctx context.Context, id, tenantID uuid.UUID, fields UpdateFields) (*domainUser.User, error) {
	m.updateCalls++
	m.gotFields = fields
	return m.user, m.updateErr
}

func (m *mockRepo) Deactivate(ctx context.Context, id, tenantID uuid.UUID) error {
	m.deactCalls++
	return m.deactErr
}

type mockHasher struct {
	hash       string
	hashErr    error
	compareErr error
	gotPwd     string
}

func (m *mockHasher) Hash(password string) (string, error) {
	m.gotPwd = password
	return m.hash, m.hashErr
}

func (m *mockHasher) Compare(hash, password string) error {
	return m.compareErr
}

// --- List tests ---

func TestService_List_Delegates(t *testing.T) {
	want := []domainUser.User{{ID: uuid.New(), Name: "a"}}
	repo := &mockRepo{users: want}
	svc := NewService(repo, &mockHasher{})

	got, err := svc.List(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d", len(got))
	}
}

func TestService_List_Error(t *testing.T) {
	repo := &mockRepo{err: errors.New("db")}
	svc := NewService(repo, &mockHasher{})
	if _, err := svc.List(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

// --- Get tests ---

func TestService_Get_Delegates(t *testing.T) {
	u := &domainUser.User{ID: uuid.New()}
	repo := &mockRepo{user: u}
	svc := NewService(repo, &mockHasher{})

	got, err := svc.Get(context.Background(), u.ID, uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != u.ID {
		t.Errorf("got id %v", got.ID)
	}
}

func TestService_Get_NotFound(t *testing.T) {
	repo := &mockRepo{err: domainUser.ErrNotFound}
	svc := NewService(repo, &mockHasher{})
	_, err := svc.Get(context.Background(), uuid.New(), uuid.New())
	if !errors.Is(err, domainUser.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// --- Create tests ---

func TestService_Create_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &mockRepo{}
	hasher := &mockHasher{hash: "hashed"}
	svc := NewService(repo, hasher)

	u, err := svc.Create(context.Background(), CreateInput{
		TenantID: tenant, Name: "Jane", Email: "jane@x.com", Role: "agent", Password: "secret",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Name != "Jane" || u.Email != "jane@x.com" || u.Role != domainUser.RoleAgent {
		t.Errorf("user = %+v", u)
	}
	if u.TenantID != tenant {
		t.Errorf("tenant = %v", u.TenantID)
	}
	if repo.gotHash != "hashed" {
		t.Errorf("hash = %q", repo.gotHash)
	}
	if hasher.gotPwd != "secret" {
		t.Errorf("hasher got pwd %q", hasher.gotPwd)
	}
	if repo.createCalls != 1 {
		t.Errorf("create calls = %d", repo.createCalls)
	}
}

func TestService_Create_DefaultRoleEndUser(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo, &mockHasher{})

	u, err := svc.Create(context.Background(), CreateInput{Name: "Jane", Email: "jane@x.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Role != domainUser.RoleEndUser {
		t.Errorf("role = %q, want end_user", u.Role)
	}
}

func TestService_Create_NoPasswordEmptyHash(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo, &mockHasher{})

	if _, err := svc.Create(context.Background(), CreateInput{Name: "Jane", Email: "jane@x.com"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.gotHash != "" {
		t.Errorf("hash = %q, want empty", repo.gotHash)
	}
}

func TestService_Create_EmptyName(t *testing.T) {
	svc := NewService(&mockRepo{}, &mockHasher{})
	_, err := svc.Create(context.Background(), CreateInput{Name: "", Email: "a@b.com"})
	if !errors.Is(err, domainUser.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Create_EmptyEmail(t *testing.T) {
	svc := NewService(&mockRepo{}, &mockHasher{})
	_, err := svc.Create(context.Background(), CreateInput{Name: "Jane", Email: ""})
	if !errors.Is(err, domainUser.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Create_InvalidRole(t *testing.T) {
	svc := NewService(&mockRepo{}, &mockHasher{})
	_, err := svc.Create(context.Background(), CreateInput{Name: "Jane", Email: "a@b.com", Role: "superuser"})
	if !errors.Is(err, domainUser.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Create_HashError(t *testing.T) {
	hasher := &mockHasher{hashErr: errors.New("bcrypt fail")}
	svc := NewService(&mockRepo{}, hasher)
	_, err := svc.Create(context.Background(), CreateInput{Name: "Jane", Email: "a@b.com", Password: "x"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestService_Create_RepoError(t *testing.T) {
	repo := &mockRepo{createErr: domainUser.ErrAlreadyExists}
	svc := NewService(repo, &mockHasher{})
	_, err := svc.Create(context.Background(), CreateInput{Name: "Jane", Email: "a@b.com"})
	if !errors.Is(err, domainUser.ErrAlreadyExists) {
		t.Errorf("expected ErrAlreadyExists, got %v", err)
	}
}

// --- Update tests ---

func TestService_Update_Delegates(t *testing.T) {
	u := &domainUser.User{ID: uuid.New(), Name: "new"}
	repo := &mockRepo{user: u}
	svc := NewService(repo, &mockHasher{})

	name := "new"
	fields := UpdateFields{Name: &name}
	got, err := svc.Update(context.Background(), u.ID, uuid.New(), fields)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != u.ID {
		t.Errorf("got id %v", got.ID)
	}
	if repo.gotFields.Name == nil || *repo.gotFields.Name != "new" {
		t.Errorf("fields = %+v", repo.gotFields)
	}
	if repo.updateCalls != 1 {
		t.Errorf("update calls = %d", repo.updateCalls)
	}
}

func TestService_Update_Error(t *testing.T) {
	repo := &mockRepo{updateErr: errors.New("nope")}
	svc := NewService(repo, &mockHasher{})
	_, err := svc.Update(context.Background(), uuid.New(), uuid.New(), UpdateFields{})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Deactivate tests ---

func TestService_Deactivate_Delegates(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo, &mockHasher{})
	if err := svc.Deactivate(context.Background(), uuid.New(), uuid.New()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.deactCalls != 1 {
		t.Errorf("deact calls = %d", repo.deactCalls)
	}
}

func TestService_Deactivate_Error(t *testing.T) {
	repo := &mockRepo{deactErr: errors.New("nope")}
	svc := NewService(repo, &mockHasher{})
	if err := svc.Deactivate(context.Background(), uuid.New(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}
