//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/tests/helpers"
	"github.com/stretchr/testify/assert"
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

	assert.Equal(s.T(), http.StatusOK, recorder.Code)
	assert.Equal(s.T(), "application/json", recorder.Header().Get("Content-Type"))

	response := helpers.DecodeCategoriesResponse(s.T(), recorder)
	assert.NotEmpty(s.T(), response)
}

func (s *CategoriesTestSuite) TestGetCategories_ReturnsCorrectCategoryStructure() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/categories", nil)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	response := helpers.DecodeCategoriesResponse(s.T(), recorder)

	for _, category := range response {
		assert.NotEmpty(s.T(), category.Code)
		assert.NotEmpty(s.T(), category.Name)
	}
}

func (s *CategoriesTestSuite) TestGetCategories_ContainsExpectedCategories() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/categories", nil)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	response := helpers.DecodeCategoriesResponse(s.T(), recorder)

	// Verify expected categories exist
	categoryCodes := make(map[string]string)
	for _, cat := range response {
		categoryCodes[cat.Code] = cat.Name
	}

	assert.Contains(s.T(), categoryCodes, "clothing")
	assert.Contains(s.T(), categoryCodes, "shoes")
	assert.Equal(s.T(), "Clothing", categoryCodes["clothing"])
	assert.Equal(s.T(), "Shoes", categoryCodes["shoes"])
}

func (s *CategoriesTestSuite) TestGetCategories_ReturnsConsistentData() {
	// First request
	recorder1 := helpers.MakeRequest(s.T(), s.mux, "GET", "/categories", nil)
	response1 := helpers.DecodeCategoriesResponse(s.T(), recorder1)

	// Second request
	recorder2 := helpers.MakeRequest(s.T(), s.mux, "GET", "/categories", nil)
	response2 := helpers.DecodeCategoriesResponse(s.T(), recorder2)

	assert.Equal(s.T(), len(response1), len(response2))
	assert.Equal(s.T(), response1, response2)
}

func (s *CategoriesTestSuite) TestCreateCategory_CreatesNewCategorySuccessfully() {
	req := categories.CreateCategoryRequest{Code: "electronics", Name: "Electronics"}
	body, _ := json.Marshal(req)
	recorder := helpers.MakeRequest(s.T(), s.mux, "POST", "/categories", body)

	assert.Equal(s.T(), http.StatusCreated, recorder.Code)
	assert.Equal(s.T(), "application/json", recorder.Header().Get("Content-Type"))

	response := helpers.DecodeCategoryResponse(s.T(), recorder)
	assert.Equal(s.T(), "electronics", response.Code)
	assert.Equal(s.T(), "Electronics", response.Name)
}

func (s *CategoriesTestSuite) TestCreateCategory_ReturnsErrorForDuplicateCode() {
	// Create first category
	req := categories.CreateCategoryRequest{Code: "duplicate", Name: "Duplicate"}
	body, _ := json.Marshal(req)
	recorder1 := helpers.MakeRequest(s.T(), s.mux, "POST", "/categories", body)
	assert.Equal(s.T(), http.StatusCreated, recorder1.Code)

	// Try to create with same code
	recorder2 := helpers.MakeRequest(s.T(), s.mux, "POST", "/categories", body)
	assert.Equal(s.T(), http.StatusInternalServerError, recorder2.Code)
}

func (s *CategoriesTestSuite) TestCreateCategory_ReturnsErrorForMissingCode() {
	req := map[string]string{"name": "NoCode"}
	body, _ := json.Marshal(req)
	recorder := helpers.MakeRequest(s.T(), s.mux, "POST", "/categories", body)

	assert.Equal(s.T(), http.StatusBadRequest, recorder.Code)
}

func (s *CategoriesTestSuite) TestCreateCategory_ReturnsErrorForMissingName() {
	req := map[string]string{"code": "noname"}
	body, _ := json.Marshal(req)
	recorder := helpers.MakeRequest(s.T(), s.mux, "POST", "/categories", body)

	assert.Equal(s.T(), http.StatusBadRequest, recorder.Code)
}

func (s *CategoriesTestSuite) TestCreateCategory_ReturnsErrorForInvalidJSON() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "POST", "/categories", []byte("invalid json"))
	assert.Equal(s.T(), http.StatusBadRequest, recorder.Code)
}
