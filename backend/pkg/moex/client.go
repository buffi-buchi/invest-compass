package moex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var (
	ErrNotFound = errors.New("not found")
)

// Client is an HTTP client for Moscow Exchange (MOEX).
type Client struct {
	client  *http.Client
	baseURL string
}

// NewClient returns a newly initialized Client for Moscow Exchange (MOEX).
func NewClient(baseURL string) *Client {
	return &Client{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *Client) doRequest(req *http.Request, result any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Response is always a JSON array at the top level with two elements when the parameter
	// iss.json=extended is set.
	// First, parse the response as a slice of JSON objects and validate the length.
	// Second, parse the second element of the slice as the result.
	var parts []json.RawMessage

	if err = json.NewDecoder(resp.Body).Decode(&parts); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if len(parts) != 2 {
		return fmt.Errorf("unexpected response format")
	}

	if err = json.Unmarshal(parts[1], result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

// baseRequestParams defines base request parameters.
type baseRequestParams struct {
	start int64
	limit int64
}

// buildURLValues constructs base URL query parameters of request. Parameters
// iss.json, iss.meta and iss.version are constant because they define response
// format.
func (p baseRequestParams) buildURLValues() url.Values {
	values := make(url.Values)

	if p.start > 0 {
		values.Add("start", strconv.FormatInt(p.start, 10))
	}

	if p.limit > 0 {
		values.Add("limit", strconv.FormatInt(p.limit, 10))
	}

	values.Add("iss.json", "extended")
	values.Add("iss.meta", "off")
	values.Add("iss.version", "on")

	return values
}

type baseCursor struct {
	Index    int64 `json:"INDEX"`
	Total    int64 `json:"TOTAL"`
	PageSize int64 `json:"PAGESIZE"`
}
