package catalog_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/app/products"
	"github.com/mytheresa/go-hiring-challenge/tests/helpers"
	"github.com/mytheresa/go-hiring-challenge/tests/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestListHandlerHandle(t *testing.T) {
	t.Run("returns 200 with catalog response on success", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return([]products.Product{}, nil)
		mockRepo.On("Total").Return(int64(0), nil)

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("returns 400 when limit exceeds 100", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?limit=200", nil)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("returns 400 when offset is negative", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?offset=-1", nil)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("returns 400 when limit is negative", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?limit=-5", nil)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("applies category filter", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "clothing", PriceLessThan: nil}).Return([]products.Product{}, nil)
		mockRepo.On("Total").Return(int64(0), nil)

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?category=clothing", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("applies price filter", func(t *testing.T) {
		price := 50.0
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: &price}).Return([]products.Product{}, nil)
		mockRepo.On("Total").Return(int64(0), nil)

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?priceLessThan=50", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("returns products with correct mapping", func(t *testing.T) {
		mockProducts := []products.Product{
			{
				Code:  "PROD001",
				Price: decimal.NewFromFloat(99.99),
				Category: &categories.Category{
					Code: "clothing",
					Name: "Clothing",
				},
			},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(1), nil)

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("handles nil category", func(t *testing.T) {
		mockProducts := []products.Product{
			{Code: "PROD001", Price: decimal.NewFromFloat(99.99), Category: nil},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(1), nil)

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("returns error from repository", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(nil, errors.New("database error"))

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog", nil)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func setupListRoutes(handler *catalog.ListHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", handler.Handle)
	return mux
}
