package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	authapp "github.com/trickreport/backend/internal/application/auth"
)

const (
	// jwtIssuer is the expected issuer claim value.
	jwtIssuer = "trickreport"
	// jwtAudience is the expected audience claim value.
	jwtAudience = "trickreport-api"
	// refreshExpHours is the lifetime of refresh tokens in days (7 days).
	refreshExpDays = 7
)

// jwtClaims is the internal claims struct embedding jwt.RegisteredClaims
// with the custom fields used by the application.
type jwtClaims struct {
	UserID    string `json:"sub"`
	TenantID  string `json:"tenant_id"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"` // "access" or "refresh"
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

// buildClaims constructs the internal jwtClaims with all standard registered
// claims populated: jti (UUID), iss, aud, iat, and exp.
func (g *JWTGenerator) buildClaims(userID, tenantID uuid.UUID, role, tokenType string, exp time.Time) jwtClaims {
	now := time.Now()
	return jwtClaims{
		UserID:    userID.String(),
		TenantID:  tenantID.String(),
		Role:      role,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Issuer:    jwtIssuer,
			Audience:  []string{jwtAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
}

// Generate creates a new HS256-signed access JWT for the given user, tenant,
// and role. The token includes jti, iss, aud, and iat claims.
func (g *JWTGenerator) Generate(userID, tenantID uuid.UUID, role string) (string, error) {
	now := time.Now()
	claims := g.buildClaims(userID, tenantID, role, "access", now.Add(time.Duration(g.expHours)*time.Hour))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(g.secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// GenerateRefresh creates a longer-lived refresh JWT (7 days) for the given
// user, tenant, and role. The token includes jti, iss, aud, and iat claims.
func (g *JWTGenerator) GenerateRefresh(userID, tenantID uuid.UUID, role string) (string, error) {
	now := time.Now()
	claims := g.buildClaims(userID, tenantID, role, "refresh", now.Add(refreshExpDays*24*time.Hour))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(g.secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

// Validate parses and validates the given JWT, returning the application
// claims. It verifies the signing method, issuer, and audience.
func (g *JWTGenerator) Validate(token string) (*authapp.Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(g.secret), nil
	}, jwt.WithIssuer(jwtIssuer), jwt.WithAudience(jwtAudience))
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

	result := &authapp.Claims{
		UserID:    userID,
		TenantID:  tenantID,
		Role:      claims.Role,
		JTI:       claims.ID,
		TokenType: claims.TokenType,
	}
	if claims.ExpiresAt != nil {
		result.ExpiresAt = claims.ExpiresAt.Time
	}
	return result, nil
}

// Compile-time interface assertion.
var _ authapp.TokenGenerator = (*JWTGenerator)(nil)
