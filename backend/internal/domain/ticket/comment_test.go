package ticket

import "testing"

func TestComment_IsVisibleTo(t *testing.T) {
	internal := &Comment{IsInternal: true}
	public := &Comment{IsInternal: false}

	tests := []struct {
		name string
		c    *Comment
		role string
		want bool
	}{
		{"admin sees internal", internal, "admin", true},
		{"admin sees public", public, "admin", true},
		{"agent sees internal", internal, "agent", true},
		{"agent sees public", public, "agent", true},
		{"end_user cannot see internal", internal, "end_user", false},
		{"end_user sees public", public, "end_user", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.IsVisibleTo(tt.role); got != tt.want {
				t.Errorf("IsVisibleTo(%q) = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}
