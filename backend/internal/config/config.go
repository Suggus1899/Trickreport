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
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Force COOKIE_SECURE=true in production unless explicitly disabled
	if cfg.IsProduction() && !cfg.SecureCookie {
		log.Println("WARNING: COOKIE_SECURE is false in production — forcing true")
		cfg.SecureCookie = true
	}

	// Fail fast in production if JWT_SECRET is not set or is a placeholder
	if cfg.IsProduction() {
		if cfg.JWTSecret == "" || strings.Contains(strings.ToLower(cfg.JWTSecret), "change") {
			log.Fatal("JWT_SECRET must be set to a strong random value in production")
		}
		if cfg.AdminPassword == "" || strings.EqualFold(cfg.AdminPassword, "change_me") {
			log.Fatal("ADMIN_PASSWORD must be set in production")
		}
	}

	return cfg
}

// IsProductionEnv is a package-level helper for early checks.
func IsProductionEnv() bool {
	return strings.EqualFold(os.Getenv("ENV"), "production")
}
