package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/app/products"
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
		mockRepo.On("GetProductsWithFilters", 0, 10, "", (*float64)(nil)).Return(mockProducts, int64(2), nil)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var response Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
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
		mockRepo.On("GetProductsWithFilters", 0, 10, "", (*float64)(nil)).Return(mockProducts, int64(1), nil)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, "PROD001", response.Products[0].Code)
		assert.Nil(t, response.Products[0].Category)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns empty products list when no products exist", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetProductsWithFilters", 0, 10, "", (*float64)(nil)).Return([]products.Product{}, int64(0), nil)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Empty(t, response.Products)
		assert.Equal(t, int64(0), response.Total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error response when service fails", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetProductsWithFilters", 0, 10, "", (*float64)(nil)).Return(nil, int64(0), errors.New("database connection failed"))

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var errorResponse map[string]string
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		assert.NoError(t, err)
		assert.Equal(t, "database connection failed", errorResponse["error"])

		mockRepo.AssertExpectations(t)
	})

	t.Run("applies custom offset and limit parameters", func(t *testing.T) {
		mockProducts := []products.Product{
			{Code: "PROD003", Price: decimal.NewFromFloat(29.99)},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetProductsWithFilters", 2, 1, "", (*float64)(nil)).Return(mockProducts, int64(8), nil)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog?offset=2&limit=1", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(8), response.Total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("enforces maximum limit of 100", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog?limit=200", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var errorResponse map[string]string
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		assert.NoError(t, err)
		assert.Equal(t, "limit cannot be greater than 100", errorResponse["error"])
	})

	t.Run("enforces minimum limit of 1", func(t *testing.T) {
		mockProducts := []products.Product{}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetProductsWithFilters", 0, 0, "", (*float64)(nil)).Return(mockProducts, int64(0), nil)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog?limit=0", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("defaults to offset 0 and limit 10 when parameters are missing", func(t *testing.T) {
		mockProducts := []products.Product{}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetProductsWithFilters", 0, 10, "", (*float64)(nil)).Return(mockProducts, int64(0), nil)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when offset is negative", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog?offset=-1", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var errorResponse map[string]string
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		assert.NoError(t, err)
		assert.Equal(t, "limit and offset need to be positive values", errorResponse["error"])
	})

	t.Run("returns error when limit is negative", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog?limit=-5", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var errorResponse map[string]string
		err := json.NewDecoder(recorder.Body).Decode(&errorResponse)
		assert.NoError(t, err)
		assert.Equal(t, "limit and offset need to be positive values", errorResponse["error"])
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
		mockRepo.On("GetProductsWithFilters", 0, 10, "clothing", (*float64)(nil)).Return(mockProducts, int64(1), nil)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog?category=clothing", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
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
		mockRepo.On("GetProductsWithFilters", 0, 10, "", &price).Return(mockProducts, int64(2), nil)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog?priceLessThan=50", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
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
		mockRepo.On("GetProductsWithFilters", 0, 10, "shoes", &price).Return(mockProducts, int64(1), nil)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog?category=shoes&priceLessThan=50", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, "PROD002", response.Products[0].Code)
		assert.Equal(t, "shoes", response.Products[0].Category.Code)

		mockRepo.AssertExpectations(t)
	})
}
