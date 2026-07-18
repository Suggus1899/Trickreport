package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validate is a shared validator instance used across HTTP handlers.
var Validate = validator.New()

// ValidationError maps field validation errors to a simple string slice.
func ValidationError(err error) []string {
	if err == nil {
		return nil
	}

	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		msgs := make([]string, 0, len(verr))
		for _, e := range verr {
			msgs = append(msgs, fmt.Sprintf("%s: %s", strings.ToLower(e.Field()), e.Tag()))
		}
		return msgs
	}
	return []string{err.Error()}
}
