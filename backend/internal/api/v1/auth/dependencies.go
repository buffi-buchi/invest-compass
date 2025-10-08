package auth

import (
	"context"
)

type Service interface {
	Login(ctx context.Context, email string, password string) (string, error)
}
