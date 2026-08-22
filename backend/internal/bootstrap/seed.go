// Seed populates the database with demo data for local development: a
// handful of users across the three roles, a realistic spread of tickets
// (for the dashboard and analytics charts to have something to show),
// comments, history entries, knowledge base articles, SLA policies for
// every priority, and one automation rule.
package bootstrap

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

// SeedPassword is the shared password for every demo user created by Seed.
const SeedPassword = "Demo1234!"

type seedUserSpec struct {
	name  string
	email string
	role  string
}

var seedUserSpecs = []seedUserSpec{
	{"Alex Agent", "alex.agent@demo.trickreport.local", "agent"},
	{"Blair Agent", "blair.agent@demo.trickreport.local", "agent"},
	{"Casey Customer", "casey.customer@demo.trickreport.local", "end_user"},
	{"Dana Customer", "dana.customer@demo.trickreport.local", "end_user"},
}

var ticketTemplates = []struct {
	title       string
	description string
	category    string
}{
	{"Printer on floor 2 won't connect", "The printer keeps dropping off the network every few minutes.", "hardware"},
	{"Cannot reset VPN password", "The self-service reset link times out before I can set a new password.", "access"},
	{"Email client crashes on startup", "The mail client crashes immediately after the last update.", "software"},
	{"Request access to shared drive", "Need read/write access to the Marketing shared drive.", "access"},
	{"Laptop battery draining fast", "Battery goes from 100% to 20% in about an hour.", "hardware"},
	{"Invoice export missing rows", "The monthly invoice CSV export is missing the last 5 rows.", "software"},
	{"New hire onboarding accounts", "Need accounts provisioned for 3 new hires starting Monday.", "access"},
	{"Conference room screen not working", "HDMI input on the 3rd floor conference room shows no signal.", "hardware"},
	{"Slow VPN connection", "VPN throughput dropped to under 1Mbps since yesterday.", "network"},
	{"Password policy question", "Does the new password policy apply to service accounts too?", "general"},
	{"Duplicate customer records", "Two records exist for the same customer in the CRM.", "software"},
	{"Mobile app login loop", "The mobile app keeps returning to the login screen after signing in.", "software"},
	{"Request additional monitor", "Would like a second monitor for the new desk setup.", "hardware"},
	{"Firewall blocking vendor API", "Outbound calls to the vendor API are being blocked.", "network"},
	{"Backup job failed overnight", "The nightly backup job for the file server failed with a timeout.", "hardware"},
	{"Update billing address", "Need to update the billing address on the account.", "general"},
	{"SSO login fails intermittently", "About 1 in 10 SSO logins fail with a generic error.", "access"},
	{"Report export takes too long", "Generating the quarterly report now takes over 10 minutes.", "software"},
}

var seedStatuses = []string{"open", "in_progress", "waiting_client", "resolved", "closed"}
var seedPriorities = []string{"low", "medium", "high", "critical"}

var commentTemplates = []string{
	"Thanks for reporting this — looking into it now.",
	"Can you confirm which device/browser this happens on?",
	"This should be resolved now, please confirm on your end.",
	"Escalating this to the infrastructure team.",
	"Following up — any update from your side?",
}

// Seed populates demo data. Safe to run against a database that already has
// data: users/SLA policies/articles/the automation rule are upserted by
// their natural key; the bulk ticket data (which has no natural key worth
// upserting on) is only generated once — if the tenant already has tickets,
// that step is skipped so re-running Seed doesn't pile up duplicates.
func Seed(ctx context.Context, pool *pgxpool.Pool, isProduction bool) error {
	if isProduction {
		return fmt.Errorf("seed: refusing to run with ENV=production")
	}

	userIDs, err := seedDemoUsers(ctx, pool)
	if err != nil {
		return fmt.Errorf("seed: users: %w", err)
	}

	if err := seedSLAPolicies(ctx, pool); err != nil {
		return fmt.Errorf("seed: sla policies: %w", err)
	}

	if err := seedArticles(ctx, pool, userIDs); err != nil {
		return fmt.Errorf("seed: articles: %w", err)
	}

	if err := seedAutomationRule(ctx, pool); err != nil {
		return fmt.Errorf("seed: automation rule: %w", err)
	}

	seededTickets, err := seedTickets(ctx, pool, userIDs)
	if err != nil {
		return fmt.Errorf("seed: tickets: %w", err)
	}

	log.Info().
		Int("users", len(userIDs)).
		Bool("tickets_seeded", seededTickets).
		Str("demo_password", SeedPassword).
		Msg("seed: demo data ready")
	return nil
}

