//go:build integration

package postgres

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
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
				store := IndexStore{
					db: db,
				}

				// Act.
				_, err := db.Exec(ctx, createTestIndexQuery)
				require.NoError(t, err)

				gotIndex, gotErr := store.GetByTicker(ctx, "MOEXBC")

				// Check.
				require.NoError(t, gotErr)

				gotIndex.CreateTime = gotIndex.CreateTime.UTC()

				assert.Equal(t, model.Index{
					Ticker:     "MOEXBC",
					Name:       "MOEXBC",
					CreateTime: now,
				}, gotIndex)

				// Cleanup.
				_, err = db.Exec(ctx, `TRUNCATE TABLE "indexes" CASCADE`)
				require.NoError(t, err)
			},
		},
		{
			name: "List success",
			run: func(t *testing.T) {
				store := &IndexStore{
					db: db,
				}

				// Act.
				_, err := db.Exec(ctx, createTestIndexQuery)
				require.NoError(t, err)

				gotIndexes, gotErr := store.List(ctx, 2, 1)

				// Check.
				require.NoError(t, gotErr)

				for i := range gotIndexes {
					gotIndexes[i].CreateTime = gotIndexes[i].CreateTime.UTC()
				}

				assert.ElementsMatch(t, []model.Index{
					{
						Ticker:     "IMOEX1",
						Name:       "IMOEX",
						CreateTime: now,
					},
					{
						Ticker:     "MOEXBC",
						Name:       "MOEXBC",
						CreateTime: now,
					},
				}, gotIndexes)

				// Cleanup.
				_, err = db.Exec(ctx, `TRUNCATE TABLE "indexes" CASCADE`)
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, tc.run)

	}
}
