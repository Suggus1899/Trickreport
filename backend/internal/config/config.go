package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"PORT"`
	DatabaseURL  string `mapstructure:"DATABASE_URL"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	JWTExpHours  int    `mapstructure:"JWT_EXP_HOURS"`
	SecureCookie bool   `mapstructure:"COOKIE_SECURE"`
	CORSOrigin   string `mapstructure:"CORS_ORIGIN"`

	LDAPEnabled bool   `mapstructure:"LDAP_ENABLED"`
	LDAPHost    string `mapstructure:"LDAP_HOST"`
	LDAPPort    int    `mapstructure:"LDAP_PORT"`
	LDAPBindDN  string `mapstructure:"LDAP_BIND_DN"`
	LDAPBindPW  string `mapstructure:"LDAP_BIND_PW"`
	LDAPBaseDN  string `mapstructure:"LDAP_BASE_DN"`
	LDAPUIDAttr string `mapstructure:"LDAP_UID_ATTR"`
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")

	// Allow environment variables to override .env values
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("JWT_EXP_HOURS", 8)
	viper.SetDefault("COOKIE_SECURE", false)
	viper.SetDefault("CORS_ORIGIN", "http://localhost:4321")
	viper.SetDefault("LDAP_ENABLED", false)
	viper.SetDefault("LDAP_PORT", 389)
	viper.SetDefault("LDAP_UID_ATTR", "uid")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	return cfg
}
