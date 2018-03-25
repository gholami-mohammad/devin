package helpers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message   string              `json:"message"`
	ErrorCode int                 `json:"error_code"`
	Errors    map[string][]string `json:"errors"`
}

func NewErrorResponse(w http.ResponseWriter, err *ErrorResponse) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(err.ErrorCode)

	return json.NewEncoder(w).Encode(&err)
}
