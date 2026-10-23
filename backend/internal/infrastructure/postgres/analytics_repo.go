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
	const q = `WITH scoped AS (
	SELECT *, (status NOT IN ('resolved','closed')) AS active FROM tickets WHERE tenant_id=$1 AND ($2::uuid IS NULL OR created_by=$2)
	), assessed AS (
	SELECT *, (sla_breached OR (sla_deadline IS NOT NULL AND (resolved_at IS NOT NULL OR (active AND sla_deadline<NOW())))) AS resolution_assessed,
	(response_sla_breached OR (response_deadline IS NOT NULL AND (first_responded_at IS NOT NULL OR (active AND response_deadline<NOW())))) AS response_assessed
	FROM scoped)
	SELECT COUNT(*), COUNT(*) FILTER(WHERE active), COUNT(*) FILTER(WHERE resolved_at IS NOT NULL),
	COUNT(*) FILTER(WHERE resolution_assessed AND (sla_breached OR COALESCE(resolved_at,NOW())>sla_deadline)), COUNT(*) FILTER(WHERE resolution_assessed),
	COUNT(*) FILTER(WHERE response_assessed AND (response_sla_breached OR COALESCE(first_responded_at,NOW())>response_deadline)), COUNT(*) FILTER(WHERE response_assessed)
	FROM assessed`

	var s analytics.Summary
	if err := r.db.QueryRow(ctx, q, tenantID, analytics.Viewer(ctx)).
		Scan(&s.TotalTickets, &s.OpenTickets, &s.ResolvedTickets, &s.SLABreached, &s.SLAEligible, &s.ResponseSLABreached, &s.ResponseSLAEligible); err != nil {
		return analytics.Summary{}, fmt.Errorf("analytics_repo.GetSummary: %w", err)
	}
	return s, nil
}

// GetVolume returns daily ticket volume for the last 30 days.
func (r *AnalyticsRepo) GetVolume(ctx context.Context, tenantID uuid.UUID) ([]analytics.VolumePoint, error) {
	return r.volume(ctx, tenantID, "created_at")
}

func (r *AnalyticsRepo) GetResolvedVolume(ctx context.Context, tenantID uuid.UUID) ([]analytics.VolumePoint, error) {
	return r.volume(ctx, tenantID, "resolved_at")
}

func (r *AnalyticsRepo) volume(ctx context.Context, tenantID uuid.UUID, column string) ([]analytics.VolumePoint, error) {
	q := fmt.Sprintf(`WITH days AS(SELECT generate_series((NOW() AT TIME ZONE 'UTC')::date-29,(NOW() AT TIME ZONE 'UTC')::date,'1 day'::interval)::date AS day),counts AS(
	SELECT (%s AT TIME ZONE 'UTC')::date AS day,COUNT(*) AS count FROM tickets WHERE tenant_id=$1 AND ($2::uuid IS NULL OR created_by=$2) AND %s IS NOT NULL GROUP BY day)
	SELECT TO_CHAR(days.day,'YYYY-MM-DD'),COALESCE(counts.count,0) FROM days LEFT JOIN counts USING(day) ORDER BY days.day`, column, column)

	rows, err := r.db.Query(ctx, q, tenantID, analytics.Viewer(ctx))
	if err != nil {
		return nil, fmt.Errorf("analytics_repo.GetVolume: query: %w", err)
	}
	defer rows.Close()

	points := make([]analytics.VolumePoint, 0, 30)
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
	const q = `SELECT status, COUNT(*) FROM tickets WHERE tenant_id = $1 AND ($2::uuid IS NULL OR created_by=$2) GROUP BY status ORDER BY status`

	rows, err := r.db.Query(ctx, q, tenantID, analytics.Viewer(ctx))
	if err != nil {
		return nil, fmt.Errorf("analytics_repo.GetStatusDistribution: query: %w", err)
	}
	defer rows.Close()

	dist := make([]analytics.StatusDistribution, 0)
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
	return r.averageTime(ctx, tenantID, "resolved_at")
}

func (r *AnalyticsRepo) GetFirstResponseTime(ctx context.Context, tenantID uuid.UUID) ([]analytics.ResolutionMetrics, error) {
	return r.averageTime(ctx, tenantID, "first_responded_at")
}

func (r *AnalyticsRepo) averageTime(ctx context.Context, tenantID uuid.UUID, column string) ([]analytics.ResolutionMetrics, error) {
	q := fmt.Sprintf(`SELECT priority,AVG(EXTRACT(EPOCH FROM (%s-created_at))/3600) FROM tickets WHERE tenant_id=$1 AND ($2::uuid IS NULL OR created_by=$2) AND %s IS NOT NULL GROUP BY priority ORDER BY priority`, column, column)

	rows, err := r.db.Query(ctx, q, tenantID, analytics.Viewer(ctx))
	if err != nil {
		return nil, fmt.Errorf("analytics_repo.GetResolutionTime: query: %w", err)
	}
	defer rows.Close()

	metrics := make([]analytics.ResolutionMetrics, 0)
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
