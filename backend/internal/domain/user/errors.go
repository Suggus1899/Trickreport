package user

import "errors"

var (
	// ErrNotFound is returned when a user does not exist.
	ErrNotFound = errors.New("user not found")

	// ErrAlreadyExists is returned when a user with the same email already exists in the tenant.
	ErrAlreadyExists = errors.New("user already exists")

	// ErrValidation is returned when input validation fails.
	ErrValidation = errors.New("validation error")
)
