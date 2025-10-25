package profile

import (
	_ "embed"
	"encoding/json"
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
	//go:embed testdata/get_profiles_response.json
	getProfilesResponse []byte
)

func TestImplementation_GetProfiles(t *testing.T) {
	t.Parallel()

	userID := uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21")
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

				profiles := []model.Profile{
					{
						ID:         uuid.MustParse("463d4cc6-023a-4d54-9da5-e6445367bf21"),
						UserID:     userID,
						Name:       "IMOEX",
						CreateTime: now,
					},
				}

				service.GeByUserIDMock.
					When(minimock.AnyContext, userID, 10, 0).
					Then(profiles, nil)

				return &Implementation{
					service: service,
					logger:  zap.NewNop(),
				}
			},
			wantResp: getProfilesResponse,
			wantCode: http.StatusOK,
		},
		{
			name: "get profiles error",
			handler: func(mc *minimock.Controller) *Implementation {
				service := NewServiceMock(mc)

				service.GeByUserIDMock.
					When(minimock.AnyContext, userID, 10, 0).
					Then(nil, assert.AnError)

				return &Implementation{
					service: service,
					logger:  zap.NewNop(),
				}
			},
			wantResp: json.RawMessage(`{ "message": "Get profiles by user ID" }`),
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)

			handler := tc.handler(mc)

			c := httptest.Case{
				Handler: func(w http.ResponseWriter, r *http.Request) {
					ctx := model.WithAuthClaims(r.Context(), model.AuthClaims{
						UserID: userID,
					})
					r = r.WithContext(ctx)
					handler.GetProfiles(w, r)
				},
				ReqBody: tc.req,
			}

			gotResp, gotStatusCode := c.Do(t)
			assert.Equal(t, tc.wantCode, gotStatusCode)
			assert.JSONEq(t, string(tc.wantResp), string(gotResp))
		})
	}
}
