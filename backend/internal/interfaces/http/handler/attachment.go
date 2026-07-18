package handler

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appTicket "github.com/trickreport/backend/internal/application/ticket"
	"github.com/trickreport/backend/internal/domain/ticket"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
)

// AttachmentHandler handles HTTP requests for ticket file attachments.
type AttachmentHandler struct {
	svc *appTicket.AttachmentService
}

// NewAttachmentHandler creates a new AttachmentHandler.
func NewAttachmentHandler(svc *appTicket.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{svc: svc}
}

// AttachmentDTO is the HTTP representation of an attachment (no file data).
type AttachmentDTO struct {
	ID          uuid.UUID `json:"id"`
	TicketID    uuid.UUID `json:"ticket_id"`
	UserID      uuid.UUID `json:"user_id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	FileSize    int64     `json:"file_size"`
	CreatedAt   string    `json:"created_at"`
	UserName    string    `json:"user_name,omitempty"`
}

// MaxUploadSize is the maximum allowed attachment size (10 MB).
const MaxUploadSize = 10 << 20

func toAttachmentDTO(a *ticket.Attachment) AttachmentDTO {
	created := ""
	if !a.CreatedAt.IsZero() {
		created = a.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")
	}
	return AttachmentDTO{
		ID:          a.ID,
		TicketID:    a.TicketID,
		UserID:      a.UserID,
		Filename:    a.Filename,
		ContentType: a.ContentType,
		FileSize:    a.FileSize,
		CreatedAt:   created,
		UserName:    a.UserName,
	}
}

// Upload handles multipart file uploads for a ticket.
func (h *AttachmentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())

	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid ticket id")
		return
	}

	// Limit request body size.
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		response.Error(w, http.StatusBadRequest, "file too large or invalid multipart form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.Error(w, http.StatusBadRequest, "missing 'file' field in multipart form")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to read uploaded file")
		return
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	att, err := h.svc.Upload(r.Context(), appTicket.UploadInput{
		TenantID:    tenantID,
		TicketID:    ticketID,
		UserID:      claims.UserID,
		Role:        claims.Role,
		Filename:    header.Filename,
		ContentType: contentType,
		FileData:    data,
	})
	if err != nil {
		switch {
		case errors.Is(err, ticket.ErrNotFound):
			response.Error(w, http.StatusNotFound, "ticket not found")
		case errors.Is(err, ticket.ErrForbidden):
			response.Error(w, http.StatusForbidden, "not allowed to access this ticket")
		case errors.Is(err, appTicket.ErrFileTooLarge):
			response.Error(w, http.StatusRequestEntityTooLarge, err.Error())
		case errors.Is(err, appTicket.ErrInvalidContent):
			response.Error(w, http.StatusUnsupportedMediaType, err.Error())
		case errors.Is(err, appTicket.ErrEmptyFilename):
			response.Error(w, http.StatusBadRequest, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, "failed to upload attachment")
		}
		return
	}

	response.JSON(w, http.StatusCreated, toAttachmentDTO(att))
}

// List returns metadata for all attachments on a ticket.
func (h *AttachmentHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())

	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid ticket id")
		return
	}

	attachments, err := h.svc.List(r.Context(), appTicket.ListInput{
		TenantID: tenantID,
		TicketID: ticketID,
		UserID:   claims.UserID,
		Role:     claims.Role,
	})
	if err != nil {
		switch {
		case errors.Is(err, ticket.ErrNotFound):
			response.Error(w, http.StatusNotFound, "ticket not found")
		case errors.Is(err, ticket.ErrForbidden):
			response.Error(w, http.StatusForbidden, "not allowed to access this ticket")
		default:
			response.Error(w, http.StatusInternalServerError, "failed to list attachments")
		}
		return
	}

	dtos := make([]AttachmentDTO, len(attachments))
	for i, a := range attachments {
		dtos[i] = toAttachmentDTO(&a)
	}
	response.JSON(w, http.StatusOK, dtos)
}

// Download returns the raw file content for an attachment.
func (h *AttachmentHandler) Download(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())

	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid ticket id")
		return
	}

	attachmentID, err := uuid.Parse(chi.URLParam(r, "aid"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid attachment id")
		return
	}

	att, data, err := h.svc.Download(r.Context(), appTicket.DownloadInput{
		TenantID:     tenantID,
		TicketID:     ticketID,
		AttachmentID: attachmentID,
		UserID:       claims.UserID,
		Role:         claims.Role,
	})
	if err != nil {
		switch {
		case errors.Is(err, ticket.ErrNotFound):
			response.Error(w, http.StatusNotFound, "ticket or attachment not found")
		case errors.Is(err, ticket.ErrForbidden):
			response.Error(w, http.StatusForbidden, "not allowed to access this ticket")
		default:
			response.Error(w, http.StatusInternalServerError, "failed to download attachment")
		}
		return
	}

	w.Header().Set("Content-Type", att.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+att.Filename+"\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
