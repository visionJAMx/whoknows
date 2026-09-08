package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/visionJAMx/whoknows/src/domain"
)

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
) (domain.User, error) {
	var user domain.User

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
		return domain.User{}, fmt.Errorf("find user by username: %w", err)
	}

	return user, nil
}
