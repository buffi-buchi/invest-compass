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
	//go:embed testdata/get_securities_response.json
	getSecuritiesResponse []byte
)

func TestClient_GetSecuritiesByIndex(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		server         func(t *testing.T) *httptest.Server
		client         func(t *testing.T, server *httptest.Server) *Client
		wantSecurities []Security
		wantMarketData []MarketData
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			server: func(t *testing.T) *httptest.Server {
				const url = "/iss/engines/stock/markets/shares/securities.json?index=IMOEX&iss.json=extended&iss.meta=off&iss.version=on&marketprice_board=1"

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, url, r.URL.String())

					w.WriteHeader(http.StatusOK)
					w.Write(getSecuritiesResponse)
				}))
			},
			client: func(t *testing.T, server *httptest.Server) *Client {
				return &Client{
					client:  &http.Client{},
					baseURL: server.URL,
				}
			},
			wantSecurities: []Security{
				{
					ID:        "AFLT",
					Name:      "Аэрофлот-росс.авиалин(ПАО)ао",
					ShortName: "Аэрофлот",
					BoardID:   "TQBR",
					BoardName: "Т+: Акции и ДР - безадрес.",
					FaceUnit:  "SUR",
					ISIN:      "RU0009062285",
				},
				{
					ID:        "GAZP",
					Name:      "\"Газпром\" (ПАО) ао",
					ShortName: "ГАЗПРОМ ао",
					BoardID:   "TQBR",
					BoardName: "Т+: Акции и ДР - безадрес.",
					FaceUnit:  "SUR",
					ISIN:      "RU0007661625",
				},
			},
			wantMarketData: []MarketData{
				{
					SecurityID: "AFLT",
					BoardID:    "TQBR",
					Open:       55.37,
					Low:        53.53,
					High:       55.5,
					Last:       55.02,
				},
				{
					SecurityID: "GAZP",
					BoardID:    "TQBR",
					Open:       128.19,
					Low:        125.35,
					High:       128.82,
					Last:       127.56,
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

			gotSecurities, gotMarketData, gotErr := client.GetSecuritiesByIndex(context.Background(), "IMOEX")

			tc.wantErr(t, gotErr)
			assert.Equal(t, tc.wantSecurities, gotSecurities)
			assert.Equal(t, tc.wantMarketData, gotMarketData)
		})
	}
}
