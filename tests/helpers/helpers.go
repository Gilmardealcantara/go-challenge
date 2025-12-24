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

// DecodeCatalogResponse decodes a catalog response from the recorder and returns it
func DecodeCatalogResponse(t *testing.T, recorder *httptest.ResponseRecorder) catalog.Response {
	var response catalog.Response
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err)
	return response
}

// DecodeErrorResponse decodes a JSON error response
func DecodeErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder) api.ErrorDataResponse {
	var errorResponse api.ErrorDataResponse
	err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
	assert.NoError(t, err)
	return errorResponse
}

// DecodeProductDetailsResponse decodes a product details response from the recorder and returns it
func DecodeProductDetailsResponse(t *testing.T, recorder *httptest.ResponseRecorder) catalog.ProductDetailsResponse {
	var response catalog.ProductDetailsResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err)
	return response
}

// DecodeCategoriesResponse decodes a categories response from the recorder and returns it
func DecodeCategoriesResponse(t *testing.T, recorder *httptest.ResponseRecorder) []categories.CategoryResponse {
	var response []categories.CategoryResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err)
	return response
}

// DecodeCategoryResponse decodes a category response from the recorder and returns it
func DecodeCategoryResponse(t *testing.T, recorder *httptest.ResponseRecorder) categories.CategoryResponse {
	var response categories.CategoryResponse
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(t, err)
	return response
}
