package mocks

import (
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/stretchr/testify/mock"
)

// ProductRepository is a mock implementation of catalog.ProductRepository using testify/mock.
type ProductRepository struct {
	mock.Mock
}

func (m *ProductRepository) GetAll() ([]catalog.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]catalog.Product), args.Error(1)
}

func (m *ProductRepository) GetWithFilters(params catalog.FilterParams) ([]catalog.Product, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]catalog.Product), args.Error(1)
}

func (m *ProductRepository) Total() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *ProductRepository) GetByCode(code string) (*catalog.Product, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*catalog.Product), args.Error(1)
}
