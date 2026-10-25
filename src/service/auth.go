package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/visionJAMx/whoknows/src/crypto"
	"github.com/visionJAMx/whoknows/src/models"
	"github.com/visionJAMx/whoknows/src/repositories"
)

const (
	UsernameMinLength = 3
	UsernameMaxLength = 50
	PasswordMinLength = 10
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (int64, error)
	FindByUsername(ctx context.Context, username string) (models.User, error)
}

type AuthService struct {
	users UserRepository
}

type RegisterInput struct {
	Username        string
	Email           string
	Password        string
	PasswordConfirm string // optional; empty means not sent
}

type LoginInput struct {
	Username string
	Password string
}

func NewAuthService(users UserRepository) *AuthService {
	return &AuthService{users: users}
}

// Register validates the input, hashes the password and stores the user.
// Invalid input, including a taken username or email, is returned as a *ValidationError.
func (s *AuthService) Register(ctx context.Context, input RegisterInput) error {
	input = normalizeRegistration(input)

	if fields := validateRegistration(input); len(fields) > 0 {
		return &ValidationError{Fields: fields}
	}

	hash, err := crypto.HashPassword(input.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	_, err = s.users.Create(ctx, models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hash,
	})

	switch {
	case errors.Is(err, repositories.ErrUsernameTaken):
		return &ValidationError{Fields: []FieldError{
			{Field: "username", Problem: ProblemTaken, Message: "Username is already taken"},
		}}
	case errors.Is(err, repositories.ErrEmailTaken):
		return &ValidationError{Fields: []FieldError{
			{Field: "email", Problem: ProblemTaken, Message: "Email is already taken"},
		}}
	case err != nil:
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (string, error) {

	// TODO: call repository layer

	return "POST /api/login", nil
}

func (s *AuthService) Logout(ctx context.Context) (string, error) {
	return "GET /api/logout", nil
}

func normalizeRegistration(input RegisterInput) RegisterInput {
	input.Username = strings.TrimSpace(input.Username)
	input.Email = strings.TrimSpace(input.Email)

	return input
}

// validateRegistration returns every problem with the input, not just the first.
func validateRegistration(input RegisterInput) []FieldError {
	var fields []FieldError

	usernameLength := utf8.RuneCountInString(input.Username)
	switch {
	case usernameLength == 0:
		fields = append(fields, FieldError{
			Field:   "username",
			Problem: ProblemMissing,
			Message: "Username is required",
		})
	case usernameLength < UsernameMinLength || usernameLength > UsernameMaxLength:
		fields = append(fields, FieldError{
			Field:   "username",
			Problem: ProblemLength,
			Message: fmt.Sprintf("Username must be between %d and %d characters", UsernameMinLength, UsernameMaxLength),
		})
	}

	switch {
	case input.Email == "":
		fields = append(fields, FieldError{
			Field:   "email",
			Problem: ProblemMissing,
			Message: "Email is required",
		})
	case !strings.Contains(input.Email, "@"):
		fields = append(fields, FieldError{
			Field:   "email",
			Problem: ProblemInvalid,
			Message: "Email must be a valid email address",
		})
	}

	switch {
	case input.Password == "":
		fields = append(fields, FieldError{
			Field:   "password",
			Problem: ProblemMissing,
			Message: "Password is required",
		})
	case utf8.RuneCountInString(input.Password) < PasswordMinLength:
		fields = append(fields, FieldError{
			Field:   "password",
			Problem: ProblemLength,
			Message: fmt.Sprintf("Password must be at least %d characters", PasswordMinLength),
		})
	}

	if input.PasswordConfirm != "" && input.PasswordConfirm != input.Password {
		fields = append(fields, FieldError{
			Field:   "password2",
			Problem: ProblemMismatch,
			Message: "Passwords do not match",
		})
	}

	return fields
}
