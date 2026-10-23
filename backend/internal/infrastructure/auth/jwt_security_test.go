package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestValidateRejectsAlgorithmConfusionAndMissingExpiry(t *testing.T) {
	gen := NewJWTGenerator("test-secret-at-least-thirty-two-bytes", 8)
	for _, tc := range []struct {
		name   string
		method jwt.SigningMethod
		expiry *jwt.NumericDate
	}{
		{"wrong HMAC algorithm", jwt.SigningMethodHS384, jwt.NewNumericDate(time.Now().Add(time.Hour))},
		{"missing expiry", jwt.SigningMethodHS256, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			claims := gen.buildClaims(uuid.New(), uuid.New(), "admin", "access", time.Now().Add(time.Hour))
			claims.ExpiresAt = tc.expiry
			raw, err := jwt.NewWithClaims(tc.method, claims).SignedString([]byte(gen.secret))
			if err != nil {
				t.Fatal(err)
			}
			if _, err := gen.Validate(raw); err == nil {
				t.Fatal("accepted structurally invalid token")
			}
		})
	}
}

func TestMFAChallengeHasSeparatePurposeAndShortExpiry(t *testing.T) {
	gen := NewJWTGenerator("test-secret-at-least-thirty-two-bytes", 8)
	issuer, ok := any(gen).(interface {
		GenerateMFA(uuid.UUID, uuid.UUID, string) (string, error)
	})
	if !ok {
		t.Fatal("MFA challenges use the unrestricted access-token generator")
	}
	raw, err := issuer.GenerateMFA(uuid.New(), uuid.New(), "admin")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := gen.Validate(raw)
	if err != nil {
		t.Fatal(err)
	}
	if claims.TokenType != "mfa" || time.Until(claims.ExpiresAt) > 5*time.Minute || time.Until(claims.ExpiresAt) < 4*time.Minute {
		t.Fatal("MFA challenge has incorrect purpose or lifetime")
	}
}
