package service

import "context"

type AuthService struct {
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

func NewAuthService() *AuthService {
	return &AuthService{}
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
