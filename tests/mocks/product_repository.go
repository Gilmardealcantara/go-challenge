package mocks

import (
	"github.com/mytheresa/go-hiring-challenge/app/products"
	"github.com/stretchr/testify/mock"
)

// ProductRepository is a mock implementation of catalog.ProductRepository using testify/mock.
type ProductRepository struct {
	mock.Mock
}

func (m *ProductRepository) GetAllProducts() ([]products.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]products.Product), args.Error(1)
}

func (m *ProductRepository) GetProductsWithPagination(offset, limit int) ([]products.Product, int64, error) {
	args := m.Called(offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]products.Product), args.Get(1).(int64), args.Error(2)
}
