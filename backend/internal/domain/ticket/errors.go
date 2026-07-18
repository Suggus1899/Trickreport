package ticket

import "errors"

var (
	// ErrNotFound is returned when a ticket does not exist or does not belong to the tenant.
	ErrNotFound = errors.New("ticket not found")

	// ErrInvalidTransition is returned when a status transition is not allowed.
	ErrInvalidTransition = errors.New("invalid status transition")

	// ErrForbidden is returned when a user lacks permission for the operation.
	ErrForbidden = errors.New("forbidden")

	// ErrValidation is returned when input validation fails.
	ErrValidation = errors.New("validation error")
)
