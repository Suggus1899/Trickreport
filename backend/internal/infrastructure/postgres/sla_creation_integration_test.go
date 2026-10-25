package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	appMSP "github.com/trickreport/backend/internal/application/msp"
	domain "github.com/trickreport/backend/internal/domain/msp"
	domainSLA "github.com/trickreport/backend/internal/domain/sla"
)

// slaSlack absorbs the microsecond precision of stored timestamps.
const slaSlack = time.Second

// TestIntegration_MSPTicketGetsContractDeadlines checks that a ticket opened for
// a client with an active contract stores the contract's clocks, so the SLA
// worker can flag it later. A client without a contract gets no clocks.
func TestIntegration_MSPTicketGetsContractDeadlines(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenantID, adminID := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenantID)

	repo := NewMSPRepo(testPool)
	svc := appMSP.NewService(repo)
	covered := seedSLAClient(t, ctx, repo, tenantID, "covered")
	uncovered := seedSLAClient(t, ctx, repo, tenantID, "uncovered")

	if _, err := svc.CreateContract(ctx, tenantID, &domain.Contract{
		TenantID: tenantID, ClientID: covered.ID, Name: "Gold", Status: "active",
		StartsAt: time.Now().Add(-time.Hour), Coverage: "business", Currency: "EUR",
		SLAResponseMinutes: 60, SLAResolutionMinutes: 240,
	}); err != nil {
		t.Fatalf("CreateContract: %v", err)
	}

	input := appMSP.MSPTicketInput{ClientID: covered.ID, Title: "Printer down", Priority: "high"}
	before := time.Now()
	ticketID, err := svc.CreateMSPTicket(ctx, tenantID, adminID, "admin", input)
	after := time.Now()
	if err != nil {
		t.Fatalf("CreateMSPTicket covered: %v", err)
	}
	got := readTicketDeadlines(t, ctx, ticketID)
	assertDeadlineNear(t, "response", got.Response, before.Add(60*time.Minute), after.Add(60*time.Minute))
	assertDeadlineNear(t, "resolution", got.Resolution, before.Add(240*time.Minute), after.Add(240*time.Minute))
	if got.Escalation != nil {
		t.Fatalf("escalation deadline = %s, want none: the contract sets no escalation", got.Escalation)
	}

	input.ClientID = uncovered.ID
	ticketID, err = svc.CreateMSPTicket(ctx, tenantID, adminID, "admin", input)
	if err != nil {
		t.Fatalf("CreateMSPTicket uncovered: %v", err)
	}
	got = readTicketDeadlines(t, ctx, ticketID)
	if got.Response != nil || got.Resolution != nil || got.Escalation != nil {
		t.Fatalf("uncovered deadlines = %+v, want none", got)
	}
}

// TestIntegration_MSPTicketCreationIsAtomic checks that a failing context row
// rolls back the ticket row written in the same transaction, so no ticket is
// left without its MSP context.
func TestIntegration_MSPTicketCreationIsAtomic(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenantID, adminID := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenantID)

	title := "orphan context " + tenantID.String()[:8]
	// The unknown client makes the context insert fail after the ticket insert.
	_, err := NewMSPRepo(testPool).CreateTicketWithContext(ctx, tenantID, adminID, title, "", "general", "high",
		&domain.TicketContext{TenantID: tenantID, ClientID: uuid.New()}, domainSLA.Deadlines{})
	if err == nil {
		t.Fatal("CreateTicketWithContext with an unknown client: want an error")
	}

	var n int
	if err := testPool.QueryRow(ctx,
		`SELECT COUNT(*) FROM tickets WHERE tenant_id = $1 AND title = $2`, tenantID, title,
	).Scan(&n); err != nil {
		t.Fatalf("count tickets: %v", err)
	}
	if n != 0 {
		t.Fatalf("tickets left behind = %d, want 0", n)
	}
}

// TestIntegration_ITSMCatalogTicketPersistsDeadlines checks that a catalog
// request stores the deadlines it was given instead of leaving the ticket
// without clocks.
func TestIntegration_ITSMCatalogTicketPersistsDeadlines(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenantID, adminID := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenantID)

	created := time.Now().UTC().Truncate(time.Second)
	response := created.Add(30 * time.Minute)
	resolution := created.Add(4 * time.Hour)
	want := domainSLA.Deadlines{Response: &response, Resolution: &resolution}

	ticketID, err := NewITSMRepo(testPool).CreateTicketFromService(ctx, tenantID, adminID,
		"VPN access", "", "Access", "high", want)
	if err != nil {
		t.Fatalf("CreateTicketFromService: %v", err)
	}
	got := readTicketDeadlines(t, ctx, ticketID)
	assertSameDeadline(t, "response", got.Response, want.Response)
	assertSameDeadline(t, "resolution", got.Resolution, want.Resolution)
	assertSameDeadline(t, "escalation", got.Escalation, nil)
}

func seedSLAClient(t *testing.T, ctx context.Context, repo *MSPRepo, tenantID uuid.UUID, label string) *domain.Client {
	t.Helper()
	c := &domain.Client{
		TenantID: tenantID, Name: label, Status: "active",
		Domain: label + "-" + tenantID.String()[:8] + ".test.local",
	}
	if err := repo.CreateClient(ctx, c); err != nil {
		t.Fatalf("CreateClient %s: %v", label, err)
	}
	return c
}

// readTicketDeadlines loads the three SLA clocks stored on a ticket. NULL reads as nil.
func readTicketDeadlines(t *testing.T, ctx context.Context, ticketID uuid.UUID) domainSLA.Deadlines {
	t.Helper()
	var dl domainSLA.Deadlines
	if err := testPool.QueryRow(ctx,
		`SELECT response_deadline, sla_deadline, escalation_deadline FROM tickets WHERE id = $1`, ticketID,
	).Scan(&dl.Response, &dl.Resolution, &dl.Escalation); err != nil {
		t.Fatalf("read deadlines of ticket %s: %v", ticketID, err)
	}
	return dl
}

func assertDeadlineNear(t *testing.T, name string, got *time.Time, lo, hi time.Time) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s deadline is nil, want one between %s and %s", name, lo, hi)
	}
	if got.Before(lo.Add(-slaSlack)) || got.After(hi.Add(slaSlack)) {
		t.Fatalf("%s deadline = %s, want between %s and %s", name, got, lo, hi)
	}
}

func assertSameDeadline(t *testing.T, name string, got, want *time.Time) {
	t.Helper()
	if (got == nil) != (want == nil) || (got != nil && !got.Equal(*want)) {
		t.Fatalf("%s deadline = %v, want %v", name, got, want)
	}
}
