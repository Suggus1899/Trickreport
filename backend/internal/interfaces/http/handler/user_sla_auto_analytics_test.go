package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appAnalytics "github.com/trickreport/backend/internal/application/analytics"
	appAuto "github.com/trickreport/backend/internal/application/automation"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	appSLA "github.com/trickreport/backend/internal/application/sla"
	appUser "github.com/trickreport/backend/internal/application/user"
	domainAuto "github.com/trickreport/backend/internal/domain/automation"
	domainSLA "github.com/trickreport/backend/internal/domain/sla"
	domainUser "github.com/trickreport/backend/internal/domain/user"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
)

// ============= USER HANDLER TESTS =============

type uMockRepo struct {
	users     []domainUser.User
	user      *domainUser.User
	err       error
	createErr error
	updateErr error
	deactErr  error
}

func (m *uMockRepo) List(ctx context.Context, tenantID uuid.UUID) ([]domainUser.User, error) {
	return m.users, m.err
}
func (m *uMockRepo) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*domainUser.User, error) {
	return m.user, m.err
}
func (m *uMockRepo) GetByEmail(ctx context.Context, email string) (*domainUser.User, error) {
	return m.user, m.err
}
func (m *uMockRepo) Create(ctx context.Context, u *domainUser.User, passwordHash string) error {
	return m.createErr
}
func (m *uMockRepo) Update(ctx context.Context, id, tenantID uuid.UUID, fields appUser.UpdateFields) (*domainUser.User, error) {
	return m.user, m.updateErr
}
func (m *uMockRepo) Deactivate(ctx context.Context, id, tenantID uuid.UUID) error {
	return m.deactErr
}

type uMockHasher struct{}

func (m *uMockHasher) Hash(password string) (string, error) { return "hashed", nil }
func (m *uMockHasher) Compare(hash, password string) error  { return nil }

func newUserHandler(repo appUser.Repository) *UserHandler {
	return NewUserHandler(appUser.NewService(repo, &uMockHasher{}))
}

func setupUserRouter(h *UserHandler, tenantID uuid.UUID) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Tenant(&mockTenantResolver{id: tenantID}))
	r.Get("/api/users", h.List)
	r.Post("/api/users", h.Create)
	r.Get("/api/users/{id}", h.Get)
	r.Put("/api/users/{id}", h.Update)
	r.Delete("/api/users/{id}", h.Delete)
	return r
}

