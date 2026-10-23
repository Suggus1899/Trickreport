package worker

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/trickreport/backend/internal/db"
	"os"
	"strings"
	"sync"
	"testing"
)

func slaTestDB(t *testing.T) (*pgxpool.Pool, uuid.UUID, uuid.UUID) {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	if err := db.RunMigrations(dsn); err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatal(err)
	}
	tenant, actor := uuid.New(), uuid.New()
	ctx := context.Background()
	if _, err := pool.Exec(ctx, `INSERT INTO tenants(id,name,slug)VALUES($1,'SLA test',$2)`, tenant, tenant.String()); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO users(id,tenant_id,name,email,role)VALUES($1,$2,'SLA actor',$3,'admin')`, actor, tenant, actor.String()+"@test.local"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DELETE FROM ticket_history WHERE user_id=$1`, actor)
		_, _ = pool.Exec(ctx, `DELETE FROM ticket_comments WHERE user_id=$1`, actor)
		_, _ = pool.Exec(ctx, `DELETE FROM tickets WHERE tenant_id=$1`, tenant)
		_, _ = pool.Exec(ctx, `DELETE FROM tenants WHERE id=$1`, tenant)
		pool.Close()
	})
	return pool, tenant, actor
}

func TestSLAWorker_PersistedDeadlinesConcurrentAndAtomic(t *testing.T) {
	pool, tenant, actor := slaTestDB(t)
	ctx := context.Background()
	ticket := uuid.New()
	if _, err := pool.Exec(ctx, `INSERT INTO sla_policies(tenant_id,priority,response_time_minutes,resolution_time_minutes,escalation_minutes)VALUES($1,'high',5000,5000,0)`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO tickets(id,tenant_id,created_by,title,description,priority,sla_deadline,response_deadline,escalation_deadline)VALUES($1,$2,$3,'deadline','test','high',NOW()-INTERVAL '1 minute',NOW()-INTERVAL '1 minute',NOW()-INTERVAL '1 minute')`, ticket, tenant, actor); err != nil {
		t.Fatal(err)
	}
	w1, w2 := New(pool, WithSystemUserID(actor)), New(pool, WithSystemUserID(actor))
	var wg sync.WaitGroup
	errs := make(chan error, 4)
	for _, w := range []*Worker{w1, w2} {
		wg.Go(func() { errs <- w.checkSLABreaches(ctx); errs <- w.checkEscalations(ctx) })
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	var breached, response bool
	var escalated *string
	var count int
	if err := pool.QueryRow(ctx, `SELECT sla_breached,response_sla_breached,escalated_at::text,(SELECT COUNT(*) FROM ticket_history WHERE ticket_id=$1)FROM tickets WHERE id=$1`, ticket).Scan(&breached, &response, &escalated, &count); err != nil {
		t.Fatal(err)
	}
	if !breached || !response || escalated == nil || count != 3 {
		t.Fatalf("resolution=%t response=%t escalated=%v history=%d", breached, response, escalated, count)
	}
	if err := w1.checkSLABreaches(ctx); err != nil {
		t.Fatal(err)
	}
	if err := w1.checkEscalations(ctx); err != nil {
		t.Fatal(err)
	}
	var comments int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM ticket_comments WHERE ticket_id=$1`, ticket).Scan(&comments); err != nil {
		t.Fatal(err)
	}
	if comments != 3 {
		t.Fatalf("duplicate SLA comments=%d", comments)
	}
}

func TestSLAWorker_HistoryFailureRollsBackAndRetry(t *testing.T) {
	pool, tenant, actor := slaTestDB(t)
	ctx := context.Background()
	ticket := uuid.New()
	if _, err := pool.Exec(ctx, `INSERT INTO tickets(id,tenant_id,created_by,title,description,sla_deadline)VALUES($1,$2,$3,'deadline','test',NOW()-INTERVAL '1 minute')`, ticket, tenant, actor); err != nil {
		t.Fatal(err)
	}
	// A scoped trigger rejects only this ticket's history, leaving other concurrent suites independent.
	function := "fail_sla_" + strings.ReplaceAll(ticket.String(), "-", "")
	if _, err := pool.Exec(ctx, `CREATE FUNCTION `+function+`() RETURNS trigger LANGUAGE plpgsql AS $$BEGIN IF NEW.ticket_id='`+ticket.String()+`'::uuid THEN RAISE EXCEPTION 'history unavailable'; END IF; RETURN NEW; END$$; CREATE TRIGGER `+function+` BEFORE INSERT ON ticket_history FOR EACH ROW EXECUTE FUNCTION `+function+`()`); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = pool.Exec(ctx, `DROP TRIGGER IF EXISTS `+function+` ON ticket_history; DROP FUNCTION IF EXISTS `+function+`()`)
	})
	w := New(pool, WithSystemUserID(actor))
	if err := w.checkSLABreaches(ctx); err == nil {
		t.Error("history failure must be returned")
	}
	var marked bool
	var count int
	if err := pool.QueryRow(ctx, `SELECT sla_breached,(SELECT COUNT(*) FROM ticket_comments WHERE ticket_id=$1)FROM tickets WHERE id=$1`, ticket).Scan(&marked, &count); err != nil {
		t.Fatal(err)
	}
	if marked || count != 0 {
		t.Fatalf("partial SLA write mark=%t comments=%d", marked, count)
	}
	if _, err := pool.Exec(ctx, `DROP TRIGGER `+function+` ON ticket_history`); err != nil {
		t.Fatal(err)
	}
	if err := w.checkSLABreaches(ctx); err != nil {
		t.Fatal(err)
	}
	if err := pool.QueryRow(ctx, `SELECT sla_breached FROM tickets WHERE id=$1`, ticket).Scan(&marked); err != nil {
		t.Fatal(err)
	}
	if !marked {
		t.Fatal("retry did not claim ticket")
	}
}

func TestSLAWorker_ClosedAndDisabledEscalation(t *testing.T) {
	pool, tenant, actor := slaTestDB(t)
	ctx := context.Background()
	if _, err := pool.Exec(ctx, `INSERT INTO sla_policies(tenant_id,priority,response_time_minutes,resolution_time_minutes,escalation_minutes)VALUES($1,'low',1,1,0)`, tenant); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO tickets(tenant_id,created_by,title,description,status,created_at,sla_deadline,response_deadline)VALUES($1,$2,'closed','test','closed',NOW()-INTERVAL '2 hours',NOW()-INTERVAL '1 hour',NOW()-INTERVAL '1 hour'),($1,$2,'disabled','test','open',NOW()-INTERVAL '2 hours',NULL,NULL)`, tenant, actor); err != nil {
		t.Fatal(err)
	}
	w := New(pool, WithSystemUserID(actor))
	if err := w.checkSLABreaches(ctx); err != nil {
		t.Fatal(err)
	}
	if err := w.checkEscalations(ctx); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM tickets WHERE tenant_id=$1 AND (sla_breached OR response_sla_breached OR escalated_at IS NOT NULL)`, tenant).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("closed/disabled tickets mutated=%d", count)
	}
}
