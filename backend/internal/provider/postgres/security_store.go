package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

var (
	//go:embed queries/create_security.sql
	createSecurityQuery string

	//go:embed queries/get_security_by_id.sql
	getSecurityByIDQuery string

	//go:embed queries/get_security_by_sec_id.sql
	getSecurityBySecIDQuery string

	//go:embed queries/list_securities.sql
	listSecuritiesQuery string
)

type SecurityStore struct {
	db  *pgxpool.Pool
	id  func() (uuid.UUID, error)
	now func() time.Time
}

func NewSecurityStore(db *pgxpool.Pool) *SecurityStore {
	return &SecurityStore{
		db:  db,
		id:  func() (uuid.UUID, error) { return uuid.NewV7() },
		now: func() time.Time { return time.Now().UTC() },
	}
}

func (s *SecurityStore) Create(ctx context.Context, security model.Security) (model.Security, error) {
	id, err := s.id()
	if err != nil {
		return model.Security{}, fmt.Errorf("create security ID: %w", err)
	}

	security.ID = id
	security.CreateTime = s.now()

	extraJSON, err := json.Marshal(security.Extra)
	if err != nil {
		return model.Security{}, fmt.Errorf("marshal extra: %w", err)
	}

	_, err = s.db.Exec(ctx, createSecurityQuery, security.ID, security.SecID, security.Ticker,
		security.ShortName, security.Type, extraJSON, security.CreateTime)
	if err != nil {
		return model.Security{}, fmt.Errorf("insert security: %w", err)
	}

	return security, nil
}

func (s *SecurityStore) GetByID(ctx context.Context, id uuid.UUID) (model.Security, error) {
	row := s.db.QueryRow(ctx, getSecurityByIDQuery, id)

	var security model.Security
	var extraJSON []byte

	err := row.Scan(&security.ID, &security.SecID, &security.Ticker, &security.ShortName,
		&security.Type, &extraJSON, &security.CreateTime)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Security{}, model.ErrNotFound
	}
	if err != nil {
		return model.Security{}, fmt.Errorf("select security by ID: %w", err)
	}

	if len(extraJSON) > 0 {
		if err := json.Unmarshal(extraJSON, &security.Extra); err != nil {
			return model.Security{}, fmt.Errorf("unmarshal extra: %w", err)
		}
	}

	return security, nil
}

func (s *SecurityStore) GetBySecID(ctx context.Context, secID string) (model.Security, error) {
	row := s.db.QueryRow(ctx, getSecurityBySecIDQuery, secID)

	var security model.Security
	var extraJSON []byte

	err := row.Scan(&security.ID, &security.SecID, &security.Ticker, &security.ShortName,
		&security.Type, &extraJSON, &security.CreateTime)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Security{}, model.ErrNotFound
	}
	if err != nil {
		return model.Security{}, fmt.Errorf("select security by sec_id: %w", err)
	}

	if len(extraJSON) > 0 {
		if err := json.Unmarshal(extraJSON, &security.Extra); err != nil {
			return model.Security{}, fmt.Errorf("unmarshal extra: %w", err)
		}
	}

	return security, nil
}

func (s *SecurityStore) List(ctx context.Context, limit, offset int64) ([]model.Security, error) {
	rows, err := s.db.Query(ctx, listSecuritiesQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select securities: %w", err)
	}

	securities, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Security, error) {
		var security model.Security
		var extraJSON []byte

		err := row.Scan(&security.ID, &security.SecID, &security.Ticker, &security.ShortName,
			&security.Type, &extraJSON, &security.CreateTime)
		if err != nil {
			return model.Security{}, err
		}

		if len(extraJSON) > 0 {
			if err := json.Unmarshal(extraJSON, &security.Extra); err != nil {
				return model.Security{}, fmt.Errorf("unmarshal extra: %w", err)
			}
		}

		return security, nil
	})
	if err != nil {
		return nil, fmt.Errorf("collect securities: %w", err)
	}

	return securities, nil
}
