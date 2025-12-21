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
		mockRepo.On("GetAllProducts").Return(mockProducts, nil)

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
		mockRepo.On("GetAllProducts").Return(mockProducts, nil)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, "PROD001", response.Products[0].Code)
		assert.Nil(t, response.Products[0].Category)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns empty products list when no products exist", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetAllProducts").Return([]products.Product{}, nil)

		handler := NewHandler(mockRepo)

		req := httptest.NewRequest("GET", "/catalog", nil)
		recorder := httptest.NewRecorder()

		handler.HandleGet(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		var response Response
		err := json.NewDecoder(recorder.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Empty(t, response.Products)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error response when service fails", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetAllProducts").Return(nil, errors.New("database connection failed"))

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
}
