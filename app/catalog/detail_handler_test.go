package catalog_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/app/products"
	"github.com/mytheresa/go-hiring-challenge/app/variants"
	"github.com/mytheresa/go-hiring-challenge/tests/helpers"
	"github.com/mytheresa/go-hiring-challenge/tests/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestDetailHandlerHandle(t *testing.T) {
	t.Run("returns 200 with product details on success", func(t *testing.T) {
		mockProduct := &products.Product{
			Code:  "PROD001",
			Price: decimal.NewFromFloat(99.99),
			Category: &categories.Category{
				Code: "clothing",
				Name: "Clothing",
			},
			Variants: []variants.Variant{},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "PROD001").Return(mockProduct, nil)

		mux := setupDetailRoutes(catalog.NewDetailHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/PROD001", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		response := helpers.DecodeProductDetailsResponse(t, recorder)
		assert.Equal(t, "PROD001", response.Code)
		assert.Equal(t, 99.99, response.Price)
		assert.NotNil(t, response.Category)
		assert.Equal(t, "clothing", response.Category.Code)
		assert.Equal(t, "Clothing", response.Category.Name)
		assert.Empty(t, response.Variants)
	})

	t.Run("returns 404 when product not found", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "INVALID").Return(nil, errors.New("record not found"))

		mux := setupDetailRoutes(catalog.NewDetailHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/INVALID", nil)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		errorResponse := helpers.DecodeErrorResponse(t, recorder)
		assert.Equal(t, "product not found", errorResponse.Error)
	})

	t.Run("returns product details with variants", func(t *testing.T) {
		mockProduct := &products.Product{
			Code:  "PROD001",
			Price: decimal.NewFromFloat(99.99),
			Category: &categories.Category{
				Code: "clothing",
				Name: "Clothing",
			},
			Variants: []variants.Variant{
				{
					Name:  "Size M",
					SKU:   "PROD001-M",
					Price: decimal.NewFromFloat(99.99),
				},
				{
					Name:  "Size L",
					SKU:   "PROD001-L",
					Price: decimal.NewFromFloat(109.99),
				},
			},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "PROD001").Return(mockProduct, nil)

		mux := setupDetailRoutes(catalog.NewDetailHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/PROD001", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		response := helpers.DecodeProductDetailsResponse(t, recorder)
		assert.Equal(t, "PROD001", response.Code)
		assert.Len(t, response.Variants, 2)
		assert.Equal(t, "Size M", response.Variants[0].Name)
		assert.Equal(t, "PROD001-M", response.Variants[0].SKU)
		assert.Equal(t, 99.99, response.Variants[0].Price)
		assert.Equal(t, "Size L", response.Variants[1].Name)
		assert.Equal(t, "PROD001-L", response.Variants[1].SKU)
		assert.Equal(t, 109.99, response.Variants[1].Price)
	})

	t.Run("inherits product price for variants with zero price", func(t *testing.T) {
		mockProduct := &products.Product{
			Code:  "PROD002",
			Price: decimal.NewFromFloat(49.99),
			Variants: []variants.Variant{
				{
					Name:  "Size 10",
					SKU:   "PROD002-10",
					Price: decimal.Zero,
				},
				{
					Name:  "Size 11",
					SKU:   "PROD002-11",
					Price: decimal.NewFromFloat(59.99),
				},
			},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "PROD002").Return(mockProduct, nil)

		mux := setupDetailRoutes(catalog.NewDetailHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/PROD002", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		response := helpers.DecodeProductDetailsResponse(t, recorder)
		assert.Equal(t, "PROD002", response.Code)
		assert.Equal(t, 49.99, response.Price)
		assert.Len(t, response.Variants, 2)
		assert.Equal(t, 49.99, response.Variants[0].Price, "variant with zero price should inherit product price")
		assert.Equal(t, 59.99, response.Variants[1].Price, "variant with explicit price should keep its own price")
	})
}

func setupDetailRoutes(handler *catalog.DetailHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog/{code}", handler.Handle)
	return mux
}
