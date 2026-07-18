package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appTicket "github.com/trickreport/backend/internal/application/ticket"
	"github.com/trickreport/backend/internal/domain/ticket"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/interfaces/http/validator"
)

// TicketHandler handles HTTP requests for tickets.
type TicketHandler struct {
	svc *appTicket.UserService
}

func NewTicketHandler(svc *appTicket.UserService) *TicketHandler {
	return &TicketHandler{svc: svc}
}

// DTOs — HTTP-specific, separate from domain entities.

type TicketDTO struct {
	ID           uuid.UUID  `json:"id"`
	TenantID     uuid.UUID  `json:"tenant_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	Priority     string     `json:"priority"`
	Category     string     `json:"category"`
	CreatedBy    uuid.UUID  `json:"created_by"`
	AssignedTo   *uuid.UUID `json:"assigned_to,omitempty"`
	SLADeadline  *time.Time `json:"sla_deadline,omitempty"`
	SLABreached  bool       `json:"sla_breached"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	CreatorName  string     `json:"creator_name,omitempty"`
	AssigneeName *string    `json:"assignee_name,omitempty"`
}

type CommentDTO struct {
	ID         uuid.UUID `json:"id"`
	TicketID   uuid.UUID `json:"ticket_id"`
	UserID     uuid.UUID `json:"user_id"`
	Content    string    `json:"content"`
	IsInternal bool      `json:"is_internal"`
	CreatedAt  time.Time `json:"created_at"`
	UserName   string    `json:"user_name,omitempty"`
}

type HistoryDTO struct {
	ID        uuid.UUID `json:"id"`
	TicketID  uuid.UUID `json:"ticket_id"`
	UserID    uuid.UUID `json:"user_id"`
	Field     string    `json:"field"`
	OldValue  *string   `json:"old_value,omitempty"`
	NewValue  *string   `json:"new_value,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"user_name,omitempty"`
}

type CreateTicketReq struct {
	Title       string `json:"title" validate:"required,min=3,max=255"`
	Description string `json:"description" validate:"required,min=10,max=5000"`
	Priority    string `json:"priority" validate:"required,oneof=low medium high critical"`
	Category    string `json:"category" validate:"required,max=100"`
}

type UpdateStatusReq struct {
	Status string `json:"status" validate:"required,oneof=open in_progress waiting_client resolved closed"`
}

type AssignReq struct {
	AssignedTo *uuid.UUID `json:"assigned_to"`
}

type CreateCommentReq struct {
	Content    string `json:"content" validate:"required,max=2000"`
	IsInternal bool   `json:"is_internal"`
}

func toTicketDTO(t *ticket.Ticket) TicketDTO {
	return TicketDTO{
		ID: t.ID, TenantID: t.TenantID, Title: t.Title, Description: t.Description,
		Status: string(t.Status), Priority: string(t.Priority), Category: t.Category,
		CreatedBy: t.CreatedBy, AssignedTo: t.AssignedTo, SLADeadline: t.SLADeadline,
		SLABreached: t.SLABreached, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
		CreatorName: t.CreatorName, AssigneeName: t.AssigneeName,
	}
}

func toCommentDTO(c *ticket.Comment) CommentDTO {
	return CommentDTO{
		ID: c.ID, TicketID: c.TicketID, UserID: c.UserID, Content: c.Content,
		IsInternal: c.IsInternal, CreatedAt: c.CreatedAt, UserName: c.UserName,
	}
}

