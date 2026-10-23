package handler

import (
	"bytes"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"github.com/google/uuid"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	domain "github.com/trickreport/backend/internal/domain/ticket"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
)

func TestTicketHandler_ConflictAndInvalidAssignee(t *testing.T) {
	actor := uuid.New()
	tk := ownedTicket(actor)
	repo := &tMockRepo{ticket: tk, updateErr: domain.ErrConcurrentModification, assignErr: domain.ErrValidation}
	router := setupTicketRouter(newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{}), tk.TenantID)
	claims := &appAuth.Claims{UserID: actor, TenantID: tk.TenantID, Role: "admin"}
	if rec := doRequest(router, "PUT", "/api/tickets/"+tk.ID.String()+"/status", `{"status":"in_progress","version":0}`, claims); rec.Code != 409 {
		t.Fatalf("stale version status=%d body=%s", rec.Code, rec.Body)
	}
	if rec := doRequest(router, "POST", "/api/tickets/"+tk.ID.String()+"/assign", `{"assigned_to":"`+uuid.New().String()+`"}`, claims); rec.Code != 400 {
		t.Fatalf("ineligible assignee status=%d", rec.Code)
	}
}

func TestTicketHandler_InvalidIdempotencyKey(t *testing.T) {
	repo := &tMockRepo{}
	tenant := uuid.New()
	router := setupTicketRouter(newTicketHandler(repo, &tMockCommentRepo{}, &tMockHistoryRepo{}), tenant)
	req := httptest.NewRequest("POST", "/api/tickets", strings.NewReader(`{"title":"valid title","description":"valid description","priority":"high","category":"general"}`))
	req.Header.Set("X-Tenant-ID", tenant.String())
	req.Header.Set("Idempotency-Key", "not-a-uuid")
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &appAuth.Claims{TenantID: tenant, UserID: uuid.New(), Role: "admin"}))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != 400 || repo.createCalls != 0 {
		t.Fatalf("status=%d writes=%d", rec.Code, repo.createCalls)
	}
}

func TestAttachmentHandler_MultipartFileLimitsAndSniffing(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	for _, test := range []struct {
		name, mime string
		size       int
		want       int
	}{{"exact file limit", "text/plain", MaxUploadSize, 201}, {"file over limit", "text/plain", MaxUploadSize + 1, 413}, {"forged MIME", "image/png", 4, 415}, {"empty", "text/plain", 0, 415}} {
		t.Run(test.name, func(t *testing.T) {
			repo := &attMockRepo{attachment: &domain.Attachment{ID: uuid.New()}}
			h := newAttachmentHandler(repo, &tMockRepo{ticket: tk})
			router := setupAttachmentRouter(h, tk.TenantID)
			var body bytes.Buffer
			writer := multipart.NewWriter(&body)
			header := textproto.MIMEHeader{}
			header.Set("Content-Disposition", `form-data; name="file"; filename="test.txt"`)
			header.Set("Content-Type", test.mime)
			part, err := writer.CreatePart(header)
			if err != nil {
				t.Fatal(err)
			}
			_, _ = part.Write(bytes.Repeat([]byte("a"), test.size))
			_ = writer.Close()
			req := httptest.NewRequest("POST", "/api/tickets/"+tk.ID.String()+"/attachments", &body)
			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("X-Tenant-ID", tk.TenantID.String())
			req = req.WithContext(middleware.ContextWithClaims(req.Context(), &appAuth.Claims{TenantID: tk.TenantID, UserID: creator, Role: "end_user"}))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != test.want {
				t.Fatalf("status=%d want=%d body=%s", rec.Code, test.want, rec.Body)
			}
		})
	}
}

func TestAttachmentHandler_DownloadSafeDisposition(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	id := uuid.New()
	repo := &attMockRepo{attachment: &domain.Attachment{ID: id, Filename: "../../unsafe\r\nname.pdf", ContentType: "application/pdf"}, data: []byte("%PDF")}
	router := setupAttachmentRouter(newAttachmentHandler(repo, &tMockRepo{ticket: tk}), tk.TenantID)
	rec := doRequest(router, http.MethodGet, "/api/tickets/"+tk.ID.String()+"/attachments/"+id.String(), "", &appAuth.Claims{TenantID: tk.TenantID, UserID: creator, Role: "end_user"})
	_, params, err := mime.ParseMediaType(rec.Header().Get("Content-Disposition"))
	if err != nil {
		t.Fatal(err)
	}
	if rec.Code != 200 || params["filename"] != "unsafe__name.pdf" || rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("status=%d filename=%q headers=%v", rec.Code, params["filename"], rec.Header())
	}
}
