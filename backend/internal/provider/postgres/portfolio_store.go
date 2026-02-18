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
	//go:embed queries/get_portfolios_by_user_id.sql
	getPortfoliosByUserIDQuery string
	//go:embed queries/create_portfolio.sql
	createPortfolioQuery string
)

type PortfolioStore struct {
	db  *pgxpool.Pool
	id  func() (uuid.UUID, error)
	now func() time.Time
}

func NewPortfolioStore(db *pgxpool.Pool) *PortfolioStore {
	return &PortfolioStore{
		db:  db,
		id:  func() (uuid.UUID, error) { return uuid.NewV7() },
		now: func() time.Time { return time.Now().UTC() },
	}
}

func (s *PortfolioStore) Create(ctx context.Context, ps model.Portfolio) (model.Portfolio, error) {
	id, err := s.id()
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("create portfolio_security ID: %w", err)
	}

	ps.ID = id
	ps.CreateTime = s.now()

	_, err = s.db.Exec(ctx, createPortfolioQuery, ps.ID, ps.UserID, ps.Name, ps.CreateTime)
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("upsert portfolio_security: %w", err)
	}

	return ps, nil
}

func (s *PortfolioStore) GetByUserID(
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
