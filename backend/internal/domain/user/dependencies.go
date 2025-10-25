package user

import (
	"context"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

//go:generate go tool minimock -g -i Store

type Store interface {
	Create(ctx context.Context, user model.User) (model.User, error)
}
