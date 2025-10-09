package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileStore struct {
	db *pgxpool.Pool
}

func NewProfileStore(db *pgxpool.Pool) *ProfileStore {
	return &ProfileStore{
		db: db,
	}
}
