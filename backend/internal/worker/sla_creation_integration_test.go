package worker

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	appMSP "github.com/trickreport/backend/internal/application/msp"
	domain "github.com/trickreport/backend/internal/domain/msp"
	"github.com/trickreport/backend/internal/infrastructure/postgres"
)

// TestSLABreachFlagsMSPTicketAfterDeadline covers the path customers depend on:
// a ticket opened under a contract must carry the contract's clocks, and the
// worker must flag it once they pass. Before the MSP creation path stored
// deadlines, the worker had nothing to flag for these tickets.
func TestSLABreachFlagsMSPTicketAfterDeadline(t *testing.T) {
	pool, tenantID, actorID := slaTestDB(t)
	ctx := context.Background()
	repo := postgres.NewMSPRepo(pool)
	svc := appMSP.NewService(repo)

	client := &domain.Client{
		TenantID: tenantID, Name: "Client", Status: "active",
		Domain: "sla-" + tenantID.String()[:8] + ".test.local",
	}
	if err := repo.CreateClient(ctx, client); err != nil {
		t.Fatalf("CreateClient: %v", err)
	}
	if _, err := svc.CreateContract(ctx, tenantID, &domain.Contract{
		TenantID: tenantID, ClientID: client.ID, Name: "Gold", Status: "active",
		StartsAt: time.Now().Add(-time.Hour), Coverage: "business", Currency: "EUR",
		SLAResponseMinutes: 60, SLAResolutionMinutes: 240,
	}); err != nil {
		t.Fatalf("CreateContract: %v", err)
	}

	input := appMSP.MSPTicketInput{ClientID: client.ID, Title: "Server down", Priority: "high"}
	expired, err := svc.CreateMSPTicket(ctx, tenantID, actorID, "admin", input)
	if err != nil {
		t.Fatalf("CreateMSPTicket expired: %v", err)
	}
	fresh, err := svc.CreateMSPTicket(ctx, tenantID, actorID, "admin", input)
	if err != nil {
		t.Fatalf("CreateMSPTicket fresh: %v", err)
	}

	// The creation path must store the contract clocks; otherwise there is nothing to check.
	var stored int
	if err := pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tickets WHERE id IN ($1, $2) AND response_deadline IS NOT NULL AND sla_deadline IS NOT NULL`,
		expired, fresh,
	).Scan(&stored); err != nil {
		t.Fatalf("count stored deadlines: %v", err)
	}
	if stored != 2 {
		t.Fatalf("tickets with stored response and resolution deadlines = %d, want 2", stored)
	}

	// Move only the first ticket past its clocks instead of waiting an hour.
	if _, err := pool.Exec(ctx,
		`UPDATE tickets SET response_deadline = NOW() - INTERVAL '1 minute', sla_deadline = NOW() - INTERVAL '1 minute' WHERE id = $1`,
		expired,
	); err != nil {
		t.Fatalf("expire ticket deadlines: %v", err)
	}

	if err := New(pool, WithSystemUserID(actorID)).checkSLABreaches(ctx); err != nil {
		t.Fatalf("checkSLABreaches: %v", err)
	}

	for _, tc := range []struct {
		name   string
		ticket uuid.UUID
		want   bool
	}{
		{name: "expired", ticket: expired, want: true},
		{name: "fresh", ticket: fresh, want: false},
	} {
		var resolution, response bool
		if err := pool.QueryRow(ctx,
			`SELECT sla_breached, response_sla_breached FROM tickets WHERE id = $1`, tc.ticket,
		).Scan(&resolution, &response); err != nil {
			t.Fatalf("read %s ticket flags: %v", tc.name, err)
		}
		if resolution != tc.want || response != tc.want {
			t.Fatalf("%s ticket: sla_breached=%v response_sla_breached=%v, want %v", tc.name, resolution, response, tc.want)
		}
	}
}
