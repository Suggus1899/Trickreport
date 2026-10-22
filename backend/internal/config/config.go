package config

import (
	"log"
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
	viper := viper.New()
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// AutomaticEnv does not expose unset-file keys to Unmarshal.
	for _, key := range []string{
		"PORT", "ENV", "DATABASE_URL", "JWT_SECRET", "JWT_EXP_HOURS", "COOKIE_SECURE", "CORS_ORIGINS",
		"ADMIN_EMAIL", "ADMIN_PASSWORD", "LDAP_ENABLED", "LDAP_HOST", "LDAP_PORT", "LDAP_USE_TLS",
		"LDAP_BIND_DN", "LDAP_BIND_PW", "LDAP_BASE_DN", "LDAP_UID_ATTR", "SMTP_HOST", "SMTP_PORT",
		"SMTP_USERNAME", "SMTP_PASSWORD", "SMTP_FROM", "SMTP_USE_TLS",
	} {
		if err := viper.BindEnv(key); err != nil {
			log.Fatal("Failed to register configuration environment key")
		}
	}

	// Defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENV", "development")
	viper.SetDefault("JWT_EXP_HOURS", 8)
	// Default to true in production, false in dev — resolved after load
	viper.SetDefault("COOKIE_SECURE", false)
	viper.SetDefault("CORS_ORIGINS", "http://localhost:4321")
	viper.SetDefault("LDAP_ENABLED", false)
	viper.SetDefault("LDAP_PORT", 636)
	viper.SetDefault("LDAP_USE_TLS", true)
	viper.SetDefault("LDAP_UID_ATTR", "uid")
	viper.SetDefault("SMTP_PORT", 587)
	viper.SetDefault("SMTP_USE_TLS", true)
	viper.SetDefault("SMTP_FROM", "noreply@trickreport.local")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	deployed := cfg.IsProduction() || strings.EqualFold(cfg.Env, "staging")
	if deployed && !cfg.SecureCookie {
		cfg.SecureCookie = true
	}

	// Fail fast in production if JWT_SECRET is not set or is a placeholder
	if deployed {
		if len(cfg.JWTSecret) < 32 || strings.Contains(strings.ToLower(cfg.JWTSecret), "change") {
			log.Fatal("JWT_SECRET must contain at least 32 random bytes in deployed environments")
		}
		if cfg.AdminPassword == "" || strings.EqualFold(cfg.AdminPassword, "change_me") {
			log.Fatal("ADMIN_PASSWORD must be set in production")
		}
	}
	if cfg.JWTExpHours <= 0 || cfg.JWTExpHours > 24 {
		log.Fatal("JWT_EXP_HOURS must be between 1 and 24")
	}

	return cfg
}

// IsProductionEnv is a package-level helper for early checks.
func IsProductionEnv() bool {
	return strings.EqualFold(os.Getenv("ENV"), "production")
}
