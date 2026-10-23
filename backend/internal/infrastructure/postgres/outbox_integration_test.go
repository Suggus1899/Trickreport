package postgres

import (
	"context"
	infraEmail "github.com/trickreport/backend/internal/infrastructure/email"
	"testing"
)

func TestOutboxSurvivesRestartAndRollsBackWithTransaction(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	store := NewEmailOutboxRepo(testPool)
	queue := infraEmail.NewEmailQueue(store)
	tx, err := testPool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := queue.SendContext(WithTx(ctx, tx), "rollback@example.test", "rollback", "private body"); err != nil {
		t.Fatal(err)
	}
	if err := tx.Rollback(ctx); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := testPool.QueryRow(ctx, `SELECT COUNT(*) FROM email_outbox WHERE recipient='rollback@example.test'`).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("rolled-back business transaction left an email")
	}
	if err := queue.Send("restart@example.test", "restart", "persistent body"); err != nil {
		t.Fatal(err)
	}
	restarted := NewEmailOutboxRepo(testPool)
	msg, err := restarted.Claim(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if msg == nil || msg.Body != "persistent body" {
		t.Fatal("email was lost on queue restart")
	}
	if err := restarted.Finish(ctx, msg, false); err != nil {
		t.Fatal(err)
	}
	if err := testPool.QueryRow(ctx, `SELECT COUNT(*) FROM email_outbox WHERE id=$1 AND status='pending' AND attempts=1`, msg.ID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatal("SMTP failure was not retained for retry")
	}
	t.Cleanup(func() {
		_, _ = testPool.Exec(ctx, `DELETE FROM email_outbox WHERE recipient IN ('restart@example.test','rollback@example.test')`)
	})
}
