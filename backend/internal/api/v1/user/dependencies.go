package user

import (
	"context"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

type Service interface {
	CreateUser(ctx context.Context, user model.User) (model.User, error)
}
