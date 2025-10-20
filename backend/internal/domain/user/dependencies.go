package user

import (
	"context"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

type Store interface {
	Create(ctx context.Context, user model.User) (model.User, error)
}
