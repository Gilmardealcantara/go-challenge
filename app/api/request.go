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

// QueryFloat parses a query string parameter as a float64.
// If the parameter is not provided or invalid, it returns nil.
// If the parameter is provided and valid, it returns a pointer to the parsed float64.
func QueryFloat(r *http.Request, param string) *float64 {
	if paramStr := r.URL.Query().Get(param); paramStr != "" {
		if parsed, err := strconv.ParseFloat(paramStr, 64); err == nil && parsed > 0 {
			return &parsed
		}
	}

	return nil
}
