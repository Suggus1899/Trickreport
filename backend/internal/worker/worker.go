package worker

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// systemUserID is the super admin from the seed migration (001_init.sql).
// Used for worker-generated comments and history entries.
// In a future refactor this should be configurable.
const systemUserID = "00000000-0000-0000-0000-000000000002"

type Worker struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Worker {
	return &Worker{db: db}
}

// Start begins the background worker loop.
// It checks for SLA breaches and runs automations.
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Info().Msg("Background worker started")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Background worker stopped")
			return
		case <-ticker.C:
			w.checkSLABreaches(ctx)
		}
	}
}

func (w *Worker) checkSLABreaches(ctx context.Context) {
	// Find tickets that are not resolved/closed, not already flagged as breached,
	// and whose created_at + resolution_time_minutes < NOW().
	q := `
		WITH breached_tickets AS (
			SELECT t.id, t.tenant_id, t.priority, s.resolution_time_minutes
			FROM tickets t
			JOIN sla_policies s ON t.tenant_id = s.tenant_id AND t.priority = s.priority
			WHERE t.status NOT IN ('resolved', 'closed')
			  AND t.sla_breached = FALSE
			  AND t.created_at + (s.resolution_time_minutes || ' minutes')::interval < NOW()
		)
		UPDATE tickets
		SET sla_breached = TRUE, updated_at = NOW()
		WHERE id IN (SELECT id FROM breached_tickets)
		RETURNING id, tenant_id;
	`

	rows, err := w.db.Query(ctx, q)
	if err != nil {
		log.Error().Err(err).Msg("Worker failed to check SLA breaches")
		return
	}
	defer rows.Close()

	breachCount := 0
	for rows.Next() {
		var id, tenantID string
		if err := rows.Scan(&id, &tenantID); err != nil {
			log.Error().Err(err).Msg("Worker failed to scan breached ticket row")
			continue
		}
		breachCount++

		// Insert internal comment — ticket_comments has no tenant_id column,
		// it's inferred from the ticket. Column is "content", not "body".
		// user_id is NOT NULL — use the system admin from seed data.
		if _, err := w.db.Exec(ctx, `
			INSERT INTO ticket_comments (ticket_id, user_id, content, is_internal)
			VALUES ($1, $2, 'SYSTEM: SLA Resolution Time Breached', TRUE)
		`, id, systemUserID); err != nil {
			log.Error().Err(err).Str("ticket_id", id).Msg("Worker failed to insert SLA breach comment")
		}

		// Insert history entry — ticket_history columns are field/old_value/new_value,
		// not action/details. No tenant_id column either.
		if _, err := w.db.Exec(ctx, `
			INSERT INTO ticket_history (ticket_id, user_id, field, old_value, new_value)
			VALUES ($1, $2, 'sla_breached', 'false', 'true')
		`, id, systemUserID); err != nil {
			log.Error().Err(err).Str("ticket_id", id).Msg("Worker failed to insert SLA breach history")
		}
	}

	if breachCount > 0 {
		log.Info().Int("breached_tickets", breachCount).Msg("SLA breaches detected and marked")
	}
}
