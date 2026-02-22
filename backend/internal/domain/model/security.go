package model

import (
	"time"
)

type Security struct {
	Ticker     string
	ShortName  string
	CreateTime time.Time
}
