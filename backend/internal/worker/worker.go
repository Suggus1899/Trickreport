package worker

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	appAuto "github.com/trickreport/backend/internal/application/automation"
)

// Metrics holds atomic counters for worker observability.
type Metrics struct {
	BreachesDetected int64
	AutomationsRun   int64
	Errors           int64
}

// Snapshot returns a copy of the current metric values.
func (m *Metrics) Snapshot() (breaches, automations, errors int64) {
	return atomic.LoadInt64(&m.BreachesDetected),
		atomic.LoadInt64(&m.AutomationsRun),
		atomic.LoadInt64(&m.Errors)
}

// Worker checks for SLA breaches and runs automations in the background.
type Worker struct {
	db           *pgxpool.Pool
	systemUserID uuid.UUID
	interval     time.Duration
	engine       *appAuto.Engine
	maxRetries   int

	metrics Metrics

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

// WithEngine wires the automation engine so SLA breach events are evaluated
// against automation rules.
func WithEngine(engine *appAuto.Engine) Option {
	return func(w *Worker) { w.engine = engine }
}

// WithRetry sets the maximum number of retries for the SLA breach check.
// Defaults to 0 (no retries).
func WithRetry(maxRetries int) Option {
	return func(w *Worker) { w.maxRetries = maxRetries }
}

// New creates a Worker. Pass options to customize behavior.
func New(db *pgxpool.Pool, opts ...Option) *Worker {
	w := &Worker{
		db:           db,
		systemUserID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		interval:     1 * time.Minute,
		maxRetries:   0,
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
				w.runWithRetry(ctx, w.checkSLABreachesErr)
			}()
		}
	}
}

// runWithRetry executes fn with retry logic and exponential backoff
// (1s, 2s, 4s). If fn returns an error, it is retried up to maxRetries times.
func (w *Worker) runWithRetry(ctx context.Context, fn func(ctx context.Context) error) {
	backoffs := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	for attempt := 0; ; attempt++ {
		if err := fn(ctx); err == nil {
			return
		} else {
			atomic.AddInt64(&w.metrics.Errors, 1)
			if attempt >= w.maxRetries {
				log.Error().Int("attempt", attempt+1).Msg("Worker check failed after max retries")
				return
			}
			backoff := backoffs[attempt%len(backoffs)]
			log.Warn().Int("attempt", attempt+1).Dur("backoff", backoff).Msg("Worker check failed — retrying")
			select {
			case <-ctx.Done():
				return
			case <-time.After(backoff):
			}
		}
	}
}

// HealthCheck verifies database connectivity. Returns nil if the database is
// reachable.
func (w *Worker) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return w.db.Ping(ctx)
}

// checkSLABreachesErr wraps checkSLABreaches to return an error for retry logic.
func (w *Worker) checkSLABreachesErr(ctx context.Context) error {
	return w.checkSLABreaches(ctx)
}

func (w *Worker) checkSLABreaches(ctx context.Context) error {
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
		return err
	}
	defer rows.Close()

	breachCount := 0
	for rows.Next() {
		var id, tenantID pgtype.UUID
		if err := rows.Scan(&id, &tenantID); err != nil {
			atomic.AddInt64(&w.metrics.Errors, 1)
			log.Error().Err(err).Msg("Worker failed to scan breached ticket row")
			continue
		}
		breachCount++
		ticketIDStr := ""
		tenantIDStr := ""
		if id.Valid {
			ticketIDStr = uuid.UUID(id.Bytes).String()
		}
		if tenantID.Valid {
			tenantIDStr = uuid.UUID(tenantID.Bytes).String()
		}

		if _, err := w.db.Exec(ctx, `
			INSERT INTO ticket_comments (ticket_id, user_id, content, is_internal)
			VALUES ($1, $2, 'SYSTEM: SLA Resolution Time Breached', TRUE)
		`, id, w.systemUserID); err != nil {
			atomic.AddInt64(&w.metrics.Errors, 1)
			log.Error().Err(err).Str("ticket_id", ticketIDStr).Msg("Worker failed to insert SLA breach comment")
		}

		if _, err := w.db.Exec(ctx, `
			INSERT INTO ticket_history (ticket_id, user_id, field, old_value, new_value)
			VALUES ($1, $2, 'sla_breached', 'false', 'true')
		`, id, w.systemUserID); err != nil {
			atomic.AddInt64(&w.metrics.Errors, 1)
			log.Error().Err(err).Str("ticket_id", ticketIDStr).Msg("Worker failed to insert SLA breach history")
		}

		// Fire automation rules for the SLA breach event.
		if w.engine != nil && id.Valid && tenantID.Valid {
			event := appAuto.Event{
				Type:     "sla_breach",
				TenantID: uuid.UUID(tenantID.Bytes),
				TicketID: uuid.UUID(id.Bytes),
				UserID:   w.systemUserID,
				OldValue: "false",
				NewValue: "true",
			}
			if err := w.engine.Evaluate(ctx, uuid.UUID(tenantID.Bytes), event); err != nil {
				atomic.AddInt64(&w.metrics.Errors, 1)
				log.Error().Err(err).Str("ticket_id", ticketIDStr).Str("tenant_id", tenantIDStr).Msg("Worker failed to evaluate automations for SLA breach")
			} else {
				atomic.AddInt64(&w.metrics.AutomationsRun, 1)
			}
		}
	}

	if breachCount > 0 {
		atomic.AddInt64(&w.metrics.BreachesDetected, int64(breachCount))
		log.Info().Int("breached_tickets", breachCount).Msg("SLA breaches detected and marked")
	}
	return nil
}
