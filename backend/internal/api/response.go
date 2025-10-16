package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func EncodeErrorf(w http.ResponseWriter, code int, format string, a ...any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Message: fmt.Sprintf(format, a...),
	})
}

func EncodeSuccess[T any](w http.ResponseWriter, code int, response T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_ = json.NewEncoder(w).Encode(response)
}
