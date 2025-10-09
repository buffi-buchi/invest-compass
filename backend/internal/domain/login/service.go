package login

import (
	"context"
	"errors"
	"fmt"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

type Service struct {
	userStore   UserStore
	jwtProvider JWTProvider
}

func NewService(userStore UserStore, jwtProvider JWTProvider) *Service {
	return &Service{
		userStore:   userStore,
		jwtProvider: jwtProvider,
	}
}

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userStore.GetByEmail(ctx, email)
	if errors.Is(err, model.ErrNotFound) {
		return "", fmt.Errorf("user not found: %w", model.ErrNotAuthorized)
	}
	if err != nil {
		return "", fmt.Errorf("get user by email: %w", err)
	}

	if !user.CheckPassword(password) {
		return "", fmt.Errorf("invalid password: %w", model.ErrNotAuthorized)
	}

	token, err := s.jwtProvider.Generate(user)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return token, nil
}
