package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	appTicket "github.com/trickreport/backend/internal/application/ticket"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
)

// --- Mocks for ticket service dependencies ---

type tMockRepo struct {
	tickets   []domainTicket.Ticket
	ticket    *domainTicket.Ticket
	err       error
	updateErr error
	assignErr error
	gotFilter appTicket.Filter
	createCalls int
}

func (m *tMockRepo) List(ctx context.Context, tenantID uuid.UUID, filter appTicket.Filter, role string, userID uuid.UUID) ([]domainTicket.Ticket, error) {
	m.gotFilter = filter
	return m.tickets, m.err
}
func (m *tMockRepo) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*domainTicket.Ticket, error) {
	return m.ticket, m.err
}
func (m *tMockRepo) Create(ctx context.Context, t *domainTicket.Ticket) error {
	m.createCalls++
	return m.err
}
func (m *tMockRepo) UpdateStatus(ctx context.Context, id, tenantID uuid.UUID, status domainTicket.Status, userID uuid.UUID) (*domainTicket.Ticket, error) {
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return m.ticket, nil
}
func (m *tMockRepo) Assign(ctx context.Context, id, tenantID, assignedTo, userID uuid.UUID) error {
	return m.assignErr
}
func (m *tMockRepo) GetCreator(ctx context.Context, id, tenantID uuid.UUID) (uuid.UUID, error) {
	return uuid.Nil, nil
}

type tMockCommentRepo struct {
	comments  []domainTicket.Comment
	err       error
	createErr error
}

func (m *tMockCommentRepo) List(ctx context.Context, ticketID uuid.UUID, role string) ([]domainTicket.Comment, error) {
	return m.comments, m.err
}
func (m *tMockCommentRepo) Create(ctx context.Context, c *domainTicket.Comment) error {
	return m.createErr
}

type tMockHistoryRepo struct {
	entries []domainTicket.HistoryEntry
	err     error
}

func (m *tMockHistoryRepo) List(ctx context.Context, ticketID uuid.UUID) ([]domainTicket.HistoryEntry, error) {
	return m.entries, m.err
}

type tMockHub struct{ broadcasts int }

func (m *tMockHub) BroadcastEvent(tenantID uuid.UUID, eventType string, data any) {
	m.broadcasts++
}

// --- Helpers ---

func newTicketHandler(repo appTicket.Repository, comments appTicket.CommentRepository, history appTicket.HistoryRepository) *TicketHandler {
	svc := appTicket.NewService(repo, comments, history, &tMockHub{}, nil)
	return NewTicketHandler(svc)
}

// tenantCtxSetter is a minimal middleware that injects the tenant UUID into
// the request context using the same key the middleware package uses. Since
// the key is unexported, we rely on the public Tenant middleware with a mock
// resolver instead.
type mockTenantResolver struct {
	id  uuid.UUID
	err error
}

func (m *mockTenantResolver) Resolve(ctx context.Context, slugOrID string) (uuid.UUID, error) {
	return m.id, m.err
}

func setupTicketRouter(h *TicketHandler, tenantID uuid.UUID) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Tenant(&mockTenantResolver{id: tenantID}))
	r.Get("/api/tickets", h.List)
	r.Post("/api/tickets", h.Create)
	r.Get("/api/tickets/{id}", h.Get)
	r.Put("/api/tickets/{id}/status", h.UpdateStatus)
	r.Post("/api/tickets/{id}/assign", h.Assign)
	r.Get("/api/tickets/{id}/comments", h.ListComments)
	r.Post("/api/tickets/{id}/comments", h.AddComment)
	r.Get("/api/tickets/{id}/history", h.ListHistory)
	return r
}

