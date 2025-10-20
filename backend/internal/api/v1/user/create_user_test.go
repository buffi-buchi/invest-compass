package user

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/buffi-buchi/invest-compass/backend/internal/domain/model"
	"github.com/buffi-buchi/invest-compass/backend/pkg/httptest"
)

var (
	//go:embed testdata/create_user_invalid_request.xml
	createUserInvalidRequest []byte

	//go:embed testdata/create_user_request.json
	createUserRequest []byte

	//go:embed testdata/create_user_response.json
	createUserResponse []byte
)

func TestImplementation_CreateUser(t *testing.T) {
	t.Parallel()

	now := time.Date(2025, time.September, 10, 0, 0, 0, 0, time.UTC)

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

				user := model.User{
					ID:         uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"),
					Email:      "user@example.com",
					Password:   "password123",
					CreateTime: now,
				}

				service.CreateMock.
					When(minimock.AnyContext, model.User{
						Email:    user.Email,
						Password: user.Password,
					}).
					Then(user, nil)

				return &Implementation{
					service: service,
					logger:  zap.NewNop(),
				}
			},
			req:      createUserRequest,
			wantResp: createUserResponse,
			wantCode: http.StatusCreated,
		},
		{
			name: "already exists",
			handler: func(mc *minimock.Controller) *Implementation {
				service := NewServiceMock(mc)

				user := model.User{
					ID:         uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"),
					Email:      "user@example.com",
					Password:   "password123",
					CreateTime: now,
				}

				service.CreateMock.
					When(minimock.AnyContext, model.User{
						Email:    user.Email,
						Password: user.Password,
					}).Then(model.User{}, fmt.Errorf("some error: %w", model.ErrAlreadyExists))

				return &Implementation{
					service: service,
					logger:  zap.NewNop(),
				}
			},
			req:      createUserRequest,
			wantResp: json.RawMessage(`{ "message": "User already exists" }`),
			wantCode: http.StatusConflict,
		},
		{
			name: "create user error",
			handler: func(mc *minimock.Controller) *Implementation {
				service := NewServiceMock(mc)

				user := model.User{
					ID:         uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"),
					Email:      "user@example.com",
					Password:   "password123",
					CreateTime: now,
				}

				service.CreateMock.
					When(minimock.AnyContext, model.User{
						Email:    user.Email,
						Password: user.Password,
					}).Then(model.User{}, errors.New("some error"))

				return &Implementation{
					service: service,
					logger:  zap.NewNop(),
				}
			},
			req:      createUserRequest,
			wantResp: json.RawMessage(`{ "message": "Create user error" }`),
			wantCode: http.StatusInternalServerError,
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
			req:      createUserInvalidRequest,
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
				Handler: handler.CreateUser,
				ReqBody: tc.req,
			}

			gotResp, gotStatusCode := c.Do(t)
			assert.Equal(t, tc.wantCode, gotStatusCode)
			assert.JSONEq(t, string(tc.wantResp), string(gotResp))
		})
	}
}
