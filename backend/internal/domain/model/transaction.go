package model

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeBuy  TransactionType = "buy"
	TransactionTypeSell TransactionType = "sell"
)

type Transaction struct {
	ID          uuid.UUID
	PortfolioID uuid.UUID
	SecurityID  uuid.UUID
	Amount      int
	Price       float64
	TradeDate   time.Time
	Type        TransactionType
	Note        string
}
