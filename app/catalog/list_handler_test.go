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
	t.Run("returns 200 with empty catalog on success", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return([]products.Product{}, nil)
		mockRepo.On("Total").Return(int64(0), nil)

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		response := helpers.DecodeCatalogResponse(t, recorder)
		assert.Empty(t, response.Products)
		assert.Equal(t, int64(0), response.Total)
	})

	t.Run("returns 400 when limit exceeds 100", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?limit=200", nil)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		errorResponse := helpers.DecodeErrorResponse(t, recorder)
		assert.Equal(t, "limit cannot be greater than 100", errorResponse.Error)
	})

	t.Run("returns 400 when offset is negative", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?offset=-1", nil)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		errorResponse := helpers.DecodeErrorResponse(t, recorder)
		assert.Equal(t, "limit and offset need to be positive values", errorResponse.Error)
	})

	t.Run("returns 400 when limit is negative", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?limit=-5", nil)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		errorResponse := helpers.DecodeErrorResponse(t, recorder)
		assert.Equal(t, "limit and offset need to be positive values", errorResponse.Error)
	})

	t.Run("applies category filter", func(t *testing.T) {
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
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "clothing", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(1), nil)

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?category=clothing", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		response := helpers.DecodeCatalogResponse(t, recorder)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, "PROD001", response.Products[0].Code)
		assert.Equal(t, "clothing", response.Products[0].Category.Code)
	})

	t.Run("applies price filter", func(t *testing.T) {
		price := 50.0
		mockProducts := []products.Product{
			{
				Code:  "PROD001",
				Price: decimal.NewFromFloat(49.99),
				Category: &categories.Category{
					Code: "shoes",
					Name: "Shoes",
				},
			},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: &price}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(1), nil)

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?priceLessThan=50", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		response := helpers.DecodeCatalogResponse(t, recorder)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, 49.99, response.Products[0].Price)
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
			{
				Code:  "PROD002",
				Price: decimal.NewFromFloat(49.99),
				Category: &categories.Category{
					Code: "shoes",
					Name: "Shoes",
				},
			},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(2), nil)

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		response := helpers.DecodeCatalogResponse(t, recorder)
		assert.Len(t, response.Products, 2)
		assert.Equal(t, int64(2), response.Total)
		assert.Equal(t, "PROD001", response.Products[0].Code)
		assert.Equal(t, 99.99, response.Products[0].Price)
		assert.Equal(t, "PROD002", response.Products[1].Code)
		assert.Equal(t, 49.99, response.Products[1].Price)
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
		response := helpers.DecodeCatalogResponse(t, recorder)
		assert.Len(t, response.Products, 1)
		assert.Nil(t, response.Products[0].Category)
	})

	t.Run("applies pagination offset and limit", func(t *testing.T) {
		mockProducts := []products.Product{
			{Code: "PROD003", Price: decimal.NewFromFloat(29.99), Category: nil},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 20, Limit: 5, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(100), nil)

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?offset=20&limit=5", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		response := helpers.DecodeCatalogResponse(t, recorder)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(100), response.Total)
	})

	t.Run("returns error from repository", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(nil, errors.New("database error"))

		mux := setupListRoutes(catalog.NewListHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog", nil)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		errorResponse := helpers.DecodeErrorResponse(t, recorder)
		assert.Equal(t, "database error", errorResponse.Error)
	})
}

func setupListRoutes(handler *catalog.ListHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", handler.Handle)
	return mux
}
