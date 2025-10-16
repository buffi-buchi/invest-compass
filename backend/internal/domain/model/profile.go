package model

import (
	"github.com/google/uuid"
)

type Profile struct {
	UserID uuid.UUID
	Ticker string
}
