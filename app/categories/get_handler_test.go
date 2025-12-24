package categories_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/tests/helpers"
	"github.com/mytheresa/go-hiring-challenge/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetHandlerHandle(t *testing.T) {
	t.Run("returns 200 with categories on success", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("GetAll").Return([]categories.Category{
			{Code: "clothing", Name: "Clothing"},
			{Code: "shoes", Name: "Shoes"},
		}, nil)

		handler := categories.NewGetHandler(mockRepo)
		recorder := helpers.MakeRequest(t, setupGetRoutes(handler), "GET", "/categories", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("returns 200 with empty list when no categories", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("GetAll").Return([]categories.Category{}, nil)

		handler := categories.NewGetHandler(mockRepo)
		recorder := helpers.MakeRequest(t, setupGetRoutes(handler), "GET", "/categories", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("returns 500 when repository returns error", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("GetAll").Return(nil, errors.New("database error"))

		handler := categories.NewGetHandler(mockRepo)
		recorder := helpers.MakeRequest(t, setupGetRoutes(handler), "GET", "/categories", nil)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func setupGetRoutes(handler *categories.GetHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories", handler.Handle)
	return mux
}
