package categories_test

import (
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/api/server"
	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/tests/helpers"
	"github.com/mytheresa/go-hiring-challenge/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandlerHandleGet(t *testing.T) {
	t.Run("returns 200 with categories on success", func(t *testing.T) {
		mockService := new(mocks.CategoriesService)
		mockService.On("GetCategories").Return([]categories.CategoryResponse{
			{Code: "clothing", Name: "Clothing"},
			{Code: "shoes", Name: "Shoes"},
		}, nil)

		handler := categories.NewHandler(mockService)
		recorder := helpers.MakeRequest(t, setupCategoriesRoutes(handler), "GET", "/categories")

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("returns 200 with empty list when no categories", func(t *testing.T) {
		mockService := new(mocks.CategoriesService)
		mockService.On("GetCategories").Return([]categories.CategoryResponse{}, nil)

		handler := categories.NewHandler(mockService)
		recorder := helpers.MakeRequest(t, setupCategoriesRoutes(handler), "GET", "/categories")

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("returns 500 when service returns error", func(t *testing.T) {
		mockService := new(mocks.CategoriesService)
		mockService.On("GetCategories").Return(nil, assert.AnError)

		handler := categories.NewHandler(mockService)
		recorder := helpers.MakeRequest(t, setupCategoriesRoutes(handler), "GET", "/categories")

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func setupCategoriesRoutes(handler *categories.Handler) *http.ServeMux {
	return server.SetupRoutes(nil, handler)
}
