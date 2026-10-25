package config

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"PORT"`
	Env          string `mapstructure:"ENV"`
	DatabaseURL  string `mapstructure:"DATABASE_URL"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	JWTExpHours  int    `mapstructure:"JWT_EXP_HOURS"`
	SecureCookie bool   `mapstructure:"COOKIE_SECURE"`
	CORSOrigins  string `mapstructure:"CORS_ORIGINS"`
	AppPublicURL string `mapstructure:"APP_PUBLIC_URL"`

	// Bootstrap admin (created on first startup if no users exist)
	AdminEmail    string `mapstructure:"ADMIN_EMAIL"`
	AdminPassword string `mapstructure:"ADMIN_PASSWORD"`

	LDAPEnabled bool   `mapstructure:"LDAP_ENABLED"`
	LDAPHost    string `mapstructure:"LDAP_HOST"`
	LDAPPort    int    `mapstructure:"LDAP_PORT"`
	LDAPUseTLS  bool   `mapstructure:"LDAP_USE_TLS"`
	LDAPBindDN  string `mapstructure:"LDAP_BIND_DN"`
	LDAPBindPW  string `mapstructure:"LDAP_BIND_PW"`
	LDAPBaseDN  string `mapstructure:"LDAP_BASE_DN"`
	LDAPUIDAttr string `mapstructure:"LDAP_UID_ATTR"`

	// Email (SMTP)
	Email EmailConfig `mapstructure:",squash"`
}

// EmailConfig holds SMTP settings for sending transactional email.
type EmailConfig struct {
	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPUsername string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom     string `mapstructure:"SMTP_FROM"`
	SMTPUseTLS   bool   `mapstructure:"SMTP_USE_TLS"`
}

// IsProduction returns true when ENV=production.
func (c *Config) IsProduction() bool {
	return strings.EqualFold(c.Env, "production")
}

// AllowedOrigins parses CORS_ORIGINS (comma-separated) into a slice.
// Returns a normalized set (trimmed, lowercased).
func (c *Config) AllowedOrigins() []string {
	raw := c.CORSOrigins
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, strings.ToLower(p))
		}
	}
	return out
}

func Load() *Config {
	v := newViper()

	if err := v.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	deployed := cfg.IsProduction() || strings.EqualFold(cfg.Env, "staging")
	if deployed && !cfg.SecureCookie {
		cfg.SecureCookie = true
	}
	if deployed {
		validateDeployedSettings(cfg)
	}

	if cfg.JWTExpHours <= 0 || cfg.JWTExpHours > 24 {
		log.Fatal("JWT_EXP_HOURS must be between 1 and 24")
	}

	return cfg
}

// newViper returns a reader that resolves configuration from .env and the process environment.
func newViper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// AutomaticEnv does not expose unset-file keys to Unmarshal.
	for _, key := range []string{
		"PORT", "ENV", "DATABASE_URL", "JWT_SECRET", "JWT_EXP_HOURS", "COOKIE_SECURE", "CORS_ORIGINS", "APP_PUBLIC_URL",
		"ADMIN_EMAIL", "ADMIN_PASSWORD", "LDAP_ENABLED", "LDAP_HOST", "LDAP_PORT", "LDAP_USE_TLS",
		"LDAP_BIND_DN", "LDAP_BIND_PW", "LDAP_BASE_DN", "LDAP_UID_ATTR", "SMTP_HOST", "SMTP_PORT",
		"SMTP_USERNAME", "SMTP_PASSWORD", "SMTP_FROM", "SMTP_USE_TLS",
	} {
		if err := v.BindEnv(key); err != nil {
			log.Fatal("Failed to register configuration environment key")
		}
	}

	// Defaults
	v.SetDefault("PORT", "8080")
	v.SetDefault("ENV", "development")
	v.SetDefault("JWT_EXP_HOURS", 8)
	// Default to true in production, false in dev — resolved after load
	v.SetDefault("COOKIE_SECURE", false)
	v.SetDefault("CORS_ORIGINS", "http://localhost:3000")
	v.SetDefault("APP_PUBLIC_URL", "http://localhost:3000")
	v.SetDefault("LDAP_ENABLED", false)
	v.SetDefault("LDAP_PORT", 636)
	v.SetDefault("LDAP_USE_TLS", true)
	v.SetDefault("LDAP_UID_ATTR", "uid")
	v.SetDefault("SMTP_PORT", 587)
	v.SetDefault("SMTP_USE_TLS", true)
	v.SetDefault("SMTP_FROM", "noreply@trickreport.local")
	return v
}

// validateDeployedSettings fails fast when a deployed environment has an unsafe public URL, JWT secret or admin password.
func validateDeployedSettings(cfg *Config) {
	publicURL, err := url.Parse(cfg.AppPublicURL)
	if err != nil || publicURL.Scheme != "https" || publicURL.Host == "" || publicURL.User != nil || publicURL.RawQuery != "" || publicURL.Fragment != "" {
		log.Fatal("APP_PUBLIC_URL must be an absolute HTTPS origin in deployed environments")
	}
	if len(cfg.JWTSecret) < 32 || strings.Contains(strings.ToLower(cfg.JWTSecret), "change") {
		log.Fatal("JWT_SECRET must contain at least 32 random bytes in deployed environments")
	}
	if cfg.AdminPassword == "" || strings.EqualFold(cfg.AdminPassword, "change_me") {
		log.Fatal("ADMIN_PASSWORD must be set in production")
	}
}

// IsProductionEnv is a package-level helper for early checks.
func IsProductionEnv() bool {
	return strings.EqualFold(os.Getenv("ENV"), "production")
}
