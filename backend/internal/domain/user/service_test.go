package user

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

func TestService_Create(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		service  func(mc *minimock.Controller) *Service
		user     model.User
		wantUser model.User
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)

				store.CreateMock.Set(func(ctx context.Context, u model.User) (model.User, error) {
					assert.Equal(mc, "user@example.com", u.Email)
					assert.True(mc, u.CheckPassword("password123"))

					return model.User{
						ID:         uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"),
						Email:      u.Email,
						Password:   u.Password,
						CreateTime: now,
					}, nil
				})

				return &Service{
					store: store,
				}
			},
			user: model.User{
				Email:    "user@example.com",
				Password: "password123",
			},
			wantUser: model.User{
				ID:         uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"),
				Email:      "user@example.com",
				CreateTime: now,
			},
			wantErr: assert.NoError,
		},
		{
			name: "hash password error",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)

				return &Service{
					store: store,
				}
			},
			user: model.User{
				Password: "password123password123password123password123password123password123password123",
			},
			wantUser: model.User{},
			wantErr:  assert.Error,
		},
		{
			name: "create user error",
			service: func(mc *minimock.Controller) *Service {
				store := NewStoreMock(mc)

				store.CreateMock.Set(func(ctx context.Context, u model.User) (model.User, error) {
					assert.Equal(mc, "user@example.com", u.Email)
					assert.True(mc, u.CheckPassword("password123"))

					return model.User{}, assert.AnError
				})

				return &Service{
					store: store,
				}
			},
			user: model.User{
				Email:    "user@example.com",
				Password: "password123",
			},
			wantUser: model.User{},
			wantErr:  assert.Error,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			service := tc.service(mc)

			gotUser, gotErr := service.Create(context.Background(), tc.user)

			tc.wantUser.Password = gotUser.Password

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantUser, gotUser)
		})
	}
}
