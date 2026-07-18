package auth

import (
	"golang.org/x/crypto/bcrypt"

	authapp "github.com/trickreport/backend/internal/application/auth"
	userapp "github.com/trickreport/backend/internal/application/user"
)

// BcryptHasher implements both auth.PasswordHasher and user.PasswordHasher
// using bcrypt.
type BcryptHasher struct{}

// NewBcryptHasher creates a new BcryptHasher.
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

// Hash hashes the given password using bcrypt with the default cost.
func (h *BcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Compare compares a bcrypt hash with a plaintext password.
func (h *BcryptHasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Compile-time interface assertions.
var (
	_ authapp.PasswordHasher = (*BcryptHasher)(nil)
	_ userapp.PasswordHasher = (*BcryptHasher)(nil)
)
