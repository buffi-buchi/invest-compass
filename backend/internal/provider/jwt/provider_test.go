package jwt

import (
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
)

func TestProvider_Generate(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name      string
		provider  func(mc *minimock.Controller) *Provider
		user      model.User
		wantToken string
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			provider: func(mc *minimock.Controller) *Provider {
				return &Provider{
					secretKey:           []byte("secret"),
					issuer:              "issuer",
					accessTokenDuration: time.Minute,
					now:                 func() time.Time { return now },
				}
			},
			user: model.User{
				ID: uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"),
			},
			wantToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI0NjNkNGNjNi0wMjNhLTRkNTQtOWRhNS1lNjQ0NTM2N2JmMjEiLCJFbWFpbCI6IiIsImlzcyI6Imlzc3VlciIsInN1YiI6IjQ2M2Q0Y2M2LTAyM2EtNGQ1NC05ZGE1LWU2NDQ1MzY3YmYyMSIsImV4cCI6MTc1NzQ2MjQ2MCwibmJmIjoxNzU3NDYyNDAwLCJpYXQiOjE3NTc0NjI0MDB9.aDPHBbTAKsUUvmRRlJJBSa_1g6FQTBmQSzmUKGhiShQ",
			wantErr:   assert.NoError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			provider := tc.provider(mc)

			gotToken, gotErr := provider.Generate(tc.user)

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantToken, gotToken)
		})
	}
}

func TestProvider_Validate(t *testing.T) {
	t.Parallel()

	now := time.Now()
	secret := []byte("secret")
	userID := uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21")
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		AuthClaims: model.AuthClaims{
			UserID: userID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}).SignedString(secret)

	cases := []struct {
		name       string
		provider   func(mc *minimock.Controller) *Provider
		token      string
		wantClaims model.AuthClaims
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			provider: func(mc *minimock.Controller) *Provider {
				return &Provider{
					secretKey: secret,
				}
			},
			token: token,
			wantClaims: model.AuthClaims{
				UserID: userID,
			},
			wantErr: assert.NoError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			provider := tc.provider(mc)

			gotClaims, gotErr := provider.Validate(tc.token)

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantClaims, gotClaims)
		})
	}
}
