package model

import (
	"time"

	"github.com/google/uuid"
)

type Index struct {
	ID         uuid.UUID
	IndexCode  string
	Name       string
	CreateTime time.Time
}
