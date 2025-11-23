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
	//go:embed testdata/get_securities_by_id_response.json
	getSecuritiesByIdResponse []byte
)

func TestClient_GetSecurity(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		server         func(t *testing.T) *httptest.Server
		client         func(t *testing.T, server *httptest.Server) *Client
		wantSecurity   Security
		wantMarketData MarketData
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			server: func(t *testing.T) *httptest.Server {
				const url = "/iss/engines/stock/markets/shares/securities/GAZP.json?iss.json=extended&iss.meta=off&iss.version=on&marketprice_board=1"

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, url, r.URL.String())

					w.WriteHeader(http.StatusOK)
					w.Write(getSecuritiesByIdResponse)
				}))
			},
			client: func(t *testing.T, server *httptest.Server) *Client {
				return &Client{
					client:  &http.Client{},
					baseURL: server.URL,
				}
			},
			wantSecurity: Security{
				ID:        "GAZP",
				Name:      "\"Газпром\" (ПАО) ао",
				ShortName: "ГАЗПРОМ ао",
				BoardID:   "TQBR",
				BoardName: "Т+: Акции и ДР - безадрес.",
				FaceUnit:  "SUR",
				ISIN:      "RU0007661625",
			},
			wantMarketData: MarketData{
				SecurityID: "GAZP",
				BoardID:    "TQBR",
				Open:       124.23,
				Low:        122.64,
				High:       129.29,
				Last:       128.11,
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

			gotSecurity, gotMarketData, gotErr := client.GetSecurity(context.Background(), "GAZP")

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantSecurity, gotSecurity)
			assert.Equal(t, tc.wantMarketData, gotMarketData)
		})
	}
}
