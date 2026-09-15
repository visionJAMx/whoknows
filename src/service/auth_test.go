package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/visionJAMx/whoknows/src/internal/repository"
)

// TestAuthenticate viser, at en bruger med korrekt password kan autentificeres.
func TestAuthenticate(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "test.db")

	db, err := repository.Open(databasePath)
	if err != nil {
		t.Fatalf("could not open test database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := repository.Initialize(ctx, db); err != nil {
		t.Fatalf("could not initialize test database: %v", err)
	}

	passwordHash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("could not hash password: %v", err)
	}

	_, err = repository.CreateUser(
		ctx,
		db,
		"osman",
		"osman@example.com",
		passwordHash,
	)
	if err != nil {
		t.Fatalf("could not create test user: %v", err)
	}

	user, err := Authenticate(ctx, db, "osman", "correct-password")
	if err != nil {
		t.Fatalf("expected successful login, got: %v", err)
	}

	if user.Username != "osman" {
		t.Errorf("expected username osman, got %s", user.Username)
	}
}

// TestAuthenticateRejectsInvalidCredentials dækker alle almindelige ugyldige loginforsøg.
func TestAuthenticateRejectsInvalidCredentials(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "test.db")

	db, err := repository.Open(databasePath)
	if err != nil {
		t.Fatalf("could not open test database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := repository.Initialize(ctx, db); err != nil {
		t.Fatalf("could not initialize test database: %v", err)
	}

	passwordHash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("could not hash password: %v", err)
	}

	_, err = repository.CreateUser(
		ctx,
		db,
		"osman",
		"osman@example.com",
		passwordHash,
	)
	if err != nil {
		t.Fatalf("could not create test user: %v", err)
	}

	testCases := []struct {
		name     string
		username string
		password string
	}{
		{"wrong password", "osman", "wrong-password"},
		{"unknown username", "unknown", "correct-password"},
		{"empty username", "", "correct-password"},
		{"empty password", "osman", ""},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := Authenticate(
				ctx,
				db,
				testCase.username,
				testCase.password,
			)

			if !errors.Is(err, ErrInvalidCredentials) {
				t.Errorf(
					"expected ErrInvalidCredentials, got %v",
					err,
				)
			}
		})
	}
}
