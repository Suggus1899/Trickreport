package article

import "testing"

func TestArticle_IsVisibleTo(t *testing.T) {
	published := &Article{Published: true}
	draft := &Article{Published: false}

	tests := []struct {
		name string
		a    *Article
		role string
		want bool
	}{
		{"admin sees published", published, "admin", true},
		{"admin sees draft", draft, "admin", true},
		{"agent sees published", published, "agent", true},
		{"agent sees draft", draft, "agent", true},
		{"end_user sees published", published, "end_user", true},
		{"end_user cannot see draft", draft, "end_user", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.IsVisibleTo(tt.role); got != tt.want {
				t.Errorf("IsVisibleTo(%q) = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}
