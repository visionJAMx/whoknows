package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/visionJAMx/whoknows/src/models"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user models.User) (int64, error) {
	result, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (username, email, password_hash)
		 VALUES (?, ?, ?)`,
		user.Username,
		user.Email,
		user.PasswordHash,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return 0, fmt.Errorf("create user: %w: %w", ErrDuplicate, err)
		}

		return 0, fmt.Errorf("create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get created user ID: %w", err)
	}

	return id, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	err := r.db.QueryRowContext(
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
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("find user by username: %w", ErrNotFound)
		}

		return models.User{}, fmt.Errorf("find user by username: %w", err)
	}

	return user, nil
}

func isUniqueViolation(err error) bool {
	var sqliteErr *sqlite.Error
	return errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE
}
