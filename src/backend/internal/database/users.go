package database

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
}

func CreateUser(
	ctx context.Context,
	db *sql.DB,
	username string,
	email string,
	passwordHash string,
) (int64, error) {
	result, err := db.ExecContext(
		ctx,
		`INSERT INTO users (username, email, password_hash)
		 VALUES (?, ?, ?)`,
		username,
		email,
		passwordHash,
	)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get created user ID: %w", err)
	}

	return id, nil
}

func FindUserByUsername(
	ctx context.Context,
	db *sql.DB,
	username string,
) (User, error) {
	var user User

	err := db.QueryRowContext(
		ctx,
		`SELECT id, username, email, password_hash
		 FROM users
		 WHERE username = ?`,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		return User{}, fmt.Errorf("find user by username: %w", err)
	}

	return user, nil
}
