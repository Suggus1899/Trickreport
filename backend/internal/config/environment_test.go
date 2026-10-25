package config

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestLoadEnvironmentWithoutDotenv(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("ENV", "development")
	t.Setenv("DATABASE_URL", "postgres://test@127.0.0.1:55432/trickreport_test?sslmode=disable")
	t.Setenv("JWT_SECRET", strings.Repeat("x", 48))
	t.Setenv("ADMIN_EMAIL", "admin@example.test")
	t.Setenv("ADMIN_PASSWORD", "TestAdminPassword123!")
	t.Setenv("SMTP_HOST", "smtp.example.test")
	t.Setenv("LDAP_HOST", "ldap.example.test")
	cfg := Load()
	if cfg.DatabaseURL == "" || cfg.JWTSecret == "" || cfg.AdminEmail != "admin@example.test" || cfg.AdminPassword == "" || cfg.Email.SMTPHost != "smtp.example.test" || cfg.LDAPHost != "ldap.example.test" {
		t.Fatal("environment-only deployment lost required settings")
	}
}

func TestLoadRejectsWeakDeployedSecret(t *testing.T) {
	if os.Getenv("TRICKREPORT_CONFIG_TEST_CHILD") == "1" {
		Load()
		return
	}
	for _, env := range []string{"production", "staging"} {
		t.Run(env, func(t *testing.T) {
			exe, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(exe, "-test.run=^TestLoadRejectsWeakDeployedSecret$")
			cmd.Dir = t.TempDir()
			cmd.Env = append(os.Environ(), "TRICKREPORT_CONFIG_TEST_CHILD=1", "ENV="+env, "JWT_SECRET=short", "ADMIN_PASSWORD=TestAdminPassword123!")
			out, err := cmd.CombinedOutput()
			if err == nil {
				t.Fatal("deployed environment accepted a weak JWT secret")
			}
			if strings.Contains(string(out), "TestAdminPassword123!") {
				t.Fatal("configuration failure exposed a credential")
			}
		})
	}
}
