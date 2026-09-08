package database

import (
	"context"
	"path/filepath"
	"testing"
)

func TestInitializeCreatesTables(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "test.db")

	db, err := Open(databasePath)
	if err != nil {
		t.Fatalf("Open() returned an error: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := Initialize(ctx, db); err != nil {
		t.Fatalf("Initialize() returned an error: %v", err)
	}

	expectedTables := []string{"users", "pages"}

	for _, table := range expectedTables {
		var name string

		err := db.QueryRowContext(
			ctx,
			`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`,
			table,
		).Scan(&name)

		if err != nil {
			t.Fatalf("expected table %q to exist: %v", table, err)
		}

		if name != table {
			t.Errorf("expected table %q, got %q", table, name)
		}
	}
}

func TestCreateAndFindUser(t *testing.T) {
	databasePath := filepath.Join(t.TempDir(), "test.db")

	db, err := Open(databasePath)
	if err != nil {
		t.Fatalf("Open() returned an error: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := Initialize(ctx, db); err != nil {
		t.Fatalf("Initialize() returned an error: %v", err)
	}

	username := `osman'); DROP TABLE users; --`
	email := "osman@example.com"
	passwordHash := "hashed-password"

	userID, err := CreateUser(ctx, db, username, email, passwordHash)
	if err != nil {
		t.Fatalf("CreateUser() returned an error: %v", err)
	}

	user, err := FindUserByUsername(ctx, db, username)
	if err != nil {
		t.Fatalf("FindUserByUsername() returned an error: %v", err)
	}

	if user.ID != userID {
		t.Errorf("expected user ID %d, got %d", userID, user.ID)
	}

	if user.Username != username {
		t.Errorf("expected username %q, got %q", username, user.Username)
	}

	if user.Email != email {
		t.Errorf("expected email %q, got %q", email, user.Email)
	}

	if user.PasswordHash != passwordHash {
		t.Errorf(
			"expected password hash %q, got %q",
			passwordHash,
			user.PasswordHash,
		)
	}
}
