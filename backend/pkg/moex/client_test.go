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

func TestClient_GetIndex(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name       string
		server     func(t *testing.T) *httptest.Server
		client     func(t *testing.T, server *httptest.Server) *Client
		ticker     string
		wantStocks []IndexStock
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			server: func(t *testing.T) *httptest.Server {
				var (
					i   int
					url = "/iss/statistics/engines/stock/markets/index/analytics/IMOEX.json?iss.json=extended&iss.meta=off&iss.version=on"
				)

				handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if i == 1 {
						assert.Fail(t, "should not be called more than once")

						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					assert.Equal(t, url, r.URL.String())

					w.WriteHeader(http.StatusOK)
					w.Write(getAnalyticsResponse)
					i++
				})
				return httptest.NewServer(handler)
			},
			client: func(t *testing.T, server *httptest.Server) *Client {
				return &Client{
					baseURL: server.URL,
					client:  &http.Client{},
				}
			},
			ticker: "IMOEX",
			wantStocks: []IndexStock{
				{
					IndexID:          "IMOEX",
					TradeDate:        date.NewDate(2025, time.November, 6),
					Ticker:           "AFLT",
					ShortNames:       "Аэрофлот",
					SecIDs:           "AFLT",
					Weight:           0.69,
					TradingSession:   3,
					TradeSessionDate: date.NewDate(2025, time.November, 6),
				},
				{
					IndexID:          "IMOEX",
					TradeDate:        date.NewDate(2025, time.November, 6),
					Ticker:           "GAZP",
					ShortNames:       "ГАЗПРОМ ао",
					SecIDs:           "GAZP",
					Weight:           11.19,
					TradingSession:   3,
					TradeSessionDate: date.NewDate(2025, time.November, 6),
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

			gotStocks, gotErr := client.GetIndex(context.Background(), tc.ticker)
			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantStocks, gotStocks)
		})
	}
}
