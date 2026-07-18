package article

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/domain/article"
)

// --- Mocks ---

type mockRepo struct {
	articles    []article.Article
	article     *article.Article
	err         error
	createErr   error
	updateErr   error
	deleteErr   error
	gotFilter   Filter
	gotRole     string
	gotCreated  *article.Article
	gotUpdated  *article.Article
	createCalls int
	updateCalls int
	deleteCalls int
}

func (m *mockRepo) List(ctx context.Context, tenantID uuid.UUID, filter Filter, role string) ([]article.Article, error) {
	m.gotFilter = filter
	m.gotRole = role
	return m.articles, m.err
}

func (m *mockRepo) GetByID(ctx context.Context, id, tenantID uuid.UUID, role string) (*article.Article, error) {
	m.gotRole = role
	return m.article, m.err
}

func (m *mockRepo) Create(ctx context.Context, a *article.Article) error {
	m.createCalls++
	m.gotCreated = a
	return m.createErr
}

func (m *mockRepo) Update(ctx context.Context, a *article.Article) error {
	m.updateCalls++
	m.gotUpdated = a
	return m.updateErr
}

func (m *mockRepo) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
	m.deleteCalls++
	return m.deleteErr
}

// --- List tests ---

func TestService_List_Delegates(t *testing.T) {
	want := []article.Article{{ID: uuid.New(), Title: "a"}}
	repo := &mockRepo{articles: want}
	svc := NewService(repo)

	got, err := svc.List(context.Background(), uuid.New(), Filter{Search: "x"}, "agent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d", len(got))
	}
	if repo.gotFilter.Search != "x" || repo.gotRole != "agent" {
		t.Errorf("repo got filter=%+v role=%q", repo.gotFilter, repo.gotRole)
	}
}

func TestService_List_Error(t *testing.T) {
	repo := &mockRepo{err: errors.New("db")}
	svc := NewService(repo)
	if _, err := svc.List(context.Background(), uuid.New(), Filter{}, "agent"); err == nil {
		t.Fatal("expected error")
	}
}

// --- Get tests ---

func TestService_Get_Delegates(t *testing.T) {
	a := &article.Article{ID: uuid.New(), Title: "x"}
	repo := &mockRepo{article: a}
	svc := NewService(repo)

	got, err := svc.Get(context.Background(), a.ID, uuid.New(), "end_user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != a.ID {
		t.Errorf("got id %v", got.ID)
	}
	if repo.gotRole != "end_user" {
		t.Errorf("role = %q", repo.gotRole)
	}
}

func TestService_Get_NotFound(t *testing.T) {
	repo := &mockRepo{err: article.ErrNotFound}
	svc := NewService(repo)
	_, err := svc.Get(context.Background(), uuid.New(), uuid.New(), "agent")
	if !errors.Is(err, article.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// --- Create tests ---

func TestService_Create_Success(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	repo := &mockRepo{}
	svc := NewService(repo)

	a, err := svc.Create(context.Background(), CreateInput{
		TenantID: tenant, Title: "How to reset", Content: "Steps here",
		Category: "guides", Tags: []string{"reset"}, Published: true, CreatedBy: creator,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Title != "How to reset" || a.Category != "guides" || !a.Published {
		t.Errorf("article = %+v", a)
	}
	if a.TenantID != tenant || a.CreatedBy != creator {
		t.Errorf("tenant=%v creator=%v", a.TenantID, a.CreatedBy)
	}
	if repo.createCalls != 1 {
		t.Errorf("create calls = %d", repo.createCalls)
	}
}

func TestService_Create_DefaultsCategoryAndTags(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo)

	a, err := svc.Create(context.Background(), CreateInput{Title: "t", Content: "c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Category != "general" {
		t.Errorf("category = %q, want general", a.Category)
	}
	if a.Tags == nil || len(a.Tags) != 0 {
		t.Errorf("tags = %v, want empty slice", a.Tags)
	}
}

func TestService_Create_EmptyTitle(t *testing.T) {
	svc := NewService(&mockRepo{})
	_, err := svc.Create(context.Background(), CreateInput{Title: "", Content: "c"})
	if !errors.Is(err, article.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Create_EmptyContent(t *testing.T) {
	svc := NewService(&mockRepo{})
	_, err := svc.Create(context.Background(), CreateInput{Title: "t", Content: ""})
	if !errors.Is(err, article.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Create_RepoError(t *testing.T) {
	repo := &mockRepo{createErr: errors.New("insert")}
	svc := NewService(repo)
	_, err := svc.Create(context.Background(), CreateInput{Title: "t", Content: "c"})
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

	a, err := svc.Update(context.Background(), UpdateInput{
		ID: id, TenantID: tenant, Title: "Updated", Content: "new content",
		Category: "guides", Tags: []string{"x"}, Published: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.ID != id || a.Title != "Updated" || a.TenantID != tenant {
		t.Errorf("article = %+v", a)
	}
	if repo.updateCalls != 1 {
		t.Errorf("update calls = %d", repo.updateCalls)
	}
}

func TestService_Update_DefaultsTags(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo)
	a, err := svc.Update(context.Background(), UpdateInput{ID: uuid.New(), TenantID: uuid.New(), Title: "t", Content: "c"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Tags == nil || len(a.Tags) != 0 {
		t.Errorf("tags = %v, want empty slice", a.Tags)
	}
}

func TestService_Update_EmptyTitle(t *testing.T) {
	svc := NewService(&mockRepo{})
	_, err := svc.Update(context.Background(), UpdateInput{Title: "", Content: "c"})
	if !errors.Is(err, article.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Update_EmptyContent(t *testing.T) {
	svc := NewService(&mockRepo{})
	_, err := svc.Update(context.Background(), UpdateInput{Title: "t", Content: ""})
	if !errors.Is(err, article.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Update_RepoError(t *testing.T) {
	repo := &mockRepo{updateErr: errors.New("update")}
	svc := NewService(repo)
	_, err := svc.Update(context.Background(), UpdateInput{Title: "t", Content: "c"})
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
