package mocks

import (
	"github.com/Gilmardealcantara/go-challenge/app/products"
	"github.com/stretchr/testify/mock"
)

// ProductRepository is a mock implementation of products.Repository using testify/mock.
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

func (m *ProductRepository) GetWithFilters(params products.FilterParams) ([]products.Product, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]products.Product), args.Error(1)
}

func (m *ProductRepository) Total() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *ProductRepository) GetByCode(code string) (*products.Product, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*products.Product), args.Error(1)
}
