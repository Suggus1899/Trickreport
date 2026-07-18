package user

import "errors"

// Password complexity requirements.
const (
	minPasswordLength = 12
)

// Password complexity errors.
var (
	ErrPasswordTooShort       = errors.New("password must be at least 12 characters long")
	ErrPasswordNoUppercase    = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLowercase    = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber       = errors.New("password must contain at least one number")
	ErrPasswordNoSymbol       = errors.New("password must contain at least one symbol")
)

// ValidatePasswordComplexity checks that a password meets the minimum
// complexity requirements: at least 12 characters, with at least one
// uppercase letter, one lowercase letter, one number, and one symbol.
// It returns a descriptive error when a requirement is not met.
func ValidatePasswordComplexity(password string) error {
	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSymbol  bool
	)
	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasNumber = true
		default:
			// Any character that is not a letter or digit counts as a symbol.
			hasSymbol = true
		}
	}

	if !hasUpper {
		return ErrPasswordNoUppercase
	}
	if !hasLower {
		return ErrPasswordNoLowercase
	}
	if !hasNumber {
		return ErrPasswordNoNumber
	}
	if !hasSymbol {
		return ErrPasswordNoSymbol
	}

	return nil
}
