package worker

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Job represents a unit of background work.
type Job struct {
	ID         uuid.UUID
	Type       string
	Payload    any
	Attempts   int
	MaxRetries int
	NextRunAt  time.Time
}

// JobHandler processes a single job. Returning an error triggers a retry
// (if attempts remain).
type JobHandler func(ctx context.Context, job Job) error

// JobQueue is a simple in-memory job queue with retry and backoff support.
type JobQueue struct {
	mu       sync.Mutex
	queue    []Job
	handlers map[string]JobHandler
}

// NewJobQueue creates a new JobQueue.
func NewJobQueue() *JobQueue {
	return &JobQueue{
		handlers: make(map[string]JobHandler),
	}
}

// Register associates a handler with a job type.
func (q *JobQueue) Register(jobType string, handler JobHandler) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.handlers[jobType] = handler
}

// Enqueue adds a job to the queue. If NextRunAt is zero, the job is eligible
// immediately.
func (q *JobQueue) Enqueue(job Job) {
	if job.ID == uuid.Nil {
		job.ID = uuid.New()
	}
	if job.NextRunAt.IsZero() {
		job.NextRunAt = time.Now()
	}
	q.mu.Lock()
	q.queue = append(q.queue, job)
	q.mu.Unlock()
}

// Dequeue removes and returns the next eligible job (NextRunAt <= now).
// Returns false if no job is ready.
func (q *JobQueue) Dequeue() (Job, bool) {
	now := time.Now()
	q.mu.Lock()
	defer q.mu.Unlock()

	for i, job := range q.queue {
		if !job.NextRunAt.After(now) {
			q.queue = append(q.queue[:i], q.queue[i+1:]...)
			return job, true
		}
	}
	return Job{}, false
}

// Len returns the number of jobs currently in the queue.
func (q *JobQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue)
}

// Process runs the job processing loop until ctx is cancelled. Jobs are
// dequeued, executed via their registered handler, and retried with
// exponential backoff (1s, 2s, 4s) up to MaxRetries times.
func (q *JobQueue) Process(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			q.processReadyJobs(ctx)
		}
	}
}

func (q *JobQueue) processReadyJobs(ctx context.Context) {
	for {
		job, ok := q.Dequeue()
		if !ok {
			return
		}
		q.executeJob(ctx, job)
	}
}

func (q *JobQueue) executeJob(ctx context.Context, job Job) {
	q.mu.Lock()
	handler, ok := q.handlers[job.Type]
	q.mu.Unlock()

	if !ok {
		log.Error().Str("job_type", job.Type).Str("job_id", job.ID.String()).Msg("No handler registered for job type")
		return
	}

	job.Attempts++
	if err := handler(ctx, job); err != nil {
		log.Error().Err(err).Str("job_type", job.Type).Int("attempt", job.Attempts).Msg("Job failed")
		if job.Attempts > job.MaxRetries {
			log.Error().Str("job_type", job.Type).Str("job_id", job.ID.String()).Msg("Job exhausted retries — dropping")
			return
		}
		backoff := backoffForAttempt(job.Attempts)
		job.NextRunAt = time.Now().Add(backoff)
		q.Enqueue(job)
	}
}

func backoffForAttempt(attempt int) time.Duration {
	backoffs := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	return backoffs[(attempt-1)%len(backoffs)]
}