func seedDemoUsers(ctx context.Context, pool *pgxpool.Pool) (map[string]uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(SeedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash demo password: %w", err)
	}

	ids := make(map[string]uuid.UUID, len(seedUserSpecs))
	for _, u := range seedUserSpecs {
		var id uuid.UUID
		err := pool.QueryRow(ctx, `
			INSERT INTO users (tenant_id, name, email, role, password, active)
			VALUES ($1, $2, $3, $4, $5, TRUE)
			ON CONFLICT (tenant_id, email) DO UPDATE SET name = EXCLUDED.name
			RETURNING id`,
			defaultTenantID, u.name, u.email, u.role, string(hash),
		).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("upsert user %s: %w", u.email, err)
		}
		ids[u.email] = id
	}
	return ids, nil
}

func seedSLAPolicies(ctx context.Context, pool *pgxpool.Pool) error {
	policies := []struct {
		priority                                  string
		responseMin, resolutionMin, escalationMin int
	}{
		{"low", 480, 4320, 1440},
		{"medium", 240, 1440, 480},
		{"high", 60, 480, 120},
		{"critical", 15, 240, 30},
	}
	for _, p := range policies {
		_, err := pool.Exec(ctx, `
			INSERT INTO sla_policies (tenant_id, priority, response_time_minutes, resolution_time_minutes, escalation_minutes)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (tenant_id, priority) DO UPDATE SET
				response_time_minutes = EXCLUDED.response_time_minutes,
				resolution_time_minutes = EXCLUDED.resolution_time_minutes,
				escalation_minutes = EXCLUDED.escalation_minutes`,
			defaultTenantID, p.priority, p.responseMin, p.resolutionMin, p.escalationMin,
		)
		if err != nil {
			return fmt.Errorf("upsert sla policy %s: %w", p.priority, err)
		}
	}
	return nil
}

