package service

import "context"

type AuthService struct {
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) error {

	// TODO: call repository layer

	return nil
}
