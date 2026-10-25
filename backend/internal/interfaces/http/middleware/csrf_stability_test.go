package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCSRFSafeReadsPreserveTokenForConcurrentWrites(t *testing.T) {
	const token = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	h := CSRF(false)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	read := httptest.NewRequest(http.MethodGet, "/api/v1/auth/csrf", nil)
	read.AddCookie(&http.Cookie{Name: csrfCookieName, Value: token})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, read)
	if got := w.Header().Get(csrfHeaderName); got != token {
		t.Fatal("GET rotated the CSRF token used by another request")
	}
	for _, cookie := range w.Result().Cookies() {
		if cookie.Name == csrfCookieName && cookie.Value != token {
			t.Fatal("GET replaced existing CSRF cookie")
		}
	}
	write := httptest.NewRequest(http.MethodPost, "/api/v1/tickets/", strings.NewReader("{}"))
	write.AddCookie(&http.Cookie{Name: csrfCookieName, Value: token})
	write.Header.Set(csrfHeaderName, token)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, write)
	if w.Code != http.StatusNoContent {
		t.Fatalf("SSR double-submit write returned %d", w.Code)
	}
}