func toHistoryDTO(h *ticket.HistoryEntry) HistoryDTO {
	return HistoryDTO{
		ID: h.ID, TicketID: h.TicketID, UserID: h.UserID, Field: h.Field,
		OldValue: h.OldValue, NewValue: h.NewValue, CreatedAt: h.CreatedAt,
		UserName: h.UserName,
	}
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	filter := appTicket.Filter{
		Status:     r.URL.Query().Get("status"),
		Priority:   r.URL.Query().Get("priority"),
		AssignedTo: r.URL.Query().Get("assigned_to"),
		Search:     r.URL.Query().Get("q"),
		Limit:      limit,
		Offset:     offset,
	}

	tickets, err := h.svc.List(r.Context(), tenantID, filter, claims.Role, claims.UserID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list tickets")
		return
	}

	dtos := make([]TicketDTO, len(tickets))
	for i, t := range tickets {
		dtos[i] = toTicketDTO(&t)
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *TicketHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())
	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	t, err := h.svc.Get(r.Context(), ticketID, tenantID, claims.Role, claims.UserID)
	if err != nil {
		if errors.Is(err, ticket.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "ticket not found")
			return
		}
		if errors.Is(err, ticket.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "not allowed to view this ticket")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to get ticket")
		return
	}

	response.JSON(w, http.StatusOK, toTicketDTO(t))
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())

	var req CreateTicketReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	t, err := h.svc.Create(r.Context(), appTicket.CreateInput{
		TenantID:    tenantID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Category:    req.Category,
		CreatedBy:   claims.UserID,
	})
	if err != nil {
		if errors.Is(err, ticket.ErrValidation) {
			response.Error(w, http.StatusBadRequest, "title and description are required")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create ticket")
		return
	}

	response.JSON(w, http.StatusCreated, toTicketDTO(t))
}

func (h *TicketHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())
	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateStatusReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	err = h.svc.ChangeStatus(r.Context(), ticketID, tenantID, claims.UserID, claims.Role, req.Status)
	if err != nil {
		switch {
		case errors.Is(err, ticket.ErrValidation):
			response.Error(w, http.StatusBadRequest, "invalid status")
		case errors.Is(err, ticket.ErrNotFound):
			response.Error(w, http.StatusNotFound, "ticket not found")
		case errors.Is(err, ticket.ErrForbidden):
			response.Error(w, http.StatusForbidden, "not allowed to update this ticket")
		case errors.Is(err, ticket.ErrInvalidTransition):
			response.Error(w, http.StatusBadRequest, "invalid status transition")
		default:
			response.Error(w, http.StatusInternalServerError, "failed to update ticket")
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": req.Status})
}

func (h *TicketHandler) Assign(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())
	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req AssignReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	assignedTo := uuid.Nil
	if req.AssignedTo != nil {
		assignedTo = *req.AssignedTo
	}

	if err := h.svc.Assign(r.Context(), ticketID, tenantID, assignedTo, claims.UserID); err != nil {
		if errors.Is(err, ticket.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "ticket not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to assign ticket")
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{"assigned_to": req.AssignedTo})
}

func (h *TicketHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())
	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	comments, err := h.svc.ListComments(r.Context(), ticketID, tenantID, claims.Role, claims.UserID)
	if err != nil {
		if errors.Is(err, ticket.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "ticket not found")
			return
		}
		if errors.Is(err, ticket.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "not allowed to view this ticket")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to list comments")
		return
	}

	dtos := make([]CommentDTO, len(comments))
	for i, c := range comments {
		dtos[i] = toCommentDTO(&c)
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *TicketHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())
	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req CreateCommentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	c, err := h.svc.AddComment(r.Context(), appTicket.AddCommentInput{
		TicketID:   ticketID,
		TenantID:   tenantID,
		UserID:     claims.UserID,
		Role:       claims.Role,
		Content:    req.Content,
		IsInternal: req.IsInternal,
	})
	if err != nil {
		switch {
		case errors.Is(err, ticket.ErrValidation):
			response.Error(w, http.StatusBadRequest, "content is required")
		case errors.Is(err, ticket.ErrNotFound):
			response.Error(w, http.StatusNotFound, "ticket not found")
		case errors.Is(err, ticket.ErrForbidden):
			response.Error(w, http.StatusForbidden, "not allowed to comment on this ticket")
		default:
			response.Error(w, http.StatusInternalServerError, "failed to create comment")
		}
		return
	}

	response.JSON(w, http.StatusCreated, toCommentDTO(c))
}

func (h *TicketHandler) ListHistory(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())
	ticketID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	entries, err := h.svc.ListHistory(r.Context(), ticketID, tenantID, claims.Role, claims.UserID)
	if err != nil {
		if errors.Is(err, ticket.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "ticket not found")
			return
		}
		if errors.Is(err, ticket.ErrForbidden) {
			response.Error(w, http.StatusForbidden, "not allowed to view this ticket")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to list history")
		return
	}

	dtos := make([]HistoryDTO, len(entries))
	for i, e := range entries {
		dtos[i] = toHistoryDTO(&e)
	}
	response.JSON(w, http.StatusOK, dtos)
}
