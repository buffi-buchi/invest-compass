package moex

import (
	"context"
	"fmt"
	"net/http"

	"github.com/buffi-buchi/invest-compass/backend/pkg/date"
)

func (c *Client) GetIndexSecurities(ctx context.Context, id string) ([]IndexSecurity, error) {
	params := getIndexSecuritiesRequest{
		ID: id,
	}

	securities := make([]IndexSecurity, 0)

	for ctx.Err() == nil {
		req, err := params.buildRequest(ctx, c.baseURL)
		if err != nil {
			return nil, err
		}

		var resp getIndexSecuritiesResponse

		err = c.doRequest(req, &resp)
		if err != nil {
			return nil, err
		}

		securities = append(securities, resp.Securities...)

		if len(resp.Cursors) != 1 {
			return nil, fmt.Errorf("invalid cursor")
		}

		cursor := resp.Cursors[0]

		if cursor.Index+cursor.PageSize > cursor.Total {
			break
		}

		params.start += cursor.PageSize
	}

	return securities, nil
}

type getIndexSecuritiesRequest struct {
	baseRequestParams
	ID string
}

func (r getIndexSecuritiesRequest) buildRequest(ctx context.Context, baseURL string) (*http.Request, error) {
	const pattern = "/iss/statistics/engines/stock/markets/index/analytics/%s.json"

	url := baseURL + fmt.Sprintf(pattern, r.ID)
	values := r.buildURLValues()

	if len(values) > 0 {
		url += "?" + values.Encode()
	}

	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

type getIndexSecuritiesResponse struct {
	Securities []IndexSecurity
	Cursors    []getIndexSecuritiesCursor `json:"analytics.cursor"`
	Dates      []getIndexSecuritiesDate   `json:"analytics.dates"`
}

type getIndexSecuritiesCursor struct {
	baseCursor
	PrevDate date.Date `json:"PREV_DATE,omitempty"`
	NextDate date.Date `json:"NEXT_DATE,omitempty"`
}

type getIndexSecuritiesDate struct {
	From date.Date `json:"from"`
	Till date.Date `json:"till"`
}

type IndexSecurity struct {
	ID        string  `json:"secids"`
	IndexID   string  `json:"indexid"`
	Ticker    string  `json:"ticker"`
	ShortName string  `json:"shortnames"`
	Weight    float64 `json:"weight"`
}
