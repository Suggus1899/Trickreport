package handler

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	appTicket "github.com/trickreport/backend/internal/application/ticket"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
)

// --- Mocks ---

type attMockRepo struct {
	attachment  *domainTicket.Attachment
	data        []byte
	attachments []domainTicket.Attachment
	err         error
	createErr   error
}

func (m *attMockRepo) Create(ctx context.Context, input appTicket.AttachmentUploadInput) (*domainTicket.Attachment, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return m.attachment, nil
}
func (m *attMockRepo) List(ctx context.Context, ticketID uuid.UUID) ([]domainTicket.Attachment, error) {
	return m.attachments, m.err
}
func (m *attMockRepo) GetByID(ctx context.Context, id, ticketID uuid.UUID) (*domainTicket.Attachment, []byte, error) {
	return m.attachment, m.data, m.err
}

func newAttachmentHandler(attRepo appTicket.AttachmentRepository, ticketRepo appTicket.Repository) *AttachmentHandler {
	svc := appTicket.NewAttachmentService(attRepo, ticketRepo)
	return NewAttachmentHandler(svc)
}

func setupAttachmentRouter(h *AttachmentHandler, tenantID uuid.UUID) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Tenant(&mockTenantResolver{id: tenantID}))
	r.Post("/api/tickets/{id}/attachments", h.Upload)
	r.Get("/api/tickets/{id}/attachments", h.List)
	r.Get("/api/tickets/{id}/attachments/{aid}", h.Download)
	return r
}

func multipartBody(t *testing.T, filename, contentType string) (*bytes.Buffer, string) {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	// Create a custom part with the specified content type so the handler
	// receives a whitelisted MIME type instead of the default octet-stream.
	h := make(map[string][]string)
	h["Content-Disposition"] = []string{`form-data; name="file"; filename="` + filename + `"`}
	h["Content-Type"] = []string{contentType}
	fw, err := w.CreatePart(h)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	fw.Write([]byte("file-content"))
	w.Close()
	return &buf, w.FormDataContentType()
}

// --- Upload tests ---

func TestAttachmentHandler_Upload_Success(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &tMockRepo{ticket: tk}
	attRepo := &attMockRepo{attachment: &domainTicket.Attachment{ID: uuid.New(), Filename: "test.txt"}}
	h := newAttachmentHandler(attRepo, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	buf, ct := multipartBody(t, "test.txt", "text/plain")
	req := httptest.NewRequest(http.MethodPost, "/api/tickets/"+tk.ID.String()+"/attachments", buf)
	req.Header.Set("X-Tenant-ID", "test-tenant")
	req.Header.Set("Content-Type", ct)
	claims := &appAuth.Claims{UserID: creator, TenantID: tenant, Role: "end_user"}
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), claims))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestAttachmentHandler_Upload_InvalidTicketID(t *testing.T) {
	tenant := uuid.New()
	h := newAttachmentHandler(&attMockRepo{}, &tMockRepo{})
	router := setupAttachmentRouter(h, tenant)

	buf, ct := multipartBody(t, "test.txt", "text/plain")
	req := httptest.NewRequest(http.MethodPost, "/api/tickets/bad/attachments", buf)
	req.Header.Set("X-Tenant-ID", "test-tenant")
	req.Header.Set("Content-Type", ct)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &appAuth.Claims{TenantID: tenant, Role: "agent"}))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_Upload_InvalidMultipart(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &tMockRepo{ticket: tk}
	h := newAttachmentHandler(&attMockRepo{}, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	req := httptest.NewRequest(http.MethodPost, "/api/tickets/"+tk.ID.String()+"/attachments", bytes.NewReader([]byte("not multipart")))
	req.Header.Set("X-Tenant-ID", "test-tenant")
	req.Header.Set("Content-Type", "text/plain")
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &appAuth.Claims{UserID: creator, TenantID: tenant, Role: "end_user"}))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_Upload_MissingFileField(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &tMockRepo{ticket: tk}
	h := newAttachmentHandler(&attMockRepo{}, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, _ := w.CreateFormFile("wrong_field", "test.txt")
	fw.Write([]byte("x"))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/tickets/"+tk.ID.String()+"/attachments", &buf)
	req.Header.Set("X-Tenant-ID", "test-tenant")
	req.Header.Set("Content-Type", w.FormDataContentType())
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &appAuth.Claims{UserID: creator, TenantID: tenant, Role: "end_user"}))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_Upload_Forbidden(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &tMockRepo{ticket: tk}
	h := newAttachmentHandler(&attMockRepo{}, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	buf, ct := multipartBody(t, "test.txt", "text/plain")
	req := httptest.NewRequest(http.MethodPost, "/api/tickets/"+tk.ID.String()+"/attachments", buf)
	req.Header.Set("X-Tenant-ID", "test-tenant")
	req.Header.Set("Content-Type", ct)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "end_user"}))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_Upload_TicketNotFound(t *testing.T) {
	tenant := uuid.New()
	ticketRepo := &tMockRepo{err: domainTicket.ErrNotFound}
	h := newAttachmentHandler(&attMockRepo{}, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	buf, ct := multipartBody(t, "test.txt", "text/plain")
	req := httptest.NewRequest(http.MethodPost, "/api/tickets/"+uuid.New().String()+"/attachments", buf)
	req.Header.Set("X-Tenant-ID", "test-tenant")
	req.Header.Set("Content-Type", ct)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "agent"}))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_Upload_CreateError(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &tMockRepo{ticket: tk}
	attRepo := &attMockRepo{createErr: errors.New("storage")}
	h := newAttachmentHandler(attRepo, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	buf, ct := multipartBody(t, "test.txt", "text/plain")
	req := httptest.NewRequest(http.MethodPost, "/api/tickets/"+tk.ID.String()+"/attachments", buf)
	req.Header.Set("X-Tenant-ID", "test-tenant")
	req.Header.Set("Content-Type", ct)
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &appAuth.Claims{UserID: creator, TenantID: tenant, Role: "end_user"}))

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

