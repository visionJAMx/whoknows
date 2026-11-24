package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/visionJAMx/whoknows/src/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user models.User) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO users (username, email, password_hash)
		 VALUES ($1, $2, $3) RETURNING id`,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(&id)
	if err != nil {
		if dupErr := duplicateUserError(err); dupErr != nil {
			return 0, fmt.Errorf("create user: %w: %w", dupErr, err)
		}

		return 0, fmt.Errorf("create user: %w", err)
	}

	return id, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, username, email, password_hash
		 FROM users
		 WHERE username = $1`,
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

// duplicateUserError returns ErrUsernameTaken or ErrEmailTaken if err is a UNIQUE violation, otherwise nil.
func duplicateUserError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if strings.Contains(pqErr.Message, "users_username_key") || strings.Contains(pqErr.Message, "username") {
			return ErrUsernameTaken
		}
		if strings.Contains(pqErr.Message, "users_email_key") || strings.Contains(pqErr.Message, "email") {
			return ErrEmailTaken
		}
	}

	return nil
}
