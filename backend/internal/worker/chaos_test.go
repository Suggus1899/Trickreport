package worker

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// TestWorker_RetryLogic_SimulatesDBDisconnection verifies that when the SLA
// check fails (simulating a DB disconnection error), the worker retries up to
// maxRetries times and increments the error metric — without crashing.
func TestWorker_RetryLogic_SimulatesDBDisconnection(t *testing.T) {
	w := New(nil, WithRetry(2))

	var attempts int32
	failFn := func(ctx context.Context) error {
		atomic.AddInt32(&attempts, 1)
		return errors.New("connection refused: db disconnected")
	}

	// runWithRetry should call failFn 3 times (1 initial + 2 retries).
	// Use a short context so the backoff sleeps don't make the test slow.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	w.runWithRetry(ctx, failFn)

	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Errorf("expected 3 attempts (1 + 2 retries), got %d", got)
	}

	breaches, _, automations, errs := w.metrics.Snapshot()
	if breaches != 0 {
		t.Errorf("breaches = %d, want 0", breaches)
	}
	if automations != 0 {
		t.Errorf("automations = %d, want 0", automations)
	}
	if errs != 3 {
		t.Errorf("errors = %d, want 3", errs)
	}
}

// TestWorker_RetryLogic_SucceedsOnRetry verifies that a transient failure
// followed by success does not exhaust the retry budget.
func TestWorker_RetryLogic_SucceedsOnRetry(t *testing.T) {
	w := New(nil, WithRetry(3))

	var attempts int32
	recoverFn := func(ctx context.Context) error {
		n := atomic.AddInt32(&attempts, 1)
		if n < 2 {
			return errors.New("transient error")
		}
		return nil // success on second attempt
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	w.runWithRetry(ctx, recoverFn)

	if got := atomic.LoadInt32(&attempts); got != 2 {
		t.Errorf("expected 2 attempts, got %d", got)
	}

	_, _, _, errs := w.metrics.Snapshot()
	if errs != 1 {
		t.Errorf("errors = %d, want 1 (only the first failure)", errs)
	}
}

// TestWorker_RetryLogic_ContextCancellationDuringBackoff verifies that
// cancelling the context during a backoff sleep exits promptly.
func TestWorker_RetryLogic_ContextCancellationDuringBackoff(t *testing.T) {
	w := New(nil, WithRetry(10))

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	w.runWithRetry(ctx, func(ctx context.Context) error {
		return errors.New("always fail")
	})
	elapsed := time.Since(start)

	// Should exit shortly after cancel, well before the 4s backoff would complete.
	if elapsed > 2*time.Second {
		t.Errorf("runWithRetry took %v, expected quick exit on cancel", elapsed)
	}
}

// TestWorker_Metrics_Snapshot verifies the metrics snapshot returns the
// current values of all counters.
func TestWorker_Metrics_Snapshot(t *testing.T) {
	var m Metrics
	b, esc, a, e := m.Snapshot()
	if b != 0 || esc != 0 || a != 0 || e != 0 {
		t.Errorf("initial snapshot = (%d, %d, %d, %d), want all 0", b, esc, a, e)
	}

	atomic.AddInt64(&m.BreachesDetected, 5)
	atomic.AddInt64(&m.EscalationsDetected, 4)
	atomic.AddInt64(&m.AutomationsRun, 3)
	atomic.AddInt64(&m.Errors, 2)

	b, esc, a, e = m.Snapshot()
	if b != 5 || esc != 4 || a != 3 || e != 2 {
		t.Errorf("snapshot = (%d, %d, %d, %d), want (5, 4, 3, 2)", b, esc, a, e)
	}
}

// TestWorker_WithRetryOption verifies the WithRetry option sets maxRetries.
func TestWorker_WithRetryOption(t *testing.T) {
	w := New(nil, WithRetry(5))
	if w.maxRetries != 5 {
		t.Errorf("maxRetries = %d, want 5", w.maxRetries)
	}

	w2 := New(nil)
	if w2.maxRetries != 0 {
		t.Errorf("default maxRetries = %d, want 0", w2.maxRetries)
	}
}
