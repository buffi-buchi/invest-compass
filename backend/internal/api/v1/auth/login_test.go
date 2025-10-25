package auth

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/buffi-buchi/invest-compass/backend/pkg/httptest"
)

var (
	//go:embed testdata/login_invalid_request.xml
	loginInvalidRequest []byte

	//go:embed testdata/login_request.json
	loginRequest []byte
)

func TestImplementation_Login(t *testing.T) {
	t.Parallel()

	token, _ := jwt.New(jwt.SigningMethodHS256).SignedString([]byte("secret"))

	cases := []struct {
		name     string
		handler  func(mc *minimock.Controller) *Implementation
		req      []byte
		wantResp []byte
		wantCode int
	}{
		{
			name: "success",
			handler: func(mc *minimock.Controller) *Implementation {
				service := NewServiceMock(mc)

				service.LoginMock.
					When(minimock.AnyContext, "user@example.com", "password123").
					Then(token, nil)

				return &Implementation{
					service: service,
					logger:  zap.NewNop(),
				}
			},
			req:      loginRequest,
			wantResp: json.RawMessage(fmt.Sprintf(`{ "token": "%s" }`, token)),
			wantCode: http.StatusOK,
		},
		{
			name: "login error",
			handler: func(mc *minimock.Controller) *Implementation {
				service := NewServiceMock(mc)

				service.LoginMock.
					When(minimock.AnyContext, "user@example.com", "password123").
					Then("", assert.AnError)

				return &Implementation{
					service: service,
					logger:  zap.NewNop(),
				}
			},
			req:      loginRequest,
			wantResp: json.RawMessage(`{ "message": "Unauthenticated" }`),
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "invalid request",
			handler: func(mc *minimock.Controller) *Implementation {
				service := NewServiceMock(mc)

				return &Implementation{
					service: service,
					logger:  zap.NewNop(),
				}
			},
			req:      loginInvalidRequest,
			wantResp: json.RawMessage(`{ "message": "Invalid request: invalid character '<' looking for beginning of value" }`),
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			handler := tc.handler(mc)

			c := httptest.Case{
				Handler: handler.Login,
				ReqBody: tc.req,
			}

			gotResp, gotStatusCode := c.Do(t)
			assert.Equal(t, tc.wantCode, gotStatusCode)
			assert.JSONEq(t, string(tc.wantResp), string(gotResp))
		})
	}
}
