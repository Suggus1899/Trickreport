package article

import "time"

// Article is a knowledge base entry.
type Article struct {
	ID        string
	TenantID  string
	Title     string
	Content   string
	Category  string
	Tags      []string
	Published bool
	CreatedBy string
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
