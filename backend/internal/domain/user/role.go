package user

import "fmt"

// Role represents a user's permission level.
type Role string

const (
	RoleAdmin   Role = "admin"
	RoleAgent   Role = "agent"
	RoleEndUser Role = "end_user"
)

// IsValid returns true if the role is a recognized value.
func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleAgent, RoleEndUser:
		return true
	}
	return false
}

// ParseRole converts a string to a Role, returning an error if invalid.
func ParseRole(s string) (Role, error) {
	rl := Role(s)
	if !rl.IsValid() {
		return "", fmt.Errorf("invalid user role: %q", s)
	}
	return rl, nil
}

// CanManageUsers returns true if the role can manage users (admin only).
func (r Role) CanManageUsers() bool {
	return r == RoleAdmin
}

// CanWriteArticles returns true if the role can create/edit articles.
func (r Role) CanWriteArticles() bool {
	return r == RoleAdmin || r == RoleAgent
}

// CanAssignTickets returns true if the role can assign tickets.
func (r Role) CanAssignTickets() bool {
	return r == RoleAdmin || r == RoleAgent
}

// CanViewAllTickets returns true if the role can see all tenant tickets.
func (r Role) CanViewAllTickets() bool {
	return r == RoleAdmin || r == RoleAgent
}
