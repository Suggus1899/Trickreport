package worker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	appAuto "github.com/trickreport/backend/internal/application/automation"
)

// Metrics holds atomic counters for worker observability.
type Metrics struct {
	BreachesDetected    int64
	EscalationsDetected int64
	AutomationsRun      int64
	Errors              int64
}

// Snapshot returns a copy of the current metric values.
func (m *Metrics) Snapshot() (breaches, escalations, automations, errors int64) {
	return atomic.LoadInt64(&m.BreachesDetected),
		atomic.LoadInt64(&m.EscalationsDetected),
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

	metrics     Metrics
	lastSuccess atomic.Int64
	lastFailed  atomic.Bool

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
	if ctx.Err() != nil {
		return
	}
	// Run immediately so readiness reflects a completed SLA cycle after startup.
	if w.db != nil {
		w.runCycle(ctx)
	}

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
				w.runCycle(ctx)
			}()
		}
	}
}

// runWithRetry executes fn with retry logic and exponential backoff
// (1s, 2s, 4s). If fn returns an error, it is retried up to maxRetries times.
func (w *Worker) runWithRetry(ctx context.Context, fn func(ctx context.Context) error) bool {
	backoffs := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	for attempt := 0; ; attempt++ {
		if err := fn(ctx); err == nil {
			return true
		}
		atomic.AddInt64(&w.metrics.Errors, 1)
		if attempt >= w.maxRetries {
			log.Error().Int("attempt", attempt+1).Msg("Worker check failed after max retries")
			return false
		}
		backoff := backoffs[attempt%len(backoffs)]
		log.Warn().Int("attempt", attempt+1).Dur("backoff", backoff).Msg("Worker check failed — retrying")
		select {
		case <-ctx.Done():
			return false
		case <-time.After(backoff):
		}
	}
}

func (w *Worker) runCycle(ctx context.Context) {
	breaches := w.runWithRetry(ctx, w.checkSLABreachesErr)
	escalations := w.runWithRetry(ctx, w.checkEscalationsErr)
	w.lastFailed.Store(!breaches || !escalations)
	if breaches && escalations {
		w.lastSuccess.Store(time.Now().Unix())
	}
}

func (w *Worker) LastSuccessUnix() int64 { return w.lastSuccess.Load() }

func (w *Worker) ReadyHealthy() bool {
	last := w.LastSuccessUnix()
	return last > 0 && !w.lastFailed.Load() && time.Since(time.Unix(last, 0)) <= 3*w.interval
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
	if err := w.processDeadline(ctx, "sla_breached", "sla_deadline", "SLA Resolution Time Breached", "sla_breach"); err != nil {
		return err
	}
	return w.processDeadline(ctx, "response_sla_breached", "response_deadline", "SLA First Response Time Breached", "sla_breach")
}
func (w *Worker) checkEscalationsErr(ctx context.Context) error { return w.checkEscalations(ctx) }
func (w *Worker) checkEscalations(ctx context.Context) error {
	return w.processDeadline(ctx, "escalated_at", "escalation_deadline", "SLA Escalation Threshold Reached", "escalation")
}

// processDeadline commits the claim, comment and history together before dispatching events.
func (w *Worker) processDeadline(ctx context.Context, field, deadline, message, eventType string) error {
	tx, err := w.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	condition := field + "=FALSE"
	update := field + "=TRUE"
	oldValue, newValue := "false", "true"
	if field == "escalated_at" {
		condition = field + " IS NULL"
		update = field + "=NOW()"
		oldValue, newValue = "", "escalated"
	}
	if field == "response_sla_breached" {
		condition += " AND first_responded_at IS NULL"
	}
	query := fmt.Sprintf(`SELECT id,tenant_id FROM tickets WHERE status NOT IN ('resolved','closed') AND %s AND %s < NOW() ORDER BY id FOR UPDATE SKIP LOCKED`, condition, deadline)
	rows, err := tx.Query(ctx, query)
	if err != nil {
		return err
	}
	type claimed struct{ id, tenant uuid.UUID }
	var tickets []claimed
	for rows.Next() {
		var ticket claimed
		if err := rows.Scan(&ticket.id, &ticket.tenant); err != nil {
			rows.Close()
			return err
		}
		tickets = append(tickets, ticket)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	for _, ticket := range tickets {
		if _, err := tx.Exec(ctx, `UPDATE tickets SET `+update+`,version=version+1,updated_at=NOW() WHERE id=$1`, ticket.id); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `INSERT INTO ticket_comments(ticket_id,user_id,content,is_internal)VALUES($1,$2,$3,TRUE)`, ticket.id, w.systemUserID, "SYSTEM: "+message); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `INSERT INTO ticket_history(ticket_id,user_id,field,old_value,new_value)VALUES($1,$2,$3,$4,$5)`, ticket.id, w.systemUserID, field, oldValue, newValue); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	if field == "escalated_at" {
		atomic.AddInt64(&w.metrics.EscalationsDetected, int64(len(tickets)))
	} else {
		atomic.AddInt64(&w.metrics.BreachesDetected, int64(len(tickets)))
	}
	for _, ticket := range tickets {
		if w.engine != nil {
			event := appAuto.Event{Type: eventType, TenantID: ticket.tenant, TicketID: ticket.id, UserID: w.systemUserID, OldValue: oldValue, NewValue: newValue}
			if err := w.engine.Evaluate(ctx, ticket.tenant, event); err != nil {
				return err
			}
			atomic.AddInt64(&w.metrics.AutomationsRun, 1)
		}
	}
	return nil
}
