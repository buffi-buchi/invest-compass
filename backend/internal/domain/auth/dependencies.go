package auth

import (
	"context"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

type UserStore interface {
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type JWTProvider interface {
	Generate(user model.User) (string, error)
}
