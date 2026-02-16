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

	//go:embed testdata/list_indexes.sql
	listIndexesQuery string
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
		{
			name: "List success",
			run: func(t *testing.T) {
				store := &IndexStore{
					db: db,
				}

				// Act.
				_, err = db.Exec(ctx, listIndexexQuery)
				require.NoError(t, err)

				gotIndexes, gotErr := store.List(ctx, 5, 3)

				// Check.
				require.NoError(t, gotErr)

				for i := range gotIndexes {
					gotIndexes[i].CreateTime = gotIndexes[i].CreateTime.UTC()
				}

				assert.ElementsMatch(t, []model.Index{
					{
						ID:         uuid.MustParse("a6d2e5b9-3c41-4f7a-8e2d-5b9f1c3a7e64"),
						Ticker:     "IMOEX2",
						Name:       "IMOEX",
						CreateTime: now,
					},
					{
						ID:         uuid.MustParse("9b1e4c72-6d3f-4a8b-b2e7-1c5f9a3d8e20"),
						Ticker:     "IMOEX3",
						Name:       "IMOEX",
						CreateTime: now,
					},
					{
						ID:         uuid.MustParse("4e8a2d1c-7f35-4b9e-9a61-3d7c2f5b8e14"),
						Ticker:     "IMOEX4",
						Name:       "IMOEX",
						CreateTime: now,
					},
					{
						ID:         uuid.MustParse("c2f7a9d4-5e31-4c8b-8d2f-6a1e3b9c4d75"),
						Ticker:     "IMOEX5",
						Name:       "IMOEX",
						CreateTime: now,
					},
					{
						ID:         uuid.MustParse("7a5d3c9e-1b42-4e8f-a6d3-9c2b7e1f4a68"),
						Ticker:     "IMOEX6",
						Name:       "IMOEX",
						CreateTime: now,
					},
				}, gotIndexes)

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
