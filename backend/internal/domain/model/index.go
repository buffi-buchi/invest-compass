package model

import (
	"time"
)

type Index struct {
	Ticker     string
	Name       string
	CreateTime time.Time
}
