package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/trickreport/backend/internal/application/article"
	domainarticle "github.com/trickreport/backend/internal/domain/article"
)

// ArticleRepo implements article.Repository.
type ArticleRepo struct {
	db *pgxpool.Pool
}

// NewArticleRepo creates a new ArticleRepo.
func NewArticleRepo(db *pgxpool.Pool) *ArticleRepo {
	return &ArticleRepo{db: db}
}

// Compile-time assertion that ArticleRepo implements article.Repository.
var _ article.Repository = (*ArticleRepo)(nil)

const articleSelectColumns = `a.id, a.tenant_id, a.title, a.content, a.category, a.tags, a.published, a.created_by, a.created_at, a.updated_at, u.name AS author_name`

// List returns articles for a tenant, optionally filtered by search and role.
func (r *ArticleRepo) List(ctx context.Context, tenantID uuid.UUID, filter article.Filter, role string) ([]domainarticle.Article, error) {
	q := `SELECT ` + articleSelectColumns + ` FROM articles a JOIN users u ON a.created_by = u.id WHERE a.tenant_id = $1`
	args := []any{tenantID}
	argIdx := 2

	if role == "end_user" {
		q += ` AND a.published = TRUE`
	}
	if filter.Search != "" {
		q += fmt.Sprintf(` AND a.title ILIKE $%d`, argIdx)
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}
	q += ` ORDER BY a.created_at DESC`

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("article_repo.List: query: %w", err)
	}
	defer rows.Close()

	var articles []domainarticle.Article
	for rows.Next() {
		a, err := scanArticle(rows)
		if err != nil {
			return nil, err
		}
		articles = append(articles, *a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("article_repo.List: rows: %w", err)
	}
	return articles, nil
}

// GetByID returns a single article by id within a tenant.
func (r *ArticleRepo) GetByID(ctx context.Context, id, tenantID uuid.UUID, role string) (*domainarticle.Article, error) {
	q := `SELECT ` + articleSelectColumns + ` FROM articles a JOIN users u ON a.created_by = u.id WHERE a.id = $1 AND a.tenant_id = $2`
	if role == "end_user" {
		q += ` AND a.published = TRUE`
	}

	row := r.db.QueryRow(ctx, q, id, tenantID)
	a, err := scanArticleRow(row)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Create inserts a new article.
func (r *ArticleRepo) Create(ctx context.Context, a *domainarticle.Article) error {
	const q = `INSERT INTO articles (tenant_id, title, content, category, tags, published, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`

	var id pgtype.UUID
	if err := r.db.QueryRow(ctx, q, a.TenantID, a.Title, a.Content, a.Category, a.Tags, a.Published, a.CreatedBy).
		Scan(&id, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return fmt.Errorf("article_repo.Create: %w", err)
	}
	a.ID = pgToUUID(id)
	return nil
}

// Update updates an existing article.
func (r *ArticleRepo) Update(ctx context.Context, a *domainarticle.Article) error {
	const q = `UPDATE articles SET title = $1, content = $2, category = $3, tags = $4, published = $5, updated_at = NOW() WHERE id = $6 AND tenant_id = $7 RETURNING id, created_at, updated_at`

	var id pgtype.UUID
	if err := r.db.QueryRow(ctx, q, a.Title, a.Content, a.Category, a.Tags, a.Published, a.ID, a.TenantID).
		Scan(&id, &a.CreatedAt, &a.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domainarticle.ErrNotFound
		}
		return fmt.Errorf("article_repo.Update: %w", err)
	}
	a.ID = pgToUUID(id)
	return nil
}

// Delete removes an article by id within a tenant.
func (r *ArticleRepo) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
	const q = `DELETE FROM articles WHERE id = $1 AND tenant_id = $2`

	ct, err := r.db.Exec(ctx, q, id, tenantID)
	if err != nil {
		return fmt.Errorf("article_repo.Delete: %w", err)
	}
	if ct.RowsAffected() == 0 {
		return domainarticle.ErrNotFound
	}
	return nil
}

// scanArticle scans an article from a pgx.Row-like scanner.
func scanArticle(row pgx.Row) (*domainarticle.Article, error) {
	var a domainarticle.Article
	var id, tid, createdBy pgtype.UUID
	if err := row.Scan(&id, &tid, &a.Title, &a.Content, &a.Category, &a.Tags, &a.Published, &createdBy, &a.CreatedAt, &a.UpdatedAt, &a.AuthorName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainarticle.ErrNotFound
		}
		return nil, fmt.Errorf("article_repo: scan: %w", err)
	}
	a.ID = pgToUUID(id)
	a.TenantID = pgToUUID(tid)
	a.CreatedBy = pgToUUID(createdBy)
	return &a, nil
}

// scanArticleRow is a thin wrapper for single-row queries.
func scanArticleRow(row pgx.Row) (*domainarticle.Article, error) {
	return scanArticle(row)
}
