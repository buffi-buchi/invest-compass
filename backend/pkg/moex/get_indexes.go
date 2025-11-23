package moex

import (
	"context"
	"net/http"

	"github.com/buffi-buchi/invest-compass/backend/pkg/date"
)

func (c *Client) GetIndexes(ctx context.Context) ([]Index, error) {
	req, err := getIndexesRequest{}.buildRequest(ctx, c.baseURL)
	if err != nil {
		return nil, err
	}

	var resp getIndexesResponse

	err = c.doRequest(req, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Indexes, nil
}

type getIndexesRequest struct {
	baseRequestParams
}

func (r getIndexesRequest) buildRequest(ctx context.Context, baseURL string) (*http.Request, error) {
	const pattern = "/iss/statistics/engines/stock/markets/index/analytics"

	url := baseURL + pattern
	values := r.buildURLValues()

	if len(values) > 0 {
		url += "?" + values.Encode()
	}

	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

type getIndexesResponse struct {
	Indexes []Index `json:"indices"`
}

type Index struct {
	ID        string    `json:"indexid"`
	ShortName string    `json:"shortname"`
	From      date.Date `json:"from"`
	Till      date.Date `json:"till"`
}
