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

type Client struct {
	baseURL string
	client  *http.Client
}

// NewClient returns a newly initialized Client for Moscow Exchange (MOEX).
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *Client) getJSON(req *http.Request, result any) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("execute get index request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var parts []json.RawMessage

	if err = json.NewDecoder(resp.Body).Decode(&parts); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	// TODO: Validate version header.

	if len(parts) != 2 {
		return errors.New("unexpected response format")
	}

	if err = json.Unmarshal(parts[1], result); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

type baseQuery struct {
	only               []string
	extendedJsonFormat bool
	disableMeta        bool
	disableData        bool
	version            bool
	start              int64
	limit              int64
}

func (q baseQuery) values() url.Values {
	values := make(url.Values)

	for _, block := range q.only {
		values.Add("iss.only", block)
	}

	if q.extendedJsonFormat {
		values.Add("iss.json", "extended")
	}

	if q.disableMeta {
		values.Add("iss.meta", "off")
	}

	if q.disableData {
		values.Add("iss.data", "off")
	}

	if q.version {
		values.Add("iss.version", "on")
	}

	if q.start != 0 {
		values.Add("start", strconv.FormatInt(q.start, 10))
	}

	if q.limit != 0 {
		values.Add("limit", strconv.FormatInt(q.limit, 10))
	}

	return values
}

type baseCursor struct {
	Index    int64 `json:"INDEX"`
	Total    int64 `json:"TOTAL"`
	PageSize int64 `json:"PAGESIZE"`
}
