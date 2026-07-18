package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appAuto "github.com/trickreport/backend/internal/application/automation"
	"github.com/trickreport/backend/internal/domain/automation"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/interfaces/http/validator"
)

type AutomationHandler struct {
	svc *appAuto.Service
}

func NewAutomationHandler(svc *appAuto.Service) *AutomationHandler {
	return &AutomationHandler{svc: svc}
}

type RuleDTO struct {
	ID          uuid.UUID       `json:"id"`
	TenantID    uuid.UUID       `json:"tenant_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	TriggerType string          `json:"trigger_type"`
	Conditions  map[string]any  `json:"conditions"`
	Actions     []any           `json:"actions"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type createRuleReq struct {
	Name        string         `json:"name" validate:"required,min=2,max=255"`
	Description string         `json:"description" validate:"max=1000"`
	TriggerType string         `json:"trigger_type" validate:"required,oneof=ticket_created status_changed sla_breach priority_changed"`
	Conditions  map[string]any `json:"conditions"`
	Actions     []any          `json:"actions"`
	IsActive    bool           `json:"is_active"`
}

type updateRuleReq struct {
	Name        string         `json:"name" validate:"required,min=2,max=255"`
	Description string         `json:"description" validate:"max=1000"`
	TriggerType string         `json:"trigger_type" validate:"required,oneof=ticket_created status_changed sla_breach priority_changed"`
	Conditions  map[string]any `json:"conditions"`
	Actions     []any          `json:"actions"`
	IsActive    bool           `json:"is_active"`
}

func toRuleDTO(r *automation.Rule) RuleDTO {
	return RuleDTO{
		ID: r.ID, TenantID: r.TenantID, Name: r.Name, Description: r.Description,
		TriggerType: r.TriggerType, Conditions: r.Conditions, Actions: r.Actions,
		IsActive: r.IsActive, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func (h *AutomationHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	rules, err := h.svc.List(r.Context(), tenantID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list automation rules")
		return
	}

	dtos := make([]RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = toRuleDTO(&rule)
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *AutomationHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	var req createRuleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	rule, err := h.svc.Create(r.Context(), appAuto.CreateInput{
		TenantID: tenantID, Name: req.Name, Description: req.Description,
		TriggerType: req.TriggerType, Conditions: req.Conditions,
		Actions: req.Actions, IsActive: req.IsActive,
	})
	if err != nil {
		if errors.Is(err, automation.ErrValidation) {
			response.Error(w, http.StatusBadRequest, "name and trigger_type are required")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create automation rule")
		return
	}

	response.JSON(w, http.StatusCreated, toRuleDTO(rule))
}

func (h *AutomationHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateRuleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	rule, err := h.svc.Update(r.Context(), appAuto.UpdateInput{
		ID: id, TenantID: tenantID, Name: req.Name, Description: req.Description,
		TriggerType: req.TriggerType, Conditions: req.Conditions,
		Actions: req.Actions, IsActive: req.IsActive,
	})
	if err != nil {
		if errors.Is(err, automation.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "rule not found")
			return
		}
		if errors.Is(err, automation.ErrValidation) {
			response.Error(w, http.StatusBadRequest, "name and trigger_type are required")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update automation rule")
		return
	}

	response.JSON(w, http.StatusOK, toRuleDTO(rule))
}

func (h *AutomationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Delete(r.Context(), id, tenantID); err != nil {
		if errors.Is(err, automation.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "rule not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete rule")
		return
	}

	response.NoContent(w)
}
