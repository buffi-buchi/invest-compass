//go:build integration

package postgres

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

var (
	//go:embed testdata/create_test_profiles.sql
	createTestProfilesQuery string
)

func TestProfileStore_GeByUserID(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "success",
			run: func(t *testing.T) {
				store := &ProfileStore{
					db: db,
				}

				// Act.
				_, err := db.Exec(ctx, createTestUserQuery)
				require.NoError(t, err)
				_, err = db.Exec(ctx, createTestProfilesQuery)
				require.NoError(t, err)

				gotProfiles, gotErr := store.GeByUserID(ctx, uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"), 10, 0)

				// Check.
				require.NoError(t, gotErr)

				for i := range gotProfiles {
					gotProfiles[i].CreateTime = gotProfiles[i].CreateTime.UTC()
				}

				assert.ElementsMatch(t, []model.Profile{
					{
						ID:         uuid.MustParse("a3e8015d-2ef2-4490-b8af-6f51f2f0e038"),
						UserID:     uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"),
						Name:       "IMOEX",
						CreateTime: now,
					},
					{
						ID:         uuid.MustParse("3147b569-1c64-41ff-8f88-329fe2bf8a6c"),
						UserID:     uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"),
						Name:       "MOEXBC",
						CreateTime: now,
					},
				}, gotProfiles)

				// Cleanup.
				_, err = db.Exec(ctx, "TRUNCATE TABLE users, profiles CASCADE")
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}
