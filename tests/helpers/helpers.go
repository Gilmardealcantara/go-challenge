package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/stretchr/testify/assert"
)

// MakeRequest makes an HTTP request and records the response
func MakeRequest(t *testing.T, handler http.HandlerFunc, method, url string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, nil)
	recorder := httptest.NewRecorder()
	handler(recorder, req)
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
