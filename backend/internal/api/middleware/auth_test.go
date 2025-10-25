package middleware

import (
	"encoding/json"
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

	now := time.Now()
	userID := uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21")
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
		header     string
		wantResp   []byte
		wantCode   int
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
			header:   "Bearer " + token,
			wantResp: json.RawMessage(`{ }`),
			wantCode: http.StatusOK,
		},
		{
			name: "no authorization header",
			middleware: func(mc *minimock.Controller) Middleware {
				jwtProvider := NewJWTProviderMock(mc)

				return NewAuthMiddleware(jwtProvider)
			},
			header:   "",
			wantResp: json.RawMessage(`{ "message": "Authorization header is required" }`),
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "invalid token",
			middleware: func(mc *minimock.Controller) Middleware {
				jwtProvider := NewJWTProviderMock(mc)

				return NewAuthMiddleware(jwtProvider)
			},
			header:   "Hello " + token,
			wantResp: json.RawMessage(`{ "message": "Authorization header format is invalid" }`),
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "validate token error",
			middleware: func(mc *minimock.Controller) Middleware {
				jwtProvider := NewJWTProviderMock(mc)

				jwtProvider.ValidateMock.
					When(token).
					Then(model.AuthClaims{}, assert.AnError)

				return NewAuthMiddleware(jwtProvider)
			},
			header:   "Bearer " + token,
			wantResp: json.RawMessage(`{ "message": "Invalid token" }`),
			wantCode: http.StatusUnauthorized,
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

				_ = json.NewEncoder(w).Encode(struct{}{})
			}))

			header := http.Header{}
			header.Set("Authorization", tc.header)

			c := httptest.Case{
				Handler: handler.ServeHTTP,
				Headers: header,
			}

			gotResp, gotStatusCode := c.Do(t)
			assert.Equal(t, tc.wantCode, gotStatusCode)
			assert.JSONEq(t, string(tc.wantResp), string(gotResp))
		})
	}
}
