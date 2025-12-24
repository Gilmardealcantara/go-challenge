package categories_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/tests/helpers"
	"github.com/mytheresa/go-hiring-challenge/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateHandlerHandle(t *testing.T) {
	t.Run("returns 201 with created category on success", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("Create", mock.MatchedBy(func(c *categories.Category) bool {
			return c.Code == "electronics" && c.Name == "Electronics"
		})).Return(nil)

		handler := categories.NewCreateHandler(mockRepo)
		req := categories.CreateCategoryRequest{Code: "electronics", Name: "Electronics"}
		body, _ := json.Marshal(req)
		recorder := helpers.MakeRequest(t, setupCreateRoutes(handler), "POST", "/categories", body)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var response categories.CategoryResponse
		json.NewDecoder(recorder.Body).Decode(&response)
		assert.Equal(t, "electronics", response.Code)
		assert.Equal(t, "Electronics", response.Name)
	})

	t.Run("returns 400 when code is missing", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		handler := categories.NewCreateHandler(mockRepo)

		req := map[string]string{"name": "Electronics"}
		body, _ := json.Marshal(req)
		recorder := helpers.MakeRequest(t, setupCreateRoutes(handler), "POST", "/categories", body)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("returns 400 when name is missing", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		handler := categories.NewCreateHandler(mockRepo)

		req := map[string]string{"code": "electronics"}
		body, _ := json.Marshal(req)
		recorder := helpers.MakeRequest(t, setupCreateRoutes(handler), "POST", "/categories", body)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("returns 400 when request body is invalid", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		handler := categories.NewCreateHandler(mockRepo)

		recorder := helpers.MakeRequest(t, setupCreateRoutes(handler), "POST", "/categories", []byte("invalid json"))

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("returns 500 when repository returns error", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		mockRepo.On("Create", mock.MatchedBy(func(c *categories.Category) bool {
			return c.Code == "electronics" && c.Name == "Electronics"
		})).Return(errors.New("database error"))

		handler := categories.NewCreateHandler(mockRepo)
		req := categories.CreateCategoryRequest{Code: "electronics", Name: "Electronics"}
		body, _ := json.Marshal(req)
		recorder := helpers.MakeRequest(t, setupCreateRoutes(handler), "POST", "/categories", body)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func setupCreateRoutes(handler *categories.CreateHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /categories", handler.Handle)
	return mux
}
