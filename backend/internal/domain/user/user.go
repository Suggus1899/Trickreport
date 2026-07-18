package user

import "time"

// User is a person who interacts with the help desk system.
type User struct {
	ID           string
	TenantID     string
	Name         string
	Email        string
	Role         Role
	PasswordHash string // empty when using LDAP/SSO only
	LDAPDN       string
	AvatarURL    *string
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
