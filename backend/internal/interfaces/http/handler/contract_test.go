package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	appAuth "github.com/trickreport/backend/internal/application/auth"
	"github.com/trickreport/backend/internal/db"
	"github.com/trickreport/backend/internal/email"
	infraAuth "github.com/trickreport/backend/internal/infrastructure/auth"
	infraRealtime "github.com/trickreport/backend/internal/infrastructure/realtime"
	httpServer "github.com/trickreport/backend/internal/interfaces/http"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/realtime"
)

// contractPool is the shared DB pool for contract tests.
var contractPool *pgxpool.Pool

// TestMain is not defined here (the handler package already has tests). We
// initialize the pool lazily in a sync.Once-free helper since contract tests
// are skipped when TEST_DATABASE_URL is unset.
func contractSetup(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if contractPool == nil {
		dsn := os.Getenv("TEST_DATABASE_URL")
		if dsn == "" {
			t.Skip("TEST_DATABASE_URL not set — skipping contract test")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := db.RunMigrations(dsn); err != nil {
			t.Fatalf("migrations: %v", err)
		}
		pool, err := db.New(ctx, dsn)
		if err != nil {
			t.Fatalf("connect: %v", err)
		}
		contractPool = pool
	}
	return contractPool
}

// contractTenant creates a fresh tenant + admin user with a known password and
// returns the tenant ID, admin ID, and the plain-text password.
func contractTenant(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (uuid.UUID, uuid.UUID, string) {
	t.Helper()
	tenantID := uuid.New()
	_, err := pool.Exec(ctx,
		`INSERT INTO tenants (id, name, slug, settings) VALUES ($1, $2, $3, '{}')`,
		tenantID, "Contract Tenant "+tenantID.String()[:8], "c-"+tenantID.String()[:8],
	)
	if err != nil {
		t.Fatalf("seed tenant: %v", err)
	}

	password := "contract-pass-123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}

	userID := uuid.New()
	_, err = pool.Exec(ctx,
		`INSERT INTO users (id, tenant_id, name, email, role, password) VALUES ($1, $2, $3, $4, 'admin', $5)`,
		userID, tenantID, "Contract Admin", "contract-"+tenantID.String()[:8]+"@test.local", string(hash),
	)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return tenantID, userID, password
}

// buildContractServer wires the actual server components (repos, services,
// handlers) and returns a test server plus the admin credentials and tenant ID.
func buildContractServer(t *testing.T) (*httptest.Server, string, string, uuid.UUID) {
	ctx := context.Background()
	pool := contractSetup(t)
	tenantID, _, password := contractTenant(t, ctx, pool)
	adminEmail := "contract-" + tenantID.String()[:8] + "@test.local"

	// Wire using the same components as http.New.
	repos := httpServer.NewRepos(pool)
	hasher := infraAuth.NewBcryptHasher()
	tokenGen := infraAuth.NewJWTGenerator("contract-test-secret", 8)

	hub := realtime.NewHub()
	go hub.Run()
	hubAdapter := infraRealtime.NewHubAdapter(hub)

	sender := email.NewConsoleSender()
	authCfg := appAuth.Config{JWTSecret: "contract-test-secret", JWTExpHours: 8, SecureCookie: false}
	services := httpServer.NewServices(repos, hasher, tokenGen, hubAdapter, sender, nil, authCfg)
	handlers := httpServer.NewHandlers(services)

	// Build a router that mirrors the /api/v1 routes from server.go.
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", handlers.Auth.Login)
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.Authenticate(tokenGen))
			r.Use(middleware.Tenant(repos.TenantResolver))

			r.Route("/tickets", func(r chi.Router) {
				r.Get("/", handlers.Ticket.List)
				r.Post("/", handlers.Ticket.Create)
			})
			r.Route("/articles", func(r chi.Router) {
				r.Get("/", handlers.Article.List)
			})
			r.Route("/admin", func(r chi.Router) {
				r.Use(middleware.RequireRole("admin"))
				r.Route("/analytics", func(r chi.Router) {
					r.Get("/charts", handlers.Analytics.GetCharts)
				})
			})
		})
	})

	ts := httptest.NewServer(r)
	t.Cleanup(func() {
		ts.Close()
		_, _ = pool.Exec(ctx, `DELETE FROM tenants WHERE id = $1`, tenantID)
	})
	return ts, adminEmail, password, tenantID
}

// loginAndGetToken performs a login and returns the JWT token.
func loginAndGetToken(t *testing.T, ts *httptest.Server, email, password string) string {
	t.Helper()
	body := `{"email":"` + email + `","password":"` + password + `"}`
	resp, err := http.Post(ts.URL+"/api/v1/auth/login", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("login request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("login status = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, raw)
	}
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	token, ok := result["token"].(string)
	if !ok {
		t.Fatal("expected token in login response")
	}
	return token
}

// authedRequest performs an authenticated request with the tenant header.
func authedRequest(t *testing.T, method, url, token, tenantID, body string) *http.Response {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	return resp
}

func TestContract_AuthLogin_Returns200WithToken(t *testing.T) {
	ts, email, password, _ := buildContractServer(t)
	token := loginAndGetToken(t, ts, email, password)
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestContract_GetTickets_Returns200WithArray(t *testing.T) {
	ts, email, password, tenantID := buildContractServer(t)
	token := loginAndGetToken(t, ts, email, password)

	resp := authedRequest(t, http.MethodGet, ts.URL+"/api/v1/tickets", token, tenantID.String(), "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("GET /tickets status = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, raw)
	}
	var arr []any
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		t.Fatalf("expected JSON array, got: %v", err)
	}
}

func TestContract_CreateTicket_Returns201WithTicket(t *testing.T) {
	ts, email, password, tenantID := buildContractServer(t)
	token := loginAndGetToken(t, ts, email, password)

	body := `{"title":"Contract test ticket","description":"Created via contract test","priority":"high","category":"bug"}`
	resp := authedRequest(t, http.MethodPost, ts.URL+"/api/v1/tickets", token, tenantID.String(), body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("POST /tickets status = %d, want %d; body=%s", resp.StatusCode, http.StatusCreated, raw)
	}
	var tk map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&tk); err != nil {
		t.Fatalf("expected JSON ticket, got: %v", err)
	}
	if tk["title"] != "Contract test ticket" {
		t.Errorf("title = %v", tk["title"])
	}
}

func TestContract_GetArticles_Returns200WithArray(t *testing.T) {
	ts, email, password, tenantID := buildContractServer(t)
	token := loginAndGetToken(t, ts, email, password)

	resp := authedRequest(t, http.MethodGet, ts.URL+"/api/v1/articles", token, tenantID.String(), "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("GET /articles status = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, raw)
	}
	var arr []any
	if err := json.NewDecoder(resp.Body).Decode(&arr); err != nil {
		t.Fatalf("expected JSON array, got: %v", err)
	}
}

func TestContract_GetAnalyticsCharts_Returns200WithData(t *testing.T) {
	ts, email, password, tenantID := buildContractServer(t)
	token := loginAndGetToken(t, ts, email, password)

	resp := authedRequest(t, http.MethodGet, ts.URL+"/api/v1/admin/analytics/charts?range=7d", token, tenantID.String(), "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		t.Fatalf("GET /analytics/charts status = %d, want %d; body=%s", resp.StatusCode, http.StatusOK, raw)
	}
	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("expected JSON object, got: %v", err)
	}
	// Charts response should contain volume, status, and resolution keys.
	if data["volume"] == nil {
		t.Error("expected 'volume' key in charts response")
	}
}
