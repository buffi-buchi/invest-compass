package moex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) GetSecurity(ctx context.Context, id string) error {
	return nil
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
}
