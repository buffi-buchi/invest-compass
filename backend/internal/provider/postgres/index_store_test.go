//go:build integration

package postgres

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed testdata/create_test_index.sql
	createTestIndexQuery string
)

func TestIndexStore_GeByTicker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "GetByTicker success",
			run: func(t *testing.T) {
				store := &IndexStore{
					db: db,
				}

				// Act.
				_, err = db.Exec(ctx, createTestIndexQuery)
				require.NoError(t, err)

				gotIndex, gotErr := store.GetByCode(ctx, "MOEXBC")

				// Check.
				require.NoError(t, gotErr)

				gotIndex.CreateTime = gotIndex.CreateTime.UTC()

				assert.Equal(t, model.Index{
					ID:         uuid.MustParse("6f1b2a6e-9c3e-4a2d-8b7f-3e5f1c9a4d21"),
					Ticker:     "MOEXBC",
					Name:       "MOEXBC",
					CreateTime: now,
				}, gotIndex)

				// Cleanup.
				_, err = db.Exec(ctx, `TRUNCATE TABLE "indexes"`)
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, tc.run)
	}
}
