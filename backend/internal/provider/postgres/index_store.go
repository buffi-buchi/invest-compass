package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed queries/create_index.sql
	createIndexQuery string
	//go:embed queries/get_index_by_code.sql
	getIndexByCodeQuery string
	//go:embed queries/list_indexes.sql
	getAllIndexesQuery string
)

type IndexStore struct {
	db  *pgxpool.Pool
	id  func() (uuid.UUID, error)
	now func() time.Time
}

func NewIndexStore(db *pgxpool.Pool) *IndexStore {
	return &IndexStore{
		db:  db,
		id:  func() (uuid.UUID, error) { return uuid.NewV7() },
		now: func() time.Time { return time.Now().UTC() },
	}
}
func (s *IndexStore) Create(ctx context.Context, index model.Index) (model.Index, error) {
	id, err := s.id()

	if err != nil {
		return model.Index{}, fmt.Errorf("create index ID: %w", err)
	}

	index.ID = id
	index.CreateTime = s.now()

	_, err = s.db.Exec(ctx, createIndexQuery, index.ID, index.Ticker, index.Name, index.CreateTime)

	if err != nil {
		return model.Index{}, fmt.Errorf("insert index: %w", err)
	}

	return index, nil
}
func (s *IndexStore) GetByTicker(ctx context.Context, code string) (model.Index, error) {
	row := s.db.QueryRow(ctx, getIndexByCodeQuery, code)
	var index model.Index
	err := row.Scan(&index.ID, &index.Ticker, &index.Name, &index.CreateTime)

	if errors.Is(err, sql.ErrNoRows) {
		return model.Index{}, model.ErrNotFound
	}

	if err != nil {
		return model.Index{}, fmt.Errorf("select index by ticker: %w", err)
	}

	return index, nil
}
func (s *IndexStore) List(ctx context.Context, limit int64,
	offset int64) ([]model.Index, error) {
	rows, err := s.db.Query(ctx, getAllIndexesQuery, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("select indexes: %w", err)
	}

	indexes, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Index, error) {
		var index model.Index
		return index, row.Scan(&index.ID, &index.Ticker, &index.Name, &index.CreateTime)
	})

	if err != nil {
		return nil, fmt.Errorf("select indexes: %w", err)
	}

	return indexes, nil
}
