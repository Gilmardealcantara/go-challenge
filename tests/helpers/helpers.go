package helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/stretchr/testify/assert"
)

// MakeRequest makes an HTTP request with optional body and records the response
func MakeRequest(t *testing.T, mux *http.ServeMux, method, url string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)
	return recorder
}

// DecodeResponse decodes a JSON response from the recorder and returns it
func DecodeResponse[T any](t *testing.T, recorder *httptest.ResponseRecorder) T {
	var response T
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err)
	return response
}

// Convenience functions for backward compatibility
func DecodeCatalogResponse(t *testing.T, recorder *httptest.ResponseRecorder) catalog.Response {
	return DecodeResponse[catalog.Response](t, recorder)
}

func DecodeErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder) api.ErrorDataResponse {
	return DecodeResponse[api.ErrorDataResponse](t, recorder)
}

func DecodeProductDetailsResponse(t *testing.T, recorder *httptest.ResponseRecorder) catalog.ProductDetailsResponse {
	return DecodeResponse[catalog.ProductDetailsResponse](t, recorder)
}

func DecodeCategoriesResponse(t *testing.T, recorder *httptest.ResponseRecorder) []categories.CategoryResponse {
	return DecodeResponse[[]categories.CategoryResponse](t, recorder)
}

func DecodeCategoryResponse(t *testing.T, recorder *httptest.ResponseRecorder) categories.CategoryResponse {
	return DecodeResponse[categories.CategoryResponse](t, recorder)
}
