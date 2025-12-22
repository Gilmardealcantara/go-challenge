package catalog_test

import (
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/app/variants"
	"github.com/mytheresa/go-hiring-challenge/tests/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestServiceGetProducts(t *testing.T) {
	t.Run("returns products with correct mapping", func(t *testing.T) {
		mockProducts := []catalog.Product{
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
		mockRepo.On("GetWithFilters", catalog.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(1), nil)

		service := catalog.NewService(mockRepo)
		response, err := service.GetProducts(catalog.FilterParams{Offset: 0, Limit: 10})

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response.Products, 1)
		assert.Equal(t, int64(1), response.Total)
		assert.Equal(t, "PROD001", response.Products[0].Code)
		assert.Equal(t, 99.99, response.Products[0].Price)
		assert.NotNil(t, response.Products[0].Category)
		assert.Equal(t, "clothing", response.Products[0].Category.Code)
	})

	t.Run("handles nil category", func(t *testing.T) {
		mockProducts := []catalog.Product{
			{Code: "PROD001", Price: decimal.NewFromFloat(99.99), Category: nil},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", catalog.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(1), nil)

		service := catalog.NewService(mockRepo)
		response, err := service.GetProducts(catalog.FilterParams{Offset: 0, Limit: 10})

		assert.NoError(t, err)
		assert.Nil(t, response.Products[0].Category)
	})

	t.Run("returns error from repository", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", catalog.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: nil}).Return(nil, errors.New("database error"))

		service := catalog.NewService(mockRepo)
		_, err := service.GetProducts(catalog.FilterParams{Offset: 0, Limit: 10})

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})

	t.Run("applies category filter", func(t *testing.T) {
		mockProducts := []catalog.Product{
			{Code: "PROD001", Price: decimal.NewFromFloat(99.99)},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", catalog.FilterParams{Offset: 0, Limit: 10, CategoryCode: "clothing", PriceLessThan: nil}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(1), nil)

		service := catalog.NewService(mockRepo)
		response, err := service.GetProducts(catalog.FilterParams{Offset: 0, Limit: 10, CategoryCode: "clothing"})

		assert.NoError(t, err)
		assert.Len(t, response.Products, 1)
	})

	t.Run("applies price filter", func(t *testing.T) {
		price := 50.0
		mockProducts := []catalog.Product{
			{Code: "PROD001", Price: decimal.NewFromFloat(49.99)},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetWithFilters", catalog.FilterParams{Offset: 0, Limit: 10, CategoryCode: "", PriceLessThan: &price}).Return(mockProducts, nil)
		mockRepo.On("Total").Return(int64(1), nil)

		service := catalog.NewService(mockRepo)
		response, err := service.GetProducts(catalog.FilterParams{Offset: 0, Limit: 10, PriceLessThan: &price})

		assert.NoError(t, err)
		assert.Len(t, response.Products, 1)
	})
}

func TestServiceGetProductByCode(t *testing.T) {
	t.Run("returns product details with variants", func(t *testing.T) {
		mockProduct := &catalog.Product{
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

		service := catalog.NewService(mockRepo)
		response, err := service.GetProductByCode("PROD001")

		assert.NoError(t, err)
		assert.Equal(t, "PROD001", response.Code)
		assert.Equal(t, 99.99, response.Price)
		assert.NotNil(t, response.Category)
		assert.Equal(t, "clothing", response.Category.Code)
		assert.Len(t, response.Variants, 2)
		assert.Equal(t, "Size M", response.Variants[0].Name)
		assert.Equal(t, 99.99, response.Variants[0].Price)
		assert.Equal(t, "Size L", response.Variants[1].Name)
		assert.Equal(t, 109.99, response.Variants[1].Price)
	})

	t.Run("inherits product price for variants with zero price", func(t *testing.T) {
		mockProduct := &catalog.Product{
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

		service := catalog.NewService(mockRepo)
		response, err := service.GetProductByCode("PROD002")

		assert.NoError(t, err)
		assert.Equal(t, 49.99, response.Variants[0].Price)
		assert.Equal(t, 59.99, response.Variants[1].Price)
	})

	t.Run("handles nil category", func(t *testing.T) {
		mockProduct := &catalog.Product{
			Code:     "PROD003",
			Price:    decimal.NewFromFloat(29.99),
			Category: nil,
			Variants: []variants.Variant{},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "PROD003").Return(mockProduct, nil)

		service := catalog.NewService(mockRepo)
		response, err := service.GetProductByCode("PROD003")

		assert.NoError(t, err)
		assert.Nil(t, response.Category)
	})

	t.Run("handles empty variants", func(t *testing.T) {
		mockProduct := &catalog.Product{
			Code:     "PROD004",
			Price:    decimal.NewFromFloat(19.99),
			Variants: []variants.Variant{},
		}

		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "PROD004").Return(mockProduct, nil)

		service := catalog.NewService(mockRepo)
		response, err := service.GetProductByCode("PROD004")

		assert.NoError(t, err)
		assert.Empty(t, response.Variants)
	})

	t.Run("returns error when product not found", func(t *testing.T) {
		mockRepo := new(mocks.ProductRepository)
		mockRepo.On("GetByCode", "INVALID").Return(nil, errors.New("record not found"))

		service := catalog.NewService(mockRepo)
		_, err := service.GetProductByCode("INVALID")

		assert.Error(t, err)
		assert.Equal(t, "product not found", err.Error())
	})
}
