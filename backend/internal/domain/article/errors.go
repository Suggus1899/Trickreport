package article

import "errors"

var (
	// ErrNotFound is returned when an article does not exist.
	ErrNotFound = errors.New("article not found")

	// ErrValidation is returned when input validation fails.
	ErrValidation = errors.New("validation error")
)
