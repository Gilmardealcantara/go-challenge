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
		response := helpers.DecodeCategoriesResponse(t, recorder)
		assert.Len(t, response, 2)
		assert.Equal(t, "clothing", response[0].Code)
		assert.Equal(t, "Clothing", response[0].Name)
		assert.Equal(t, "shoes", response[1].Code)
		assert.Equal(t, "Shoes", response[1].Name)
	})

	t.Run("returns 200 with empty list when no categories", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("GetAll").Return([]categories.Category{}, nil)

		handler := categories.NewGetHandler(mockRepo)
		recorder := helpers.MakeRequest(t, setupGetRoutes(handler), "GET", "/categories", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		response := helpers.DecodeCategoriesResponse(t, recorder)
		assert.Empty(t, response)
	})

	t.Run("returns 500 when repository returns error", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("GetAll").Return(nil, errors.New("database error"))

		handler := categories.NewGetHandler(mockRepo)
		recorder := helpers.MakeRequest(t, setupGetRoutes(handler), "GET", "/categories", nil)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		errorResponse := helpers.DecodeErrorResponse(t, recorder)
		assert.Equal(t, "database error", errorResponse.Error)
	})
}

func setupGetRoutes(handler *categories.GetHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories", handler.Handle)
	return mux
}
