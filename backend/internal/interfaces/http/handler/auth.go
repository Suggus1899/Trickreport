package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	"github.com/trickreport/backend/internal/domain/user"
	"github.com/trickreport/backend/internal/interfaces/http/response"
	"github.com/trickreport/backend/internal/interfaces/http/validator"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	svc *appAuth.Service
}

func NewAuthHandler(svc *appAuth.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

type loginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type userInfo struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	TenantID uuid.UUID `json:"tenant_id"`
}

type loginResp struct {
	Token string   `json:"token"`
	User  userInfo `json:"user"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	result, err := h.svc.Login(r.Context(), appAuth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, appAuth.ErrInvalidCredentials) || errors.Is(err, appAuth.ErrPasswordNotConfigured) {
			response.Error(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		response.Error(w, http.StatusInternalServerError, "authentication error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "trickreport_token",
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   result.SecureCookie,
		SameSite: http.SameSiteStrictMode,
		Expires:  result.ExpiresAt,
	})

	response.JSON(w, http.StatusOK, loginResp{
		Token: result.Token,
		User: userInfo{
			ID:       result.User.ID,
			Name:     result.User.Name,
			Email:    result.User.Email,
			Role:     string(result.User.Role),
			TenantID: result.User.TenantID,
		},
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Claims are already validated by middleware; we need the userID
	// We'll use the auth service to fetch the full profile
	// But we need the userID from context — use the middleware's ClaimsFromContext
	claims, ok := getClaimsFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	u, err := h.svc.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	response.JSON(w, http.StatusOK, userInfo{
		ID:       u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Role:     string(u.Role),
		TenantID: u.TenantID,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "trickreport_token",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})
	response.NoContent(w)
}

// Ensure user import is used (domain types referenced via service)
var _ = user.ErrNotFound
