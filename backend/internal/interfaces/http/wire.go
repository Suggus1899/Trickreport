package http

// wire.go — dependency wiring helpers

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	appAnalytics "github.com/trickreport/backend/internal/application/analytics"
	appArticle "github.com/trickreport/backend/internal/application/article"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	appAuto "github.com/trickreport/backend/internal/application/automation"
	appSLA "github.com/trickreport/backend/internal/application/sla"
	appTicket "github.com/trickreport/backend/internal/application/ticket"
	appUser "github.com/trickreport/backend/internal/application/user"
	"github.com/trickreport/backend/internal/email"
	infraAuth "github.com/trickreport/backend/internal/infrastructure/auth"
	"github.com/trickreport/backend/internal/infrastructure/postgres"
	infraRealtime "github.com/trickreport/backend/internal/infrastructure/realtime"
	"github.com/trickreport/backend/internal/interfaces/http/handler"
)

// Repos holds all repository implementations.
type Repos struct {
	Ticket         *postgres.TicketRepo
	Comment        *postgres.CommentRepo
	History        *postgres.HistoryRepo
	User           *postgres.UserRepo
	Article        *postgres.ArticleRepo
	SLA            *postgres.SLARepo
	Automation     *postgres.AutomationRepo
	Analytics      *postgres.AnalyticsRepo
	TenantResolver *postgres.TenantResolver
	Attachment     *postgres.AttachmentRepo
	PasswordReset  *postgres.PasswordResetRepo
	Session        *postgres.SessionRepo
	TxManager      *postgres.TxManager
	Notification   *postgres.NotificationRepo
}

// NewRepos creates all repositories from a pgxpool.
func NewRepos(pool *pgxpool.Pool) *Repos {
	return &Repos{
		Ticket:         postgres.NewTicketRepo(pool),
		Comment:        postgres.NewCommentRepo(pool),
		History:        postgres.NewHistoryRepo(pool),
		User:           postgres.NewUserRepo(pool),
		Article:        postgres.NewArticleRepo(pool),
		SLA:            postgres.NewSLARepo(pool),
		Automation:     postgres.NewAutomationRepo(pool),
		Analytics:      postgres.NewAnalyticsRepo(pool),
		TenantResolver: postgres.NewTenantResolver(pool),
		Attachment:     postgres.NewAttachmentRepo(pool),
		PasswordReset:  postgres.NewPasswordResetRepo(pool),
		Session:        postgres.NewSessionRepo(pool),
		TxManager:      postgres.NewTxManager(pool),
		Notification:   postgres.NewNotificationRepo(pool),
	}
}

// Services holds all application services.
type Services struct {
	Auth          *appAuth.Service
	Ticket        *appTicket.UserService
	User          *appUser.Service
	Article       *appArticle.Service
	SLA           *appSLA.Service
	Automation    *appAuto.Service
	Analytics     *appAnalytics.Service
	Engine        *appAuto.Engine
	Attachment    *appTicket.AttachmentService
	MFA           *appAuth.MFAService
	PasswordReset *appAuth.PasswordResetService
	TokenStore    *infraAuth.MemoryTokenStore
	Notification  *appTicket.NotificationService
}

// NewServices creates all application services from repos and adapters.
func NewServices(
	repos *Repos,
	hasher *infraAuth.BcryptHasher,
	tokenGen *infraAuth.JWTGenerator,
	hubAdapter *infraRealtime.HubAdapter,
	sender email.Sender,
	ldapAuth appAuth.LDAPAuthenticator,
	authCfg appAuth.Config,
) *Services {
	// Build the automation engine: executor uses direct SQL on the pool,
	// engine uses the automation repo + executor.
	executor := postgres.NewAutomationExecutor(repos.Automation.DB())
	engine := appAuto.NewEngine(repos.Automation, executor, log.Logger)

	ticketSvc := appTicket.NewService(repos.Ticket, repos.Comment, repos.History, hubAdapter, sender)
	ticketSvc.SetEngine(engine)
	ticketSvc.SetTxManager(repos.TxManager)
	ticketSvc.SetSLAPolicyFetcher(repos.SLA)

	// Notification service — wired into the ticket service so events create
	// persistent notifications and broadcast them in real-time via WebSocket.
	notifSvc := appTicket.NewNotificationService(repos.Notification)
	ticketSvc.SetNotificationService(notifSvc)

	// Auth service with token store, session repo, account lockout, MFA repo,
	// and full password hasher (for password reset and MFA flows).
	tokenStore := infraAuth.NewMemoryTokenStore()
	authSvc := appAuth.NewService(repos.User, hasher, ldapAuth, tokenGen, authCfg)
	authSvc.SetTokenStore(tokenStore)
	authSvc.SetSessionRepository(repos.Session)
	authSvc.SetAccountLockoutRepository(repos.User)
	authSvc.SetPasswordHasherFull(hasher)
	authSvc.SetMFARepository(repos.User)
	authSvc.SetUserCreator(repos.User)

	// MFA service
	mfaSvc := appAuth.NewMFAService(repos.User, appAuth.MFAConfig{Issuer: "Trickreport"})

	// Password reset service
	resetSvc := appAuth.NewPasswordResetService(
		repos.User,
		repos.PasswordReset,
		repos.User,
		hasher,
		sender,
		appAuth.PasswordResetConfig{
			TokenExpiry: time.Hour,
			ResetURL:    "/reset-password",
		},
	)

	return &Services{
		Ticket:        ticketSvc,
		User:          appUser.NewService(repos.User, hasher),
		Article:       appArticle.NewService(repos.Article),
		SLA:           appSLA.NewService(repos.SLA),
		Automation:    appAuto.NewService(repos.Automation),
		Analytics:     appAnalytics.NewService(repos.Analytics),
		Auth:          authSvc,
		Engine:        engine,
		Attachment:    appTicket.NewAttachmentService(repos.Attachment, repos.Ticket),
		MFA:           mfaSvc,
		PasswordReset: resetSvc,
		TokenStore:    tokenStore,
		Notification:  notifSvc,
	}
}

// Handlers holds all HTTP handlers.
type Handlers struct {
	Auth         *handler.AuthHandler
	Ticket       *handler.TicketHandler
	User         *handler.UserHandler
	Article      *handler.ArticleHandler
	SLA          *handler.SLAHandler
	Automation   *handler.AutomationHandler
	Analytics    *handler.AnalyticsHandler
	Attachment   *handler.AttachmentHandler
	Notification *handler.NotificationHandler
}

// NewHandlers creates all HTTP handlers from services.
func NewHandlers(services *Services) *Handlers {
	authH := handler.NewAuthHandler(services.Auth)
	authH.SetMFAService(services.MFA)
	authH.SetPasswordResetService(services.PasswordReset)

	return &Handlers{
		Auth:         authH,
		User:         handler.NewUserHandler(services.User),
		Ticket:       handler.NewTicketHandler(services.Ticket),
		Article:      handler.NewArticleHandler(services.Article),
		SLA:          handler.NewSLAHandler(services.SLA),
		Automation:   handler.NewAutomationHandler(services.Automation),
		Analytics:    handler.NewAnalyticsHandler(services.Analytics),
		Attachment:   handler.NewAttachmentHandler(services.Attachment),
		Notification: handler.NewNotificationHandler(services.Notification),
	}
}
