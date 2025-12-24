package categories_test

import (
	"encoding/json"
	"net/http"
	"testing"

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

		handler := categories.NewHandlerWithService(mockService)
		recorder := helpers.MakeRequest(t, setupCategoriesRoutes(handler), "GET", "/categories", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	})

	t.Run("returns 200 with empty list when no categories", func(t *testing.T) {
		mockService := new(mocks.CategoriesService)
		mockService.On("GetCategories").Return([]categories.CategoryResponse{}, nil)

		handler := categories.NewHandlerWithService(mockService)
		recorder := helpers.MakeRequest(t, setupCategoriesRoutes(handler), "GET", "/categories", nil)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("returns 500 when service returns error", func(t *testing.T) {
		mockService := new(mocks.CategoriesService)
		mockService.On("GetCategories").Return(nil, assert.AnError)

		handler := categories.NewHandlerWithService(mockService)
		recorder := helpers.MakeRequest(t, setupCategoriesRoutes(handler), "GET", "/categories", nil)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func TestHandlerHandleCreate(t *testing.T) {
	t.Run("returns 201 with created category on success", func(t *testing.T) {
		mockService := new(mocks.CategoriesService)
		req := categories.CreateCategoryRequest{Code: "electronics", Name: "Electronics"}
		mockService.On("CreateCategory", req).Return(&categories.CategoryResponse{Code: "electronics", Name: "Electronics"}, nil)

		handler := categories.NewHandlerWithService(mockService)
		body, _ := json.Marshal(req)
		recorder := helpers.MakeRequest(t, setupCategoriesRoutes(handler), "POST", "/categories", body)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

		var response categories.CategoryResponse
		json.NewDecoder(recorder.Body).Decode(&response)
		assert.Equal(t, "electronics", response.Code)
		assert.Equal(t, "Electronics", response.Name)
	})

	t.Run("returns 400 when code is missing", func(t *testing.T) {
		mockService := new(mocks.CategoriesService)
		handler := categories.NewHandlerWithService(mockService)

		req := map[string]string{"name": "Electronics"}
		body, _ := json.Marshal(req)
		recorder := helpers.MakeRequest(t, setupCategoriesRoutes(handler), "POST", "/categories", body)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("returns 400 when name is missing", func(t *testing.T) {
		mockService := new(mocks.CategoriesService)
		handler := categories.NewHandlerWithService(mockService)

		req := map[string]string{"code": "electronics"}
		body, _ := json.Marshal(req)
		recorder := helpers.MakeRequest(t, setupCategoriesRoutes(handler), "POST", "/categories", body)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("returns 400 when request body is invalid", func(t *testing.T) {
		mockService := new(mocks.CategoriesService)
		handler := categories.NewHandlerWithService(mockService)

		recorder := helpers.MakeRequest(t, setupCategoriesRoutes(handler), "POST", "/categories", []byte("invalid json"))

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
	})

	t.Run("returns 500 when service returns error", func(t *testing.T) {
		mockService := new(mocks.CategoriesService)
		req := categories.CreateCategoryRequest{Code: "electronics", Name: "Electronics"}
		mockService.On("CreateCategory", req).Return(nil, assert.AnError)

		handler := categories.NewHandlerWithService(mockService)
		body, _ := json.Marshal(req)
		recorder := helpers.MakeRequest(t, setupCategoriesRoutes(handler), "POST", "/categories", body)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
}

func setupCategoriesRoutes(handler *categories.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories", handler.HandleGet)
	mux.HandleFunc("POST /categories", handler.HandleCreate)
	return mux
}
