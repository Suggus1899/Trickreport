package auth

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	"github.com/trickreport/backend/internal/config"
)

// LDAPAuthenticator handles authentication against an LDAP directory.
// It implements the application/auth.LDAPAuthenticator interface.
type LDAPAuthenticator struct {
	host    string
	port    int
	useTLS  bool
	bindDN  string
	bindPW  string
	baseDN  string
	uidAttr string
}

// NewLDAPAuthenticator creates a new LDAPAuthenticator from config.
func NewLDAPAuthenticator(cfg *config.Config) *LDAPAuthenticator {
	return &LDAPAuthenticator{
		host:    cfg.LDAPHost,
		port:    cfg.LDAPPort,
		useTLS:  cfg.LDAPUseTLS,
		bindDN:  cfg.LDAPBindDN,
		bindPW:  cfg.LDAPBindPW,
		baseDN:  cfg.LDAPBaseDN,
		uidAttr: cfg.LDAPUIDAttr,
	}
}

// Authenticate verifies credentials against the LDAP server.
// Returns the user's DN on success.
// When useTLS is true, it connects via LDAPS (TLS) on the configured port.
// Otherwise it falls back to plaintext LDAP (not recommended for production).
func (a *LDAPAuthenticator) Authenticate(username, password string) (string, error) {
	scheme := "ldap"
	if a.useTLS {
		scheme = "ldaps"
	}
	addr := fmt.Sprintf("%s://%s:%s", scheme, a.host, strconv.Itoa(a.port))

	conn, err := ldap.DialURL(addr)
	if err != nil {
		return "", fmt.Errorf("ldap dial failed: %w", err)
	}
	defer conn.Close()

	// For plaintext ldap://, attempt StartTLS upgrade if TLS is requested
	// but the URL scheme was not ldaps:// (e.g. port 389 with StartTLS).
	if a.useTLS && scheme == "ldap" {
		if err := conn.StartTLS(&tls.Config{ServerName: a.host}); err != nil {
			return "", fmt.Errorf("ldap StartTLS failed: %w", err)
		}
	}

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