func TestUserHandler_List_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &uMockRepo{users: []domainUser.User{{ID: uuid.New(), Name: "a"}}}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/users", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestUserHandler_List_Error(t *testing.T) {
	tenant := uuid.New()
	repo := &uMockRepo{err: errors.New("db")}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/users", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Get_Success(t *testing.T) {
	tenant := uuid.New()
	u := &domainUser.User{ID: uuid.New(), Name: "Jane"}
	repo := &uMockRepo{user: u}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/users/"+u.ID.String(), "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestUserHandler_Get_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newUserHandler(&uMockRepo{})
	router := setupUserRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/users/bad", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Get_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &uMockRepo{err: domainUser.ErrNotFound}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/users/"+uuid.New().String(), "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Create_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &uMockRepo{user: &domainUser.User{ID: uuid.New(), Name: "Jane"}}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	body := `{"name":"Jane Doe","email":"jane@x.com","role":"agent","password":"secret123"}`
	rec := doRequest(router, http.MethodPost, "/api/users", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestUserHandler_Create_InvalidJSON(t *testing.T) {
	tenant := uuid.New()
	h := newUserHandler(&uMockRepo{})
	router := setupUserRouter(h, tenant)

	rec := doRequest(router, http.MethodPost, "/api/users", "{bad", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Create_ValidationError(t *testing.T) {
	tenant := uuid.New()
	h := newUserHandler(&uMockRepo{})
	router := setupUserRouter(h, tenant)

	body := `{"name":"a","email":"jane@x.com","role":"agent","password":"secret123"}`
	rec := doRequest(router, http.MethodPost, "/api/users", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Create_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &uMockRepo{createErr: domainUser.ErrAlreadyExists}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	body := `{"name":"Jane Doe","email":"jane@x.com","role":"agent","password":"secret123"}`
	rec := doRequest(router, http.MethodPost, "/api/users", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Update_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &uMockRepo{user: &domainUser.User{ID: uuid.New(), Name: "updated"}}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	body := `{"name":"Updated Name"}`
	rec := doRequest(router, http.MethodPut, "/api/users/"+uuid.New().String(), body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestUserHandler_Update_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newUserHandler(&uMockRepo{})
	router := setupUserRouter(h, tenant)

	body := `{"name":"Updated"}`
	rec := doRequest(router, http.MethodPut, "/api/users/bad", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Update_InvalidJSON(t *testing.T) {
	tenant := uuid.New()
	h := newUserHandler(&uMockRepo{})
	router := setupUserRouter(h, tenant)

	rec := doRequest(router, http.MethodPut, "/api/users/"+uuid.New().String(), "{bad", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Update_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &uMockRepo{updateErr: domainUser.ErrNotFound}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	body := `{"name":"Updated Name"}`
	rec := doRequest(router, http.MethodPut, "/api/users/"+uuid.New().String(), body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Update_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &uMockRepo{updateErr: errors.New("db")}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	body := `{"name":"Updated Name"}`
	rec := doRequest(router, http.MethodPut, "/api/users/"+uuid.New().String(), body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Delete_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &uMockRepo{}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	rec := doRequest(router, http.MethodDelete, "/api/users/"+uuid.New().String(), "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Delete_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newUserHandler(&uMockRepo{})
	router := setupUserRouter(h, tenant)

	rec := doRequest(router, http.MethodDelete, "/api/users/bad", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Delete_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &uMockRepo{deactErr: domainUser.ErrNotFound}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	rec := doRequest(router, http.MethodDelete, "/api/users/"+uuid.New().String(), "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestUserHandler_Delete_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &uMockRepo{deactErr: errors.New("db")}
	h := newUserHandler(repo)
	router := setupUserRouter(h, tenant)

	rec := doRequest(router, http.MethodDelete, "/api/users/"+uuid.New().String(), "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

// ============= SLA HANDLER TESTS =============

type slaMockRepo struct {
	policies  []domainSLA.Policy
	policy    *domainSLA.Policy
	err       error
	upsertErr error
}

func (m *slaMockRepo) List(ctx context.Context, tenantID uuid.UUID) ([]domainSLA.Policy, error) {
	return m.policies, m.err
}
func (m *slaMockRepo) Upsert(ctx context.Context, p *domainSLA.Policy) error {
	return m.upsertErr
}

func newSLAHandler(repo appSLA.Repository) *SLAHandler {
	return NewSLAHandler(appSLA.NewService(repo))
}

func setupSLARouter(h *SLAHandler, tenantID uuid.UUID) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Tenant(&mockTenantResolver{id: tenantID}))
	r.Get("/api/sla", h.List)
	r.Put("/api/sla/{priority}", h.Upsert)
	return r
}

func TestSLAHandler_List_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &slaMockRepo{policies: []domainSLA.Policy{{Priority: "high"}}}
	h := newSLAHandler(repo)
	router := setupSLARouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/sla", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestSLAHandler_List_Error(t *testing.T) {
	tenant := uuid.New()
	repo := &slaMockRepo{err: errors.New("db")}
	h := newSLAHandler(repo)
	router := setupSLARouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/sla", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestSLAHandler_Upsert_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &slaMockRepo{policy: &domainSLA.Policy{Priority: "high"}}
	h := newSLAHandler(repo)
	router := setupSLARouter(h, tenant)

	body := `{"response_time_minutes":60,"resolution_time_minutes":480,"escalation_minutes":720}`
	rec := doRequest(router, http.MethodPut, "/api/sla/high", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestSLAHandler_Upsert_InvalidJSON(t *testing.T) {
	tenant := uuid.New()
	h := newSLAHandler(&slaMockRepo{})
	router := setupSLARouter(h, tenant)

	rec := doRequest(router, http.MethodPut, "/api/sla/high", "{bad", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestSLAHandler_Upsert_ValidationError(t *testing.T) {
	tenant := uuid.New()
	h := newSLAHandler(&slaMockRepo{})
	router := setupSLARouter(h, tenant)

	body := `{"response_time_minutes":0,"resolution_time_minutes":480,"escalation_minutes":720}`
	rec := doRequest(router, http.MethodPut, "/api/sla/high", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestSLAHandler_Upsert_InvalidPriority(t *testing.T) {
	tenant := uuid.New()
	h := newSLAHandler(&slaMockRepo{})
	router := setupSLARouter(h, tenant)

	body := `{"response_time_minutes":60,"resolution_time_minutes":480,"escalation_minutes":720}`
	rec := doRequest(router, http.MethodPut, "/api/sla/urgent", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestSLAHandler_Upsert_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &slaMockRepo{upsertErr: errors.New("db")}
	h := newSLAHandler(repo)
	router := setupSLARouter(h, tenant)

	body := `{"response_time_minutes":60,"resolution_time_minutes":480,"escalation_minutes":720}`
	rec := doRequest(router, http.MethodPut, "/api/sla/high", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

// ============= AUTOMATION HANDLER TESTS =============

type autoMockRepo struct {
	rules     []domainAuto.Rule
	rule      *domainAuto.Rule
	err       error
	createErr error
	updateErr error
	deleteErr error
}

func (m *autoMockRepo) List(ctx context.Context, tenantID uuid.UUID) ([]domainAuto.Rule, error) {
	return m.rules, m.err
}
func (m *autoMockRepo) Create(ctx context.Context, r *domainAuto.Rule) error {
	return m.createErr
}
func (m *autoMockRepo) Update(ctx context.Context, r *domainAuto.Rule) error {
	return m.updateErr
}
func (m *autoMockRepo) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
	return m.deleteErr
}

func newAutoHandler(repo appAuto.Repository) *AutomationHandler {
	return NewAutomationHandler(appAuto.NewService(repo))
}

func setupAutoRouter(h *AutomationHandler, tenantID uuid.UUID) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Tenant(&mockTenantResolver{id: tenantID}))
	r.Get("/api/automation/rules", h.List)
	r.Post("/api/automation/rules", h.Create)
	r.Put("/api/automation/rules/{id}", h.Update)
	r.Delete("/api/automation/rules/{id}", h.Delete)
	return r
}

func TestAutomationHandler_List_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &autoMockRepo{rules: []domainAuto.Rule{{ID: uuid.New(), Name: "r"}}}
	h := newAutoHandler(repo)
	router := setupAutoRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/automation/rules", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_List_Error(t *testing.T) {
	tenant := uuid.New()
	repo := &autoMockRepo{err: errors.New("db")}
	h := newAutoHandler(repo)
	router := setupAutoRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/automation/rules", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Create_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &autoMockRepo{rule: &domainAuto.Rule{ID: uuid.New(), Name: "r"}}
	h := newAutoHandler(repo)
	router := setupAutoRouter(h, tenant)

	body := `{"name":"Auto-assign","trigger_type":"ticket_created","is_active":true}`
	rec := doRequest(router, http.MethodPost, "/api/automation/rules", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestAutomationHandler_Create_InvalidJSON(t *testing.T) {
	tenant := uuid.New()
	h := newAutoHandler(&autoMockRepo{})
	router := setupAutoRouter(h, tenant)

	rec := doRequest(router, http.MethodPost, "/api/automation/rules", "{bad", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Create_ValidationError(t *testing.T) {
	tenant := uuid.New()
	h := newAutoHandler(&autoMockRepo{})
	router := setupAutoRouter(h, tenant)

	body := `{"name":"a","trigger_type":"ticket_created","is_active":true}`
	rec := doRequest(router, http.MethodPost, "/api/automation/rules", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Create_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &autoMockRepo{createErr: errors.New("db")}
	h := newAutoHandler(repo)
	router := setupAutoRouter(h, tenant)

	body := `{"name":"Auto-assign","trigger_type":"ticket_created","is_active":true}`
	rec := doRequest(router, http.MethodPost, "/api/automation/rules", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Update_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &autoMockRepo{rule: &domainAuto.Rule{ID: uuid.New(), Name: "r"}}
	h := newAutoHandler(repo)
	router := setupAutoRouter(h, tenant)

	body := `{"name":"Updated rule","trigger_type":"status_changed","is_active":false}`
	rec := doRequest(router, http.MethodPut, "/api/automation/rules/"+uuid.New().String(), body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestAutomationHandler_Update_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newAutoHandler(&autoMockRepo{})
	router := setupAutoRouter(h, tenant)

	body := `{"name":"Updated","trigger_type":"status_changed","is_active":false}`
	rec := doRequest(router, http.MethodPut, "/api/automation/rules/bad", body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Update_InvalidJSON(t *testing.T) {
	tenant := uuid.New()
	h := newAutoHandler(&autoMockRepo{})
	router := setupAutoRouter(h, tenant)

	rec := doRequest(router, http.MethodPut, "/api/automation/rules/"+uuid.New().String(), "{bad", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Update_ValidationError(t *testing.T) {
	tenant := uuid.New()
	h := newAutoHandler(&autoMockRepo{})
	router := setupAutoRouter(h, tenant)

	body := `{"name":"a","trigger_type":"status_changed","is_active":false}`
	rec := doRequest(router, http.MethodPut, "/api/automation/rules/"+uuid.New().String(), body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Update_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &autoMockRepo{updateErr: domainAuto.ErrNotFound}
	h := newAutoHandler(repo)
	router := setupAutoRouter(h, tenant)

	body := `{"name":"Updated rule","trigger_type":"status_changed","is_active":false}`
	rec := doRequest(router, http.MethodPut, "/api/automation/rules/"+uuid.New().String(), body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Update_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &autoMockRepo{updateErr: errors.New("db")}
	h := newAutoHandler(repo)
	router := setupAutoRouter(h, tenant)

	body := `{"name":"Updated rule","trigger_type":"status_changed","is_active":false}`
	rec := doRequest(router, http.MethodPut, "/api/automation/rules/"+uuid.New().String(), body, &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Delete_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &autoMockRepo{}
	h := newAutoHandler(repo)
	router := setupAutoRouter(h, tenant)

	rec := doRequest(router, http.MethodDelete, "/api/automation/rules/"+uuid.New().String(), "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Delete_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newAutoHandler(&autoMockRepo{})
	router := setupAutoRouter(h, tenant)

	rec := doRequest(router, http.MethodDelete, "/api/automation/rules/bad", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Delete_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &autoMockRepo{deleteErr: domainAuto.ErrNotFound}
	h := newAutoHandler(repo)
	router := setupAutoRouter(h, tenant)

	rec := doRequest(router, http.MethodDelete, "/api/automation/rules/"+uuid.New().String(), "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAutomationHandler_Delete_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &autoMockRepo{deleteErr: errors.New("db")}
	h := newAutoHandler(repo)
	router := setupAutoRouter(h, tenant)

	rec := doRequest(router, http.MethodDelete, "/api/automation/rules/"+uuid.New().String(), "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

// ============= ANALYTICS HANDLER TESTS =============

type anMockRepo struct {
	summary    appAnalytics.Summary
	summaryErr error
	volume     []appAnalytics.VolumePoint
	volumeErr  error
	dist       []appAnalytics.StatusDistribution
	distErr    error
	resolution []appAnalytics.ResolutionMetrics
	resErr     error
}

func (m *anMockRepo) GetSummary(ctx context.Context, tenantID uuid.UUID) (appAnalytics.Summary, error) {
	return m.summary, m.summaryErr
}
func (m *anMockRepo) GetVolume(ctx context.Context, tenantID uuid.UUID) ([]appAnalytics.VolumePoint, error) {
	return m.volume, m.volumeErr
}
func (m *anMockRepo) GetStatusDistribution(ctx context.Context, tenantID uuid.UUID) ([]appAnalytics.StatusDistribution, error) {
	return m.dist, m.distErr
}
func (m *anMockRepo) GetResolutionTime(ctx context.Context, tenantID uuid.UUID) ([]appAnalytics.ResolutionMetrics, error) {
	return m.resolution, m.resErr
}

func newAnalyticsHandler(repo appAnalytics.Repository) *AnalyticsHandler {
	return NewAnalyticsHandler(appAnalytics.NewService(repo))
}

func setupAnalyticsRouter(h *AnalyticsHandler, tenantID uuid.UUID) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Tenant(&mockTenantResolver{id: tenantID}))
	r.Get("/api/analytics/summary", h.GetSummary)
	r.Get("/api/analytics/volume", h.GetVolume)
	r.Get("/api/analytics/status", h.GetStatusDistribution)
	r.Get("/api/analytics/resolution", h.GetResolutionTime)
	r.Get("/api/analytics/charts", h.GetCharts)
	return r
}

func TestAnalyticsHandler_GetSummary_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &anMockRepo{summary: appAnalytics.Summary{TotalTickets: 10}}
	h := newAnalyticsHandler(repo)
	router := setupAnalyticsRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/analytics/summary", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestAnalyticsHandler_GetSummary_Error(t *testing.T) {
	tenant := uuid.New()
	repo := &anMockRepo{summaryErr: errors.New("db")}
	h := newAnalyticsHandler(repo)
	router := setupAnalyticsRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/analytics/summary", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAnalyticsHandler_GetVolume_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &anMockRepo{volume: []appAnalytics.VolumePoint{{Date: "2025-01-01", Count: 5}}}
	h := newAnalyticsHandler(repo)
	router := setupAnalyticsRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/analytics/volume", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestAnalyticsHandler_GetVolume_Error(t *testing.T) {
	tenant := uuid.New()
	repo := &anMockRepo{volumeErr: errors.New("db")}
	h := newAnalyticsHandler(repo)
	router := setupAnalyticsRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/analytics/volume", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAnalyticsHandler_GetStatusDistribution_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &anMockRepo{dist: []appAnalytics.StatusDistribution{{Status: "open", Count: 3}}}
	h := newAnalyticsHandler(repo)
	router := setupAnalyticsRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/analytics/status", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestAnalyticsHandler_GetStatusDistribution_Error(t *testing.T) {
	tenant := uuid.New()
	repo := &anMockRepo{distErr: errors.New("db")}
	h := newAnalyticsHandler(repo)
	router := setupAnalyticsRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/analytics/status", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAnalyticsHandler_GetResolutionTime_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &anMockRepo{resolution: []appAnalytics.ResolutionMetrics{{Priority: "high", AvgHours: 4.5}}}
	h := newAnalyticsHandler(repo)
	router := setupAnalyticsRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/analytics/resolution", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestAnalyticsHandler_GetResolutionTime_Error(t *testing.T) {
	tenant := uuid.New()
	repo := &anMockRepo{resErr: errors.New("db")}
	h := newAnalyticsHandler(repo)
	router := setupAnalyticsRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/analytics/resolution", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAnalyticsHandler_GetCharts_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &anMockRepo{
		summary: appAnalytics.Summary{TotalTickets: 10, SLABreached: 2},
		volume:  []appAnalytics.VolumePoint{{Date: "2025-01-01", Count: 5}},
		dist:    []appAnalytics.StatusDistribution{{Status: "open", Count: 3}},
	}
	h := newAnalyticsHandler(repo)
	router := setupAnalyticsRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/analytics/charts", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestAnalyticsHandler_GetCharts_Error(t *testing.T) {
	tenant := uuid.New()
	repo := &anMockRepo{volumeErr: errors.New("db")}
	h := newAnalyticsHandler(repo)
	router := setupAnalyticsRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/analytics/charts", "", &appAuth.Claims{TenantID: tenant, Role: "admin"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

// Ensure strings import is used.
var _ = strings.NewReader