// --- List tests ---

func TestAttachmentHandler_List_Success(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &tMockRepo{ticket: tk}
	attRepo := &attMockRepo{attachments: []domainTicket.Attachment{{ID: uuid.New(), Filename: "a.txt"}}}
	h := newAttachmentHandler(attRepo, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	claims := &appAuth.Claims{UserID: creator, TenantID: tenant, Role: "end_user"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/attachments", "", claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestAttachmentHandler_List_InvalidTicketID(t *testing.T) {
	tenant := uuid.New()
	h := newAttachmentHandler(&attMockRepo{}, &tMockRepo{})
	router := setupAttachmentRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/tickets/bad/attachments", "", &appAuth.Claims{TenantID: tenant, Role: "agent"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_List_NotFound(t *testing.T) {
	tenant := uuid.New()
	ticketRepo := &tMockRepo{err: domainTicket.ErrNotFound}
	h := newAttachmentHandler(&attMockRepo{}, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/tickets/"+uuid.New().String()+"/attachments", "", &appAuth.Claims{TenantID: tenant, Role: "agent"})
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_List_Forbidden(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &tMockRepo{ticket: tk}
	h := newAttachmentHandler(&attMockRepo{}, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/attachments", "", &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "end_user"})
	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_List_RepoError(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &tMockRepo{ticket: tk}
	attRepo := &attMockRepo{err: errors.New("db")}
	h := newAttachmentHandler(attRepo, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/attachments", "", &appAuth.Claims{UserID: creator, TenantID: tenant, Role: "end_user"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}

// --- Download tests ---

func TestAttachmentHandler_Download_Success(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	attID := uuid.New()
	ticketRepo := &tMockRepo{ticket: tk}
	attRepo := &attMockRepo{
		attachment: &domainTicket.Attachment{ID: attID, Filename: "file.pdf", ContentType: "application/pdf"},
		data:       []byte("pdf-data"),
	}
	h := newAttachmentHandler(attRepo, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	claims := &appAuth.Claims{UserID: creator, TenantID: tenant, Role: "end_user"}
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/attachments/"+attID.String(), "", claims)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if rec.Body.String() != "pdf-data" {
		t.Errorf("body = %q", rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Errorf("content-type = %q", ct)
	}
}

func TestAttachmentHandler_Download_InvalidTicketID(t *testing.T) {
	tenant := uuid.New()
	h := newAttachmentHandler(&attMockRepo{}, &tMockRepo{})
	router := setupAttachmentRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/tickets/bad/attachments/"+uuid.New().String(), "", &appAuth.Claims{TenantID: tenant, Role: "agent"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_Download_InvalidAttachmentID(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &tMockRepo{ticket: tk}
	h := newAttachmentHandler(&attMockRepo{}, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/attachments/bad", "", &appAuth.Claims{UserID: creator, TenantID: tenant, Role: "end_user"})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_Download_NotFound(t *testing.T) {
	tenant := uuid.New()
	ticketRepo := &tMockRepo{err: domainTicket.ErrNotFound}
	h := newAttachmentHandler(&attMockRepo{}, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/tickets/"+uuid.New().String()+"/attachments/"+uuid.New().String(), "", &appAuth.Claims{TenantID: tenant, Role: "agent"})
	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_Download_Forbidden(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &tMockRepo{ticket: tk}
	h := newAttachmentHandler(&attMockRepo{}, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/attachments/"+uuid.New().String(), "", &appAuth.Claims{UserID: uuid.New(), TenantID: tenant, Role: "end_user"})
	if rec.Code != http.StatusForbidden {
		t.Errorf("status = %d", rec.Code)
	}
}

func TestAttachmentHandler_Download_RepoError(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &tMockRepo{ticket: tk}
	attRepo := &attMockRepo{err: errors.New("db")}
	h := newAttachmentHandler(attRepo, ticketRepo)
	router := setupAttachmentRouter(h, tenant)

	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/attachments/"+uuid.New().String(), "", &appAuth.Claims{UserID: creator, TenantID: tenant, Role: "end_user"})
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d", rec.Code)
	}
}
