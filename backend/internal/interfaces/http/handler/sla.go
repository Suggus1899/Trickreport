package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appSLA "github.com/trickreport/backend/internal/application/sla"
	"github.com/trickreport/backend/internal/domain/sla"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/interfaces/http/validator"
)

type SLAHandler struct {
	svc *appSLA.Service
}

func NewSLAHandler(svc *appSLA.Service) *SLAHandler {
	return &SLAHandler{svc: svc}
}

type SLAPolicyDTO struct {
	ID                    uuid.UUID `json:"id"`
	TenantID              uuid.UUID `json:"tenant_id"`
	Priority              string    `json:"priority"`
	ResponseTimeMinutes   int       `json:"response_time_minutes"`
	ResolutionTimeMinutes int       `json:"resolution_time_minutes"`
	EscalationMinutes     int       `json:"escalation_minutes"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type upsertSLAReq struct {
	ResponseTimeMinutes   int `json:"response_time_minutes" validate:"required,min=1,max=525600"`
	ResolutionTimeMinutes int `json:"resolution_time_minutes" validate:"required,min=1,max=525600"`
	EscalationMinutes     int `json:"escalation_minutes" validate:"required,min=1,max=525600"`
}

func toSLADTO(p *sla.Policy) SLAPolicyDTO {
	return SLAPolicyDTO{
		ID: p.ID, TenantID: p.TenantID, Priority: p.Priority,
		ResponseTimeMinutes: p.ResponseTimeMinutes,
		ResolutionTimeMinutes: p.ResolutionTimeMinutes,
		EscalationMinutes: p.EscalationMinutes,
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func (h *SLAHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	policies, err := h.svc.List(r.Context(), tenantID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list sla policies")
		return
	}

	dtos := make([]SLAPolicyDTO, len(policies))
	for i, p := range policies {
		dtos[i] = toSLADTO(&p)
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *SLAHandler) Upsert(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	priority := chi.URLParam(r, "priority")

	var req upsertSLAReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	p, err := h.svc.Upsert(r.Context(), appSLA.UpsertInput{
		TenantID:              tenantID,
		Priority:              priority,
		ResponseTimeMinutes:   req.ResponseTimeMinutes,
		ResolutionTimeMinutes: req.ResolutionTimeMinutes,
		EscalationMinutes:     req.EscalationMinutes,
	})
	if err != nil {
		if errors.Is(err, sla.ErrInvalidPriority) {
			response.Error(w, http.StatusBadRequest, "invalid priority")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to upsert sla policy")
		return
	}

	response.JSON(w, http.StatusOK, toSLADTO(p))
}
