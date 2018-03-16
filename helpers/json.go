package helpers

import (
	"net/http"
	"strings"
)

// HasJSONRequest check for application/json header.
func HasJSONRequest(r *http.Request) bool {
	if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "application/json") == false {
		return false
	}

	return true
}
