package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func (g *JWTGenerator) GenerateMFA(userID, tenantID uuid.UUID, role string) (string, error) {
	claims := g.buildClaims(userID, tenantID, role, "mfa", time.Now().Add(5*time.Minute))
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(g.secret))
}
