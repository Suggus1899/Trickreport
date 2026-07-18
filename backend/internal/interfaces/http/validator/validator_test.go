package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

type testStruct struct {
	Name  string `validate:"required,min=3"`
	Email string `validate:"required,email"`
}

func TestValidationError_Nil(t *testing.T) {
	if got := ValidationError(nil); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestValidationError_WithFieldErrors(t *testing.T) {
	v := validator.New()
	err := v.Struct(testStruct{Name: "ab", Email: "not-email"})
	msgs := ValidationError(err)
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d: %v", len(msgs), msgs)
	}
	// Each message should be "field: tag" format
	for _, m := range msgs {
		if m == "" {
			t.Error("expected non-empty message")
		}
	}
}

func TestValidationError_NonValidatorError(t *testing.T) {
	msgs := ValidationError(errCustom{})
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
}

type errCustom struct{}

func (errCustom) Error() string { return "custom error" }
