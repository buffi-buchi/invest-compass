package auth

import (
	"context"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

//go:generate go tool minimock -g -i UserStore

type UserStore interface {
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

//go:generate go tool minimock -g -i JWTProvider

type JWTProvider interface {
	Generate(user model.User) (string, error)
}
