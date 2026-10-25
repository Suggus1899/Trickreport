package config

import (
	"os"
	"testing"
)

func TestIsProduction(t *testing.T) {
	tests := []struct {
		env  string
		want bool
	}{
		{"production", true},
		{"PRODUCTION", true},
		{"Production", true},
		{"development", false},
		{"staging", false},
		{"", false},
	}
	for _, tt := range tests {
		c := &Config{Env: tt.env}
		if got := c.IsProduction(); got != tt.want {
			t.Errorf("IsProduction(%q) = %v, want %v", tt.env, got, tt.want)
		}
	}
}

func TestAllowedOrigins_Empty(t *testing.T) {
	c := &Config{CORSOrigins: ""}
	if got := c.AllowedOrigins(); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestAllowedOrigins_Single(t *testing.T) {
	c := &Config{CORSOrigins: "http://localhost:3000"}
	got := c.AllowedOrigins()
	if len(got) != 1 || got[0] != "http://localhost:3000" {
		t.Errorf("got %v", got)
	}
}

func TestAllowedOrigins_Multiple(t *testing.T) {
	c := &Config{CORSOrigins: "http://localhost:3000, https://app.example.com , http://localhost:3001"}
	got := c.AllowedOrigins()
	if len(got) != 3 {
		t.Fatalf("got %d origins, want 3", len(got))
	}
	if got[0] != "http://localhost:3000" {
		t.Errorf("got[0] = %q", got[0])
	}
	if got[1] != "https://app.example.com" {
		t.Errorf("got[1] = %q", got[1])
	}
	if got[2] != "http://localhost:3001" {
		t.Errorf("got[2] = %q", got[2])
	}
}

func TestAllowedOrigins_WithEmptyParts(t *testing.T) {
	c := &Config{CORSOrigins: "http://localhost:3000, , , https://app.example.com"}
	got := c.AllowedOrigins()
	if len(got) != 2 {
		t.Fatalf("got %d origins, want 2 (empty parts should be skipped)", len(got))
	}
}

func TestIsProductionEnv(t *testing.T) {
	t.Setenv("ENV", "production")
	if !IsProductionEnv() {
		t.Error("IsProductionEnv() = false, want true for ENV=production")
	}

	t.Setenv("ENV", "development")
	if IsProductionEnv() {
		t.Error("IsProductionEnv() = true, want false for ENV=development")
	}

	t.Setenv("ENV", "")
	if IsProductionEnv() {
		t.Error("IsProductionEnv() = true, want false for ENV=''")
	}
}

func TestLoad_Defaults(t *testing.T) {
	// Clear potentially problematic env vars so viper defaults apply.
	clearEnv(t, "ENV")
	clearEnv(t, "PORT")
	clearEnv(t, "JWT_SECRET")
	clearEnv(t, "ADMIN_PASSWORD")

	cfg := Load()
	if cfg.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Port)
	}
	if cfg.Env != "development" {
		t.Errorf("Env = %q, want development", cfg.Env)
	}
	if cfg.JWTExpHours != 8 {
		t.Errorf("JWTExpHours = %d, want 8", cfg.JWTExpHours)
	}
	if cfg.SecureCookie {
		t.Error("SecureCookie = true, want false in development")
	}
}

// clearEnv unsets key for the duration of the test. t.Setenv records the previous
// value and restores it on cleanup.
func clearEnv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset %s: %v", key, err)
	}
}
