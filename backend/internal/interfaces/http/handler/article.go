package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	appArticle "github.com/trickreport/backend/internal/application/article"
	"github.com/trickreport/backend/internal/domain/article"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/interfaces/http/validator"
)

type ArticleHandler struct {
	svc *appArticle.Service
}

func NewArticleHandler(svc *appArticle.Service) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

type ArticleDTO struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	Published bool      `json:"published"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AuthorName string   `json:"author_name,omitempty"`
}

type createArticleReq struct {
	Title     string   `json:"title" validate:"required,min=3,max=255"`
	Content   string   `json:"content" validate:"required,min=10,max=20000"`
	Category  string   `json:"category" validate:"required,max=100"`
	Tags      []string `json:"tags" validate:"omitempty,dive,max=50"`
	Published bool     `json:"published"`
}

type updateArticleReq struct {
	Title     string   `json:"title" validate:"required,min=3,max=255"`
	Content   string   `json:"content" validate:"required,min=10,max=20000"`
	Category  string   `json:"category" validate:"required,max=100"`
	Tags      []string `json:"tags" validate:"omitempty,dive,max=50"`
	Published bool     `json:"published"`
}

func toArticleDTO(a *article.Article) ArticleDTO {
	return ArticleDTO{
		ID: a.ID, TenantID: a.TenantID, Title: a.Title, Content: a.Content,
		Category: a.Category, Tags: a.Tags, Published: a.Published,
		CreatedBy: a.CreatedBy, CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt,
		AuthorName: a.AuthorName,
	}
}

func (h *ArticleHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())

	articles, err := h.svc.List(r.Context(), tenantID, appArticle.Filter{
		Search: r.URL.Query().Get("q"),
	}, claims.Role)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list articles")
		return
	}

	dtos := make([]ArticleDTO, len(articles))
	for i, a := range articles {
		dtos[i] = toArticleDTO(&a)
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *ArticleHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())
	id := chi.URLParam(r, "id")

	a, err := h.svc.Get(r.Context(), id, tenantID, claims.Role)
	if err != nil {
		if errors.Is(err, article.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "article not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to get article")
		return
	}

	response.JSON(w, http.StatusOK, toArticleDTO(a))
}

func (h *ArticleHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.ClaimsFromContext(r.Context())
	tenantID := middleware.TenantFromContext(r.Context())

	var req createArticleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	a, err := h.svc.Create(r.Context(), appArticle.CreateInput{
		TenantID:  tenantID,
		Title:     req.Title,
		Content:   req.Content,
		Category:  req.Category,
		Tags:      req.Tags,
		Published: req.Published,
		CreatedBy: claims.UserID,
	})
	if err != nil {
		if errors.Is(err, article.ErrValidation) {
			response.Error(w, http.StatusBadRequest, "title and content are required")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create article")
		return
	}

	response.JSON(w, http.StatusCreated, toArticleDTO(a))
}

func (h *ArticleHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req updateArticleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	a, err := h.svc.Update(r.Context(), appArticle.UpdateInput{
		ID: id, TenantID: tenantID, Title: req.Title, Content: req.Content,
		Category: req.Category, Tags: req.Tags, Published: req.Published,
	})
	if err != nil {
		if errors.Is(err, article.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "article not found")
			return
		}
		if errors.Is(err, article.ErrValidation) {
			response.Error(w, http.StatusBadRequest, "title and content are required")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update article")
		return
	}

	response.JSON(w, http.StatusOK, toArticleDTO(a))
}

func (h *ArticleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), id, tenantID); err != nil {
		if errors.Is(err, article.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "article not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete article")
		return
	}

	response.NoContent(w)
}
