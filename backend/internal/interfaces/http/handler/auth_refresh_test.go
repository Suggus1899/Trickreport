package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	infraAuth "github.com/trickreport/backend/internal/infrastructure/auth"
)

type refreshSessions struct {
	appAuth.SessionRepository
	active string
}

func (s *refreshSessions) RotateSession(_ context.Context, old string, replacement appAuth.SessionInput) (bool, error) {
	if s.active != old {
		return false, nil
	}
	s.active = replacement.RefreshJTI
	return true, nil
}

func TestAuthHandler_RefreshCookieRotation(t *testing.T) {
	for _, secure := range []bool{false, true} {
		for _, source := range []string{"cookie", "body"} {
			t.Run(source+map[bool]string{false: "-local", true: "-secure"}[secure], func(t *testing.T) {
				u := newTestUser()
				u.Active = true
				tokens := infraAuth.NewJWTGenerator("refresh-handler-test-secret-long-enough", 1)
				old, err := tokens.GenerateRefresh(u.ID, u.TenantID, string(u.Role))
				if err != nil {
					t.Fatal(err)
				}
				claims, err := tokens.Validate(old)
				if err != nil {
					t.Fatal(err)
				}
				sessions := &refreshSessions{active: claims.JTI}
				svc := appAuth.NewService(&authMockUserRepo{u: u}, &authMockHasher{}, nil, tokens, appAuth.Config{SecureCookie: secure})
				svc.SetSessionRepository(sessions)
				h := NewAuthHandler(svc)
				body := "{}"
				if source == "body" {
					body = `{"refresh_token":"` + old + `"}`
				}
				req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", strings.NewReader(body))
				if source == "cookie" {
					req.AddCookie(&http.Cookie{Name: "trickreport_refresh", Value: old})
				}
				rec := httptest.NewRecorder()
				h.Refresh(rec, req)
				if rec.Code != http.StatusOK {
					t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
				}
				cookies := rec.Result().Cookies()
				if len(cookies) != 2 {
					t.Fatalf("expected two rotated cookies, got %v", cookies)
				}
				for _, cookie := range cookies {
					if !cookie.HttpOnly || cookie.Secure != secure || cookie.Path != "/" || cookie.SameSite != http.SameSiteStrictMode {
						t.Errorf("unsafe cookie flags: %+v", cookie)
					}
					c, err := tokens.Validate(cookie.Value)
					if err != nil {
						t.Fatal(err)
					}
					wantType := "access"
					if cookie.Name == "trickreport_refresh" {
						wantType = "refresh"
					} else if cookie.Name != "trickreport_token" {
						t.Fatalf("unexpected cookie %q", cookie.Name)
					}
					if c.TokenType != wantType || !cookie.Expires.Equal(c.ExpiresAt) {
						t.Errorf("cookie lifetime/type differs from JWT: %+v / %+v", cookie, c)
					}
					if cookie.Value == old {
						t.Fatal("refresh token was not rotated")
					}
				}
				replay := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", strings.NewReader(`{"refresh_token":"`+old+`"}`))
				replayed := httptest.NewRecorder()
				h.Refresh(replayed, replay)
				if replayed.Code != http.StatusUnauthorized || len(replayed.Result().Cookies()) != 0 {
					t.Fatalf("replay issued credentials: %d", replayed.Code)
				}
			})
		}
	}
}

func TestAuthHandler_RefreshRejectsInvalidInput(t *testing.T) {
	for _, tt := range []struct {
		name, body, cookie string
		status             int
	}{
		{"missing", "{}", "", http.StatusBadRequest},
		{"invalid-json", "{", "valid-looking-cookie", http.StatusBadRequest},
		{"invalid-cookie", "{}", "invalid", http.StatusUnauthorized},
		{"explicit-invalid-body", `{"refresh_token":"invalid"}`, uuid.NewString(), http.StatusUnauthorized},
	} {
		t.Run(tt.name, func(t *testing.T) {
			h := newAuthHandler(&authMockUserRepo{}, &authMockHasher{}, nil, &authMockTokenGen{})
			r := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", strings.NewReader(tt.body))
			if tt.cookie != "" {
				r.AddCookie(&http.Cookie{Name: "trickreport_refresh", Value: tt.cookie})
			}
			w := httptest.NewRecorder()
			h.Refresh(w, r)
			if w.Code != tt.status || len(w.Result().Cookies()) != 0 {
				t.Fatalf("status=%d want=%d cookies=%v", w.Code, tt.status, w.Result().Cookies())
			}
		})
	}
}
