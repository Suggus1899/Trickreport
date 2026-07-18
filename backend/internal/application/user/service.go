package user

import (
	"context"

	"github.com/trickreport/backend/internal/domain/user"
)

// Repository is the port for user persistence.
type Repository interface {
	List(ctx context.Context, tenantID string) ([]user.User, error)
	GetByID(ctx context.Context, id, tenantID string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	Create(ctx context.Context, u *user.User, passwordHash string) error
	Update(ctx context.Context, id, tenantID string, fields UpdateFields) (*user.User, error)
	Deactivate(ctx context.Context, id, tenantID string) error
}

// UpdateFields holds optional fields for partial updates.
type UpdateFields struct {
	Name      *string
	Role      *string
	AvatarURL *string
	Active    *bool
}

// PasswordHasher is the port for password hashing.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

// Service is the application service for user operations.
type Service struct {
	repo   Repository
	hasher PasswordHasher
}

// NewService creates a new user application service.
func NewService(repo Repository, hasher PasswordHasher) *Service {
	return &Service{repo: repo, hasher: hasher}
}

// List returns all users in the tenant.
func (s *Service) List(ctx context.Context, tenantID string) ([]user.User, error) {
	return s.repo.List(ctx, tenantID)
}

// Get returns a single user by ID.
func (s *Service) Get(ctx context.Context, id, tenantID string) (*user.User, error) {
	return s.repo.GetByID(ctx, id, tenantID)
}

// CreateInput holds the data for creating a new user.
type CreateInput struct {
	TenantID string
	Name     string
	Email    string
	Role     string
	Password string
}

// Create creates a new user within a tenant.
func (s *Service) Create(ctx context.Context, input CreateInput) (*user.User, error) {
	if input.Name == "" || input.Email == "" {
		return nil, user.ErrValidation
	}

	role := user.RoleEndUser
	if input.Role != "" {
		r, err := user.ParseRole(input.Role)
		if err != nil {
			return nil, user.ErrValidation
		}
		role = r
	}

	var passwordHash string
	if input.Password != "" {
		hash, err := s.hasher.Hash(input.Password)
		if err != nil {
			return nil, err
		}
		passwordHash = hash
	}

	u := &user.User{
		TenantID: input.TenantID,
		Name:     input.Name,
		Email:    input.Email,
		Role:     role,
	}

	if err := s.repo.Create(ctx, u, passwordHash); err != nil {
		return nil, err
	}

	return u, nil
}

// Update partially updates a user.
func (s *Service) Update(ctx context.Context, id, tenantID string, fields UpdateFields) (*user.User, error) {
	return s.repo.Update(ctx, id, tenantID, fields)
}

// Deactivate soft-deletes a user by setting active = false.
func (s *Service) Deactivate(ctx context.Context, id, tenantID string) error {
	return s.repo.Deactivate(ctx, id, tenantID)
}
