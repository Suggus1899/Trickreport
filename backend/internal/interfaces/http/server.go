package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"

	"github.com/jackc/pgx/v5/pgxpool"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	"github.com/trickreport/backend/internal/config"
	"github.com/trickreport/backend/internal/email"
	infraAuth "github.com/trickreport/backend/internal/infrastructure/auth"
	infraEmail "github.com/trickreport/backend/internal/infrastructure/email"
	infraRealtime "github.com/trickreport/backend/internal/infrastructure/realtime"
	httpMiddleware "github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/handler"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/realtime"
	"github.com/trickreport/backend/internal/worker"
)

// Server is the HTTP server with all dependencies wired.
type Server struct {
	cfg       *config.Config
	pool      *pgxpool.Pool
	router    *chi.Mux
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
	var sender email.Sender
	if cfg.Email.SMTPHost != "" {
		sender = infraEmail.NewSMTPSender(cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.SMTPUsername, cfg.Email.SMTPPassword, cfg.Email.SMTPFrom, cfg.Email.SMTPUseTLS)
		log.Info().Str("host", cfg.Email.SMTPHost).Int("port", cfg.Email.SMTPPort).Msg("SMTP email sender configured")
	} else {
		sender = email.NewConsoleSender()
		log.Info().Msg("Console email sender configured (no SMTP_HOST set)")
	}

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

	// ── Worker ────────────────────────────────────────────────────────
	wrk := worker.New(pool, worker.WithEngine(services.Engine))
	go wrk.Start(ctx)

	// ── Interface: Handlers ───────────────────────────────────────────
	handlers := NewHandlers(services)
	authHandler := handlers.Auth
	userHandler := handlers.User
	ticketHandler := handlers.Ticket
	articleHandler := handlers.Article
	slaHandler := handlers.SLA
	automationHandler := handlers.Automation
	analyticsHandler := handlers.Analytics
	attachmentHandler := handlers.Attachment
	notificationHandler := handlers.Notification

	// ── Router ────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(httpMiddleware.RequestLogger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(30 * time.Second))

	// Metrics (Prometheus)
	r.Use(httpMiddleware.Metrics)

	// Compression (gzip)
	r.Use(httpMiddleware.Compress())

	// Request body size limit (10MB max)
	r.Use(httpMiddleware.MaxBodySize(10 << 20))

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
			// Login rate limiting (per email+IP)
			loginLimiter := httpMiddleware.NewLoginRateLimiter(rate.Limit(5), 10, 15)
			loginLimiter.Start()

			r.With(loginLimiter.LoginLimit).Post("/login", authHandler.Login)
			r.With(loginLimiter.LoginLimit).Post("/register", authHandler.Register)
			r.Post("/logout", authHandler.Logout)
			r.Post("/refresh", authHandler.Refresh)
			r.Post("/password-reset", authHandler.PasswordReset)
			r.Post("/password-reset/confirm", authHandler.PasswordResetConfirm)
			r.Post("/mfa/login", authHandler.MFALogin)
			r.With(httpMiddleware.Authenticate(tokenGen)).Get("/me", authHandler.Me)

			// MFA management (requires auth)
			r.Group(func(r chi.Router) {
				r.Use(httpMiddleware.Authenticate(tokenGen))
				r.Post("/mfa/setup", authHandler.MFASetup)
				r.Post("/mfa/verify", authHandler.MFAVerify)
				r.Post("/mfa/enable", authHandler.MFAEnable)
				r.Post("/mfa/disable", authHandler.MFADisable)
			})
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
					r.Get("/charts", analyticsHandler.GetCharts)
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

				// Attachments
				r.Post("/{id}/attachments", attachmentHandler.Upload)
				r.Get("/{id}/attachments", attachmentHandler.List)
				r.Get("/{id}/attachments/{aid}", attachmentHandler.Download)

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

		// Notifications (real-time + persistent)
		r.Route("/notifications", func(r chi.Router) {
			r.Get("/", notificationHandler.List)
			r.Get("/unread-count", notificationHandler.UnreadCount)
			r.Post("/{id}/read", notificationHandler.MarkRead)
			r.Post("/read-all", notificationHandler.MarkAllRead)
		})
	})

	// Metrics endpoint (Prometheus)
	r.Handle("/metrics", handler.MetricsHandler())

	// Swagger UI
	r.Handle("/swagger/*", handler.SwaggerHandler())

	// Health check — includes database connectivity
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			httpMiddleware.LoggerFromContext(r.Context()).Error().Err(err).Msg("health check failed")
			response.Error(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Readiness check — verifies DB connectivity and worker availability.
	// The worker runs in the background; if the process is up and the DB is
	// reachable, the worker is considered running.
	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			httpMiddleware.LoggerFromContext(r.Context()).Error().Err(err).Msg("ready check failed")
			response.Error(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
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
