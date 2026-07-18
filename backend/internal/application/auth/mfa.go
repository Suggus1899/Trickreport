package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

// MFAConfig holds MFA-specific configuration.
type MFAConfig struct {
	Issuer string // e.g. "Trickreport"
}

// MFAService handles MFA/TOTP setup, verification, and management.
type MFAService struct {
	repo MFARepository
	cfg  MFAConfig
}

// NewMFAService creates a new MFAService.
func NewMFAService(repo MFARepository, cfg MFAConfig) *MFAService {
	if cfg.Issuer == "" {
		cfg.Issuer = "Trickreport"
	}
	return &MFAService{repo: repo, cfg: cfg}
}

// SetMFARepository wires the MFA repository. Allows late binding.
func (s *MFAService) SetMFARepository(repo MFARepository) { s.repo = repo }

// ErrMFANotConfigured is returned when MFA is not set up for a user.
var ErrMFANotConfigured = errors.New("mfa not configured for this user")

// ErrMFAAlreadyEnabled is returned when MFA is already enabled.
var ErrMFAAlreadyEnabled = errors.New("mfa is already enabled")

// ErrInvalidMFACode is returned when the TOTP code is invalid.
var ErrInvalidMFACode = errors.New("invalid mfa code")

// SetupResult holds the data for MFA setup.
type SetupResult struct {
	Secret  string
	QRURL   string
}

// Setup generates a new TOTP secret for the user and stores it (not yet enabled).
// It returns the secret and a QR code URL for authenticator apps.
func (s *MFAService) Setup(ctx context.Context, userID uuid.UUID, email string) (*SetupResult, error) {
	if s.repo == nil {
		return nil, errors.New("mfa repository not configured")
	}

	// Check if already enabled
	_, enabled, err := s.repo.GetMFASecret(ctx, userID)
	if err != nil {
		return nil, err
	}
	if enabled {
		return nil, ErrMFAAlreadyEnabled
	}

	// Generate a new TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      s.cfg.Issuer,
		AccountName: email,
	})
	if err != nil {
		return nil, fmt.Errorf("mfa setup: generate key: %w", err)
	}

	// Store the secret (not enabled yet)
	if err := s.repo.SetMFASecret(ctx, userID, key.Secret()); err != nil {
		return nil, err
	}

	return &SetupResult{
		Secret: key.Secret(),
		QRURL:  key.URL(),
	}, nil
}

// Verify validates a TOTP code against the user's stored secret without
// enabling MFA. Used during setup to confirm the user can generate codes.
func (s *MFAService) Verify(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	if s.repo == nil {
		return false, errors.New("mfa repository not configured")
	}

	secret, _, err := s.repo.GetMFASecret(ctx, userID)
	if err != nil {
		return false, err
	}
	if secret == "" {
		return false, ErrMFANotConfigured
	}

	return validateTOTP(secret, code), nil
}

// Enable activates MFA for the user after verifying the TOTP code matches
// the stored secret.
func (s *MFAService) Enable(ctx context.Context, userID uuid.UUID, secret, code string) error {
	if s.repo == nil {
		return errors.New("mfa repository not configured")
	}

	if !validateTOTP(secret, code) {
		return ErrInvalidMFACode
	}

	return s.repo.EnableMFA(ctx, userID)
}

// Disable deactivates and clears MFA for the user.
func (s *MFAService) Disable(ctx context.Context, userID uuid.UUID) error {
	if s.repo == nil {
		return errors.New("mfa repository not configured")
	}
	return s.repo.DisableMFA(ctx, userID)
}

// validateTOTP checks whether the given code is valid for the secret.
// This is a package-level helper used by both MFAService and the auth Service.
func validateTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}
