package http

// wire.go — dependency wiring helpers

import (
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
	"github.com/jackc/pgx/v5/pgxpool"
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
	}
}

// Services holds all application services.
type Services struct {
	Auth       *appAuth.Service
	Ticket     *appTicket.UserService
	User       *appUser.Service
	Article    *appArticle.Service
	SLA        *appSLA.Service
	Automation *appAuto.Service
	Analytics  *appAnalytics.Service
}

// NewServices creates all application services from repos and adapters.
func NewServices(
	repos *Repos,
	hasher *infraAuth.BcryptHasher,
	tokenGen *infraAuth.JWTGenerator,
	hubAdapter *infraRealtime.HubAdapter,
	sender *email.ConsoleSender,
	ldapAuth appAuth.LDAPAuthenticator,
	authCfg appAuth.Config,
) *Services {
	return &Services{
		Ticket:     appTicket.NewService(repos.Ticket, repos.Comment, repos.History, hubAdapter, sender),
		User:       appUser.NewService(repos.User, hasher),
		Article:    appArticle.NewService(repos.Article),
		SLA:        appSLA.NewService(repos.SLA),
		Automation: appAuto.NewService(repos.Automation),
		Analytics:  appAnalytics.NewService(repos.Analytics),
		Auth:       appAuth.NewService(repos.User, hasher, ldapAuth, tokenGen, authCfg),
	}
}

// Handlers holds all HTTP handlers.
type Handlers struct {
	Auth       *handler.AuthHandler
	Ticket     *handler.TicketHandler
	User       *handler.UserHandler
	Article    *handler.ArticleHandler
	SLA        *handler.SLAHandler
	Automation *handler.AutomationHandler
	Analytics  *handler.AnalyticsHandler
}

// NewHandlers creates all HTTP handlers from services.
func NewHandlers(services *Services) *Handlers {
	return &Handlers{
		Auth:       handler.NewAuthHandler(services.Auth),
		User:       handler.NewUserHandler(services.User),
		Ticket:     handler.NewTicketHandler(services.Ticket),
		Article:    handler.NewArticleHandler(services.Article),
		SLA:        handler.NewSLAHandler(services.SLA),
		Automation: handler.NewAutomationHandler(services.Automation),
		Analytics:  handler.NewAnalyticsHandler(services.Analytics),
	}
}
