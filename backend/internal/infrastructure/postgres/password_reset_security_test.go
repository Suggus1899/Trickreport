package postgres

import (
	"context"
	"crypto/sha256"
	"fmt"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	infraAuth "github.com/trickreport/backend/internal/infrastructure/auth"
	"net/url"
	"strings"
	"sync"
	"testing"
)

type resetCaptureSender struct{ body string }

func (s *resetCaptureSender) Send(to, subject, body string) error { s.body = body; return nil }

func resetTestToken(t *testing.T) (*appAuth.PasswordResetService, *appAuth.Service, *appAuth.LoginResult, string) {
	t.Helper()
	authSvc, _, input := sessionTestService(t)
	login, err := authSvc.Login(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	sender := &resetCaptureSender{}
	users := NewUserRepo(testPool)
	svc := appAuth.NewPasswordResetService(users, NewPasswordResetRepo(testPool), users, infraAuth.NewBcryptHasher(), sender, appAuth.PasswordResetConfig{ResetURL: "https://support.example.test/password-reset/confirm"})
	if err := svc.RequestReset(context.Background(), input.TenantID, input.Email); err != nil {
		t.Fatal(err)
	}
	for _, word := range strings.Fields(sender.body) {
		if strings.HasPrefix(word, "https://support.example.test/password-reset/confirm?") {
			u, err := url.Parse(word)
			if err != nil {
				t.Fatal(err)
			}
			return svc, authSvc, login, u.Query().Get("token")
		}
	}
	t.Fatal("reset email did not contain an absolute reset link")
	return nil, nil, nil, ""
}

func TestPasswordResetPersistsOnlyTokenDigest(t *testing.T) {
	_, _, _, token := resetTestToken(t)
	var rawCount, hashCount int
	ctx := context.Background()
	if err := testPool.QueryRow(ctx, `SELECT COUNT(*) FROM password_reset_tokens WHERE token=$1`, token).Scan(&rawCount); err != nil {
		t.Fatal(err)
	}
	digest := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
	if err := testPool.QueryRow(ctx, `SELECT COUNT(*) FROM password_reset_tokens WHERE token=$1`, digest).Scan(&hashCount); err != nil {
		t.Fatal(err)
	}
	if rawCount != 0 || hashCount != 1 {
		t.Fatalf("plaintext=%d digest=%d, want 0 and 1", rawCount, hashCount)
	}
}

func TestPasswordResetConcurrentConsumptionRevokesSessions(t *testing.T) {
	svc, authSvc, login, token := resetTestToken(t)
	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			results <- svc.ConfirmReset(context.Background(), token, "ChangedPassword456!")
		}()
	}
	close(start)
	wg.Wait()
	close(results)
	successes := 0
	for err := range results {
		if err == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("reset successes=%d, want 1", successes)
	}
	if _, err := authSvc.AuthenticateAccess(context.Background(), login.Token); err == nil {
		t.Fatal("old session survived password reset")
	}
	if _, err := authSvc.Refresh(context.Background(), appAuth.RefreshInput{RefreshToken: login.RefreshToken}); err == nil {
		t.Fatal("old refresh survived password reset")
	}
}
