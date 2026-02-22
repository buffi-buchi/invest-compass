package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed queries/create_security.sql
	createSecurityQuery string
	//go:embed queries/get_security_by_ticker.sql
	getSecurityByTickerQuery string
	//go:embed queries/list_security.sql
	listSecurityQuery string
)

type SecurityStore struct {
	db  *pgxpool.Pool
	now func() time.Time
}

func NewSecurityStore(db *pgxpool.Pool) *SecurityStore {
	return &SecurityStore{
		db:  db,
		now: func() time.Time { return time.Now().UTC() },
	}
}
func (s *SecurityStore) Create(ctx context.Context, security model.Security) (model.Security, error) {

	security.CreateTime = s.now()

	_, err := s.db.Exec(ctx, createSecurityQuery, security.Ticker, security.ShortName, security.CreateTime)
	if err != nil {
		return model.Security{}, fmt.Errorf("insert security: %w", err)
	}

	return security, nil
}
func (s *SecurityStore) GetByTicker(ctx context.Context, code string) (model.Security, error) {
	row := s.db.QueryRow(ctx, getSecurityByTickerQuery, code)
	var security model.Security
	err := row.Scan(&security.Ticker, &security.ShortName, &security.CreateTime)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Security{}, model.ErrNotFound
	}
	if err != nil {
		return model.Security{}, fmt.Errorf("select security by ticker: %w", err)
	}

	return security, nil
}
func (s *SecurityStore) List(
	ctx context.Context,
	limit int64,
	offset int64,
) ([]model.Security, error) {
	rows, err := s.db.Query(ctx, listSecurityQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select securities: %w", err)
	}

	securities, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Security, error) {
		var security model.Security
		return security, row.Scan(&security.Ticker, &security.ShortName, &security.CreateTime)
	})
	if err != nil {
		return nil, fmt.Errorf("select indexes: %w", err)
	}

	return securities, nil
}
