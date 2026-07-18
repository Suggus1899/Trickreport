package user

import "testing"

func TestRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role Role
		want bool
	}{
		{"admin", RoleAdmin, true},
		{"agent", RoleAgent, true},
		{"end_user", RoleEndUser, true},
		{"invalid", Role("superuser"), false},
		{"empty", Role(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.IsValid(); got != tt.want {
				t.Errorf("Role(%q).IsValid() = %v, want %v", tt.role, got, tt.want)
			}
		})
	}
}

func TestParseRole(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Role
		wantErr bool
	}{
		{"admin", "admin", RoleAdmin, false},
		{"invalid", "foo", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRole(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseRole(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParseRole(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestRole_Permissions(t *testing.T) {
	tests := []struct {
		name       string
		role       Role
		manage     bool
		writeArt   bool
		assign     bool
		viewAll    bool
	}{
		{"admin", RoleAdmin, true, true, true, true},
		{"agent", RoleAgent, false, true, true, true},
		{"end_user", RoleEndUser, false, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.CanManageUsers(); got != tt.manage {
				t.Errorf("CanManageUsers() = %v, want %v", got, tt.manage)
			}
			if got := tt.role.CanWriteArticles(); got != tt.writeArt {
				t.Errorf("CanWriteArticles() = %v, want %v", got, tt.writeArt)
			}
			if got := tt.role.CanAssignTickets(); got != tt.assign {
				t.Errorf("CanAssignTickets() = %v, want %v", got, tt.assign)
			}
			if got := tt.role.CanViewAllTickets(); got != tt.viewAll {
				t.Errorf("CanViewAllTickets() = %v, want %v", got, tt.viewAll)
			}
		})
	}
}
