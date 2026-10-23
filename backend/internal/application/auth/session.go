package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/domain/user"
)

var ErrAuthUnavailable = errors.New("authentication service unavailable")

type MFAChallengeRepository interface {
	CreateChallenge(context.Context, string, uuid.UUID, time.Time) error
	ConsumeChallenge(context.Context, string, uuid.UUID, bool) (bool, error)
}

func (s *Service) SetMFAChallengeRepository(repo MFAChallengeRepository) { s.mfaChallenges = repo }

func (s *Service) currentUser(ctx context.Context, claims *Claims) (*user.User, error) {
	u, err := s.users.GetByIDNoTenant(ctx, claims.UserID)
	if errors.Is(err, user.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, fmt.Errorf("%w: user lookup", ErrAuthUnavailable)
	}
	if u == nil || !u.Active || u.TenantID != claims.TenantID || string(u.Role) != claims.Role {
		return nil, ErrInvalidCredentials
	}
	return u, nil
}

func (s *Service) AuthenticateAccess(ctx context.Context, token string) (*Claims, error) {
	claims, err := s.tokens.Validate(token)
	if err != nil || claims == nil || claims.TokenType != "access" {
		return nil, ErrInvalidCredentials
	}
	if s.sessions == nil {
		return nil, ErrAuthUnavailable
	}
	session, err := s.sessions.FindActiveByAccessJTI(ctx, claims.JTI)
	if errors.Is(err, ErrInvalidCredentials) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, ErrAuthUnavailable
	}
	if session == nil || session.UserID != claims.UserID {
		return nil, ErrInvalidCredentials
	}
	if _, err := s.currentUser(ctx, claims); err != nil {
		return nil, err
	}
	return claims, nil
}
