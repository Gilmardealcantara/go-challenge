package api

import (
	"net/http"
	"strconv"
)

// QueryInt parses a query string parameter as an integer.
// If the parameter is not provided or invalid, it returns the default value.
func QueryInt(r *http.Request, param string, defaultValue int) int {
	value := defaultValue

	if paramStr := r.URL.Query().Get(param); paramStr != "" {
		if parsed, err := strconv.Atoi(paramStr); err == nil {
			value = parsed
		}
	}

	return value
}
