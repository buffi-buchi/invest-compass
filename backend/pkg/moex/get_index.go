package moex

import (
	"context"
	"fmt"
	"net/http"

	"github.com/buffi-buchi/invest-compass/backend/pkg/date"
)

type IndexStock struct {
	IndexID          string    `json:"indexid"`
	TradeDate        date.Date `json:"tradedate"`
	Ticker           string    `json:"ticker"`
	ShortNames       string    `json:"shortnames"`
	SecIDs           string    `json:"secids"`
	Weight           float64   `json:"weight"`
	TradingSession   int32     `json:"tradingsession"`
	TradeSessionDate date.Date `json:"trade_session_date"`
}

// GetIndex returns index by its ticker.
func (c *Client) GetIndex(ctx context.Context, ticker string) ([]IndexStock, error) {
	query := getIndexQuery{
		baseQuery: baseQuery{
			extendedJsonFormat: true,
			disableMeta:        true,
			version:            true,
		},
		ticker: ticker,
	}

	req, err := query.request(ctx, c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("build get index request: %w", err)
	}

	stocks := make([]IndexStock, 0)

	for ctx.Err() == nil {
		var resp getIndexResponse

		err = c.getJSON(req, &resp)
		if err != nil {
			return nil, err
		}

		stocks = append(stocks, resp.Stocks...)

		// TODO: Validate cursor.
		if resp.Cursors[0].Index+resp.Cursors[0].PageSize > resp.Cursors[0].Total {
			break
		}

		query.start += resp.Cursors[0].PageSize
		req.URL.RawQuery = query.values().Encode()
	}

	return stocks, nil
}

type getIndexQuery struct {
	baseQuery
	ticker string
}

func (q getIndexQuery) request(ctx context.Context, baseURL string) (*http.Request, error) {
	const pattern = "/iss/statistics/engines/stock/markets/index/analytics/%s.json"

	url := baseURL + fmt.Sprintf(pattern, q.ticker)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	values := q.values()
	req.URL.RawQuery = values.Encode()

	return req, nil
}

type getIndexResponse struct {
	Stocks  []IndexStock      `json:"analytics"`
	Cursors []analyticsCursor `json:"analytics.cursor"`
	Dates   []analyticsDate   `json:"analytics.dates"`
}

type analyticsCursor struct {
	baseCursor
	PrevDate date.Date `json:"PREV_DATE,omitempty"`
	NextDate date.Date `json:"NEXT_DATE,omitempty"`
}

type analyticsDate struct {
	From date.Date `json:"from"`
	Till date.Date `json:"till"`
}
