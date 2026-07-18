package handler

import (
	"net/http"

	appAnalytics "github.com/trickreport/backend/internal/application/analytics"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
)

type AnalyticsHandler struct {
	svc *appAnalytics.Service
}

func NewAnalyticsHandler(svc *appAnalytics.Service) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

func (h *AnalyticsHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	summary, err := h.svc.GetSummary(r.Context(), tenantID)
	if err != nil {
		middleware.LoggerFromContext(r.Context()).Error().Err(err).Str("tenant_id", tenantID.String()).Msg("analytics summary failed")
		response.Error(w, http.StatusInternalServerError, "failed to get metrics")
		return
	}

	response.JSON(w, http.StatusOK, summary)
}

func (h *AnalyticsHandler) GetVolume(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	volume, err := h.svc.GetVolume(r.Context(), tenantID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get volume")
		return
	}

	response.JSON(w, http.StatusOK, volume)
}

func (h *AnalyticsHandler) GetStatusDistribution(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	dist, err := h.svc.GetStatusDistribution(r.Context(), tenantID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get status distribution")
		return
	}

	response.JSON(w, http.StatusOK, dist)
}

func (h *AnalyticsHandler) GetResolutionTime(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	metrics, err := h.svc.GetResolutionTime(r.Context(), tenantID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to get resolution time")
		return
	}

	response.JSON(w, http.StatusOK, metrics)
}
