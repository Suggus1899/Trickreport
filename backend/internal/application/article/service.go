package article

import (
	"context"

	"github.com/trickreport/backend/internal/domain/article"
)

// Filter holds query parameters for listing articles.
type Filter struct {
	Search string
}

// Repository is the port for article persistence.
type Repository interface {
	List(ctx context.Context, tenantID string, filter Filter, role string) ([]article.Article, error)
	GetByID(ctx context.Context, id, tenantID string, role string) (*article.Article, error)
	Create(ctx context.Context, a *article.Article) error
	Update(ctx context.Context, a *article.Article) error
	Delete(ctx context.Context, id, tenantID string) error
}

// Service is the application service for article operations.
type Service struct {
	repo Repository
}

// NewService creates a new article application service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// List returns articles for the tenant, filtered by role visibility.
func (s *Service) List(ctx context.Context, tenantID string, filter Filter, role string) ([]article.Article, error) {
	return s.repo.List(ctx, tenantID, filter, role)
}

// Get returns a single article, checking role-based visibility.
func (s *Service) Get(ctx context.Context, id, tenantID, role string) (*article.Article, error) {
	return s.repo.GetByID(ctx, id, tenantID, role)
}

// CreateInput holds the data for creating a new article.
type CreateInput struct {
	TenantID  string
	Title     string
	Content   string
	Category  string
	Tags      []string
	Published bool
	CreatedBy string
}

// Create creates a new article.
func (s *Service) Create(ctx context.Context, input CreateInput) (*article.Article, error) {
	if input.Title == "" || input.Content == "" {
		return nil, article.ErrValidation
	}
	category := input.Category
	if category == "" {
		category = "general"
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}

	a := &article.Article{
		TenantID:  input.TenantID,
		Title:     input.Title,
		Content:   input.Content,
		Category:  category,
		Tags:      tags,
		Published: input.Published,
		CreatedBy: input.CreatedBy,
	}

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

// UpdateInput holds the data for updating an article.
type UpdateInput struct {
	ID        string
	TenantID  string
	Title     string
	Content   string
	Category  string
	Tags      []string
	Published bool
}

// Update modifies an existing article.
func (s *Service) Update(ctx context.Context, input UpdateInput) (*article.Article, error) {
	if input.Title == "" || input.Content == "" {
		return nil, article.ErrValidation
	}
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}

	a := &article.Article{
		ID:        input.ID,
		TenantID:  input.TenantID,
		Title:     input.Title,
		Content:   input.Content,
		Category:  input.Category,
		Tags:      tags,
		Published: input.Published,
	}

	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

// Delete removes an article.
func (s *Service) Delete(ctx context.Context, id, tenantID string) error {
	return s.repo.Delete(ctx, id, tenantID)
}
