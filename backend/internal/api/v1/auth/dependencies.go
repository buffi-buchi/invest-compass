package auth

import (
	"context"
)

//go:generate go tool minimock -g -i Service

type Service interface {
	Login(ctx context.Context, email string, password string) (string, error)
}
