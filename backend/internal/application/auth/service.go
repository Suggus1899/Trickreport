package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/domain/user"
)

// UserRepository is the port for user lookup during authentication.
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetByIDNoTenant(ctx context.Context, id uuid.UUID) (*user.User, error)
}

// PasswordHasher is the port for password verification.
type PasswordHasher interface {
	Compare(hash, password string) error
}

// PasswordHasherFull is the port for hashing and comparing passwords.
// Used by password reset and MFA flows that need to set new passwords.
type PasswordHasherFull interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

// LDAPAuthenticator is the port for LDAP authentication.
type LDAPAuthenticator interface {
	Authenticate(username, password string) (string, error)
}

// TokenGenerator is the port for JWT generation and validation.
type TokenGenerator interface {
	Generate(userID, tenantID uuid.UUID, role string) (string, error)
	GenerateRefresh(userID, tenantID uuid.UUID, role string) (string, error)
	Validate(token string) (*Claims, error)
}

// TokenStore is the port for blacklisting JWTs by their jti claim.
type TokenStore interface {
	Blacklist(jti string, expiresAt time.Time)
	IsBlacklisted(jti string) bool
}

// AccountLockoutRepository is the port for account lockout state persistence.
type AccountLockoutRepository interface {
	IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error
	ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error
	LockAccount(ctx context.Context, userID uuid.UUID, until time.Time) error
}

// PasswordUpdater is the port for updating a user's password hash.
type PasswordUpdater interface {
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
}

// SessionRepository is the port for user session persistence.
type SessionRepository interface {
	CreateSession(ctx context.Context, input SessionInput) error
	RevokeSession(ctx context.Context, jti string) error
	RevokeByRefreshJTI(ctx context.Context, refreshJTI string) error
}

// SessionInput holds the data for creating a session record.
type SessionInput struct {
	UserID      uuid.UUID
	JTI         string
	RefreshJTI  string
	IPAddress   string
	UserAgent   string
	ExpiresAt   time.Time
}

// MFARepository is the port for MFA/TOTP secret persistence.
type MFARepository interface {
	GetMFASecret(ctx context.Context, userID uuid.UUID) (secret string, enabled bool, err error)
	SetMFASecret(ctx context.Context, userID uuid.UUID, secret string) error
	EnableMFA(ctx context.Context, userID uuid.UUID) error
	DisableMFA(ctx context.Context, userID uuid.UUID) error
}

// Claims represents the authenticated user's JWT claims.
type Claims struct {
	UserID    uuid.UUID
	TenantID  uuid.UUID
	Role      string
	JTI       string
	TokenType string // "access" or "refresh"
	ExpiresAt time.Time
}

// Config holds auth-specific configuration.
type Config struct {
	JWTSecret    string
	JWTExpHours  int
	SecureCookie bool
}

// Account lockout constants.
const (
	maxFailedAttempts  = 5
	lockoutDuration    = 15 * time.Minute
)

