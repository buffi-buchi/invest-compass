package user

import (
	"context"
	"fmt"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) Create(ctx context.Context, user model.User) (model.User, error) {
	err := user.HashPassword()
	if err != nil {
		return model.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err = s.store.Create(ctx, user)
	if err != nil {
		return model.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
