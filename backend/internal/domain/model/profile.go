package model

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Name       string
	CreateTime time.Time
}
