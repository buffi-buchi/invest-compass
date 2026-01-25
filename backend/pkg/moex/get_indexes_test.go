package moex

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/buffi-buchi/invest-compass/backend/pkg/date"
)

var (
	//go:embed testdata/get_analytics_response.json
	getAnalyticsResponse []byte
)

func TestClient_GetIndexes(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		server      func(t *testing.T) *httptest.Server
		client      func(t *testing.T, server *httptest.Server) *Client
		wantIndexes []Index
		wantErr     assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			server: func(t *testing.T) *httptest.Server {
				const url = "/iss/statistics/engines/stock/markets/index/analytics?iss.json=extended&iss.meta=off&iss.version=on"

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, url, r.URL.String())

					w.WriteHeader(http.StatusOK)
					w.Write(getAnalyticsResponse)
				}))
			},
			client: func(t *testing.T, server *httptest.Server) *Client {
				return &Client{
					baseURL: server.URL,
					client:  &http.Client{},
				}
			},
			wantIndexes: []Index{
				{
					ID:        "IMOEX",
					ShortName: "Индекс МосБиржи",
					From:      date.NewDate(2001, time.January, 3),
					Till:      date.NewDate(2025, time.November, 14),
				},
				{
					ID:        "MOEXBC",
					ShortName: "Индекс голубых фишек",
					From:      date.NewDate(2010, time.January, 14),
					Till:      date.NewDate(2025, time.November, 14),
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := tc.server(t)
			defer server.Close()

			client := tc.client(t, server)

			gotIndexes, gotErr := client.GetIndexes(context.Background())

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantIndexes, gotIndexes)
		})
	}
}
