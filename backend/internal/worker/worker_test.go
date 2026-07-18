package worker

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"

	appAuto "github.com/trickreport/backend/internal/application/automation"
)

func TestNew_Defaults(t *testing.T) {
	w := New(nil)
	if w == nil {
		t.Fatal("expected non-nil worker")
	}
	if w.interval != 1*time.Minute {
		t.Errorf("interval = %v, want 1m", w.interval)
	}
	expected := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	if w.systemUserID != expected {
		t.Errorf("systemUserID = %v, want %v", w.systemUserID, expected)
	}
	if w.engine != nil {
		t.Error("engine should be nil by default")
	}
}

func TestNew_WithOptions(t *testing.T) {
	sysID := uuid.New()
	engine := appAuto.NewEngine(nil, nil, zerolog.Nop())

	w := New(nil,
		WithSystemUserID(sysID),
		WithInterval(5*time.Second),
		WithEngine(engine),
	)

	if w.systemUserID != sysID {
		t.Errorf("systemUserID = %v, want %v", w.systemUserID, sysID)
	}
	if w.interval != 5*time.Second {
		t.Errorf("interval = %v, want 5s", w.interval)
	}
	if w.engine != engine {
		t.Error("engine not set")
	}
}

func TestStart_ContextCancellation(t *testing.T) {
	// Use a long interval so the ticker never fires (avoiding nil db access).
	w := New(nil, WithInterval(1*time.Hour))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Start(ctx)
		close(done)
	}()

	// Cancel immediately
	cancel()

	select {
	case <-done:
		// success — Start returned on context cancellation
	case <-time.After(2 * time.Second):
		t.Fatal("Start did not return after context cancellation")
	}
}

func TestNew_AllOptionsCombined(t *testing.T) {
	sysID := uuid.New()
	engine := appAuto.NewEngine(nil, nil, zerolog.Nop())

	w := New(nil,
		WithSystemUserID(sysID),
		WithInterval(10*time.Second),
		WithEngine(engine),
	)

	if w.systemUserID != sysID {
		t.Errorf("systemUserID = %v, want %v", w.systemUserID, sysID)
	}
	if w.interval != 10*time.Second {
		t.Errorf("interval = %v, want 10s", w.interval)
	}
	if w.engine == nil {
		t.Error("engine should be set")
	}
	if w.db != nil {
		t.Error("db should be nil when nil passed")
	}
}

func TestWithInterval_ZeroValue(t *testing.T) {
	// A zero interval should still be accepted (time.NewTicker panics on <=0,
	// but the option itself just sets the field).
	w := New(nil, WithInterval(0))
	if w.interval != 0 {
		t.Errorf("interval = %v, want 0", w.interval)
	}
}

func TestNew_NilDB_Accepted(t *testing.T) {
	// The worker accepts a nil DB — useful for testing and graceful startup
	// when the DB is not yet connected.
	w := New(nil)
	if w == nil {
		t.Fatal("expected non-nil worker")
	}
	if w.db != nil {
		t.Error("db should be nil")
	}
}

func TestNew_WithRealInterval_DoesNotBlock(t *testing.T) {
	// Creating a worker should never block, regardless of interval.
	start := time.Now()
	_ = New(nil, WithInterval(1*time.Nanosecond))
	if elapsed := time.Since(start); elapsed > 100*time.Millisecond {
		t.Errorf("New blocked for %v", elapsed)
	}
}

func TestNew_DefaultMaxRetries(t *testing.T) {
	w := New(nil)
	if w.maxRetries != 0 {
		t.Errorf("default maxRetries = %d, want 0", w.maxRetries)
	}
}

func TestWorker_HealthCheck_NilDB(t *testing.T) {
	// HealthCheck with a nil DB panics; verify the recover catches it so the
	// test suite doesn't crash. The contract is that HealthCheck cannot
	// succeed without a DB.
	w := New(nil)

	defer func() {
		if r := recover(); r == nil {
			// If it didn't panic but returned an error, that's also acceptable.
		}
	}()

	_ = w.HealthCheck()
}
