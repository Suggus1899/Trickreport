package user

import (
	"time"

	"github.com/google/uuid"
)

// User is a person who interacts with the help desk system.
type User struct {
	ID                   uuid.UUID
	TenantID             uuid.UUID
	Name                 string
	Email                string
	Role                 Role
	PasswordHash         string // empty when using LDAP/SSO only
	LDAPDN               string
	AvatarURL            *string
	Active               bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
	FailedLoginAttempts  int       // number of consecutive failed login attempts
	LockedUntil          time.Time // zero value means not locked
	MFASecret            string    // TOTP secret; empty when MFA is not configured
	MFAEnabled           bool      // true when MFA is active
}
