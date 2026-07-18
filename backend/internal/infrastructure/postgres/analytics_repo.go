package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/trickreport/backend/internal/application/analytics"
)

// AnalyticsRepo implements analytics.Repository.
type AnalyticsRepo struct {
	db *pgxpool.Pool
}

// NewAnalyticsRepo creates a new AnalyticsRepo.
func NewAnalyticsRepo(db *pgxpool.Pool) *AnalyticsRepo {
	return &AnalyticsRepo{db: db}
}

// Compile-time assertion that AnalyticsRepo implements analytics.Repository.
var _ analytics.Repository = (*AnalyticsRepo)(nil)

// GetSummary returns aggregate ticket counts for a tenant.
func (r *AnalyticsRepo) GetSummary(ctx context.Context, tenantID uuid.UUID) (analytics.Summary, error) {
	const q = `SELECT COUNT(*), COUNT(*) FILTER (WHERE status IN ('open', 'in_progress', 'waiting_client')), COUNT(*) FILTER (WHERE status = 'resolved'), COUNT(*) FILTER (WHERE sla_breached = TRUE) FROM tickets WHERE tenant_id = $1`

	var s analytics.Summary
	if err := r.db.QueryRow(ctx, q, tenantID).
		Scan(&s.TotalTickets, &s.OpenTickets, &s.ResolvedTickets, &s.SLABreached); err != nil {
		return analytics.Summary{}, fmt.Errorf("analytics_repo.GetSummary: %w", err)
	}
	return s, nil
}

// GetVolume returns daily ticket volume for the last 30 days.
func (r *AnalyticsRepo) GetVolume(ctx context.Context, tenantID uuid.UUID) ([]analytics.VolumePoint, error) {
	const q = `SELECT TO_CHAR(created_at, 'YYYY-MM-DD') as date, COUNT(*) FROM tickets WHERE tenant_id = $1 AND created_at >= NOW() - INTERVAL '30 days' GROUP BY date ORDER BY date ASC`

	rows, err := r.db.Query(ctx, q, tenantID)
	if err != nil {
		return nil, fmt.Errorf("analytics_repo.GetVolume: query: %w", err)
	}
	defer rows.Close()

	var points []analytics.VolumePoint
	for rows.Next() {
		var p analytics.VolumePoint
		if err := rows.Scan(&p.Date, &p.Count); err != nil {
			return nil, fmt.Errorf("analytics_repo.GetVolume: scan: %w", err)
		}
		points = append(points, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("analytics_repo.GetVolume: rows: %w", err)
	}
	return points, nil
}

// GetStatusDistribution returns ticket counts grouped by status.
func (r *AnalyticsRepo) GetStatusDistribution(ctx context.Context, tenantID uuid.UUID) ([]analytics.StatusDistribution, error) {
	const q = `SELECT status, COUNT(*) FROM tickets WHERE tenant_id = $1 GROUP BY status`

	rows, err := r.db.Query(ctx, q, tenantID)
	if err != nil {
		return nil, fmt.Errorf("analytics_repo.GetStatusDistribution: query: %w", err)
	}
	defer rows.Close()

	var dist []analytics.StatusDistribution
	for rows.Next() {
		var d analytics.StatusDistribution
		if err := rows.Scan(&d.Status, &d.Count); err != nil {
			return nil, fmt.Errorf("analytics_repo.GetStatusDistribution: scan: %w", err)
		}
		dist = append(dist, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("analytics_repo.GetStatusDistribution: rows: %w", err)
	}
	return dist, nil
}

// GetResolutionTime returns average resolution hours grouped by priority.
func (r *AnalyticsRepo) GetResolutionTime(ctx context.Context, tenantID uuid.UUID) ([]analytics.ResolutionMetrics, error) {
	const q = `SELECT priority, COALESCE(AVG(EXTRACT(EPOCH FROM (updated_at - created_at)) / 3600), 0) AS avg_hours FROM tickets WHERE tenant_id = $1 AND status IN ('resolved', 'closed') GROUP BY priority`

	rows, err := r.db.Query(ctx, q, tenantID)
	if err != nil {
		return nil, fmt.Errorf("analytics_repo.GetResolutionTime: query: %w", err)
	}
	defer rows.Close()

	var metrics []analytics.ResolutionMetrics
	for rows.Next() {
		var m analytics.ResolutionMetrics
		if err := rows.Scan(&m.Priority, &m.AvgHours); err != nil {
			return nil, fmt.Errorf("analytics_repo.GetResolutionTime: scan: %w", err)
		}
		metrics = append(metrics, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("analytics_repo.GetResolutionTime: rows: %w", err)
	}
	return metrics, nil
}
