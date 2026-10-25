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
	"github.com/trickreport/backend/internal/interfaces/http/handler"
	httpMiddleware "github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/realtime"
	"github.com/trickreport/backend/internal/worker"
)

// Server is the HTTP server with all dependencies wired.
type Server struct {
	cfg            *config.Config
	pool           *pgxpool.Pool
	router         *chi.Mux
	workerCtx      context.Context
	backgroundDone []<-chan struct{}
}

// New creates a new Server with all dependencies wired.
func New(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool) *Server {
	srv := &Server{cfg: cfg, pool: pool, workerCtx: ctx}
	startBackground := func(run func()) {
		done := make(chan struct{})
		srv.backgroundDone = append(srv.backgroundDone, done)
		go func() {
			defer close(done)
			run()
		}()
	}

	// ── Realtime ──────────────────────────────────────────────────────
	hub := realtime.NewHub()
	startBackground(func() { hub.Run(ctx) })

	// Configure WebSocket origin whitelist from CORS config
	realtime.ConfigureAllowedOrigins(cfg.AllowedOrigins())

	// ── Email ─────────────────────────────────────────────────────────
	var sender email.Sender
	if cfg.Email.SMTPHost != "" {
		sender = infraEmail.NewSMTPSender(cfg.Email.SMTPHost, cfg.Email.SMTPPort, cfg.Email.SMTPUsername, cfg.Email.SMTPPassword, cfg.Email.SMTPFrom, cfg.Email.SMTPUseTLS)
		log.Info().Str("host", cfg.Email.SMTPHost).Int("port", cfg.Email.SMTPPort).Msg("SMTP email sender configured")
	} else if cfg.IsProduction() || strings.EqualFold(cfg.Env, "staging") {
		sender = email.UnavailableSender{}
		log.Warn().Msg("SMTP is not configured; email will remain in the outbox")
	} else {
		sender = email.NewConsoleSender()
		log.Info().Msg("Console email sender configured (no SMTP_HOST set)")
	}

	// ── Infrastructure: Repositories ──────────────────────────────────
	repos := NewRepos(pool)

	// ── Infrastructure: Adapters ──────────────────────────────────────
	hasher := infraAuth.NewBcryptHasher()
	tokenGen := infraAuth.NewJWTGenerator(cfg.JWTSecret, cfg.JWTExpHours)
	hubAdapter := infraRealtime.NewHubAdapter(hub, repos.User)

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
		PublicURL:    cfg.AppPublicURL,
	}
	services := NewServices(repos, hasher, tokenGen, hubAdapter, sender, ldapAuth, authCfg)

	// ── Worker ────────────────────────────────────────────────────────
	wrk := worker.New(pool, worker.WithEngine(services.Engine))
	if pool != nil {
		startBackground(func() { services.EmailQueue.Process(ctx, sender) })
		startBackground(func() { wrk.Start(ctx) })
	}

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
	itsmHandler := handlers.ITSM
	mspHandler := handlers.MSP

	// ── Router ────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
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
	srv.backgroundDone = append(srv.backgroundDone, globalLimiter.Start(ctx))

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
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID, X-CSRF-Token, Idempotency-Key")
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

	// Root redirect — send browser to the frontend (first allowed origin)
	// or fall back to a JSON health message when no origin is configured.
	r.With(globalLimiter.LimitByIP).Get("/", func(w http.ResponseWriter, r *http.Request) {
		if len(allowedOrigins) > 0 {
			http.Redirect(w, r, allowedOrigins[0], http.StatusSeeOther)
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// ── API routes (JSON) ─────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(httpMiddleware.CSRF(cfg.SecureCookie))

		// Public: auth (no tenant required)
		r.Route("/auth", func(r chi.Router) {
			r.With(globalLimiter.LimitByIP).Get("/csrf", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
			// Login rate limiting (per email+IP)
			loginLimiter := httpMiddleware.NewLoginRateLimiter(rate.Limit(5), 10, 15)
			srv.backgroundDone = append(srv.backgroundDone, loginLimiter.Start(ctx))

			// These run before any JWT exists, so the tenant can only come from
			// the X-Tenant-ID header or the subdomain — hence the Tenant
			// middleware here. Without it these lookups would be global, and an
			// email that exists in two organizations would resolve arbitrarily.
			tenant := httpMiddleware.Tenant(repos.TenantResolver)

			r.With(globalLimiter.LimitByIP, loginLimiter.LoginLimit, tenant).Post("/login", authHandler.Login)
			r.With(globalLimiter.LimitByIP, loginLimiter.LoginLimit, tenant).Post("/register", authHandler.Register)
			r.With(globalLimiter.LimitByIP).Post("/logout", authHandler.Logout)
			r.With(globalLimiter.LimitByIP).Post("/refresh", authHandler.Refresh)
			r.With(globalLimiter.LimitByIP, tenant).Post("/password-reset", authHandler.PasswordReset)
			r.With(globalLimiter.LimitByIP).Post("/password-reset/confirm", authHandler.PasswordResetConfirm)
			r.With(globalLimiter.LimitByIP, tenant).Post("/mfa/login", authHandler.MFALogin)
			r.With(httpMiddleware.Authenticate(services.Auth), globalLimiter.LimitByUser).Get("/me", authHandler.Me)

			// MFA management (requires auth)
			r.Group(func(r chi.Router) {
				r.Use(httpMiddleware.Authenticate(services.Auth))
				r.Use(globalLimiter.LimitByUser)
				r.Post("/mfa/setup", authHandler.MFASetup)
				r.Post("/mfa/verify", authHandler.MFAVerify)
				r.Post("/mfa/enable", authHandler.MFAEnable)
				r.Post("/mfa/disable", authHandler.MFADisable)
			})
		})

		// Zabbix webhook (Stage 4): no user session exists for a
		// monitoring system, so this route sits outside the
		// authenticated group. It still resolves the tenant from the
		// X-Tenant-ID header and authenticates with the per-source
		// webhook token inside the handler; CSRF is exempt for this
		// path only (see middleware CSRF).
		r.With(httpMiddleware.Tenant(repos.TenantResolver)).Post("/msp/monitoring/webhook", mspHandler.Webhook)

		// Protected routes (auth + tenant required)
		r.Group(func(r chi.Router) {
			r.Use(httpMiddleware.Authenticate(services.Auth))
			r.Use(httpMiddleware.Tenant(repos.TenantResolver))
			// Authenticated identities isolate users sharing the frontend/proxy connection.
			r.Use(globalLimiter.LimitByUser)

			// WebSockets
			r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
				tenantID := httpMiddleware.TenantFromContext(r.Context())
				claims, _ := httpMiddleware.ClaimsFromContext(r.Context())
				rawToken := ""
				if c, err := r.Cookie("trickreport_token"); err == nil {
					rawToken = c.Value
				}
				if header := r.Header.Get("Authorization"); strings.HasPrefix(header, "Bearer ") {
					rawToken = strings.TrimPrefix(header, "Bearer ")
				}
				if rawToken == "" {
					rawToken = r.URL.Query().Get("token")
				}
				realtime.ServeWs(hub, w, r, tenantID, claims.UserID, func(ctx context.Context) error { _, err := services.Auth.AuthenticateAccess(ctx, rawToken); return err })
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
					r.Get("/itsm", itsmHandler.Report)
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

			// Notifications (real-time + persistent). Every handler here reads the
			// caller from the request context, so this has to stay inside the
			// authenticated group.
			r.Route("/notifications", func(r chi.Router) {
				r.Get("/", notificationHandler.List)
				r.Get("/unread-count", notificationHandler.UnreadCount)
				r.Post("/{id}/read", notificationHandler.MarkRead)
				r.Post("/read-all", notificationHandler.MarkAllRead)
			})

			// ITSM (Stage 3): catalog, approvals, assets/CMDB, problems,
			// changes. Reads are open to any authenticated user of the
			// tenant; writes are staff-only except requesting a catalog
			// service and deciding an approval, whose per-record rule
			// lives in the application service.
			r.Route("/catalog", func(r chi.Router) {
				r.Get("/categories", itsmHandler.ListCategories)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/categories", itsmHandler.CreateCategory)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Put("/categories/{id}", itsmHandler.UpdateCategory)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Delete("/categories/{id}", itsmHandler.DeleteCategory)

				r.Get("/services", itsmHandler.ListServices)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/services", itsmHandler.CreateService)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Put("/services/{id}", itsmHandler.UpdateService)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Delete("/services/{id}", itsmHandler.DeleteService)
				r.Post("/services/{id}/request", itsmHandler.RequestService)
			})

			r.Route("/approvals", func(r chi.Router) {
				r.Get("/", itsmHandler.ListApprovals)
				r.Get("/{id}/history", itsmHandler.ApprovalHistory)
				r.Post("/{id}/decide", itsmHandler.DecideApproval)
			})

			r.Route("/ci-types", func(r chi.Router) {
				r.Get("/", itsmHandler.ListCITypes)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", itsmHandler.CreateCIType)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Delete("/{id}", itsmHandler.DeleteCIType)
			})

			r.Route("/assets", func(r chi.Router) {
				r.Get("/", itsmHandler.ListAssets)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", itsmHandler.CreateAsset)
				r.Get("/{id}", itsmHandler.GetAsset)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Put("/{id}", itsmHandler.UpdateAsset)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Delete("/{id}", itsmHandler.DeleteAsset)
				r.Get("/{id}/relationships", itsmHandler.ListRelationships)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/relationships", itsmHandler.CreateRelationship)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Delete("/relationships/{id}", itsmHandler.DeleteRelationship)
			})

			r.Route("/ticket-assets", func(r chi.Router) {
				r.Get("/{ticketId}", itsmHandler.AssetsForTicket)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", itsmHandler.LinkTicketAsset)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Delete("/{ticketId}/{assetId}", itsmHandler.UnlinkTicketAsset)
			})

			r.Route("/problems", func(r chi.Router) {
				r.Get("/", itsmHandler.ListProblems)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", itsmHandler.CreateProblem)
				r.Get("/{id}", itsmHandler.GetProblem)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Patch("/{id}", itsmHandler.UpdateProblem)
				r.Get("/{id}/incidents", itsmHandler.IncidentsForProblem)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/incidents", itsmHandler.LinkIncident)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Delete("/{id}/incidents/{ticketId}", itsmHandler.UnlinkIncident)
			})

			r.Route("/changes", func(r chi.Router) {
				r.Get("/", itsmHandler.ListChanges)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", itsmHandler.CreateChange)
				r.Get("/{id}", itsmHandler.GetChange)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Patch("/{id}/status", itsmHandler.UpdateChangeStatus)
				r.Get("/{id}/tickets", itsmHandler.TicketsForChange)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/tickets", itsmHandler.LinkChangeTicket)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Delete("/{id}/tickets/{ticketId}", itsmHandler.UnlinkChangeTicket)
			})

			// MSP (Stage 4): multi-client operation. Reads are open to
			// any authenticated user of the tenant, but the service
			// restricts client contacts to their own clients; writes
			// are staff-only except creating a multi-client ticket
			// and requesting a remote session, whose per-record rules
			// live in the application service.
			r.Route("/msp/clients", func(r chi.Router) {
				r.Get("/", mspHandler.ListClients)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", mspHandler.CreateClient)
				r.Get("/{id}", mspHandler.GetClient)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Put("/{id}", mspHandler.UpdateClient)
				r.With(httpMiddleware.RequireRole("admin")).Delete("/{id}", mspHandler.DeleteClient)
				r.Get("/{id}/contacts", mspHandler.ListContacts)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/contacts", mspHandler.CreateContact)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Delete("/{id}/contacts/{contactId}", mspHandler.DeleteContact)
				r.Get("/{id}/assets", mspHandler.ListClientAssets)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/assets", mspHandler.LinkClientAsset)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Delete("/{id}/assets/{assetId}", mspHandler.UnlinkClientAsset)
				r.Get("/{id}/discovered-assets", mspHandler.ListDiscoveredAssets)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/discovered-assets", mspHandler.CreateDiscoveredAsset)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/discovered-assets/{discoveredId}/import", mspHandler.ImportDiscoveredAsset)
			})

			r.Route("/msp/contracts", func(r chi.Router) {
				r.Get("/", mspHandler.ListContracts)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", mspHandler.CreateContract)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Put("/{id}", mspHandler.UpdateContract)
			})

			r.Route("/msp/tickets", func(r chi.Router) {
				r.Get("/", mspHandler.ListMSPTickets)
				r.Post("/", mspHandler.CreateMSPTicket)
				r.Get("/{id}/context", mspHandler.GetTicketContext)
			})

			r.Route("/msp/technicians", func(r chi.Router) {
				r.Get("/", mspHandler.ListTechnicians)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", mspHandler.CreateTechnician)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Put("/{id}", mspHandler.UpdateTechnician)
				r.Get("/{id}/availability", mspHandler.ListAvailability)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/availability", mspHandler.CreateAvailability)
			})

			r.Route("/msp/interventions", func(r chi.Router) {
				r.Get("/", mspHandler.ListInterventions)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", mspHandler.CreateIntervention)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Patch("/{id}/status", mspHandler.UpdateInterventionStatus)
			})

			r.Route("/msp/time-entries", func(r chi.Router) {
				r.Get("/", mspHandler.ListTimeEntries)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", mspHandler.CreateTimeEntry)
			})

			r.Route("/msp/expenses", func(r chi.Router) {
				r.Get("/", mspHandler.ListExpenses)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", mspHandler.CreateExpense)
			})

			r.Route("/msp/projects", func(r chi.Router) {
				r.Get("/", mspHandler.ListProjects)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", mspHandler.CreateProject)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Patch("/{id}/status", mspHandler.UpdateProjectStatus)
				r.Get("/{id}/tasks", mspHandler.ListTasks)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/tasks", mspHandler.CreateTask)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Patch("/{id}/tasks/{taskId}/status", mspHandler.UpdateTaskStatus)
			})

			r.Route("/msp/billing-drafts", func(r chi.Router) {
				r.Get("/", mspHandler.ListBillingDrafts)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", mspHandler.GenerateBillingDraft)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Patch("/{id}/status", mspHandler.UpdateBillingDraftStatus)
			})

			r.Get("/msp/profitability", mspHandler.Profitability)
			r.With(httpMiddleware.RequireRole("admin")).Get("/admin/msp", mspHandler.Report)

			r.Route("/msp/monitoring", func(r chi.Router) {
				r.Route("/sources", func(r chi.Router) {
					r.Get("/", mspHandler.ListMonitoringSources)
					r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/", mspHandler.CreateMonitoringSource)
					r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/rotate-token", mspHandler.RotateSourceToken)
				})
				r.Get("/alerts", mspHandler.ListMonitoringAlerts)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/alerts/{id}/resolve", mspHandler.ResolveMonitoringAlert)
			})

			r.Route("/msp/remote-sessions", func(r chi.Router) {
				r.Get("/", mspHandler.ListRemoteSessions)
				r.Post("/", mspHandler.RequestRemoteSession)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/authorize", mspHandler.AuthorizeRemoteSession)
				r.Post("/{id}/start", mspHandler.StartRemoteSession)
				r.With(httpMiddleware.RequireRole("admin", "agent")).Post("/{id}/end", mspHandler.EndRemoteSession)
			})
		})
	})

	// Metrics endpoint (Prometheus)
	r.Handle("/metrics", handler.MetricsHandler())

	// Swagger UI
	r.Handle("/swagger/*", handler.SwaggerHandler())

	// Health check — includes database connectivity
	r.With(globalLimiter.LimitByIP).Get("/health", func(w http.ResponseWriter, r *http.Request) {
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
	r.With(globalLimiter.LimitByIP).Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			httpMiddleware.LoggerFromContext(r.Context()).Error().Err(err).Msg("ready check failed")
			response.Error(w, http.StatusServiceUnavailable, "database unavailable")
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"status": "ready"})
	})

	srv.router = r
	return srv
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
	for _, done := range s.backgroundDone {
		<-done
	}
	log.Info().Msg("Server stopped")
}
