package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Errorf("expected no error hashing password, got: %v", err)
	}
	if hash == "" {
		t.Error("expected non-empty hash, got empty string")
	}
	if hash == "secret123" {
		t.Error("hash should not equal the original password")
	}
}

func TestHashPassword_DifferentEachTime(t *testing.T) {
	// bcrypt generates a different hash every time due to salting
	// both should still be valid hashes of the same password
	hash1, _ := HashPassword("secret123")
	hash2, _ := HashPassword("secret123")

	if hash1 == hash2 {
		t.Error("expected different hashes for same password due to bcrypt salting")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash, _ := HashPassword("secret123")

	if !CheckPasswordHash("secret123", hash) {
		t.Error("expected correct password to match hash")
	}

	if CheckPasswordHash("wrongpassword", hash) {
		t.Error("expected wrong password to not match hash")
	}
}
