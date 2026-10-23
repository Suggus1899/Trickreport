package postgres

import (
	"context"
	"github.com/pquerna/otp/totp"
	"sync"
	"testing"
	"time"

	appAuth "github.com/trickreport/backend/internal/application/auth"
	infraAuth "github.com/trickreport/backend/internal/infrastructure/auth"
)

func sessionTestService(t *testing.T) (*appAuth.Service, *infraAuth.JWTGenerator, appAuth.LoginInput) {
	t.Helper()
	skipIfNoDB(t)
	ctx := context.Background()
	tenantID, userID := seedTenantUser(t, ctx)
	t.Cleanup(func() { cleanupTenant(ctx, tenantID) })
	hasher := infraAuth.NewBcryptHasher()
	hash, err := hasher.Hash("SessionPassword123!")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := testPool.Exec(ctx, `UPDATE users SET password=$1 WHERE id=$2`, hash, userID); err != nil {
		t.Fatal(err)
	}
	gen := infraAuth.NewJWTGenerator("session-test-secret-at-least-32-bytes", 8)
	svc := appAuth.NewService(NewUserRepo(testPool), hasher, nil, gen, appAuth.Config{JWTExpHours: 8})
	svc.SetSessionRepository(NewSessionRepo(testPool))
	return svc, gen, appAuth.LoginInput{TenantID: tenantID, Email: "admin-" + tenantID.String()[:8] + "@test.local", Password: "SessionPassword123!"}
}

func TestSessionLogoutRevokesRefreshAfterServiceRestart(t *testing.T) {
	svc, gen, input := sessionTestService(t)
	ctx := context.Background()
	login, err := svc.Login(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.Logout(ctx, login.Token); err != nil {
		t.Fatal(err)
	}
	restarted := appAuth.NewService(NewUserRepo(testPool), infraAuth.NewBcryptHasher(), nil, gen, appAuth.Config{JWTExpHours: 8})
	restarted.SetSessionRepository(NewSessionRepo(testPool))
	if _, err := restarted.Refresh(ctx, appAuth.RefreshInput{RefreshToken: login.RefreshToken}); err == nil {
		t.Fatal("revoked session regained access after service restart")
	}
}

func TestRefreshConcurrentSingleWinner(t *testing.T) {
	svc, _, input := sessionTestService(t)
	login, err := svc.Login(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := svc.Refresh(context.Background(), appAuth.RefreshInput{RefreshToken: login.RefreshToken})
			results <- err
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
		t.Fatalf("refresh successes=%d, want exactly one", successes)
	}
}

func TestAuthenticateAccessRejectsWrongPurposeMissingSessionAndChangedRole(t *testing.T) {
	svc, gen, input := sessionTestService(t)
	ctx := context.Background()
	login, err := svc.Login(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := svc.AuthenticateAccess(ctx, login.Token)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AuthenticateAccess(ctx, login.RefreshToken); err == nil {
		t.Fatal("refresh token authorized API access")
	}
	untracked, err := gen.Generate(claims.UserID, claims.TenantID, claims.Role)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AuthenticateAccess(ctx, untracked); err == nil {
		t.Fatal("unpersisted JWT authorized API access")
	}
	restarted := appAuth.NewService(NewUserRepo(testPool), infraAuth.NewBcryptHasher(), nil, gen, appAuth.Config{JWTExpHours: 8})
	restarted.SetSessionRepository(NewSessionRepo(testPool))
	if _, err := restarted.AuthenticateAccess(ctx, login.Token); err != nil {
		t.Fatalf("active session lost on restart: %v", err)
	}
	if _, err := testPool.Exec(ctx, `UPDATE users SET role='end_user' WHERE id=$1`, claims.UserID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AuthenticateAccess(ctx, login.Token); err == nil {
		t.Fatal("stale role authorized API access")
	}
	if _, err := svc.Refresh(ctx, appAuth.RefreshInput{RefreshToken: login.RefreshToken}); err == nil {
		t.Fatal("stale role regained access through refresh")
	}
}

func TestMFAChallengeSingleUseAndAttemptLimit(t *testing.T) {
	svc, _, input := sessionTestService(t)
	ctx := context.Background()
	users := NewUserRepo(testPool)
	svc.SetMFARepository(users)
	svc.SetMFAChallengeRepository(NewMFAChallengeRepo(testPool))
	u, err := users.GetByEmail(ctx, input.TenantID, input.Email)
	if err != nil {
		t.Fatal(err)
	}
	const secret = "JBSWY3DPEHPK3PXP"
	if _, err := testPool.Exec(ctx, `UPDATE users SET mfa_secret=$1,mfa_enabled=TRUE WHERE id=$2`, secret, u.ID); err != nil {
		t.Fatal(err)
	}
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	login, err := svc.Login(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	if !login.RequiresMFA || login.Token != "" {
		t.Fatal("MFA login leaked access token")
	}
	if _, err := svc.AuthenticateAccess(ctx, login.MFAToken); err == nil {
		t.Fatal("MFA challenge authorized API access")
	}
	result, err := svc.MFALogin(ctx, appAuth.MFALoginInput{MFAToken: login.MFAToken, Code: code})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AuthenticateAccess(ctx, result.Token); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.MFALogin(ctx, appAuth.MFALoginInput{MFAToken: login.MFAToken, Code: code}); err == nil {
		t.Fatal("MFA challenge replay succeeded")
	}
	challenge, err := svc.Login(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	for range 5 {
		if _, err := svc.MFALogin(ctx, appAuth.MFALoginInput{MFAToken: challenge.MFAToken, Code: "invalid"}); err == nil {
			t.Fatal("invalid MFA code accepted")
		}
	}
	if _, err := svc.MFALogin(ctx, appAuth.MFALoginInput{MFAToken: challenge.MFAToken, Code: code}); err == nil {
		t.Fatal("MFA challenge accepted after five failures")
	}
}

func TestMFAEnableRejectsClientSuppliedSecret(t *testing.T) {
	svc, _, input := sessionTestService(t)
	_ = svc
	ctx := context.Background()
	users := NewUserRepo(testPool)
	u, err := users.GetByEmail(ctx, input.TenantID, input.Email)
	if err != nil {
		t.Fatal(err)
	}
	if err := users.SetMFASecret(ctx, u.ID, "JBSWY3DPEHPK3PXP"); err != nil {
		t.Fatal(err)
	}
	const attackerSecret = "KRSXG5DSNFXGOIDB"
	code, err := totp.GenerateCode(attackerSecret, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	mfa := appAuth.NewMFAService(users, appAuth.MFAConfig{})
	if err := mfa.Enable(ctx, u.ID, attackerSecret, code); err == nil {
		t.Fatal("enabled MFA using a secret that was never stored")
	}
}