func doRequest(router http.Handler, method, target string, body string, claims *appAuth.Claims) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("X-Tenant-ID", "test-tenant")
	if claims != nil {
		req = req.WithContext(middleware.ContextWithClaims(req.Context(), claims))
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func decodeJSON(t *testing.T, body string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.NewDecoder(strings.NewReader(body)).Decode(&m); err != nil {
		t.Fatalf("decode: %v (body=%s)", err, body)
	}
	return m
}

func ownedTicket(creator uuid.UUID) *domainTicket.Ticket {
	return &domainTicket.Ticket{
		ID:        uuid.New(),
		TenantID:  uuid.New(),
		Title:     "Broken login",
		Status:    domainTicket.StatusOpen,
		Priority:  domainTicket.PriorityMedium,
		Category:  "general",
		CreatedBy: creator,
	}
}

// --- List tests ---

func TestTicketHandler_List_Success(t *testing.T) {
	tenant := uuid.New()
	tickets := []domainTicket.Ticket{{ID: uuid.New(), Title: "t1"}, {ID: uuid.New(), Title: "t2"}}
	repo := &tMockRepo{tickets: tickets}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets?status=open&limit=10", "", claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var got []map[string]any
	if err := json.NewDecoder(strings.NewReader(rec.Body.String())).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("got %d tickets", len(got))
	}
	if repo.gotFilter.Status != "open" || repo.gotFilter.Limit != 10 {
		t.Errorf("filter = %+v", repo.gotFilter)
	}
}

func TestTicketHandler_List_Unauthorized(t *testing.T) {
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, uuid.New())

	// No X-Tenant-ID header and no subdomain (localhost has no dot) ->
	// tenant middleware returns 400.
	req := httptest.NewRequest(http.MethodGet, "/api/tickets", nil)
	req.Host = "localhost"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d (missing tenant)", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_List_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &tMockRepo{err: errors.New("db down")}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets", "", claims)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- Create tests ---

func TestTicketHandler_Create_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &tMockRepo{}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "end_user"}
	body := `{"title":"My ticket","description":"Something is broken","priority":"high","category":"bug"}`
	rec := doRequest(router, http.MethodPost, "/api/tickets", body, claims)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	resp := decodeJSON(t, rec.Body.String())
	if resp["title"] != "My ticket" {
		t.Errorf("title = %v", resp["title"])
	}
	if resp["status"] != "open" {
		t.Errorf("status = %v", resp["status"])
	}
	if repo.createCalls != 1 {
		t.Errorf("create calls = %d", repo.createCalls)
	}
}

