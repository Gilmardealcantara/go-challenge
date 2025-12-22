package catalog_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/api/server"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/app/products"
	"github.com/mytheresa/go-hiring-challenge/app/variants"
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

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog")
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

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog")
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

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog")
		response := helpers.DecodeCatalogResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Empty(t, response.Products)
		assert.Equal(t, int64(0), response.Total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error response when service fails", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(nil, errors.New("database connection failed"))

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog")
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

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?offset=2&limit=1")
		response := helpers.DecodeCatalogResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(8), response.Total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("enforces maximum limit of 100", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?limit=200")
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

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?limit=0")

		assert.Equal(t, http.StatusOK, recorder.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("defaults to offset 0 and limit 10 when parameters are missing", func(t *testing.T) {
		mockProducts := []products.Product{}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", products.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(0), nil)

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog")

		assert.Equal(t, http.StatusOK, recorder.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns error when offset is negative", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?offset=-1")
		errorResponse := helpers.DecodeErrorResponse(t, recorder)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Equal(t, "limit and offset need to be positive values", errorResponse.Error)
	})

	t.Run("returns error when limit is negative", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?limit=-5")
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

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?category=clothing")
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

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?priceLessThan=50")
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

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog?category=shoes&priceLessThan=50")
		response := helpers.DecodeCatalogResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, "PROD002", response.Products[0].Code)
		assert.Equal(t, "shoes", response.Products[0].Category.Code)

		mockRepo.AssertExpectations(t)
	})
}

func TestHandlerHandleGetByCode(t *testing.T) {
	t.Run("returns product details with variants on success", func(t *testing.T) {
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

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/PROD001")
		response := helpers.DecodeProductDetailsResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Equal(t, "PROD001", response.Code)
		assert.Equal(t, 99.99, response.Price)
		assert.NotNil(t, response.Category)
		assert.Equal(t, "clothing", response.Category.Code)
		assert.Equal(t, "Clothing", response.Category.Name)
		assert.Len(t, response.Variants, 2)
		assert.Equal(t, "Size M", response.Variants[0].Name)
		assert.Equal(t, "PROD001-M", response.Variants[0].SKU)
		assert.Equal(t, 99.99, response.Variants[0].Price)
		assert.Equal(t, "Size L", response.Variants[1].Name)
		assert.Equal(t, "PROD001-L", response.Variants[1].SKU)
		assert.Equal(t, 109.99, response.Variants[1].Price)

		mockRepo.AssertExpectations(t)
	})

	t.Run("inherits product price for variants without specific price", func(t *testing.T) {
		mockProduct := &products.Product{
			Code:  "PROD002",
			Price: decimal.NewFromFloat(49.99),
			Category: &categories.Category{
				Code: "shoes",
				Name: "Shoes",
			},
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

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/PROD002")
		response := helpers.DecodeProductDetailsResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Len(t, response.Variants, 2)
		assert.Equal(t, 49.99, response.Variants[0].Price)
		assert.Equal(t, 59.99, response.Variants[1].Price)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns product with no category when category is not set", func(t *testing.T) {
		mockProduct := &products.Product{
			Code:     "PROD003",
			Price:    decimal.NewFromFloat(29.99),
			Category: nil,
			Variants: []variants.Variant{
				{
					Name:  "Default",
					SKU:   "PROD003-DEFAULT",
					Price: decimal.NewFromFloat(29.99),
				},
			},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "PROD003").Return(mockProduct, nil)

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/PROD003")
		response := helpers.DecodeProductDetailsResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "PROD003", response.Code)
		assert.Nil(t, response.Category)
		assert.Len(t, response.Variants, 1)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns product with no variants", func(t *testing.T) {
		mockProduct := &products.Product{
			Code:  "PROD004",
			Price: decimal.NewFromFloat(19.99),
			Category: &categories.Category{
				Code: "accessories",
				Name: "Accessories",
			},
			Variants: []variants.Variant{},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "PROD004").Return(mockProduct, nil)

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/PROD004")
		response := helpers.DecodeProductDetailsResponse(t, recorder)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "PROD004", response.Code)
		assert.NotNil(t, response.Category)
		assert.Empty(t, response.Variants)

		mockRepo.AssertExpectations(t)
	})

	t.Run("returns 404 when product not found", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "INVALID").Return(nil, errors.New("record not found"))

		mux := server.SetupRoutes(catalog.NewHandler(mockRepo))
		recorder := helpers.MakeRequest(t, mux, "GET", "/catalog/INVALID")
		errorResponse := helpers.DecodeErrorResponse(t, recorder)

		assert.Equal(t, http.StatusNotFound, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
		assert.Equal(t, "product not found", errorResponse.Error)

		mockRepo.AssertExpectations(t)
	})
}
