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

func NewSuccessResponse(w http.ResponseWriter, message string) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var ok struct {
		Message string
	}
	ok.Message = message

	return json.NewEncoder(w).Encode(&ok)
}

// IsRequestBodyNil check request body to being not nil
func IsRequestBodyNil(w http.ResponseWriter, r *http.Request) bool {
	// Check request boby
	if r.Body == nil {
		err := ErrorResponse{
			ErrorCode: http.StatusInternalServerError,
			Message:   "Request body cant be empty",
		}
		NewErrorResponse(w, &err)

		return true
	}

	return false
}