func TestTicketHandler_Create_InvalidJSON(t *testing.T) {
	tenant := uuid.New()
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "end_user"}
	rec := doRequest(router, http.MethodPost, "/api/tickets", "{bad", claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_Create_ValidationError(t *testing.T) {
	tenant := uuid.New()
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "end_user"}
	// title too short (min=3)
	body := `{"title":"ab","description":"something broken here","priority":"high","category":"bug"}`
	rec := doRequest(router, http.MethodPost, "/api/tickets", body, claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_Create_RepoError(t *testing.T) {
	tenant := uuid.New()
	repo := &tMockRepo{err: errors.New("insert failed")}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "end_user"}
	body := `{"title":"My ticket","description":"Something is broken","priority":"high","category":"bug"}`
	rec := doRequest(router, http.MethodPost, "/api/tickets", body, claims)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- Get tests ---

func TestTicketHandler_Get_Success(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String(), "", claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	resp := decodeJSON(t, rec.Body.String())
	if resp["title"] != "Broken login" {
		t.Errorf("title = %v", resp["title"])
	}
}

func TestTicketHandler_Get_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/not-a-uuid", "", claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_Get_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &tMockRepo{err: domainTicket.ErrNotFound}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+uuid.New().String(), "", claims)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestTicketHandler_Get_Forbidden(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "end_user"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String(), "", claims)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

// --- UpdateStatus tests ---

func TestTicketHandler_UpdateStatus_Success(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"status":"in_progress"}`
	rec := doRequest(router, http.MethodPut, "/api/tickets/"+tk.ID.String()+"/status", body, claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	resp := decodeJSON(t, rec.Body.String())
	if resp["status"] != "in_progress" {
		t.Errorf("status = %v", resp["status"])
	}
}

func TestTicketHandler_UpdateStatus_InvalidJSON(t *testing.T) {
	tenant := uuid.New()
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodPut, "/api/tickets/"+uuid.New().String()+"/status", "{bad", claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_UpdateStatus_ValidationError(t *testing.T) {
	tenant := uuid.New()
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"status":"bogus"}`
	rec := doRequest(router, http.MethodPut, "/api/tickets/"+uuid.New().String()+"/status", body, claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_UpdateStatus_InvalidTransition(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	tk.Status = domainTicket.StatusClosed
	repo := &tMockRepo{ticket: tk}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"status":"open"}`
	rec := doRequest(router, http.MethodPut, "/api/tickets/"+tk.ID.String()+"/status", body, claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_UpdateStatus_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &tMockRepo{err: domainTicket.ErrNotFound}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"status":"in_progress"}`
	rec := doRequest(router, http.MethodPut, "/api/tickets/"+uuid.New().String()+"/status", body, claims)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestTicketHandler_UpdateStatus_Forbidden(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	// end_user trying to set status to in_progress (not closed) -> forbidden
	claims := &appAuth.Claims{UserID: creator, TenantID: tenant, Role: "end_user"}
	body := `{"status":"in_progress"}`
	rec := doRequest(router, http.MethodPut, "/api/tickets/"+tk.ID.String()+"/status", body, claims)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestTicketHandler_UpdateStatus_UpdateError(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk, updateErr: errors.New("db")}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"status":"in_progress"}`
	rec := doRequest(router, http.MethodPut, "/api/tickets/"+tk.ID.String()+"/status", body, claims)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- AddComment tests ---

func TestTicketHandler_AddComment_Success(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"content":"looks bad","is_internal":true}`
	rec := doRequest(router, http.MethodPost, "/api/tickets/"+tk.ID.String()+"/comments", body, claims)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	resp := decodeJSON(t, rec.Body.String())
	if resp["content"] != "looks bad" {
		t.Errorf("content = %v", resp["content"])
	}
}

func TestTicketHandler_AddComment_InvalidJSON(t *testing.T) {
	tenant := uuid.New()
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodPost, "/api/tickets/"+uuid.New().String()+"/comments", "{bad", claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_AddComment_ValidationError(t *testing.T) {
	tenant := uuid.New()
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	// content required
	body := `{"content":"","is_internal":false}`
	rec := doRequest(router, http.MethodPost, "/api/tickets/"+uuid.New().String()+"/comments", body, claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_AddComment_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &tMockRepo{err: domainTicket.ErrNotFound}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"content":"hi","is_internal":false}`
	rec := doRequest(router, http.MethodPost, "/api/tickets/"+uuid.New().String()+"/comments", body, claims)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestTicketHandler_AddComment_Forbidden(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "end_user"}
	body := `{"content":"hi","is_internal":false}`
	rec := doRequest(router, http.MethodPost, "/api/tickets/"+tk.ID.String()+"/comments", body, claims)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestTicketHandler_AddComment_CreateError(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	comments := &tMockCommentRepo{createErr: errors.New("insert")}
	h := newTicketHandler(repo, comments, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"content":"hi","is_internal":false}`
	rec := doRequest(router, http.MethodPost, "/api/tickets/"+tk.ID.String()+"/comments", body, claims)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- Assign tests ---

func TestTicketHandler_Assign_Success(t *testing.T) {
	tenant := uuid.New()
	assignee := uuid.New()
	repo := &tMockRepo{}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"assigned_to":"` + assignee.String() + `"}`
	rec := doRequest(router, http.MethodPost, "/api/tickets/"+uuid.New().String()+"/assign", body, claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	resp := decodeJSON(t, rec.Body.String())
	if resp["assigned_to"] != assignee.String() {
		t.Errorf("assigned_to = %v, want %v", resp["assigned_to"], assignee.String())
	}
}

func TestTicketHandler_Assign_Unassign(t *testing.T) {
	tenant := uuid.New()
	repo := &tMockRepo{}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{}`
	rec := doRequest(router, http.MethodPost, "/api/tickets/"+uuid.New().String()+"/assign", body, claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestTicketHandler_Assign_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"assigned_to":"` + uuid.New().String() + `"}`
	rec := doRequest(router, http.MethodPost, "/api/tickets/not-a-uuid/assign", body, claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_Assign_InvalidJSON(t *testing.T) {
	tenant := uuid.New()
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodPost, "/api/tickets/"+uuid.New().String()+"/assign", "{bad", claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_Assign_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &tMockRepo{assignErr: domainTicket.ErrNotFound}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"assigned_to":"` + uuid.New().String() + `"}`
	rec := doRequest(router, http.MethodPost, "/api/tickets/"+uuid.New().String()+"/assign", body, claims)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestTicketHandler_Assign_InternalError(t *testing.T) {
	tenant := uuid.New()
	repo := &tMockRepo{assignErr: errors.New("db down")}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	body := `{"assigned_to":"` + uuid.New().String() + `"}`
	rec := doRequest(router, http.MethodPost, "/api/tickets/"+uuid.New().String()+"/assign", body, claims)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- ListComments tests ---

func TestTicketHandler_ListComments_Success(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	comments := &tMockCommentRepo{
		comments: []domainTicket.Comment{
			{ID: uuid.New(), TicketID: tk.ID, UserID: uuid.New(), Content: "hello", IsInternal: false},
			{ID: uuid.New(), TicketID: tk.ID, UserID: uuid.New(), Content: "world", IsInternal: true},
		},
	}
	h := newTicketHandler(repo, comments, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/comments", "", claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var got []map[string]any
	if err := json.NewDecoder(strings.NewReader(rec.Body.String())).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("got %d comments", len(got))
	}
}

func TestTicketHandler_ListComments_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/not-a-uuid/comments", "", claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_ListComments_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &tMockRepo{err: domainTicket.ErrNotFound}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+uuid.New().String()+"/comments", "", claims)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestTicketHandler_ListComments_Forbidden(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "end_user"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/comments", "", claims)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestTicketHandler_ListComments_RepoError(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	comments := &tMockCommentRepo{err: errors.New("db down")}
	h := newTicketHandler(repo, comments, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/comments", "", claims)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

// --- ListHistory tests ---

func TestTicketHandler_ListHistory_Success(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	history := &tMockHistoryRepo{
		entries: []domainTicket.HistoryEntry{
			{ID: uuid.New(), TicketID: tk.ID, UserID: uuid.New(), Field: "status", OldValue: strPtr("open"), NewValue: strPtr("closed")},
		},
	}
	h := newTicketHandler(repo, &tMockCommentRepo{}, history)
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/history", "", claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	var got []map[string]any
	if err := json.NewDecoder(strings.NewReader(rec.Body.String())).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d entries", len(got))
	}
	if got[0]["field"] != "status" {
		t.Errorf("field = %v", got[0]["field"])
	}
}

func TestTicketHandler_ListHistory_InvalidID(t *testing.T) {
	tenant := uuid.New()
	h := newTicketHandler(&tMockRepo{}, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/not-a-uuid/history", "", claims)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_ListHistory_NotFound(t *testing.T) {
	tenant := uuid.New()
	repo := &tMockRepo{err: domainTicket.ErrNotFound}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+uuid.New().String()+"/history", "", claims)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestTicketHandler_ListHistory_Forbidden(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	h := newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{})
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "end_user"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/history", "", claims)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestTicketHandler_ListHistory_RepoError(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &tMockRepo{ticket: tk}
	history := &tMockHistoryRepo{err: errors.New("db down")}
	h := newTicketHandler(repo, &tMockCommentRepo{}, history)
	router := setupTicketRouter(h, tenant)

	claims := &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/history", "", claims)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func strPtr(s string) *string { return &s }
