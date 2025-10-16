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
	//go:embed queries/get_profiles_by_user_id.sql
	getProfilesByUserIDQuery string
)

type ProfileStore struct {
	db *pgxpool.Pool
}

func NewProfileStore(db *pgxpool.Pool) *ProfileStore {
	return &ProfileStore{
		db: db,
	}
}

func (s *ProfileStore) GeByUserID(
	ctx context.Context,
	userID uuid.UUID,
	limit int64,
	offset int64,
) ([]model.Profile, error) {
	rows, err := s.db.Query(ctx, getProfilesByUserIDQuery, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select profiles by user id: %w", err)
	}

	profiles, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (model.Profile, error) {
		var profile model.Profile

		return profile, row.Scan(&profile.ID, &profile.UserID, &profile.Name, &profile.CreateTime)
	})
	if err != nil {
		return nil, fmt.Errorf("select profiles by user id: %w", err)
	}

	return profiles, nil
}
