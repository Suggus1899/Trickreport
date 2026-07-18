package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	authapp "github.com/trickreport/backend/internal/application/auth"
)

// jwtClaims is the internal claims struct embedding jwt.RegisteredClaims
// with the custom fields used by the application.
type jwtClaims struct {
	UserID   string `json:"sub"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTGenerator implements auth.TokenGenerator using HS256-signed JWTs.
type JWTGenerator struct {
	secret   string
	expHours int
}

// NewJWTGenerator creates a new JWTGenerator with the given signing secret
// and token expiration in hours.
func NewJWTGenerator(secret string, expHours int) *JWTGenerator {
	return &JWTGenerator{
		secret:   secret,
		expHours: expHours,
	}
}

// Generate creates a new HS256-signed JWT for the given user, tenant, and role.
func (g *JWTGenerator) Generate(userID, tenantID uuid.UUID, role string) (string, error) {
	now := time.Now()
	claims := jwtClaims{
		UserID:   userID.String(),
		TenantID: tenantID.String(),
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(g.expHours) * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(g.secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// Validate parses and validates the given JWT, returning the application
// claims. It verifies that the token is signed with the HMAC method.
func (g *JWTGenerator) Validate(token string) (*authapp.Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(g.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*jwtClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token claims")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user id in token claims")
	}
	tenantID, err := uuid.Parse(claims.TenantID)
	if err != nil {
		return nil, errors.New("invalid tenant id in token claims")
	}

	return &authapp.Claims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     claims.Role,
	}, nil
}

// Compile-time interface assertion.
var _ authapp.TokenGenerator = (*JWTGenerator)(nil)
