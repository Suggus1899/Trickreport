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
