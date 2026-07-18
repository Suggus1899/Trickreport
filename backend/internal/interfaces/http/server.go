package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"golang.org/x/time/rate"
	"github.com/rs/zerolog/log"

	appAuth "github.com/trickreport/backend/internal/application/auth"
	"github.com/trickreport/backend/internal/config"
	"github.com/trickreport/backend/internal/email"
	infraAuth "github.com/trickreport/backend/internal/infrastructure/auth"
	infraRealtime "github.com/trickreport/backend/internal/infrastructure/realtime"
	httpMiddleware "github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/realtime"
	"github.com/trickreport/backend/internal/worker"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server is the HTTP server with all dependencies wired.
type Server struct {
	cfg      *config.Config
	pool     *pgxpool.Pool
	router   *chi.Mux
	workerCtx context.Context
}

// New creates a new Server with all dependencies wired.
func New(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) *Server {
	// ── Realtime ──────────────────────────────────────────────────────
	hub := realtime.NewHub()
	go hub.Run()

	// Configure WebSocket origin whitelist from CORS config
	realtime.ConfigureAllowedOrigins(cfg.AllowedOrigins())

	// ── Email ─────────────────────────────────────────────────────────
	sender := email.NewConsoleSender()

	// ── Worker ────────────────────────────────────────────────────────
	wrk := worker.New(pool)
	go wrk.Start(ctx)

	// ── Infrastructure: Repositories ──────────────────────────────────
	repos := NewRepos(pool)

	// ── Infrastructure: Adapters ──────────────────────────────────────
	hasher := infraAuth.NewBcryptHasher()
	tokenGen := infraAuth.NewJWTGenerator(cfg.JWTSecret, cfg.JWTExpHours)
	hubAdapter := infraRealtime.NewHubAdapter(hub)

	// LDAP (optional)
	var ldapAuth appAuth.LDAPAuthenticator
	if cfg.LDAPEnabled {
		ldapAuth = infraAuth.NewLDAPAuthenticator(cfg)
	}

	// ── Application: Services ─────────────────────────────────────────
	authCfg := appAuth.Config{
		JWTSecret:    cfg.JWTSecret,
		JWTExpHours:  cfg.JWTExpHours,
		SecureCookie: cfg.SecureCookie,
	}
	services := NewServices(repos, hasher, tokenGen, hubAdapter, sender, ldapAuth, authCfg)

	// ── Interface: Handlers ───────────────────────────────────────────
	handlers := NewHandlers(services)
	authHandler := handlers.Auth
	userHandler := handlers.User
	ticketHandler := handlers.Ticket
	articleHandler := handlers.Article
	slaHandler := handlers.SLA
	automationHandler := handlers.Automation
	analyticsHandler := handlers.Analytics

	// ── Router ────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(30 * time.Second))

	// Rate limiting
	globalLimiter := httpMiddleware.NewRateLimiter(rate.Limit(100), 200, 10*time.Minute)
	globalLimiter.Start()
	r.Use(globalLimiter.LimitByIP)

	// Security headers
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			h := w.Header()
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			h.Set("X-XSS-Protection", "1; mode=block")
			if cfg.IsProduction() {
				h.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
				h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'")
			}
			next.ServeHTTP(w, req)
		})
	})

	// CORS — whitelist of allowed origins (comma-separated in CORS_ORIGINS)
	allowedOrigins := cfg.AllowedOrigins()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			origin := strings.ToLower(strings.TrimSpace(req.Header.Get("Origin")))
			allowed := false
			for _, o := range allowedOrigins {
				if o == origin {
					allowed = true
					break
				}
			}
			if allowed && origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Vary", "Origin")
			}
			if req.Method == http.MethodOptions {
				if allowed {
					w.WriteHeader(http.StatusNoContent)
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
				return
			}
			next.ServeHTTP(w, req)
		})
	})

	// Root redirect — send browser to the Astro frontend (first allowed origin)
	// or fall back to a JSON health message when no origin is configured.
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if len(allowedOrigins) > 0 {
			http.Redirect(w, r, allowedOrigins[0], http.StatusSeeOther)
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// ── API routes (JSON) ─────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		// Public: auth (no tenant required)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/logout", authHandler.Logout)
			r.With(httpMiddleware.Authenticate(tokenGen)).Get("/me", authHandler.Me)
		})

		// Protected routes (auth + tenant required)
		r.Group(func(r chi.Router) {
			r.Use(httpMiddleware.Authenticate(tokenGen))
			r.Use(httpMiddleware.Tenant(repos.TenantResolver))

			// WebSockets
			r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
				tenantID := httpMiddleware.TenantFromContext(r.Context())
				claims, _ := httpMiddleware.ClaimsFromContext(r.Context())
				realtime.ServeWs(hub, w, r, tenantID, claims.UserID)
			})

			// Admin namespace
			r.Route("/admin", func(r chi.Router) {
				r.Use(httpMiddleware.RequireRole("admin"))

				r.Route("/users", func(r chi.Router) {
					r.Get("/", userHandler.List)
					r.Get("/{id}", userHandler.Get)
					r.Post("/", userHandler.Create)
					r.Put("/{id}", userHandler.Update)
					r.Delete("/{id}", userHandler.Delete)
				})

				r.Route("/sla", func(r chi.Router) {
					r.Get("/", slaHandler.List)
					r.Put("/{priority}", slaHandler.Upsert)
				})

				r.Route("/automations", func(r chi.Router) {
					r.Get("/", automationHandler.List)
					r.Post("/", automationHandler.Create)
					r.Put("/{id}", automationHandler.Update)
					r.Delete("/{id}", automationHandler.Delete)
				})

				r.Route("/analytics", func(r chi.Router) {
					r.Get("/summary", analyticsHandler.GetSummary)
					r.Get("/volume", analyticsHandler.GetVolume)
					r.Get("/status-distribution", analyticsHandler.GetStatusDistribution)
					r.Get("/resolution-time", analyticsHandler.GetResolutionTime)
				})
			})

			// Dashboard summary available to any authenticated user
			r.Get("/dashboard/summary", analyticsHandler.GetSummary)

			// Agents: read-only user list for assignee selection
			r.With(httpMiddleware.RequireRole("agent")).Get("/users", userHandler.List)

			// Tickets
			r.Route("/tickets", func(r chi.Router) {
				r.Get("/", ticketHandler.List)
				r.Post("/", ticketHandler.Create)
				r.Get("/{id}", ticketHandler.Get)
				r.Patch("/{id}/status", ticketHandler.UpdateStatus)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/assign", ticketHandler.Assign)

				r.Get("/{id}/comments", ticketHandler.ListComments)
				r.Post("/{id}/comments", ticketHandler.AddComment)

				r.With(httpMiddleware.RequireRole("admin", "agent")).Get("/{id}/history", ticketHandler.ListHistory)
			})

			// Articles (Knowledge Base)
			r.Route("/articles", func(r chi.Router) {
				r.Get("/", articleHandler.List)
				r.Get("/{id}", articleHandler.Get)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", articleHandler.Create)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Put("/{id}", articleHandler.Update)
				r.With(httpMiddleware.RequireRole("admin")).Delete("/{id}", articleHandler.Delete)
			})
		})
	})

	// Health check — includes database connectivity
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			log.Error().Err(err).Msg("health check failed")
			response.Error(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return &Server{cfg: cfg, pool: pool, router: r, workerCtx: ctx}
}

// Start runs the HTTP server with graceful shutdown.
// It blocks until the context is cancelled.
func (s *Server) Start() {
	srv := &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      s.router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("Server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	<-s.workerCtx.Done()
	log.Info().Msg("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Graceful shutdown failed")
	}
	log.Info().Msg("Server stopped")
}
