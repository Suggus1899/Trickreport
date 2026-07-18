package article

import (
	"time"

	"github.com/google/uuid"
)

// Article is a knowledge base entry.
type Article struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Title     string
	Content   string
	Category  string
	Tags      []string
	Published bool
	CreatedBy uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	// Derived
	AuthorName string
}

// IsVisibleTo returns true if the article can be seen by the given role.
// End users only see published articles.
func (a *Article) IsVisibleTo(role string) bool {
	if role == "admin" || role == "agent" {
		return true
	}
	return a.Published
}
