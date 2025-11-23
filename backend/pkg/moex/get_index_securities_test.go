package moex

import (
	"context"
	_ "embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	//go:embed testdata/get_analytics_securities_response.json
	getAnalyticsSecuritiesResponse []byte
)

func TestClient_GetIndexSecurities(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                string
		server              func(t *testing.T) *httptest.Server
		client              func(t *testing.T, server *httptest.Server) *Client
		wantIndexSecurities []IndexSecurity
		wantErr             assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			server: func(t *testing.T) *httptest.Server {
				const url = "/iss/statistics/engines/stock/markets/index/analytics/IMOEX.json?iss.json=extended&iss.meta=off&iss.version=on"

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, url, r.URL.String())

					w.WriteHeader(http.StatusOK)
					w.Write(getAnalyticsSecuritiesResponse)
				}))
			},
			client: func(t *testing.T, server *httptest.Server) *Client {
				return &Client{
					client:  &http.Client{},
					baseURL: server.URL,
				}
			},
			wantIndexSecurities: []IndexSecurity{
				{
					ID:        "AFLT",
					IndexID:   "IMOEX",
					Ticker:    "AFLT",
					ShortName: "Аэрофлот",
					Weight:    0.7,
				},
				{
					ID:        "GAZP",
					IndexID:   "IMOEX",
					Ticker:    "GAZP",
					ShortName: "ГАЗПРОМ ао",
					Weight:    11.45,
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

			gotIndexSecurities, gotErr := client.GetIndexSecurities(context.Background(), "IMOEX")

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantIndexSecurities, gotIndexSecurities)
		})
	}
}
