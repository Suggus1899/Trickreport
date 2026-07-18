package worker

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// Worker checks for SLA breaches and runs automations in the background.
type Worker struct {
	db           *pgxpool.Pool
	systemUserID uuid.UUID
	interval     time.Duration

	wg sync.WaitGroup
}

// Option configures a Worker.
type Option func(*Worker)

// WithSystemUserID sets the user ID used for worker-generated comments and
// history entries. Defaults to the seed admin UUID.
func WithSystemUserID(id uuid.UUID) Option {
	return func(w *Worker) { w.systemUserID = id }
}

// WithInterval sets the poll interval. Defaults to 1 minute.
func WithInterval(d time.Duration) Option {
	return func(w *Worker) { w.interval = d }
}

// New creates a Worker. Pass options to customize behavior.
func New(db *pgxpool.Pool, opts ...Option) *Worker {
	w := &Worker{
		db:           db,
		systemUserID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		interval:     1 * time.Minute,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Start begins the background worker loop. It blocks until ctx is cancelled.
// Use Stop for graceful shutdown — it waits for the current check to finish.
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Info().Dur("interval", w.interval).Msg("Background worker started")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Background worker stopping — waiting for current check")
			w.wg.Wait()
			log.Info().Msg("Background worker stopped")
			return
		case <-ticker.C:
			w.wg.Add(1)
			func() {
				defer w.wg.Done()
				w.checkSLABreaches(ctx)
			}()
		}
	}
}

func (w *Worker) checkSLABreaches(ctx context.Context) {
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
		var id, tenantID pgtype.UUID
		if err := rows.Scan(&id, &tenantID); err != nil {
			log.Error().Err(err).Msg("Worker failed to scan breached ticket row")
			continue
		}
		breachCount++
		ticketIDStr := ""
		if id.Valid {
			ticketIDStr = uuid.UUID(id.Bytes).String()
		}

		if _, err := w.db.Exec(ctx, `
			INSERT INTO ticket_comments (ticket_id, user_id, content, is_internal)
			VALUES ($1, $2, 'SYSTEM: SLA Resolution Time Breached', TRUE)
		`, id, w.systemUserID); err != nil {
			log.Error().Err(err).Str("ticket_id", ticketIDStr).Msg("Worker failed to insert SLA breach comment")
		}

		if _, err := w.db.Exec(ctx, `
			INSERT INTO ticket_history (ticket_id, user_id, field, old_value, new_value)
			VALUES ($1, $2, 'sla_breached', 'false', 'true')
		`, id, w.systemUserID); err != nil {
			log.Error().Err(err).Str("ticket_id", ticketIDStr).Msg("Worker failed to insert SLA breach history")
		}
	}

	if breachCount > 0 {
		log.Info().Int("breached_tickets", breachCount).Msg("SLA breaches detected and marked")
	}
}
