package auth

import (
	"testing"

	"github.com/trickreport/backend/internal/config"
)

func TestNewLDAPAuthenticator(t *testing.T) {
	cfg := &config.Config{
		LDAPHost:    "ldap.example.com",
		LDAPPort:    636,
		LDAPUseTLS:  true,
		LDAPBindDN:  "cn=admin,dc=example,dc=com",
		LDAPBindPW:  "secret",
		LDAPBaseDN:  "dc=example,dc=com",
		LDAPUIDAttr: "uid",
	}

	a := NewLDAPAuthenticator(cfg)
	if a == nil {
		t.Fatal("expected non-nil authenticator")
	}
}

func TestLDAPAuthenticator_Authenticate_DialError(t *testing.T) {
	cfg := &config.Config{
		LDAPHost:    "nonexistent.invalid",
		LDAPPort:    389,
		LDAPUseTLS:  false,
		LDAPBindDN:  "cn=admin,dc=example,dc=com",
		LDAPBindPW:  "secret",
		LDAPBaseDN:  "dc=example,dc=com",
		LDAPUIDAttr: "uid",
	}

	a := NewLDAPAuthenticator(cfg)
	_, err := a.Authenticate("testuser", "password")
	if err == nil {
		t.Error("expected error for non-existent LDAP server")
	}
}
