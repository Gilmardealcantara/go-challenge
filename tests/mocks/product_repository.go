package mocks

import (
	"github.com/mytheresa/go-hiring-challenge/app/products"
	"github.com/stretchr/testify/mock"
)

// ProductRepository is a mock implementation of catalog.ProductRepository using testify/mock.
type ProductRepository struct {
	mock.Mock
}

func (m *ProductRepository) GetAll() ([]products.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]products.Product), args.Error(1)
}

func (m *ProductRepository) GetWithPagination(offset, limit int) ([]products.Product, error) {
	args := m.Called(offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]products.Product), args.Error(1)
}

func (m *ProductRepository) GetWithFilters(offset, limit int, categoryCode string, priceLessThan *float64) ([]products.Product, error) {
	args := m.Called(offset, limit, categoryCode, priceLessThan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]products.Product), args.Error(1)
}

func (m *ProductRepository) Total() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
