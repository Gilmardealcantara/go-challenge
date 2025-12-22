package mocks

import (
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/stretchr/testify/mock"
)

// CatalogService is a mock implementation of catalog.CatalogService using testify/mock.
type CatalogService struct {
	mock.Mock
}

func (m *CatalogService) GetProducts(params catalog.FilterParams) (*catalog.Response, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*catalog.Response), args.Error(1)
}

func (m *CatalogService) GetProductByCode(code string) (*catalog.ProductDetailsResponse, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*catalog.ProductDetailsResponse), args.Error(1)
}
