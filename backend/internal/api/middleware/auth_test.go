package middleware

import (
	"net/http"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
	jwtprovider "github.com/buffi-buchi/invest-compass/backend/internal/provider/jwt"
	"github.com/buffi-buchi/invest-compass/backend/pkg/httptest"
)

func Test_NewAuthMiddleware(t *testing.T) {
	t.Parallel()

	userID := uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21")
	now := time.Now()
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtprovider.Claims{
		AuthClaims: model.AuthClaims{
			UserID: userID,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}).SignedString([]byte("secret"))

	cases := []struct {
		name       string
		middleware func(mc *minimock.Controller) Middleware
	}{
		{
			name: "success",
			middleware: func(mc *minimock.Controller) Middleware {
				jwtProvider := NewJWTProviderMock(mc)

				jwtProvider.ValidateMock.
					When(token).
					Then(model.AuthClaims{
						UserID: userID,
					}, nil)

				return NewAuthMiddleware(jwtProvider)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			middleware := tc.middleware(mc)

			handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				claims, _ := model.AuthClaimsValue(r.Context())
				assert.Equal(t, userID, claims.UserID)
			}))

			header := http.Header{}
			header.Set("Authorization", "Bearer "+token)

			c := httptest.Case{
				Handler: handler.ServeHTTP,
				Headers: header,
			}

			gotResp, gotStatusCode := c.Do(t)
			assert.Equal(t, http.StatusOK, gotStatusCode)
			assert.Equal(t, []byte{}, gotResp)
		})
	}
}
