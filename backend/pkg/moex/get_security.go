package moex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetSecurity(ctx context.Context, id string) (Security, MarketData, error) {
	req, err := getSecurityRequest{ID: id}.buildRequest(ctx, c.baseURL)
	if err != nil {
		return Security{}, MarketData{}, err
	}

	var resp getSecurityResponse

	err = c.doRequest(req, &resp)
	if err != nil {
		return Security{}, MarketData{}, err
	}

	if len(resp.Securities) == 0 || len(resp.MarketData) == 0 {
		return Security{}, MarketData{}, ErrNotFound
	}

	if len(resp.Securities) > 1 || len(resp.MarketData) > 1 {
		return Security{}, MarketData{}, fmt.Errorf("too many records")
	}

	return resp.Securities[0], resp.MarketData[0], nil
}

type getSecurityRequest struct {
	baseRequestParams
	ID string
}

func (r getSecurityRequest) buildRequest(ctx context.Context, baseURL string) (*http.Request, error) {
	const pattern = "/iss/engines/stock/markets/shares/securities/%s.json"

	url := baseURL + fmt.Sprintf(pattern, r.ID)
	values := r.buildURLValues()

	if len(values) > 0 {
		url += "?" + values.Encode()
	}

	return http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
}

func (r getSecurityRequest) buildURLValues() url.Values {
	values := r.baseRequestParams.buildURLValues()

	// Flag "marketprice_board" tells API to return data only for the main trading mode of the security.
	// Flag "primary_board" works the same as the "marketprice_board" in out case so do not use it.
	values.Add("marketprice_board", "1")

	return values
}

type getSecurityResponse struct {
	Securities []Security   `json:"securities"`
	MarketData []MarketData `json:"marketdata"`
}

type Security struct {
	ID        string `json:"SECID"`
	Name      string `json:"SECNAME"`
	ShortName string `json:"SHORTNAME"`
	BoardID   string `json:"BOARDID"`
	BoardName string `json:"BOARDNAME"`
	FaceUnit  string `json:"FACEUNIT"`
	ISIN      string `json:"ISIN"`
}

type MarketData struct {
	SecurityID string  `json:"SECID"`
	BoardID    string  `json:"BOARDID"`
	Open       float64 `json:"OPEN"`
	Low        float64 `json:"LOW"`
	High       float64 `json:"HIGH"`
	Last       float64 `json:"LAST"`
}
