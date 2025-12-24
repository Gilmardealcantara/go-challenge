package categories_test

import (
	"errors"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestServiceCreateCategory(t *testing.T) {
	t.Run("creates category successfully", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("Create", mock.MatchedBy(func(c *categories.Category) bool {
			return c.Code == "electronics" && c.Name == "Electronics"
		})).Return(nil)

		service := categories.NewService(mockRepo)
		req := categories.CreateCategoryRequest{Code: "electronics", Name: "Electronics"}
		response, err := service.CreateCategory(req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, "electronics", response.Code)
		assert.Equal(t, "Electronics", response.Name)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("Create", mock.MatchedBy(func(c *categories.Category) bool {
			return c.Code == "electronics" && c.Name == "Electronics"
		})).Return(errors.New("duplicate key"))

		service := categories.NewService(mockRepo)
		req := categories.CreateCategoryRequest{Code: "electronics", Name: "Electronics"}
		_, err := service.CreateCategory(req)

		assert.Error(t, err)
		assert.Equal(t, "duplicate key", err.Error())
	})
}
