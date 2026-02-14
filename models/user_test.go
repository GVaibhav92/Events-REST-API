package models

import (
	"testing"
)

func TestSaveUser(t *testing.T) {
	setupTestDB(t)

	user := User{
		Email:    "test@example.com",
		Password: "secret123",
	}

	err := user.Save()
	if err != nil {
		t.Errorf("expected no error saving user, got: %v", err)
	}
	if user.ID == 0 {
		t.Error("expected user ID to be set after save, got 0")
	}
}

func TestSaveUser_DuplicateEmail(t *testing.T) {
	setupTestDB(t)

	user := User{Email: "duplicate@example.com", Password: "secret123"}
	user.Save()

	duplicate := User{Email: "duplicate@example.com", Password: "different123"}
	err := duplicate.Save()
	if err == nil {
		t.Error("expected error for duplicate email, got nil")
	}
	if err.Error() != "email already registered" {
		t.Errorf("expected 'email already registered' error, got: %v", err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	setupTestDB(t)

	user := User{Email: "findme@example.com", Password: "secret123"}
	user.Save()

	found, err := GetUserByEmail("findme@example.com")
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
	if found == nil {
		t.Fatal("expected to find user, got nil")
	}
	if found.Email != "findme@example.com" {
		t.Errorf("expected email findme@example.com, got %s", found.Email)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	setupTestDB(t)

	found, err := GetUserByEmail("ghost@example.com")
	if err != nil {
		t.Errorf("expected no error for missing user, got: %v", err)
	}
	if found != nil {
		t.Error("expected nil for non-existent user, got a result")
	}
}

func TestValidateCredentials(t *testing.T) {
	setupTestDB(t)

	user := User{Email: "auth@example.com", Password: "secret123"}
	user.Save()

	// Test correct credentials
	loginUser := User{Email: "auth@example.com", Password: "secret123"}
	err := loginUser.ValidateCredentials()
	if err != nil {
		t.Errorf("expected valid credentials to pass, got: %v", err)
	}

	// Test wrong password
	wrongPass := User{Email: "auth@example.com", Password: "wrongpassword"}
	err = wrongPass.ValidateCredentials()
	if err == nil {
		t.Error("expected error for wrong password, got nil")
	}

	// Test non-existent user
	noUser := User{Email: "nobody@example.com", Password: "secret123"}
	err = noUser.ValidateCredentials()
	if err == nil {
		t.Error("expected error for non-existent user, got nil")
	}
}
