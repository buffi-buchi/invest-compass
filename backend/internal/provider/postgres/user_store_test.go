//go:build integration

package postgres

import (
	"context"
	_ "embed"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

var (
	//go:embed testdata/create_test_user.sql
	createTestUserQuery string
)

func TestUserStore_Create(t *testing.T) {
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
				store := &UserStore{
					db:  db,
					id:  func() (uuid.UUID, error) { return uuid.Parse("463d4cc6-023a-4d54-9da5-e6445367bf21") },
					now: func() time.Time { return now },
				}

				// Act.
				gotUser, gotErr := store.Create(ctx, model.User{
					Email: "user@example.com",
				})

				// Check.
				assert.NoError(t, gotErr)

				row := db.QueryRow(ctx, getUserByIDQuery, gotUser.ID)

				var wantUser model.User
				err := row.Scan(&wantUser.ID, &wantUser.Email, &wantUser.Password, &wantUser.CreateTime)
				assert.NoError(t, err)

				wantUser.CreateTime = wantUser.CreateTime.UTC()
				assert.Equal(t, wantUser, gotUser)

				// Cleanup.
				_, err = db.Exec(ctx, "TRUNCATE TABLE users CASCADE;")
				assert.NoError(t, err)
			},
		},
		{
			name: "already exists",
			run: func(t *testing.T) {
				store := &UserStore{
					db:  db,
					id:  func() (uuid.UUID, error) { return uuid.NewV7() },
					now: func() time.Time { return now },
				}

				// Act.
				_, err := db.Exec(ctx, createTestUserQuery)
				assert.NoError(t, err)

				gotUser, gotErr := store.Create(ctx, model.User{
					Email: "user@example.com",
				})

				// Check.
				assert.Error(t, gotErr)
				assert.Zero(t, gotUser)

				// Cleanup.
				_, err = db.Exec(ctx, "TRUNCATE TABLE users CASCADE;")
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}

func TestUserStore_GetByID(t *testing.T) {
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
				store := &UserStore{
					db: db,
				}

				// Act.
				_, err := db.Exec(ctx, createTestUserQuery)
				assert.NoError(t, err)

				gotUser, gotErr := store.GetByID(ctx, uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"))

				// Check.
				assert.NoError(t, gotErr)

				gotUser.CreateTime = gotUser.CreateTime.UTC()
				assert.Equal(t, model.User{
					ID:         uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"),
					Email:      "user@example.com",
					Password:   "",
					CreateTime: now,
				}, gotUser)

				// Cleanup.
				_, err = db.Exec(ctx, "TRUNCATE TABLE users CASCADE;")
				assert.NoError(t, err)
			},
		},
		{
			name: "not found",
			run: func(t *testing.T) {
				store := &UserStore{
					db: db,
				}

				// Act.
				gotUser, gotErr := store.GetByID(ctx, uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"))

				// Check.
				assert.Error(t, gotErr)
				assert.Zero(t, gotUser)

				// Cleanup.
				_, err := db.Exec(ctx, "TRUNCATE TABLE users CASCADE;")
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}

func TestUserStore_GetByEmail(t *testing.T) {
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
				store := &UserStore{
					db: db,
				}

				// Act.
				_, err := db.Exec(ctx, createTestUserQuery)
				assert.NoError(t, err)

				gotUser, gotErr := store.GetByEmail(ctx, "user@example.com")

				// Check.
				assert.NoError(t, gotErr)

				gotUser.CreateTime = gotUser.CreateTime.UTC()
				assert.Equal(t, model.User{
					ID:         uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"),
					Email:      "user@example.com",
					Password:   "",
					CreateTime: now,
				}, gotUser)

				// Cleanup.
				_, err = db.Exec(ctx, "TRUNCATE TABLE users CASCADE;")
				assert.NoError(t, err)
			},
		},
		{
			name: "not found",
			run: func(t *testing.T) {
				store := &UserStore{
					db: db,
				}

				// Act.
				gotUser, gotErr := store.GetByEmail(ctx, "user@example.com")

				// Check.
				assert.Error(t, gotErr)
				assert.Zero(t, gotUser)

				// Cleanup.
				_, err := db.Exec(ctx, "TRUNCATE TABLE users CASCADE;")
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.run(t)
		})
	}
}
