package model

import (
	"time"
)

type Index struct {
	Ticker     string
	ShortName  string
	CreateTime time.Time
}
