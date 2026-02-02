package model

import (
	"time"

	"github.com/google/uuid"
)

type SecurityType string

const (
	SecurityTypeShare SecurityType = "share"
	SecurityTypeBond  SecurityType = "bond"
	SecurityTypeFund  SecurityType = "fund"
)

type Security struct {
	ID         uuid.UUID
	SecID      string
	Ticker     string
	ShortName  string
	Type       SecurityType
	Extra      map[string]interface{}
	CreateTime time.Time
}
