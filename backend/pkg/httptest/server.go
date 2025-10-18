package httptest

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

//TODO: Make struct for customize HTTP method, pattern and headers in request.

func DoRequest(
	t *testing.T,
	handler http.HandlerFunc,
	reqBody []byte,
) ([]byte, int) {
	mux := http.NewServeMux()
	mux.Handle("/", handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	req, err := http.NewRequest("", server.URL, bytes.NewBuffer(reqBody))
	require.NoError(t, err)

	resp, err := server.Client().Do(req)
	require.NoError(t, err)

	defer resp.Body.Close()

	gotResp, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return gotResp, resp.StatusCode
}