// Service is the application service for authentication.
type Service struct {
	users   UserRepository
	hasher  PasswordHasher
	ldap    LDAPAuthenticator
	tokens  TokenGenerator
	cfg     Config

	// Optional dependencies — wired via setters. When nil, the corresponding
	// features are disabled (graceful degradation).
	// TODO: Wire in wire.go after the auth middleware
	tokenStore   TokenStore
	sessions     SessionRepository
	lockoutRepo  AccountLockoutRepository
	hasherFull   PasswordHasherFull
	mfaRepo      MFARepository
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

// SetTokenStore wires the token blacklist store.
func (s *Service) SetTokenStore(store TokenStore) { s.tokenStore = store }

// SetSessionRepository wires the session repository.
func (s *Service) SetSessionRepository(repo SessionRepository) { s.sessions = repo }

// SetAccountLockoutRepository wires the account lockout repository.
func (s *Service) SetAccountLockoutRepository(repo AccountLockoutRepository) { s.lockoutRepo = repo }

// SetPasswordHasherFull wires a hasher that can also hash (not just compare).
func (s *Service) SetPasswordHasherFull(h PasswordHasherFull) { s.hasherFull = h }

// SetMFARepository wires the MFA repository.
func (s *Service) SetMFARepository(repo MFARepository) { s.mfaRepo = repo }

// LoginInput holds the credentials for login.
type LoginInput struct {
	Email    string
	Password string
}

// LoginResult holds the result of a successful login.
type LoginResult struct {
	Token        string
	RefreshToken string
	User         *user.User
	ExpiresAt    time.Time
	SecureCookie bool
	RequiresMFA  bool
	MFAToken     string
}

// ErrInvalidCredentials is returned when authentication fails.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrPasswordNotConfigured is returned when the user has no password (LDAP only).
var ErrPasswordNotConfigured = errors.New("password login not configured for this account")

// ErrAccountLocked is returned when the account is temporarily locked.
var ErrAccountLocked = errors.New("account is temporarily locked due to too many failed attempts")

// ErrMFARequired is returned when MFA verification is needed to complete login.
var ErrMFARequired = errors.New("mfa verification required")

// ErrTokenBlacklisted is returned when a refresh token has been blacklisted.
var ErrTokenBlacklisted = errors.New("token has been revoked")

// ErrInvalidRefreshToken is returned when a refresh token is invalid or not a refresh type.
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

// Login authenticates a user and returns a JWT.
func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	if input.Email == "" || input.Password == "" {
		return nil, ErrInvalidCredentials
	}

	u, err := s.users.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check account lockout
	if !u.LockedUntil.IsZero() && time.Now().Before(u.LockedUntil) {
		return nil, ErrAccountLocked
	}

	// Authenticate via LDAP or local password
	authFailed := false
	if s.ldap != nil {
		if _, err := s.ldap.Authenticate(input.Email, input.Password); err != nil {
			authFailed = true
		}
	} else {
		if u.PasswordHash == "" {
			return nil, ErrPasswordNotConfigured
		}
		if err := s.hasher.Compare(u.PasswordHash, input.Password); err != nil {
			authFailed = true
		}
	}

	if authFailed {
		s.handleFailedLogin(ctx, u)
		return nil, ErrInvalidCredentials
	}

	// Reset failed attempts on successful authentication
	if s.lockoutRepo != nil {
		_ = s.lockoutRepo.ResetFailedAttempts(ctx, u.ID)
	}

	// If MFA is enabled, return a temporary MFA token instead of full tokens
	if u.MFAEnabled {
		mfaToken, err := s.tokens.Generate(u.ID, u.TenantID, string(u.Role))
		if err != nil {
			return nil, err
		}
		return &LoginResult{
			RequiresMFA: true,
			MFAToken:    mfaToken,
			User:        u,
		}, nil
	}

	return s.issueTokens(ctx, u)
}

// handleFailedLogin increments the failed attempt counter and locks the
// account when the threshold is reached.
func (s *Service) handleFailedLogin(ctx context.Context, u *user.User) {
	if s.lockoutRepo == nil {
		return
	}
	_ = s.lockoutRepo.IncrementFailedAttempts(ctx, u.ID)
	u.FailedLoginAttempts++
	if u.FailedLoginAttempts >= maxFailedAttempts {
		until := time.Now().Add(lockoutDuration)
		_ = s.lockoutRepo.LockAccount(ctx, u.ID, until)
		u.LockedUntil = until
	}
}

// issueTokens generates access + refresh tokens and optionally creates a session.
func (s *Service) issueTokens(ctx context.Context, u *user.User) (*LoginResult, error) {
	token, err := s.tokens.Generate(u.ID, u.TenantID, string(u.Role))
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.tokens.GenerateRefresh(u.ID, u.TenantID, string(u.Role))
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(time.Duration(s.cfg.JWTExpHours) * time.Hour)

	return &LoginResult{
		Token:        token,
		RefreshToken: refreshToken,
		User:         u,
		ExpiresAt:    expiresAt,
		SecureCookie: s.cfg.SecureCookie,
	}, nil
}

