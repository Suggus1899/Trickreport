package automation

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/domain/automation"
)

// --- Mocks ---

type mockRepo struct {
	rules       []automation.Rule
	err         error
	createErr   error
	updateErr   error
	deleteErr   error
	gotCreated  *automation.Rule
	gotUpdated  *automation.Rule
	createCalls int
	updateCalls int
	deleteCalls int
}

func (m *mockRepo) List(ctx context.Context, tenantID uuid.UUID) ([]automation.Rule, error) {
	return m.rules, m.err
}

func (m *mockRepo) Create(ctx context.Context, r *automation.Rule) error {
	m.createCalls++
	m.gotCreated = r
	return m.createErr
}

func (m *mockRepo) Update(ctx context.Context, r *automation.Rule) error {
	m.updateCalls++
	m.gotUpdated = r
	return m.updateErr
}

func (m *mockRepo) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
	m.deleteCalls++
	return m.deleteErr
}

// --- List tests ---

func TestService_List_Delegates(t *testing.T) {
	want := []automation.Rule{{ID: uuid.New(), Name: "r"}}
	repo := &mockRepo{rules: want}
	svc := NewService(repo)

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
	svc := NewService(repo)
	if _, err := svc.List(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

// --- Create tests ---

func TestService_Create_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &mockRepo{}
	svc := NewService(repo)

	conditions := map[string]any{"priority": "high"}
	actions := []any{"notify"}
	r, err := svc.Create(context.Background(), CreateInput{
		TenantID: tenant, Name: "Auto-assign", Description: "desc",
		TriggerType: "on_create", Conditions: conditions, Actions: actions, IsActive: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Name != "Auto-assign" || r.TriggerType != "on_create" || !r.IsActive {
		t.Errorf("rule = %+v", r)
	}
	if r.TenantID != tenant {
		t.Errorf("tenant = %v", r.TenantID)
	}
	if repo.createCalls != 1 {
		t.Errorf("create calls = %d", repo.createCalls)
	}
}

func TestService_Create_EmptyName(t *testing.T) {
	svc := NewService(&mockRepo{})
	_, err := svc.Create(context.Background(), CreateInput{Name: "", TriggerType: "on_create"})
	if !errors.Is(err, automation.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Create_EmptyTriggerType(t *testing.T) {
	svc := NewService(&mockRepo{})
	_, err := svc.Create(context.Background(), CreateInput{Name: "r", TriggerType: ""})
	if !errors.Is(err, automation.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Create_RepoError(t *testing.T) {
	repo := &mockRepo{createErr: errors.New("insert")}
	svc := NewService(repo)
	_, err := svc.Create(context.Background(), CreateInput{Name: "r", TriggerType: "on_create"})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Update tests ---

func TestService_Update_Success(t *testing.T) {
	id := uuid.New()
	tenant := uuid.New()
	repo := &mockRepo{}
	svc := NewService(repo)

	r, err := svc.Update(context.Background(), UpdateInput{
		ID: id, TenantID: tenant, Name: "Updated", TriggerType: "on_update", IsActive: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.ID != id || r.TenantID != tenant || r.Name != "Updated" {
		t.Errorf("rule = %+v", r)
	}
	if repo.updateCalls != 1 {
		t.Errorf("update calls = %d", repo.updateCalls)
	}
}

func TestService_Update_EmptyName(t *testing.T) {
	svc := NewService(&mockRepo{})
	_, err := svc.Update(context.Background(), UpdateInput{Name: "", TriggerType: "on_update"})
	if !errors.Is(err, automation.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Update_EmptyTriggerType(t *testing.T) {
	svc := NewService(&mockRepo{})
	_, err := svc.Update(context.Background(), UpdateInput{Name: "r", TriggerType: ""})
	if !errors.Is(err, automation.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Update_RepoError(t *testing.T) {
	repo := &mockRepo{updateErr: errors.New("update")}
	svc := NewService(repo)
	_, err := svc.Update(context.Background(), UpdateInput{Name: "r", TriggerType: "on_update"})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Delete tests ---

func TestService_Delete_Delegates(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo)
	if err := svc.Delete(context.Background(), uuid.New(), uuid.New()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.deleteCalls != 1 {
		t.Errorf("delete calls = %d", repo.deleteCalls)
	}
}

func TestService_Delete_Error(t *testing.T) {
	repo := &mockRepo{deleteErr: errors.New("nope")}
	svc := NewService(repo)
	if err := svc.Delete(context.Background(), uuid.New(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}
