package profile

import (
	"context"

	"github.com/google/uuid"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

type Service interface {
	GetProfilesByUserID(ctx context.Context, userID uuid.UUID) ([]model.Profile, error)
}
