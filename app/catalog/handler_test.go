package catalog_test

import (
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/api/server"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/tests/helpers"
	"github.com/mytheresa/go-hiring-challenge/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandlerHandleGet(t *testing.T) {
	t.Run("returns 200 with catalog response on success", func(t *testing.T) {
		mockService := new(mocks.CatalogService)
		mockService.On("GetProducts", catalog.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(&catalog.Response{Products: []catalog.ProductResponse{}, Total: 0}, nil)

		mux := server.SetupRoutes(catalog.NewHandlerWithService(mockService))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog")

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("returns 400 when limit exceeds 100", func(t *testing.T) {
		mockService := new(mocks.CatalogService)
		mux := server.SetupRoutes(catalog.NewHandlerWithService(mockService))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?limit=200")

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("returns 400 when offset is negative", func(t *testing.T) {
		mockService := new(mocks.CatalogService)
		mux := server.SetupRoutes(catalog.NewHandlerWithService(mockService))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?offset=-1")

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("returns 400 when limit is negative", func(t *testing.T) {
		mockService := new(mocks.CatalogService)
		mux := server.SetupRoutes(catalog.NewHandlerWithService(mockService))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?limit=-5")

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("applies category filter", func(t *testing.T) {
		mockService := new(mocks.CatalogService)
		mockService.On("GetProducts", catalog.FilterParams{Offset: 0, Limit: 10, CategoryCode: "clothing", PriceLessThan: nil}).Return(&catalog.Response{Products: []catalog.ProductResponse{}, Total: 0}, nil)

		mux := server.SetupRoutes(catalog.NewHandlerWithService(mockService))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?category=clothing")

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("applies price filter", func(t *testing.T) {
		price := 50.0
		mockService := new(mocks.CatalogService)
		mockService.On("GetProducts", catalog.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: &price}).Return(&catalog.Response{Products: []catalog.ProductResponse{}, Total: 0}, nil)

		mux := server.SetupRoutes(catalog.NewHandlerWithService(mockService))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?priceLessThan=50")

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}

func TestHandlerHandleGetByCode(t *testing.T) {
	t.Run("returns 200 with product details on success", func(t *testing.T) {
		mockService := new(mocks.CatalogService)
		mockService.On("GetProductByCode", "PROD001").Return(&catalog.ProductDetailsResponse{}, nil)

		mux := server.SetupRoutes(catalog.NewHandlerWithService(mockService))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/PROD001")

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("returns 404 when product not found", func(t *testing.T) {
		mockService := new(mocks.CatalogService)
		mockService.On("GetProductByCode", "INVALID").Return(nil, assert.AnError)

		mux := server.SetupRoutes(catalog.NewHandlerWithService(mockService))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/INVALID")

		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
}
