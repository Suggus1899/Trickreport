package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestJWTGenerator_GenerateAndValidate_RoundTrip(t *testing.T) {
	gen := NewJWTGenerator("test-secret", 24)
	userID := uuid.New()
	tenantID := uuid.New()

	token, err := gen.Generate(userID, tenantID, "agent")
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}

	claims, err := gen.Validate(token)
	if err != nil {
		t.Fatalf("Validate error: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("userID = %v, want %v", claims.UserID, userID)
	}
	if claims.TenantID != tenantID {
		t.Errorf("tenantID = %v, want %v", claims.TenantID, tenantID)
	}
	if claims.Role != "agent" {
		t.Errorf("role = %q, want agent", claims.Role)
	}
}

func TestJWTGenerator_Validate_InvalidToken(t *testing.T) {
	gen := NewJWTGenerator("test-secret", 24)
	_, err := gen.Validate("not-a-jwt")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestJWTGenerator_Validate_WrongSecret(t *testing.T) {
	gen1 := NewJWTGenerator("secret-1", 24)
	gen2 := NewJWTGenerator("secret-2", 24)

	token, _ := gen1.Generate(uuid.New(), uuid.New(), "admin")
	_, err := gen2.Validate(token)
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestJWTGenerator_Validate_EmptyToken(t *testing.T) {
	gen := NewJWTGenerator("test-secret", 24)
	_, err := gen.Validate("")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestNewJWTGenerator(t *testing.T) {
	gen := NewJWTGenerator("secret", 1)
	if gen == nil {
		t.Fatal("expected non-nil generator")
	}
}

// TestJWTGenerator_Validate_NonUUIDUserID creates a token with a non-UUID
// user ID to exercise the uuid.Parse error path in Validate.
func TestJWTGenerator_Validate_NonUUIDUserID(t *testing.T) {
	gen := NewJWTGenerator("test-secret", 24)
	claims := jwtClaims{
		UserID:   "not-a-uuid",
		TenantID: uuid.New().String(),
		Role:     "agent",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("test-secret"))

	_, err := gen.Validate(signed)
	if err == nil {
		t.Fatal("expected error for non-UUID user ID")
	}
}

// TestJWTGenerator_Validate_NonUUIDTenantID creates a token with a valid
// user ID but non-UUID tenant ID to exercise that error path.
func TestJWTGenerator_Validate_NonUUIDTenantID(t *testing.T) {
	gen := NewJWTGenerator("test-secret", 24)
	claims := jwtClaims{
		UserID:   uuid.New().String(),
		TenantID: "not-a-uuid",
		Role:     "agent",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("test-secret"))

	_, err := gen.Validate(signed)
	if err == nil {
		t.Fatal("expected error for non-UUID tenant ID")
	}
}

// TestJWTGenerator_Validate_WrongSigningMethod creates a token signed with
// a non-HMAC method to exercise the signing method check.
func TestJWTGenerator_Validate_WrongSigningMethod(t *testing.T) {
	gen := NewJWTGenerator("test-secret", 24)
	claims := jwtClaims{
		UserID:   uuid.New().String(),
		TenantID: uuid.New().String(),
		Role:     "agent",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	// Use SigningMethodNone which is not HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	signed, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	_, err := gen.Validate(signed)
	if err == nil {
		t.Fatal("expected error for wrong signing method")
	}
}

// TestJWTGenerator_Validate_ExpiredToken creates an already-expired token.
func TestJWTGenerator_Validate_ExpiredToken(t *testing.T) {
	gen := NewJWTGenerator("test-secret", 24)
	claims := jwtClaims{
		UserID:   uuid.New().String(),
		TenantID: uuid.New().String(),
		Role:     "agent",
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte("test-secret"))

	_, err := gen.Validate(signed)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}
