package ticket

import "time"

// Comment is a public or internal note on a ticket.
type Comment struct {
	ID         string
	TicketID   string
	UserID     string
	Content    string
	IsInternal bool
	CreatedAt  time.Time

	// Derived
	UserName string
}

// IsVisibleTo returns true if the comment can be seen by the given role.
// End users cannot see internal comments.
func (c *Comment) IsVisibleTo(role string) bool {
	if role == "admin" || role == "agent" {
		return true
	}
	return !c.IsInternal
}
