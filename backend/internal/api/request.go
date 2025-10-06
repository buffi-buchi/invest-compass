package api

import (
	"encoding/json"
	"net/http"
)

func DecodeRequest[T any](r *http.Request) (T, error) {
	var req T
	return req, json.NewDecoder(r.Body).Decode(&req)
}
