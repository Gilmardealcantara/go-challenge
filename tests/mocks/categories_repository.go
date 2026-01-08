package mocks

import (
	"github.com/Gilmardealcantara/go-challenge/app/categories"
	"github.com/stretchr/testify/mock"
)

// CategoriesRepository is a mock implementation of categories.Repository using testify/mock.
type CategoriesRepository struct {
	mock.Mock
}

func (m *CategoriesRepository) GetAll() ([]categories.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]categories.Category), args.Error(1)
}

func (m *CategoriesRepository) Create(category *categories.Category) error {
	args := m.Called(category)
	return args.Error(0)
}
