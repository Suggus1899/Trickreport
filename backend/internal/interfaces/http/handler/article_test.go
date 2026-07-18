package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appArticle "github.com/trickreport/backend/internal/application/article"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	"github.com/trickreport/backend/internal/domain/article"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
)

// --- Mocks ---

type aMockRepo struct {
	articles   []article.Article
	article    *article.Article
	err        error
	createErr  error
	updateErr  error
	deleteErr  error
	gotRole    string
	createCalls int
}

func (m *aMockRepo) List(ctx context.Context, tenantID uuid.UUID, filter appArticle.Filter, role string) ([]article.Article, error) {
	m.gotRole = role
	return m.articles, m.err
}
func (m *aMockRepo) GetByID(ctx context.Context, id, tenantID uuid.UUID, role string) (*article.Article, error) {
	m.gotRole = role
	return m.article, m.err
}
func (m *aMockRepo) Create(ctx context.Context, a *article.Article) error {
	m.createCalls++
	return m.createErr
}
func (m *aMockRepo) Update(ctx context.Context, a *article.Article) error {
	return m.updateErr
}
func (m *aMockRepo) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
	return m.deleteErr
}

func newArticleHandler(repo appArticle.Repository) *ArticleHandler {
	return NewArticleHandler(appArticle.NewService(repo))
}

func setupArticleRouter(h *ArticleHandler, tenantID uuid.UUID) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Tenant(&mockTenantResolver{id: tenantID}))
	r.Get("/api/articles", h.List)
	r.Post("/api/articles", h.Create)
	r.Get("/api/articles/{id}", h.Get)
	r.Put("/api/articles/{id}", h.Update)
	r.Delete("/api/articles/{id}", h.Delete)
	return r
}

// --- List ---

func TestArticleHandler_List_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &aMockRepo{articles: []article.Article{{ID: uuid.New(), Title: "a1"}}}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/articles?q=test", "", claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if repo.gotRole != "agent" {
		t.Errorf("role = %q", repo.gotRole)
	}
}

func TestArticleHandler_List_Error(t *testing.T) {
	tenant := uuid.New()
	repo := &aMockRepo{err: errors.New("db")}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/articles", "", claims)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

// --- Get ---

func TestArticleHandler_Get_Success(t *testing.T) {
	tenant := uuid.New()
	a := &article.Article{ID: uuid.New(), Title: "test"}
	repo := &aMockRepo{article: a}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/articles/"+a.ID.String(), "", claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestArticleHandler_Get_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newArticleHandler(&aMockRepo{})
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/articles/bad", "", claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestArticleHandler_Get_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &aMockRepo{err: article.ErrNotFound}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/articles/"+uuid.New().String(), "", claims)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d", rec.Code)
	}
}

// --- Create ---

func TestArticleHandler_Create_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &aMockRepo{}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"title":"How to reset","content":"Steps to reset password","category":"guides","published":true}`
	rec := doRequest(router, http.MethodPost, "/api/articles", body, claims)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
	if repo.createCalls != 1 {
		t.Errorf("create calls = %d", repo.createCalls)
	}
}

func TestArticleHandler_Create_InvalidJSON(t *testing.T) {
	tenant := uuid.New()
	h := newArticleHandler(&aMockRepo{})
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodPost, "/api/articles", "{bad", claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestArticleHandler_Create_ValidationError(t *testing.T) {
	tenant := uuid.New()
	h := newArticleHandler(&aMockRepo{})
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"title":"ab","content":"too short","category":"x"}`
	rec := doRequest(router, http.MethodPost, "/api/articles", body, claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestArticleHandler_Create_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &aMockRepo{createErr: errors.New("db")}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"title":"How to reset","content":"Steps to reset password","category":"guides"}`
	rec := doRequest(router, http.MethodPost, "/api/articles", body, claims)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

// --- Update ---

func TestArticleHandler_Update_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &aMockRepo{article: &article.Article{ID: uuid.New(), Title: "updated"}}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"title":"Updated title","content":"Updated content here","category":"guides"}`
	rec := doRequest(router, http.MethodPut, "/api/articles/"+uuid.New().String(), body, claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestArticleHandler_Update_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newArticleHandler(&aMockRepo{})
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"title":"Updated","content":"Content here","category":"x"}`
	rec := doRequest(router, http.MethodPut, "/api/articles/bad", body, claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestArticleHandler_Update_InvalidJSON(t *testing.T) {
	tenant := uuid.New()
	h := newArticleHandler(&aMockRepo{})
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodPut, "/api/articles/"+uuid.New().String(), "{bad", claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestArticleHandler_Update_ValidationError(t *testing.T) {
	tenant := uuid.New()
	h := newArticleHandler(&aMockRepo{})
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"title":"ab","content":"too short","category":"x"}`
	rec := doRequest(router, http.MethodPut, "/api/articles/"+uuid.New().String(), body, claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestArticleHandler_Update_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &aMockRepo{updateErr: article.ErrNotFound}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"title":"Updated","content":"Content here ok","category":"x"}`
	rec := doRequest(router, http.MethodPut, "/api/articles/"+uuid.New().String(), body, claims)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestArticleHandler_Update_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &aMockRepo{updateErr: errors.New("db")}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"title":"Updated","content":"Content here ok","category":"x"}`
	rec := doRequest(router, http.MethodPut, "/api/articles/"+uuid.New().String(), body, claims)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

// --- Delete ---

func TestArticleHandler_Delete_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &aMockRepo{}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodDelete, "/api/articles/"+uuid.New().String(), "", claims)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestArticleHandler_Delete_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newArticleHandler(&aMockRepo{})
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodDelete, "/api/articles/bad", "", claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestArticleHandler_Delete_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &aMockRepo{deleteErr: article.ErrNotFound}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodDelete, "/api/articles/"+uuid.New().String(), "", claims)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestArticleHandler_Delete_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &aMockRepo{deleteErr: errors.New("db")}
	h := newArticleHandler(repo)
	router := setupArticleRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodDelete, "/api/articles/"+uuid.New().String(), "", claims)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

// Ensure strings import is used.
var _ = strings.NewReader
