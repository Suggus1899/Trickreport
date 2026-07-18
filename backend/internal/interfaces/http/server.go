package http

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
	"github.com/rs/zerolog/log"

	appAnalytics "github.com/trickreport/backend/internal/application/analytics"
	appArticle "github.com/trickreport/backend/internal/application/article"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	appAuto "github.com/trickreport/backend/internal/application/automation"
	appSLA "github.com/trickreport/backend/internal/application/sla"
	appTicket "github.com/trickreport/backend/internal/application/ticket"
	appUser "github.com/trickreport/backend/internal/application/user"
	"github.com/trickreport/backend/internal/config"
	"github.com/trickreport/backend/internal/domain/ticket"
	"github.com/trickreport/backend/internal/email"
	infraAuth "github.com/trickreport/backend/internal/infrastructure/auth"
	"github.com/trickreport/backend/internal/infrastructure/postgres"
	infraRealtime "github.com/trickreport/backend/internal/infrastructure/realtime"
	"github.com/trickreport/backend/internal/interfaces/http/handler"
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

	// ── Email ─────────────────────────────────────────────────────────
	sender := email.NewConsoleSender()

	// ── Worker ────────────────────────────────────────────────────────
	wrk := worker.New(pool)
	go wrk.Start(ctx)

	// ── Infrastructure: Repositories ──────────────────────────────────
	ticketRepo := postgres.NewTicketRepo(pool)
	commentRepo := postgres.NewCommentRepo(pool)
	historyRepo := postgres.NewHistoryRepo(pool)
	userRepo := postgres.NewUserRepo(pool)
	authUserRepo := postgres.NewAuthUserRepo(pool)
	articleRepo := postgres.NewArticleRepo(pool)
	slaRepo := postgres.NewSLARepo(pool)
	automationRepo := postgres.NewAutomationRepo(pool)
	analyticsRepo := postgres.NewAnalyticsRepo(pool)
	tenantResolver := postgres.NewTenantResolver(pool)

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
	ticketSvc := appTicket.NewService(ticketRepo, commentRepo, historyRepo, hubAdapter, sender)
	userSvc := appUser.NewService(userRepo, hasher)
	articleSvc := appArticle.NewService(articleRepo)
	slaSvc := appSLA.NewService(slaRepo)
	automationSvc := appAuto.NewService(automationRepo)
	analyticsSvc := appAnalytics.NewService(analyticsRepo)

	authCfg := appAuth.Config{
		JWTSecret:    cfg.JWTSecret,
		JWTExpHours:  cfg.JWTExpHours,
		SecureCookie: cfg.SecureCookie,
	}
	authSvc := appAuth.NewService(authUserRepo, hasher, ldapAuth, tokenGen, authCfg)

	// ── Interface: Handlers ───────────────────────────────────────────
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	ticketHandler := handler.NewTicketHandler(ticketSvc)
	articleHandler := handler.NewArticleHandler(articleSvc)
	slaHandler := handler.NewSLAHandler(slaSvc)
	automationHandler := handler.NewAutomationHandler(automationSvc)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsSvc)
	webHandler := handler.NewWebHandler(authSvc, userSvc, ticketSvc, articleSvc, slaSvc, automationSvc, analyticsSvc)

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

	// CORS
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			origin := cfg.CORSOrigin
			if origin == "" {
				origin = "http://localhost:4321"
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if req.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, req)
		})
	})

	// ── Static files ──────────────────────────────────────────────────
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		staticDir = "./internal/interfaces/http/static"
	}
	fileServer := http.FileServer(http.Dir(staticDir))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// ── Web pages (HTML, server-side rendered) ────────────────────────
	// Public: login (stricter rate limiting to mitigate brute force attacks)
	authLimiter := httpMiddleware.NewRateLimiter(rate.Limit(5/60.0), 5, 10*time.Minute)
	authLimiter.Start()

	r.Get("/login", webHandler.LoginPage)
	r.With(authLimiter.LimitByIP).Post("/auth/login", webHandler.LoginSubmit)
	r.Post("/auth/logout", webHandler.Logout)

	// Protected web pages (auth + tenant required)
	r.Group(func(r chi.Router) {
		r.Use(httpMiddleware.Authenticate(tokenGen))
		r.Use(httpMiddleware.Tenant(tenantResolver))

		r.Get("/dashboard", webHandler.Dashboard)

		// Tickets
		r.Get("/tickets", webHandler.TicketList)
		r.Get("/tickets/new", webHandler.TicketFormPage)
		r.Post("/tickets", webHandler.TicketCreate)
		r.Get("/tickets/{id}", webHandler.TicketDetailPage)
		r.Post("/tickets/{id}/comments", webHandler.TicketAddComment)

		// Articles (Knowledge Base)
		r.Get("/articles", webHandler.ArticleList)
		r.Get("/articles/new", webHandler.ArticleFormPage)
		r.Post("/articles", webHandler.ArticleCreate)
		r.Get("/articles/{id}", webHandler.ArticleDetailPage)

		// Admin: Users
		r.Group(func(r chi.Router) {
			r.Use(httpMiddleware.RequireRole("admin"))
			r.Get("/admin/users", webHandler.UserList)
			r.Get("/admin/users/new", webHandler.UserFormPage)
			r.Post("/admin/users", webHandler.UserCreate)

			// Admin: SLA
			r.Get("/admin/sla", webHandler.SLAList)
			r.Get("/admin/sla/{priority}/edit", webHandler.SLAFormPage)
			r.Post("/admin/sla/{priority}", webHandler.SLAUpsert)

			// Admin: Automations
			r.Get("/admin/automations", webHandler.AutomationList)
			r.Get("/admin/automations/new", webHandler.AutomationFormPage)
			r.Post("/admin/automations", webHandler.AutomationCreate)

			// Admin: Analytics
			r.Get("/admin/analytics", webHandler.AnalyticsPage)
		})
	})

	// Root redirect
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
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
			r.Use(httpMiddleware.Tenant(tenantResolver))

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

// Ensure imports are used
var (
	_ = zerolog.ConsoleWriter{}
	_ = ticket.StatusOpen
)
