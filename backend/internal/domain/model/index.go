package model

import (
	"time"

	"github.com/google/uuid"
)

type Index struct {
	ID         uuid.UUID
	Ticker     string
	Name       string
	CreateTime time.Time
}
