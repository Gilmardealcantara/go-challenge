package categories_test

import (
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestServiceGetCategories(t *testing.T) {
	t.Run("returns categories with correct mapping", func(t *testing.T) {
		mockCategories := []categories.Category{
			{ID: 1, Code: "clothing", Name: "Clothing"},
			{ID: 2, Code: "shoes", Name: "Shoes"},
		}

		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("GetAll").Return(mockCategories, nil)

		service := categories.NewService(mockRepo)
		response, err := service.GetCategories()

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response, 2)
		assert.Equal(t, "clothing", response[0].Code)
		assert.Equal(t, "Clothing", response[0].Name)
		assert.Equal(t, "shoes", response[1].Code)
		assert.Equal(t, "Shoes", response[1].Name)
	})

	t.Run("returns empty list when no categories", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("GetAll").Return([]categories.Category{}, nil)

		service := categories.NewService(mockRepo)
		response, err := service.GetCategories()

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Len(t, response, 0)
	})

	t.Run("returns error from repository", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("GetAll").Return(nil, errors.New("database error"))

		service := categories.NewService(mockRepo)
		_, err := service.GetCategories()

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}
