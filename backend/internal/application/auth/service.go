package auth

import (
	"context"
	"errors"
	"time"

	"github.com/trickreport/backend/internal/domain/user"
)

// UserRepository is the port for user lookup during authentication.
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetByID(ctx context.Context, id string) (*user.User, error)
}

// PasswordHasher is the port for password verification.
type PasswordHasher interface {
	Compare(hash, password string) error
}

// LDAPAuthenticator is the port for LDAP authentication.
type LDAPAuthenticator interface {
	Authenticate(username, password string) (string, error)
}

// TokenGenerator is the port for JWT generation.
type TokenGenerator interface {
	Generate(userID, tenantID, role string) (string, error)
	Validate(token string) (*Claims, error)
}

// Claims represents the authenticated user's JWT claims.
type Claims struct {
	UserID   string
	TenantID string
	Role     string
}

// Config holds auth-specific configuration.
type Config struct {
	JWTSecret    string
	JWTExpHours  int
	SecureCookie bool
}

// Service is the application service for authentication.
type Service struct {
	users  UserRepository
	hasher PasswordHasher
	ldap   LDAPAuthenticator
	tokens TokenGenerator
	cfg    Config
}

// NewService creates a new auth application service.
func NewService(users UserRepository, hasher PasswordHasher, ldap LDAPAuthenticator, tokens TokenGenerator, cfg Config) *Service {
	return &Service{
		users:  users,
		hasher: hasher,
		ldap:   ldap,
		tokens: tokens,
		cfg:    cfg,
	}
}

// LoginInput holds the credentials for login.
type LoginInput struct {
	Email    string
	Password string
}

// LoginResult holds the result of a successful login.
type LoginResult struct {
	Token        string
	User         *user.User
	ExpiresAt    time.Time
	SecureCookie bool
}

// ErrInvalidCredentials is returned when authentication fails.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrPasswordNotConfigured is returned when the user has no password (LDAP only).
var ErrPasswordNotConfigured = errors.New("password login not configured for this account")

// Login authenticates a user and returns a JWT.
func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	if input.Email == "" || input.Password == "" {
		return nil, ErrInvalidCredentials
	}

	u, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Authenticate via LDAP or local password
	if s.ldap != nil {
		if _, err := s.ldap.Authenticate(input.Email, input.Password); err != nil {
			return nil, ErrInvalidCredentials
		}
	} else {
		if u.PasswordHash == "" {
			return nil, ErrPasswordNotConfigured
		}
		if err := s.hasher.Compare(u.PasswordHash, input.Password); err != nil {
			return nil, ErrInvalidCredentials
		}
	}

	token, err := s.tokens.Generate(u.ID, u.TenantID, string(u.Role))
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Token:        token,
		User:         u,
		ExpiresAt:    time.Now().Add(time.Duration(s.cfg.JWTExpHours) * time.Hour),
		SecureCookie: s.cfg.SecureCookie,
	}, nil
}

// GetProfile returns the current user's profile.
func (s *Service) GetProfile(ctx context.Context, userID string) (*user.User, error) {
	return s.users.GetByID(ctx, userID)
}