// RefreshInput holds the data for refreshing tokens.
type RefreshInput struct {
	RefreshToken string
	IPAddress    string
	UserAgent    string
}

// RefreshResult holds the new tokens after a refresh.
type RefreshResult struct {
	Token        string
	RefreshToken string
	ExpiresAt    time.Time
}

// Refresh validates a refresh token, checks the blacklist, generates new
// access + refresh tokens, and blacklists the old refresh token (rotation).
func (s *Service) Refresh(ctx context.Context, input RefreshInput) (*RefreshResult, error) {
	claims, err := s.tokens.Validate(input.RefreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}
	if claims.TokenType != "refresh" {
		return nil, ErrInvalidRefreshToken
	}

	// Check blacklist
	if s.tokenStore != nil && s.tokenStore.IsBlacklisted(claims.JTI) {
		return nil, ErrTokenBlacklisted
	}

	// Blacklist the old refresh token (rotation)
	if s.tokenStore != nil {
		s.tokenStore.Blacklist(claims.JTI, claims.ExpiresAt)
	}

	// Revoke old session and create new one
	if s.sessions != nil {
		_ = s.sessions.RevokeByRefreshJTI(ctx, claims.JTI)
	}

	token, err := s.tokens.Generate(claims.UserID, claims.TenantID, claims.Role)
	if err != nil {
		return nil, err
	}
	newRefresh, err := s.tokens.GenerateRefresh(claims.UserID, claims.TenantID, claims.Role)
	if err != nil {
		return nil, err
	}

	// Create new session record
	if s.sessions != nil {
		newClaims, _ := s.tokens.Validate(newRefresh)
		var newJTI string
		if newClaims != nil {
			newJTI = newClaims.JTI
		}
		_ = s.sessions.CreateSession(ctx, SessionInput{
			UserID:     claims.UserID,
			JTI:        newJTI,
			RefreshJTI: newJTI,
			IPAddress:  input.IPAddress,
			UserAgent:  input.UserAgent,
			ExpiresAt:  time.Now().Add(7 * 24 * time.Hour),
		})
	}

	return &RefreshResult{
		Token:        token,
		RefreshToken: newRefresh,
		ExpiresAt:    time.Now().Add(time.Duration(s.cfg.JWTExpHours) * time.Hour),
	}, nil
}

// Logout revokes the session associated with the given access token's jti.
func (s *Service) Logout(ctx context.Context, accessToken string) error {
	claims, err := s.tokens.Validate(accessToken)
	if err != nil {
		return nil // already invalid, nothing to revoke
	}

	if s.tokenStore != nil {
		s.tokenStore.Blacklist(claims.JTI, claims.ExpiresAt)
	}
	if s.sessions != nil {
		_ = s.sessions.RevokeSession(ctx, claims.JTI)
	}
	return nil
}

// MFALoginInput holds the data for completing MFA login.
type MFALoginInput struct {
	MFAToken string
	Code     string
}

// MFALogin completes the login flow after MFA verification. It validates the
// temporary MFA token, verifies the TOTP code, and issues full tokens.
func (s *Service) MFALogin(ctx context.Context, input MFALoginInput) (*LoginResult, error) {
	claims, err := s.tokens.Validate(input.MFAToken)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	u, err := s.users.GetByIDNoTenant(ctx, claims.UserID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !u.MFAEnabled {
		return nil, ErrMFARequired
	}

	if s.mfaRepo == nil {
		return nil, errors.New("mfa not configured")
	}

	secret, enabled, err := s.mfaRepo.GetMFASecret(ctx, u.ID)
	if err != nil || !enabled {
		return nil, ErrInvalidCredentials
	}

	if !validateTOTP(secret, input.Code) {
		return nil, ErrInvalidCredentials
	}

	return s.issueTokens(ctx, u)
}

// GetProfile returns the current user's profile.
func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	return s.users.GetByIDNoTenant(ctx, userID)
}
