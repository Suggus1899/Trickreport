package sla

import "errors"

var (
	// ErrNotFound is returned when an SLA policy does not exist.
	ErrNotFound = errors.New("sla policy not found")

	// ErrInvalidPriority is returned when the priority is not a valid enum value.
	ErrInvalidPriority = errors.New("invalid priority")

	// ErrValidation is returned when input validation fails.
	ErrValidation = errors.New("validation error")
)
