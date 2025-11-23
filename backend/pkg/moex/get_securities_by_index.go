package moex

import (
	"context"
	"net/http"
	"net/url"
)

func (c *Client) GetSecuritiesByIndex(ctx context.Context, indexID string) ([]Security, []MarketData, error) {
	req, err := getSecuritiesByIndexRequest{IndexID: indexID}.buildRequest(ctx, c.baseURL)
	if err != nil {
		return nil, nil, err
	}

	var resp getSecuritiesByIndexResponse

	err = c.doRequest(req, &resp)
	if err != nil {
		return nil, nil, err
	}

	if len(resp.Securities) == 0 || len(resp.MarketData) == 0 {
		return nil, nil, ErrNotFound
	}

	return resp.Securities, resp.MarketData, nil
}

type getSecuritiesByIndexRequest struct {
	baseRequestParams
	IndexID string
}

func (r getSecuritiesByIndexRequest) buildRequest(ctx context.Context, baseURL string) (*http.Request, error) {
	const pattern = "/iss/engines/stock/markets/shares/securities.json"

	url := baseURL + pattern
	values := r.buildURLValues()

	if len(values) > 0 {
		url += "?" + values.Encode()
	}

	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

func (r getSecuritiesByIndexRequest) buildURLValues() url.Values {
	values := r.baseRequestParams.buildURLValues()

	// Выводить акции из базы индекса.
	// Только для фондового рынка.
	values.Add("index", r.IndexID)

	// Flag "marketprice_board" tells API to return data only for the main trading mode of the security.
	// Flag "primary_board" works the same as the "marketprice_board" in out case so do not use it.
	values.Add("marketprice_board", "1")

	return values
}

type getSecuritiesByIndexResponse struct {
	Securities []Security   `json:"securities"`
	MarketData []MarketData `json:"marketdata"`
}
