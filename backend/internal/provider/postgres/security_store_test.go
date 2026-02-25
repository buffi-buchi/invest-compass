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
	//go:embed testdata/create_test_securities.sql
	createTestSecuritiesQuery string
)

func TestSecurityStore_GetByTicker(t *testing.T) {
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
				store := SecurityStore{
					db: db,
				}

				// Act.
				_, err := db.Exec(ctx, createTestSecuritiesQuery)
				require.NoError(t, err)

				gotSecurity, gotErr := store.GetByTicker(ctx, "LKOH")

				// Check.
				require.NoError(t, gotErr)

				gotSecurity.CreateTime = gotSecurity.CreateTime.UTC()

				assert.Equal(t, model.Security{
					Ticker:     "LKOH",
					ShortName:  "ЛУКОЙЛ",
					CreateTime: now,
				}, gotSecurity)

				// Cleanup.
				_, err = db.Exec(ctx, `TRUNCATE TABLE "securities" CASCADE`)
				require.NoError(t, err)
			},
		},
		{
			name: "GetByTicker fail",
			run: func(t *testing.T) {
				store := SecurityStore{
					db: db,
				}

				// Act.
				_, err := db.Exec(ctx, createTestSecuritiesQuery)
				require.NoError(t, err)

				gotSecurity, gotErr := store.GetByTicker(ctx, "T")

				// Check.
				require.Equal(t, model.ErrNotFound, gotErr)
				require.Equal(t, model.Security{}, gotSecurity)

				// Cleanup.
				_, err = db.Exec(ctx, `TRUNCATE TABLE "securities" CASCADE`)
				require.NoError(t, err)
			},
		},
		{
			name: "List success",
			run: func(t *testing.T) {
				store := &SecurityStore{
					db: db,
				}

				// Act.
				_, err := db.Exec(ctx, createTestSecuritiesQuery)
				require.NoError(t, err)

				gotSecurities, gotErr := store.List(ctx, 2, 1, nil)

				// Check.
				require.NoError(t, gotErr)

				for i := range gotSecurities {
					gotSecurities[i].CreateTime = gotSecurities[i].CreateTime.UTC()
				}

				assert.ElementsMatch(t, []model.Security{
					{
						Ticker:     "GMKN",
						ShortName:  "ГМКНорНик",
						CreateTime: now,
					},
					{
						Ticker:     "LKOH",
						ShortName:  "ЛУКОЙЛ",
						CreateTime: now,
					},
				}, gotSecurities)

				// Cleanup.
				_, err = db.Exec(ctx, `TRUNCATE TABLE "securities" CASCADE`)
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, tc.run)

	}
}
