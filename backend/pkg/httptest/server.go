package httptest

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type Case struct {
	Handler http.HandlerFunc
	ReqBody []byte
	Pattern string
	Headers http.Header
}

func (c *Case) Do(t *testing.T) ([]byte, int) {
	if c.Pattern == "" {
		c.Pattern = "/"
	}

	mux := http.NewServeMux()
	mux.Handle(c.Pattern, c.Handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	req, err := http.NewRequest("", server.URL+c.Pattern, bytes.NewBuffer(c.ReqBody))
	require.NoError(t, err)

	for key, values := range c.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := server.Client().Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	gotResp, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return gotResp, resp.StatusCode
}
