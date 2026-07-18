package postgres

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewTicketRepo(t *testing.T) {
	r := NewTicketRepo(nil)
	if r == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewCommentRepo(t *testing.T) {
	r := NewCommentRepo(nil)
	if r == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewHistoryRepo(t *testing.T) {
	r := NewHistoryRepo(nil)
	if r == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewUserRepo(t *testing.T) {
	r := NewUserRepo(nil)
	if r == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewArticleRepo(t *testing.T) {
	r := NewArticleRepo(nil)
	if r == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewAttachmentRepo(t *testing.T) {
	r := NewAttachmentRepo(nil)
	if r == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewAnalyticsRepo(t *testing.T) {
	r := NewAnalyticsRepo(nil)
	if r == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewSLARepo(t *testing.T) {
	r := NewSLARepo(nil)
	if r == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewAutomationRepo(t *testing.T) {
	r := NewAutomationRepo(nil)
	if r == nil {
		t.Fatal("expected non-nil repo")
	}
}

func TestNewTenantResolver(t *testing.T) {
	r := NewTenantResolver(nil)
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}

func TestNewAutomationExecutor(t *testing.T) {
	e := NewAutomationExecutor(nil)
	if e == nil {
		t.Fatal("expected non-nil executor")
	}
	expected := uuid.MustParse("00000000-0000-0000-0000-000000000002")
	if e.systemUserID != expected {
		t.Errorf("systemUserID = %v, want %v", e.systemUserID, expected)
	}
}

func TestAutomationExecutor_WithSystemUserID(t *testing.T) {
	e := NewAutomationExecutor(nil)
	customID := uuid.New()
	result := e.WithSystemUserID(customID)
	if result != e {
		t.Error("expected to return same executor instance")
	}
	if e.systemUserID != customID {
		t.Errorf("systemUserID = %v, want %v", e.systemUserID, customID)
	}
}

func TestAutomationRepo_DB(t *testing.T) {
	r := NewAutomationRepo(nil)
	if r.DB() != nil {
		t.Error("expected nil DB")
	}
}
