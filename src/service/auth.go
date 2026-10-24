package service

import (
	"context"

	"github.com/visionJAMx/whoknows/src/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (int64, error)
	FindByUsername(ctx context.Context, username string) (models.User, error)
}

type AuthService struct {
	users UserRepository
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
}

type LoginInput struct {
	Username string
	Password string
}

func NewAuthService(users UserRepository) *AuthService {
	return &AuthService{users: users}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (string, error) {

	// TODO: call repository layer

	return "POST /api/register", nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (string, error) {

	// TODO: call repository layer

	return "POST /api/login", nil
}

func (s *AuthService) Logout(ctx context.Context) (string, error) {
	return "GET /api/logout", nil
}
