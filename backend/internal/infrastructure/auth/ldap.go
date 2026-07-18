package auth

import (
	"errors"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	"github.com/trickreport/backend/internal/config"
)

// LDAPAuthenticator handles authentication against an LDAP directory.
// It implements the application/auth.LDAPAuthenticator interface.
type LDAPAuthenticator struct {
	host     string
	port     int
	bindDN   string
	bindPW   string
	baseDN   string
	uidAttr  string
}

// NewLDAPAuthenticator creates a new LDAPAuthenticator from config.
func NewLDAPAuthenticator(cfg *config.Config) *LDAPAuthenticator {
	return &LDAPAuthenticator{
		host:    cfg.LDAPHost,
		port:    cfg.LDAPPort,
		bindDN:  cfg.LDAPBindDN,
		bindPW:  cfg.LDAPBindPW,
		baseDN:  cfg.LDAPBaseDN,
		uidAttr: cfg.LDAPUIDAttr,
	}
}

// Authenticate verifies credentials against the LDAP server.
// Returns the user's DN on success.
func (a *LDAPAuthenticator) Authenticate(username, password string) (string, error) {
	addr := fmt.Sprintf("%s:%d", a.host, a.port)

	conn, err := ldap.DialURL(fmt.Sprintf("ldap://%s", addr))
	if err != nil {
		return "", fmt.Errorf("ldap dial failed: %w", err)
	}
	defer conn.Close()

	// Bind as service account to search for the user
	if err := conn.Bind(a.bindDN, a.bindPW); err != nil {
		return "", fmt.Errorf("ldap bind failed: %w", err)
	}

	searchReq := ldap.NewSearchRequest(
		a.baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 0, false,
		fmt.Sprintf("(%s=%s)", a.uidAttr, ldap.EscapeFilter(username)),
		[]string{"dn", "mail", "cn"},
		nil,
	)

	result, err := conn.Search(searchReq)
	if err != nil {
		return "", fmt.Errorf("ldap search failed: %w", err)
	}

	if len(result.Entries) == 0 {
		return "", errors.New("user not found in LDAP")
	}

	userDN := result.Entries[0].DN

	// Re-bind as the actual user to verify their password
	if err := conn.Bind(userDN, password); err != nil {
		return "", errors.New("invalid LDAP credentials")
	}

	return userDN, nil
}
