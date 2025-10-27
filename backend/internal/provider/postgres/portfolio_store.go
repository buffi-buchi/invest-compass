package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

var (
	//go:embed queries/get_portfolios_by_user_id.sql
	getPortfoliosByUserIDQuery string
)

type PortfolioStore struct {
	db *pgxpool.Pool
}

func NewPortfolioStore(db *pgxpool.Pool) *PortfolioStore {
	return &PortfolioStore{
		db: db,
	}
}

func (s *PortfolioStore) GeByUserID(
	ctx context.Context,
	userID uuid.UUID,
	limit int64,
	offset int64,
) ([]model.Portfolio, error) {
	rows, err := s.db.Query(ctx, getPortfoliosByUserIDQuery, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select portfolios by user id: %w", err)
	}

	portfolios, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Portfolio, error) {
		var portfolio model.Portfolio

		return portfolio, row.Scan(&portfolio.ID, &portfolio.UserID, &portfolio.Name, &portfolio.CreateTime)
	})
	if err != nil {
		return nil, fmt.Errorf("select portfolios by user id: %w", err)
	}

	return portfolios, nil
}
