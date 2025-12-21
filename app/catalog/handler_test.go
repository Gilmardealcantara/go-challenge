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

func TestHandlerHandleGet(t *testing.T) {
	t.Run("returns catalog with products on success", func(t *testing.T) {
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
				Price: decimal.NewFromFloat(49.50),
				Category: &categories.Category{
					Code: "shoes",
					Name: "Shoes",
				},
			},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(2), nil)

		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog")
		response := helpers.DecodeCatalogResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Len(t, response.Products, 2)
		assert.Equal(t, int64(2), response.Total)
		assert.Equal(t, "PROD001", response.Products[0].Code)
		assert.Equal(t, 99.99, response.Products[0].Price)
		assert.NotNil(t, response.Products[0].Category)
		assert.Equal(t, "clothing", response.Products[0].Category.Code)
		assert.Equal(t, "Clothing", response.Products[0].Category.Name)
		assert.Equal(t, "PROD002", response.Products[1].Code)
		assert.Equal(t, 49.50, response.Products[1].Price)
		assert.NotNil(t, response.Products[1].Category)
		assert.Equal(t, "shoes", response.Products[1].Category.Code)
		assert.Equal(t, "Shoes", response.Products[1].Category.Name)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns products with null category when category is not set", func(t *testing.T) {
		mockProducts := []products.Product{
			{Code: "PROD001", Price: decimal.NewFromFloat(99.99), Category: nil},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(1), nil)

		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog")
		response := helpers.DecodeCatalogResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, "PROD001", response.Products[0].Code)
		assert.Nil(t, response.Products[0].Category)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns empty products list when no products exist", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return([]products.Product{}, nil)
		mockRepo.On("Total").Return(int64(0), nil)

		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog")
		response := helpers.DecodeCatalogResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Empty(t, response.Products)
		assert.Equal(t, int64(0), response.Total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error response when service fails", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(nil, errors.New("database connection failed"))

		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog")
		errorResponse := helpers.DecodeErrorResponse(t, recorder)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Equal(t, "database connection failed", errorResponse.Error)

		mockRepo.AssertExpectations(t)
	})

	t.Run("applies custom offset and limit parameters", func(t *testing.T) {
		mockProducts := []products.Product{
			{Code: "PROD003", Price: decimal.NewFromFloat(29.99)},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 2, Limit: 1, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(8), nil)

		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog?offset=2&limit=1")
		response := helpers.DecodeCatalogResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(8), response.Total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("enforces maximum limit of 100", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog?limit=200")
		errorResponse := helpers.DecodeErrorResponse(t, recorder)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Equal(t, "limit cannot be greater than 100", errorResponse.Error)
	})

	t.Run("enforces minimum limit of 1", func(t *testing.T) {
		mockProducts := []products.Product{}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 0, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(0), nil)

		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog?limit=0")

		assert.Equal(t, http.StatusOK, recorder.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("defaults to offset 0 and limit 10 when parameters are missing", func(t *testing.T) {
		mockProducts := []products.Product{}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(0), nil)

		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog")

		assert.Equal(t, http.StatusOK, recorder.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when offset is negative", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog?offset=-1")
		errorResponse := helpers.DecodeErrorResponse(t, recorder)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Equal(t, "limit and offset need to be positive values", errorResponse.Error)
	})

	t.Run("returns error when limit is negative", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog?limit=-5")
		errorResponse := helpers.DecodeErrorResponse(t, recorder)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Equal(t, "limit and offset need to be positive values", errorResponse.Error)
	})

	t.Run("filters products by category", func(t *testing.T) {
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

		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog?category=clothing")
		response := helpers.DecodeCatalogResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, "PROD001", response.Products[0].Code)
		assert.Equal(t, "clothing", response.Products[0].Category.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("filters products by price less than", func(t *testing.T) {
		mockProducts := []products.Product{
			{Code: "PROD002", Price: decimal.NewFromFloat(49.50)},
			{Code: "PROD003", Price: decimal.NewFromFloat(29.99)},
		}

		price := 50.0
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: &price}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(2), nil)

		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog?priceLessThan=50")
		response := helpers.DecodeCatalogResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Len(t, response.Products, 2)
		assert.Equal(t, int64(2), response.Total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("filters products by category and price", func(t *testing.T) {
		mockProducts := []products.Product{
			{
				Code:  "PROD002",
				Price: decimal.NewFromFloat(49.50),
				Category: &categories.Category{
					Code: "shoes",
					Name: "Shoes",
				},
			},
		}

		price := 50.0
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "shoes", PriceLessThan: &price}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(1), nil)

		handler := catalog.NewHandler(mockRepo)
		recorder := helpers.MakeRequest(t, handler.HandleGet, "GET", "/catalog?category=shoes&priceLessThan=50")
		response := helpers.DecodeCatalogResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, "PROD002", response.Products[0].Code)
		assert.Equal(t, "shoes", response.Products[0].Category.Code)

		mockRepo.AssertExpectations(t)
	})
}
