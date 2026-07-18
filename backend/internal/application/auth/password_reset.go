package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/domain/user"
	"github.com/trickreport/backend/internal/email"
)

// PasswordResetConfig holds password reset-specific configuration.
type PasswordResetConfig struct {
	TokenExpiry time.Duration // default: 1 hour
	ResetURL    string        // base URL for the reset link, e.g. "https://app.trickreport.com/reset"
}

// PasswordResetRepository is the port for password reset token persistence.
type PasswordResetRepository interface {
	CreateToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	FindValidToken(ctx context.Context, token string) (userID uuid.UUID, expiresAt time.Time, err error)
	MarkUsed(ctx context.Context, token string) error
}

// PasswordResetService handles the password reset flow: requesting a reset
// (generating a token and sending an email) and confirming the reset
// (validating the token and updating the password).
type PasswordResetService struct {
	users     UserRepository
	tokens    PasswordResetRepository
	updater   PasswordUpdater
	hasher    PasswordHasherFull
	sender    email.Sender
	cfg       PasswordResetConfig
}

// NewPasswordResetService creates a new PasswordResetService.
func NewPasswordResetService(
	users UserRepository,
	tokens PasswordResetRepository,
	updater PasswordUpdater,
	hasher PasswordHasherFull,
	sender email.Sender,
	cfg PasswordResetConfig,
) *PasswordResetService {
	if cfg.TokenExpiry == 0 {
		cfg.TokenExpiry = time.Hour
	}
	return &PasswordResetService{
		users:   users,
		tokens:  tokens,
		updater: updater,
		hasher:  hasher,
		sender:  sender,
		cfg:     cfg,
	}
}

// ErrPasswordResetTokenInvalid is returned when a reset token is invalid or expired.
var ErrPasswordResetTokenInvalid = errors.New("password reset token is invalid or expired")

// ErrPasswordResetTokenUsed is returned when a reset token has already been used.
var ErrPasswordResetTokenUsed = errors.New("password reset token has already been used")

// RequestReset generates a reset token for the user with the given email and
// sends a reset email. If no user exists with the email, the method returns
// nil to avoid leaking which emails are registered.
func (s *PasswordResetService) RequestReset(ctx context.Context, emailAddr string) error {
	u, err := s.users.GetByEmail(ctx, emailAddr)
	if err != nil {
		// Don't leak whether the email exists
		return nil
	}

	token := uuid.NewString()
	expiresAt := time.Now().Add(s.cfg.TokenExpiry)

	if err := s.tokens.CreateToken(ctx, u.ID, token, expiresAt); err != nil {
		return fmt.Errorf("password reset: create token: %w", err)
	}

	if s.sender != nil {
		resetLink := fmt.Sprintf("%s?token=%s", s.cfg.ResetURL, token)
		body := fmt.Sprintf(
			"Hello %s,\n\nYou requested a password reset for your Trickreport account.\n"+
				"Click the link below to reset your password (valid for %s):\n\n%s\n\n"+
				"If you did not request this reset, you can safely ignore this email.",
			u.Name, s.cfg.TokenExpiry, resetLink,
		)
		_ = s.sender.Send(u.Email, "Password Reset - Trickreport", body)
	}

	return nil
}

// ConfirmReset validates the reset token and updates the user's password.
// The new password is validated for complexity before being stored.
func (s *PasswordResetService) ConfirmReset(ctx context.Context, token, newPassword string) error {
	if err := user.ValidatePasswordComplexity(newPassword); err != nil {
		return err
	}

	userID, expiresAt, err := s.tokens.FindValidToken(ctx, token)
	if err != nil {
		return ErrPasswordResetTokenInvalid
	}

	if time.Now().After(expiresAt) {
		return ErrPasswordResetTokenInvalid
	}

	hash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("password reset: hash: %w", err)
	}

	if err := s.updater.UpdatePassword(ctx, userID, hash); err != nil {
		return fmt.Errorf("password reset: update password: %w", err)
	}

	if err := s.tokens.MarkUsed(ctx, token); err != nil {
		// Password was already updated; log but don't fail
		return nil
	}

	return nil
}
