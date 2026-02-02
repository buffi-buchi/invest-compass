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
	//go:embed queries/create_portfolio_security.sql
	createPortfolioSecurityQuery string

	//go:embed queries/get_portfolio_securities.sql
	getPortfolioSecuritiesQuery string
)

type PortfolioSecurityStore struct {
	db  *pgxpool.Pool
	id  func() (uuid.UUID, error)
	now func() time.Time
}

func NewPortfolioSecurityStore(db *pgxpool.Pool) *PortfolioSecurityStore {
	return &PortfolioSecurityStore{
		db:  db,
		id:  func() (uuid.UUID, error) { return uuid.NewV7() },
		now: func() time.Time { return time.Now().UTC() },
	}
}

func (s *PortfolioSecurityStore) Upsert(ctx context.Context, ps model.PortfolioSecurity) (model.PortfolioSecurity, error) {
	id, err := s.id()
	if err != nil {
		return model.PortfolioSecurity{}, fmt.Errorf("create portfolio_security ID: %w", err)
	}

	ps.ID = id
	ps.CreateTime = s.now()

	_, err = s.db.Exec(ctx, createPortfolioSecurityQuery, ps.ID, ps.PortfolioID, ps.SecurityID, ps.Amount, ps.CreateTime)
	if err != nil {
		return model.PortfolioSecurity{}, fmt.Errorf("upsert portfolio_security: %w", err)
	}

	return ps, nil
}

func (s *PortfolioSecurityStore) GetByPortfolioID(ctx context.Context, portfolioID uuid.UUID) ([]model.PortfolioSecurity, error) {
	rows, err := s.db.Query(ctx, getPortfolioSecuritiesQuery, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("select portfolio securities: %w", err)
	}

	portfolioSecurities, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.PortfolioSecurity, error) {
		var ps model.PortfolioSecurity

		return ps, row.Scan(&ps.ID, &ps.PortfolioID, &ps.SecurityID, &ps.Amount, &ps.CreateTime)
	})
	if err != nil {
		return nil, fmt.Errorf("collect portfolio securities: %w", err)
	}

	return portfolioSecurities, nil
}
