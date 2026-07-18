package automation

import "errors"

var (
	// ErrNotFound is returned when an automation rule does not exist.
	ErrNotFound = errors.New("automation rule not found")

	// ErrValidation is returned when input validation fails.
	ErrValidation = errors.New("validation error")
)