func seedArticles(ctx context.Context, pool *pgxpool.Pool, userIDs map[string]uuid.UUID) error {
	author := userIDs["alex.agent@demo.trickreport.local"]
	articles := []struct {
		title, content, category string
		tags                     []string
	}{
		{
			"How to reset your VPN password",
			"Go to the self-service portal, click \"Forgot password\", and follow the emailed link. Links expire after 15 minutes.",
			"access", []string{"vpn", "password"},
		},
		{
			"Requesting access to a shared drive",
			"Submit a ticket with the drive name and the access level you need (read or read/write). Your manager will be asked to approve it.",
			"access", []string{"access", "shared-drive"},
		},
		{
			"Setting up email on your phone",
			"Use the Exchange/ActiveSync profile with your full email address as the username. MFA must be enabled first.",
			"software", []string{"email", "mobile"},
		},
	}
	for _, a := range articles {
		var exists bool
		if err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM articles WHERE tenant_id = $1 AND title = $2)`, defaultTenantID, a.title).Scan(&exists); err != nil {
			return fmt.Errorf("check article %q: %w", a.title, err)
		}
		if exists {
			continue
		}
		_, err := pool.Exec(ctx, `
			INSERT INTO articles (tenant_id, title, content, category, tags, published, created_by)
			VALUES ($1, $2, $3, $4, $5, TRUE, $6)`,
			defaultTenantID, a.title, a.content, a.category, a.tags, author,
		)
		if err != nil {
			return fmt.Errorf("insert article %q: %w", a.title, err)
		}
	}
	return nil
}

func seedAutomationRule(ctx context.Context, pool *pgxpool.Pool) error {
	const name = "Escalate critical tickets"
	var exists bool
	if err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM automation_rules WHERE tenant_id = $1 AND name = $2)`, defaultTenantID, name).Scan(&exists); err != nil {
		return fmt.Errorf("check automation rule: %w", err)
	}
	if exists {
		return nil
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO automation_rules (tenant_id, name, description, trigger_type, conditions, actions, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, TRUE)`,
		defaultTenantID,
		name,
		"Automatically flags newly created critical tickets for immediate attention.",
		"ticket_created",
		`{"priority": "critical"}`,
		`[{"type": "change_status", "status": "in_progress"}]`,
	)
	if err != nil {
		return fmt.Errorf("insert automation rule: %w", err)
	}
	return nil
}

// seedTickets returns whether it actually inserted new tickets (false when
// skipped because the tenant already has some).
func seedTickets(ctx context.Context, pool *pgxpool.Pool, userIDs map[string]uuid.UUID) (bool, error) {
	var count int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM tickets WHERE tenant_id = $1`, defaultTenantID).Scan(&count); err != nil {
		return false, fmt.Errorf("count tickets: %w", err)
	}
	if count > 0 {
		return false, nil
	}

	creators := []uuid.UUID{userIDs["casey.customer@demo.trickreport.local"], userIDs["dana.customer@demo.trickreport.local"]}
	assignees := []uuid.UUID{userIDs["alex.agent@demo.trickreport.local"], userIDs["blair.agent@demo.trickreport.local"]}

	now := time.Now()
	for i, tpl := range ticketTemplates {
		status := seedStatuses[i%len(seedStatuses)]
		priority := seedPriorities[i%len(seedPriorities)]
		creator := creators[i%len(creators)]
		createdAt := now.Add(-time.Duration(rand.Intn(28*24)+1) * time.Hour)

		var assignedTo *uuid.UUID
		var updatedAt time.Time = createdAt
		if status != "open" {
			a := assignees[i%len(assignees)]
			assignedTo = &a
		}
		if status == "resolved" || status == "closed" {
			updatedAt = createdAt.Add(time.Duration(rand.Intn(70)+2) * time.Hour)
		}
		slaBreached := priority == "critical" && i%3 == 0

		var ticketID uuid.UUID
		err := pool.QueryRow(ctx, `
			INSERT INTO tickets (tenant_id, title, description, status, priority, category, created_by, assigned_to, sla_breached, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id`,
			defaultTenantID, tpl.title, tpl.description, status, priority, tpl.category,
			creator, assignedTo, slaBreached, createdAt, updatedAt,
		).Scan(&ticketID)
		if err != nil {
			return false, fmt.Errorf("insert ticket %q: %w", tpl.title, err)
		}

		if status != "open" {
			commenter := creator
			if assignedTo != nil {
				commenter = *assignedTo
			}
			_, err = pool.Exec(ctx, `
				INSERT INTO ticket_comments (ticket_id, user_id, content, is_internal, created_at)
				VALUES ($1, $2, $3, FALSE, $4)`,
				ticketID, commenter, commentTemplates[i%len(commentTemplates)], createdAt.Add(time.Hour),
			)
			if err != nil {
				return false, fmt.Errorf("insert comment for ticket %q: %w", tpl.title, err)
			}

			_, err = pool.Exec(ctx, `
				INSERT INTO ticket_history (ticket_id, user_id, field, old_value, new_value, created_at)
				VALUES ($1, $2, 'status', 'open', $3, $4)`,
				ticketID, commenter, status, updatedAt,
			)
			if err != nil {
				return false, fmt.Errorf("insert history for ticket %q: %w", tpl.title, err)
			}
		}
	}

	return true, nil
}
