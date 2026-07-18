package analytics

import "context"

// Summary holds high-level KPI metrics for a tenant.
type Summary struct {
	TotalTickets    int `json:"total_tickets"`
	OpenTickets     int `json:"open_tickets"`
	ResolvedTickets int `json:"resolved_tickets"`
	SLABreached     int `json:"sla_breached"`
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

// Repository is the port for analytics queries.
type Repository interface {
	GetSummary(ctx context.Context, tenantID string) (Summary, error)
	GetVolume(ctx context.Context, tenantID string) ([]VolumePoint, error)
	GetStatusDistribution(ctx context.Context, tenantID string) ([]StatusDistribution, error)
	GetResolutionTime(ctx context.Context, tenantID string) ([]ResolutionMetrics, error)
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
func (s *Service) GetSummary(ctx context.Context, tenantID string) (Summary, error) {
	return s.repo.GetSummary(ctx, tenantID)
}

// GetVolume returns ticket volume per day for the last 30 days.
func (s *Service) GetVolume(ctx context.Context, tenantID string) ([]VolumePoint, error) {
	return s.repo.GetVolume(ctx, tenantID)
}

// GetStatusDistribution returns the breakdown of tickets by status.
func (s *Service) GetStatusDistribution(ctx context.Context, tenantID string) ([]StatusDistribution, error) {
	return s.repo.GetStatusDistribution(ctx, tenantID)
}

// GetResolutionTime returns the average resolution time in hours, grouped by priority.
func (s *Service) GetResolutionTime(ctx context.Context, tenantID string) ([]ResolutionMetrics, error) {
	return s.repo.GetResolutionTime(ctx, tenantID)
}
