package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	domainUser "github.com/trickreport/backend/internal/domain/user"
	"testing"
)

func TestCreateRejectsWeakLocalPassword(t *testing.T) {
	repo, hasher := &mockRepo{}, &mockHasher{}
	svc := NewService(repo, hasher)
	if _, err := svc.Create(context.Background(), CreateInput{TenantID: uuid.New(), Name: "Test", Email: "test@example.test", Password: "weak"}); !errors.Is(err, domainUser.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if repo.createCalls != 0 || hasher.gotPwd != "" {
		t.Fatal("invalid password reached hashing or persistence")
	}
}

func TestUpdateRejectsUnknownRole(t *testing.T) {
	repo := &mockRepo{}
	svc := NewService(repo, &mockHasher{})
	role := "superadmin"
	if _, err := svc.Update(context.Background(), uuid.New(), uuid.New(), UpdateFields{Role: &role}); !errors.Is(err, domainUser.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if repo.updateCalls != 0 {
		t.Fatal("invalid role reached persistence")
	}
}
