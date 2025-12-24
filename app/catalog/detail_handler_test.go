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
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("returns 404 when product not found", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "INVALID").Return(nil, errors.New("record not found"))

		mux := setupDetailRoutes(catalog.NewDetailHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/INVALID", nil)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
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
	})

	t.Run("handles nil category", func(t *testing.T) {
		mockProduct := &products.Product{
			Code:     "PROD003",
			Price:    decimal.NewFromFloat(29.99),
			Category: nil,
			Variants: []variants.Variant{},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "PROD003").Return(mockProduct, nil)

		mux := setupDetailRoutes(catalog.NewDetailHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/PROD003", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("handles empty variants", func(t *testing.T) {
		mockProduct := &products.Product{
			Code:     "PROD004",
			Price:    decimal.NewFromFloat(19.99),
			Variants: []variants.Variant{},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "PROD004").Return(mockProduct, nil)

		mux := setupDetailRoutes(catalog.NewDetailHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/PROD004", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}

func setupDetailRoutes(handler *catalog.DetailHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog/{code}", handler.Handle)
	return mux
}
