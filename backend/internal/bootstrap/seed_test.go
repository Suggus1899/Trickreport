package bootstrap

import (
	"context"
	"testing"
)

func TestSeed_RefusesProduction(t *testing.T) {
	// The production guard must trip before touching the pool at all —
	// passing a nil pool here proves it never gets that far.
	err := Seed(context.Background(), nil, true)
	if err == nil {
		t.Fatal("expected an error when isProduction is true, got nil")
	}
}
