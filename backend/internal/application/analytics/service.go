package analytics

import (
	"context"

	"github.com/google/uuid"
)

// Summary holds high-level KPI metrics for a tenant.
type Summary struct {
	TotalTickets        int `json:"total_tickets"`
	OpenTickets         int `json:"open_tickets"`
	ResolvedTickets     int `json:"resolved_tickets"`
	SLABreached         int `json:"sla_breached"`
	SLAEligible         int `json:"sla_eligible"`
	ResponseSLABreached int `json:"response_sla_breached"`
	ResponseSLAEligible int `json:"response_sla_eligible"`
}

type viewerKey struct{}

// WithViewer limits end-user analytics to the same tickets they can list.
func WithViewer(ctx context.Context, role string, userID uuid.UUID) context.Context {
	if role == "end_user" {
		return context.WithValue(ctx, viewerKey{}, userID)
	}
	return ctx
}

func Viewer(ctx context.Context) any {
	if id, ok := ctx.Value(viewerKey{}).(uuid.UUID); ok {
		return id
	}
	return nil
}

// VolumePoint represents ticket count for a single day.
type VolumePoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// StatusDistribution represents ticket count per status.
type StatusDistribution struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// ResolutionMetrics represents average resolution time per priority.
type ResolutionMetrics struct {
	Priority string  `json:"priority"`
	AvgHours float64 `json:"avg_hours"`
}

// SLACompliance represents SLA breach statistics.
type SLACompliance struct {
	Total     int `json:"total"`
	Breached  int `json:"breached"`
	Compliant int `json:"compliant"`
}

// ChartsData holds aggregated chart data optimized for frontend rendering.
type ChartsData struct {
	Volume                []VolumePoint        `json:"volume"`
	StatusDistribution    []StatusDistribution `json:"status_distribution"`
	ResolutionTime        []ResolutionMetrics  `json:"resolution_time"`
	SLACompliance         SLACompliance        `json:"sla_compliance"`
	FirstResponseTime     []ResolutionMetrics  `json:"first_response_time"`
	ResponseSLACompliance SLACompliance        `json:"response_sla_compliance"`
	ResolvedVolume        []VolumePoint        `json:"resolved_volume"`
	Timezone              string               `json:"timezone"`
}

// Repository is the port for analytics queries.
type Repository interface {
	GetSummary(ctx context.Context, tenantID uuid.UUID) (Summary, error)
	GetVolume(ctx context.Context, tenantID uuid.UUID) ([]VolumePoint, error)
	GetStatusDistribution(ctx context.Context, tenantID uuid.UUID) ([]StatusDistribution, error)
	GetResolutionTime(ctx context.Context, tenantID uuid.UUID) ([]ResolutionMetrics, error)
}

// Service is the application service for analytics operations.
type Service struct {
	repo Repository
}

// NewService creates a new analytics application service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetSummary returns KPI metrics for the tenant.
func (s *Service) GetSummary(ctx context.Context, tenantID uuid.UUID) (Summary, error) {
	return s.repo.GetSummary(ctx, tenantID)
}

// GetVolume returns ticket volume per day for the last 30 days.
func (s *Service) GetVolume(ctx context.Context, tenantID uuid.UUID) ([]VolumePoint, error) {
	return s.repo.GetVolume(ctx, tenantID)
}

// GetStatusDistribution returns the breakdown of tickets by status.
func (s *Service) GetStatusDistribution(ctx context.Context, tenantID uuid.UUID) ([]StatusDistribution, error) {
	return s.repo.GetStatusDistribution(ctx, tenantID)
}

// GetResolutionTime returns the average resolution time in hours, grouped by priority.
func (s *Service) GetResolutionTime(ctx context.Context, tenantID uuid.UUID) ([]ResolutionMetrics, error) {
	return s.repo.GetResolutionTime(ctx, tenantID)
}

// GetCharts returns aggregated chart data (volume, status distribution,
// resolution time, and SLA compliance) in a single response.
func (s *Service) GetCharts(ctx context.Context, tenantID uuid.UUID) (ChartsData, error) {
	var data ChartsData
	var err error

	data.Volume, err = s.repo.GetVolume(ctx, tenantID)
	if err != nil {
		return ChartsData{}, err
	}

	data.StatusDistribution, err = s.repo.GetStatusDistribution(ctx, tenantID)
	if err != nil {
		return ChartsData{}, err
	}

	data.ResolutionTime, err = s.repo.GetResolutionTime(ctx, tenantID)
	if err != nil {
		return ChartsData{}, err
	}

	summary, err := s.repo.GetSummary(ctx, tenantID)
	if err != nil {
		return ChartsData{}, err
	}
	data.SLACompliance = SLACompliance{
		Total:     summary.SLAEligible,
		Breached:  summary.SLABreached,
		Compliant: summary.SLAEligible - summary.SLABreached,
	}
	data.ResponseSLACompliance = SLACompliance{Total: summary.ResponseSLAEligible, Breached: summary.ResponseSLABreached, Compliant: summary.ResponseSLAEligible - summary.ResponseSLABreached}
	if repo, ok := s.repo.(interface {
		GetFirstResponseTime(context.Context, uuid.UUID) ([]ResolutionMetrics, error)
		GetResolvedVolume(context.Context, uuid.UUID) ([]VolumePoint, error)
	}); ok {
		data.FirstResponseTime, err = repo.GetFirstResponseTime(ctx, tenantID)
		if err != nil {
			return ChartsData{}, err
		}
		data.ResolvedVolume, err = repo.GetResolvedVolume(ctx, tenantID)
		if err != nil {
			return ChartsData{}, err
		}
	}
	if data.Volume == nil {
		data.Volume = []VolumePoint{}
	}
	if data.StatusDistribution == nil {
		data.StatusDistribution = []StatusDistribution{}
	}
	if data.ResolutionTime == nil {
		data.ResolutionTime = []ResolutionMetrics{}
	}
	if data.FirstResponseTime == nil {
		data.FirstResponseTime = []ResolutionMetrics{}
	}
	if data.ResolvedVolume == nil {
		data.ResolvedVolume = []VolumePoint{}
	}
	data.Timezone = "UTC"

	return data, nil
}
