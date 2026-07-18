package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	appUser "github.com/trickreport/backend/internal/application/user"
	"github.com/trickreport/backend/internal/domain/user"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/interfaces/http/validator"
)

type UserHandler struct {
	svc *appUser.Service
}

func NewUserHandler(svc *appUser.Service) *UserHandler {
	return &UserHandler{svc: svc}
}

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type createUserReq struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Role     string `json:"role" validate:"required,oneof=admin agent end_user"`
	Password string `json:"password" validate:"required,min=6,max=255"`
}

type updateUserReq struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Role      *string `json:"role,omitempty" validate:"omitempty,oneof=admin agent end_user"`
	AvatarURL *string `json:"avatar_url,omitempty" validate:"omitempty,url,max=500"`
	Active    *bool   `json:"active,omitempty"`
}

func toUserDTO(u *user.User) UserDTO {
	return UserDTO{
		ID: u.ID, TenantID: u.TenantID, Name: u.Name, Email: u.Email,
		Role: string(u.Role), AvatarURL: u.AvatarURL, Active: u.Active,
		CreatedAt: u.CreatedAt,
	}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	users, err := h.svc.List(r.Context(), tenantID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to retrieve users")
		return
	}

	dtos := make([]UserDTO, len(users))
	for i, u := range users {
		dtos[i] = toUserDTO(&u)
	}
	response.JSON(w, http.StatusOK, dtos)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	u, err := h.svc.Get(r.Context(), id, tenantID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	response.JSON(w, http.StatusOK, toUserDTO(u))
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())

	var req createUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	// Validate password complexity before calling the service
	if err := user.ValidatePasswordComplexity(req.Password); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	u, err := h.svc.Create(r.Context(), appUser.CreateInput{
		TenantID: tenantID,
		Name:     req.Name,
		Email:    req.Email,
		Role:     req.Role,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, user.ErrValidation) {
			response.Error(w, http.StatusBadRequest, "name and email are required")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	response.JSON(w, http.StatusCreated, toUserDTO(u))
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req updateUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	u, err := h.svc.Update(r.Context(), id, tenantID, appUser.UpdateFields{
		Name:      req.Name,
		Role:      req.Role,
		AvatarURL: req.AvatarURL,
		Active:    req.Active,
	})
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	response.JSON(w, http.StatusOK, toUserDTO(u))
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := middleware.TenantFromContext(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.svc.Deactivate(r.Context(), id, tenantID); err != nil {
		if errors.Is(err, user.ErrNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	response.NoContent(w)
}
