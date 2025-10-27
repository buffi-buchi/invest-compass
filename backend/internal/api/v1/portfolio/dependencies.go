package portfolio

import (
	"context"

	"github.com/google/uuid"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

//go:generate go tool minimock -g -i Service

type Service interface {
	GeByUserID(
		ctx context.Context,
		userID uuid.UUID,
		limit int64,
		offset int64,
	) ([]model.Portfolio, error)
}
