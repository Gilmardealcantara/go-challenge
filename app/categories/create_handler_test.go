package categories_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/Gilmardealcantara/go-challenge/app/categories"
	"github.com/Gilmardealcantara/go-challenge/tests/helpers"
	"github.com/Gilmardealcantara/go-challenge/tests/mocks"
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
		response := helpers.DecodeCategoryResponse(t, recorder)
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
		errorResponse := helpers.DecodeErrorResponse(t, recorder)
		assert.Equal(t, "code and name are required", errorResponse.Error)
	})

	t.Run("returns 400 when name is missing", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		handler := categories.NewCreateHandler(mockRepo)

		req := map[string]string{"code": "electronics"}
		body, _ := json.Marshal(req)
		recorder := helpers.MakeRequest(t, setupCreateRoutes(handler), "POST", "/categories", body)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		errorResponse := helpers.DecodeErrorResponse(t, recorder)
		assert.Equal(t, "code and name are required", errorResponse.Error)
	})

	t.Run("returns 400 when request body is invalid", func(t *testing.T) {
		mockRepo := new(mocks.CategoriesRepository)
		handler := categories.NewCreateHandler(mockRepo)

		recorder := helpers.MakeRequest(t, setupCreateRoutes(handler), "POST", "/categories", []byte("invalid json"))

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		errorResponse := helpers.DecodeErrorResponse(t, recorder)
		assert.Equal(t, "invalid request body", errorResponse.Error)
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
		errorResponse := helpers.DecodeErrorResponse(t, recorder)
		assert.Equal(t, "database error", errorResponse.Error)
	})
}

func setupCreateRoutes(handler *categories.CreateHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /categories", handler.Handle)
	return mux
}
