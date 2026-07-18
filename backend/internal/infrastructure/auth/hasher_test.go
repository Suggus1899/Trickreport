package auth

import (
	"testing"
)

func TestBcryptHasher_HashAndCompare_Success(t *testing.T) {
	h := NewBcryptHasher()
	hash, err := h.Hash("mysecret")
	if err != nil {
		t.Fatalf("Hash error: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash")
	}
	if hash == "mysecret" {
		t.Error("hash should not equal plaintext")
	}

	if err := h.Compare(hash, "mysecret"); err != nil {
		t.Errorf("Compare with correct password: %v", err)
	}
}

func TestBcryptHasher_Compare_WrongPassword(t *testing.T) {
	h := NewBcryptHasher()
	hash, _ := h.Hash("correct-password")

	if err := h.Compare(hash, "wrong-password"); err == nil {
		t.Error("expected error for wrong password")
	}
}

func TestBcryptHasher_Compare_InvalidHash(t *testing.T) {
	h := NewBcryptHasher()
	if err := h.Compare("not-a-hash", "password"); err == nil {
		t.Error("expected error for invalid hash")
	}
}

func TestNewBcryptHasher(t *testing.T) {
	h := NewBcryptHasher()
	if h == nil {
		t.Fatal("expected non-nil hasher")
	}
}
