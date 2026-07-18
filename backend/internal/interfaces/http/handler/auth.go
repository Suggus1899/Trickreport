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
	svc      *appAuth.Service
	mfaSvc   *appAuth.MFAService
	resetSvc *appAuth.PasswordResetService
}

// NewAuthHandler creates a new AuthHandler with the core auth service.
// Use the Set* methods to wire optional services (MFA, password reset).
func NewAuthHandler(svc *appAuth.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// SetMFAService wires the MFA service.
func (h *AuthHandler) SetMFAService(svc *appAuth.MFAService) { h.mfaSvc = svc }

// SetPasswordResetService wires the password reset service.
func (h *AuthHandler) SetPasswordResetService(svc *appAuth.PasswordResetService) { h.resetSvc = svc }

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
	Token        string   `json:"token"`
	RefreshToken string   `json:"refresh_token,omitempty"`
	User         userInfo `json:"user"`
	RequiresMFA  bool     `json:"requires_mfa,omitempty"`
	MFAToken     string   `json:"mfa_token,omitempty"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// @Summary      Login
	// @Description  Authenticate a user with email and password, returning a JWT.
	// @Tags         auth
	// @Accept       json
	// @Produce      json
	// @Param        body  body      loginReq  true  "Login credentials"
	// @Success      200   {object}  loginResp
	// @Failure      400   {object}  response.ErrorBody
	// @Failure      401   {object}  response.ErrorBody
	// @Router       /auth/login [post]
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
		switch {
		case errors.Is(err, appAuth.ErrInvalidCredentials), errors.Is(err, appAuth.ErrPasswordNotConfigured):
			response.Error(w, http.StatusUnauthorized, "invalid credentials")
		case errors.Is(err, appAuth.ErrAccountLocked):
			response.Error(w, http.StatusTooManyRequests, "account is temporarily locked")
		default:
			response.Error(w, http.StatusInternalServerError, "authentication error")
		}
		return
	}

	// MFA challenge — return a temporary token instead of full credentials
	if result.RequiresMFA {
		response.JSON(w, http.StatusOK, loginResp{
			RequiresMFA: true,
			MFAToken:    result.MFAToken,
		})
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
	if result.RefreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "trickreport_refresh",
			Value:    result.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   result.SecureCookie,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})
	}

	response.JSON(w, http.StatusOK, loginResp{
		Token:        result.Token,
		RefreshToken: result.RefreshToken,
		User: userInfo{
			ID:       result.User.ID,
			Name:     result.User.Name,
			Email:    result.User.Email,
			Role:     string(result.User.Role),
			TenantID: result.User.TenantID,
		},
	})
}

// refreshReq is the request body for the refresh endpoint.
type refreshReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// refreshResp is the response body for the refresh endpoint.
type refreshResp struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// Refresh accepts a refresh token and returns new access + refresh tokens.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	// Fall back to cookie if not in body
	if req.RefreshToken == "" {
		if c, err := r.Cookie("trickreport_refresh"); err == nil {
			req.RefreshToken = c.Value
		}
	}

	result, err := h.svc.Refresh(r.Context(), appAuth.RefreshInput{
		RefreshToken: req.RefreshToken,
		IPAddress:    r.RemoteAddr,
		UserAgent:    r.UserAgent(),
	})
	if err != nil {
		switch {
		case errors.Is(err, appAuth.ErrTokenBlacklisted):
			response.Error(w, http.StatusUnauthorized, "token has been revoked")
		case errors.Is(err, appAuth.ErrInvalidRefreshToken):
			response.Error(w, http.StatusUnauthorized, "invalid refresh token")
		default:
			response.Error(w, http.StatusInternalServerError, "token refresh error")
		}
		return
	}

	response.JSON(w, http.StatusOK, refreshResp{
		Token:        result.Token,
		RefreshToken: result.RefreshToken,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// @Summary      Get current user
	// @Description  Returns the profile of the authenticated user.
	// @Tags         auth
	// @Produce      json
	// @Security     BearerAuth
	// @Success      200  {object}  userInfo
	// @Failure      401  {object}  response.ErrorBody
	// @Failure      404  {object}  response.ErrorBody
	// @Router       /auth/me [get]
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
	// @Summary      Logout
	// @Description  Clears the authentication cookie.
	// @Tags         auth
	// @Success      204
	// @Router       /auth/logout [post]
	// Try to revoke the session via the access token
	token := ""
	if authHeader := r.Header.Get("Authorization"); len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}
	if token == "" {
		if c, err := r.Cookie("trickreport_token"); err == nil {
			token = c.Value
		}
	}
	if token != "" {
		_ = h.svc.Logout(r.Context(), token)
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "trickreport_token",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})
	http.SetCookie(w, &http.Cookie{
		Name:    "trickreport_refresh",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(0, 0),
	})
	response.NoContent(w)
}

// --- Password Reset ---

type passwordResetReq struct {
	Email string `json:"email" validate:"required,email"`
}

// PasswordReset handles POST /api/v1/auth/password-reset.
func (h *AuthHandler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	if h.resetSvc == nil {
		response.Error(w, http.StatusServiceUnavailable, "password reset not configured")
		return
	}

	var req passwordResetReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	// Always return 202 to avoid leaking which emails are registered
	_ = h.resetSvc.RequestReset(r.Context(), req.Email)
	response.JSON(w, http.StatusAccepted, map[string]string{"status": "if the email exists, a reset link has been sent"})
}

type passwordResetConfirmReq struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

// PasswordResetConfirm handles POST /api/v1/auth/password-reset/confirm.
func (h *AuthHandler) PasswordResetConfirm(w http.ResponseWriter, r *http.Request) {
	if h.resetSvc == nil {
		response.Error(w, http.StatusServiceUnavailable, "password reset not configured")
		return
	}

	var req passwordResetConfirmReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	// Validate password complexity before calling the service
	if err := user.ValidatePasswordComplexity(req.NewPassword); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.resetSvc.ConfirmReset(r.Context(), req.Token, req.NewPassword); err != nil {
		switch {
		case errors.Is(err, appAuth.ErrPasswordResetTokenInvalid):
			response.Error(w, http.StatusBadRequest, "invalid or expired reset token")
		case errors.Is(err, appAuth.ErrPasswordResetTokenUsed):
			response.Error(w, http.StatusBadRequest, "reset token has already been used")
		default:
			// Complexity errors and other validation errors
			response.Error(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "password updated successfully"})
}

// --- MFA ---

type mfaSetupResp struct {
	Secret string `json:"secret"`
	QRURL  string `json:"qr_url"`
}

// MFASetup handles POST /api/v1/auth/mfa/setup.
func (h *AuthHandler) MFASetup(w http.ResponseWriter, r *http.Request) {
	if h.mfaSvc == nil {
		response.Error(w, http.StatusServiceUnavailable, "mfa not configured")
		return
	}

	claims, ok := getClaimsFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "authentication required")
		return
	}

	// Use the user's email from their profile
	u, err := h.svc.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "user not found")
		return
	}

	result, err := h.mfaSvc.Setup(r.Context(), claims.UserID, u.Email)
	if err != nil {
		if errors.Is(err, appAuth.ErrMFAAlreadyEnabled) {
			response.Error(w, http.StatusConflict, "mfa is already enabled")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to setup mfa")
		return
	}

	response.JSON(w, http.StatusOK, mfaSetupResp{Secret: result.Secret, QRURL: result.QRURL})
}

type mfaVerifyReq struct {
	Code string `json:"code" validate:"required"`
}

// MFAVerify handles POST /api/v1/auth/mfa/verify.
func (h *AuthHandler) MFAVerify(w http.ResponseWriter, r *http.Request) {
	if h.mfaSvc == nil {
		response.Error(w, http.StatusServiceUnavailable, "mfa not configured")
		return
	}

	claims, ok := getClaimsFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req mfaVerifyReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	valid, err := h.mfaSvc.Verify(r.Context(), claims.UserID, req.Code)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]bool{"valid": valid})
}

type mfaEnableReq struct {
	Secret string `json:"secret" validate:"required"`
	Code   string `json:"code" validate:"required"`
}

// MFAEnable handles POST /api/v1/auth/mfa/enable.
// This endpoint enables MFA after the user has verified they can generate codes.
func (h *AuthHandler) MFAEnable(w http.ResponseWriter, r *http.Request) {
	if h.mfaSvc == nil {
		response.Error(w, http.StatusServiceUnavailable, "mfa not configured")
		return
	}

	claims, ok := getClaimsFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req mfaEnableReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	if err := h.mfaSvc.Enable(r.Context(), claims.UserID, req.Secret, req.Code); err != nil {
		if errors.Is(err, appAuth.ErrInvalidMFACode) {
			response.Error(w, http.StatusBadRequest, "invalid mfa code")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to enable mfa")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "mfa enabled"})
}

type mfaDisableReq struct {
	Code string `json:"code" validate:"required"`
}

// MFADisable handles POST /api/v1/auth/mfa/disable.
func (h *AuthHandler) MFADisable(w http.ResponseWriter, r *http.Request) {
	if h.mfaSvc == nil {
		response.Error(w, http.StatusServiceUnavailable, "mfa not configured")
		return
	}

	claims, ok := getClaimsFromContext(r)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req mfaDisableReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	// Verify the code before disabling
	valid, err := h.mfaSvc.Verify(r.Context(), claims.UserID, req.Code)
	if err != nil || !valid {
		response.Error(w, http.StatusBadRequest, "invalid mfa code")
		return
	}

	if err := h.mfaSvc.Disable(r.Context(), claims.UserID); err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to disable mfa")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "mfa disabled"})
}

type mfaLoginReq struct {
	MFAToken string `json:"mfa_token" validate:"required"`
	Code     string `json:"code" validate:"required"`
}

// MFALogin handles POST /api/v1/auth/mfa/login.
func (h *AuthHandler) MFALogin(w http.ResponseWriter, r *http.Request) {
	var req mfaLoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{"error": "validation failed", "details": validator.ValidationError(err)})
		return
	}

	result, err := h.svc.MFALogin(r.Context(), appAuth.MFALoginInput{
		MFAToken: req.MFAToken,
		Code:     req.Code,
	})
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid mfa code or token")
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
	if result.RefreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "trickreport_refresh",
			Value:    result.RefreshToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   result.SecureCookie,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})
	}

	response.JSON(w, http.StatusOK, loginResp{
		Token:        result.Token,
		RefreshToken: result.RefreshToken,
		User: userInfo{
			ID:       result.User.ID,
			Name:     result.User.Name,
			Email:    result.User.Email,
			Role:     string(result.User.Role),
			TenantID: result.User.TenantID,
		},
	})
}

// Ensure user import is used (domain types referenced via service)
var _ = user.ErrNotFound
