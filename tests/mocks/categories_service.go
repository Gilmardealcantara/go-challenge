package mocks

import (
	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/stretchr/testify/mock"
)

// CategoriesService is a mock implementation of categories.Service using testify/mock.
type CategoriesService struct {
	mock.Mock
}

func (m *CategoriesService) GetCategories() ([]categories.CategoryResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]categories.CategoryResponse), args.Error(1)
}
