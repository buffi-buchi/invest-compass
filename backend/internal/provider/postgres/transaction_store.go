package postgres

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

var (
	//go:embed queries/create_transaction.sql
	createTransactionQuery string

	//go:embed queries/get_transactions_by_portfolio.sql
	getTransactionsByPortfolioQuery string
)

type TransactionStore struct {
	db  *pgxpool.Pool
	id  func() (uuid.UUID, error)
	now func() time.Time
}

func NewTransactionStore(db *pgxpool.Pool) *TransactionStore {
	return &TransactionStore{
		db:  db,
		id:  func() (uuid.UUID, error) { return uuid.NewV7() },
		now: func() time.Time { return time.Now().UTC() },
	}
}

func (s *TransactionStore) Create(ctx context.Context, tx model.Transaction) (model.Transaction, error) {
	id, err := s.id()
	if err != nil {
		return model.Transaction{}, fmt.Errorf("create transaction ID: %w", err)
	}

	tx.ID = id

	_, err = s.db.Exec(ctx, createTransactionQuery, tx.ID, tx.PortfolioID, tx.SecurityID,
		tx.Amount, tx.Price, tx.TradeDate, tx.Type, tx.Note)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("insert transaction: %w", err)
	}

	return tx, nil
}

func (s *TransactionStore) GetByPortfolioID(ctx context.Context, portfolioID uuid.UUID, limit, offset int64) ([]model.Transaction, error) {
	rows, err := s.db.Query(ctx, getTransactionsByPortfolioQuery, portfolioID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select transactions: %w", err)
	}

	transactions, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Transaction, error) {
		var tx model.Transaction

		return tx, row.Scan(&tx.ID, &tx.PortfolioID, &tx.SecurityID, &tx.Amount,
			&tx.Price, &tx.TradeDate, &tx.Type, &tx.Note)
	})
	if err != nil {
		return nil, fmt.Errorf("collect transactions: %w", err)
	}

	return transactions, nil
}
