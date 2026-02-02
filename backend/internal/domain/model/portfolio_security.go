package model

import (
	"time"

	"github.com/google/uuid"
)

type PortfolioSecurity struct {
	ID          uuid.UUID
	PortfolioID uuid.UUID
	SecurityID  uuid.UUID
	Amount      int
	CreateTime  time.Time
}
