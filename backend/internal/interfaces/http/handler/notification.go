package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appTicket "github.com/trickreport/backend/internal/application/ticket"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
)

// Compile-time check that appTicket is used (for the service type).
var _ = (*appTicket.NotificationService)(nil)

// NotificationHandler handles notification endpoints.
type NotificationHandler struct {
	svc *appTicket.NotificationService
}

// NotificationDTO is the JSON representation of a notification.
type NotificationDTO struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Type      string `json:"type"`
	RefID     string `json:"ref_id,omitempty"`
	RefType   string `json:"ref_type,omitempty"`
	Read      bool   `json:"read"`
	CreatedAt string `json:"created_at"`
}

// NewNotificationHandler creates a new NotificationHandler.
func NewNotificationHandler(svc *appTicket.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

// List returns notifications for the authenticated user.
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	notifications, err := h.svc.List(r.Context(), tenantID, claims.UserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list notifications")
		return
	}
	dtos := make([]NotificationDTO, 0, len(notifications))
	for _, n := range notifications {
		dto := NotificationDTO{
			ID:        n.ID.String(),
			Title:     n.Title,
			Body:      n.Body,
			Type:      n.Type,
			Read:      n.Read,
			CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			RefType:   n.RefType,
		}
		if n.RefID != nil {
			dto.RefID = n.RefID.String()
		}
		dtos = append(dtos, dto)
	}
	response.JSON(w, http.StatusOK, dtos)
}

// MarkRead marks a single notification as read.
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	notifID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid notification id")
		return
	}
	if err := h.svc.MarkRead(r.Context(), tenantID, claims.UserID, notifID); err != nil {
		response.Error(w, http.StatusNotFound, "notification not found")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// MarkAllRead marks all notifications as read for the authenticated user.
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.svc.MarkAllRead(r.Context(), tenantID, claims.UserID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to mark all read")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// UnreadCount returns the unread notification count.
func (h *NotificationHandler) UnreadCount(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	count, err := h.svc.UnreadCount(r.Context(), tenantID, claims.UserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get unread count")
		return
	}
	response.JSON(w, http.StatusOK, map[string]int{"unread": count})
}
