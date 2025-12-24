//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/tests/helpers"
	"github.com/stretchr/testify/suite"
)

type CategoriesTestSuite struct {
	IntegrationTestSuite
}

func TestCategoriesSuite(t *testing.T) {
	suite.Run(t, new(CategoriesTestSuite))
}

func (s *CategoriesTestSuite) TestGetCategories_ReturnsAllCategories() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/categories", nil)

	s.Equal(http.StatusOK, recorder.Code)
	s.Equal("application/json", recorder.Header().Get("Content-Type"))

	response := helpers.DecodeCategoriesResponse(s.T(), recorder)
	s.NotEmpty(response)
	s.Len(response, 3)
	for _, category := range response {
		s.NotEmpty(category.Code)
		s.NotEmpty(category.Name)
	}
}

func (s *CategoriesTestSuite) TestCreateCategory_CreatesNewCategorySuccessfully() {
	req := categories.CreateCategoryRequest{Code: "electronics", Name: "Electronics"}
	body, _ := json.Marshal(req)
	recorder := helpers.MakeRequest(s.T(), s.mux, "POST", "/categories", body)

	s.Equal(http.StatusCreated, recorder.Code)
	s.Equal("application/json", recorder.Header().Get("Content-Type"))

	response := helpers.DecodeCategoryResponse(s.T(), recorder)
	s.Equal("electronics", response.Code)
	s.Equal("Electronics", response.Name)
}

func (s *CategoriesTestSuite) TestCreateCategory_ReturnsErrorForDuplicateCode() {
	// Create first category
	req := categories.CreateCategoryRequest{Code: "duplicate", Name: "Duplicate"}
	body, _ := json.Marshal(req)
	recorder1 := helpers.MakeRequest(s.T(), s.mux, "POST", "/categories", body)
	s.Equal(http.StatusCreated, recorder1.Code)

	// Try to create with same code
	recorder2 := helpers.MakeRequest(s.T(), s.mux, "POST", "/categories", body)
	s.Equal(http.StatusInternalServerError, recorder2.Code)
}

func (s *CategoriesTestSuite) TestCreateCategory_ReturnsErrorForInvalidJSON() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "POST", "/categories", []byte("invalid json"))
	s.Equal(http.StatusBadRequest, recorder.Code)
}
