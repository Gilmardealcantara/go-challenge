//go:build integration

package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/products"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CatalogTestSuite struct {
	IntegrationTestSuite
}

func TestCatalogSuite(t *testing.T) {
	suite.Run(t, new(CatalogTestSuite))
}

func (s *CatalogTestSuite) TestGetCatalog_ReturnsAllProductsWithPagination() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	req := httptest.NewRequest("GET", "/catalog", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	var response catalog.Response
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), len(response.Products), 8)
	assert.Equal(s.T(), response.Total, int64(8))
}

func (s *CatalogTestSuite) TestGetCatalog_ReturnsProductsWithCategoryInformation() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	req := httptest.NewRequest("GET", "/catalog", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	var response catalog.Response
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(s.T(), err)

	for _, p := range response.Products {
		assert.NotNil(s.T(), p.Category)
		assert.NotEmpty(s.T(), p.Category.Code)
		assert.NotEmpty(s.T(), p.Category.Name)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_RespectsOffsetAndLimitParameters() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	req := httptest.NewRequest("GET", "/catalog?offset=2&limit=2", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	var response catalog.Response
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(s.T(), err)
	assert.LessOrEqual(s.T(), len(response.Products), 2)
	assert.Equal(s.T(), "PROD003", response.Products[0].Code)
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByCategory() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	req := httptest.NewRequest("GET", "/catalog?category=clothing", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	var response catalog.Response
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(s.T(), err)

	for _, p := range response.Products {
		assert.NotNil(s.T(), p.Category)
		assert.Equal(s.T(), "clothing", p.Category.Code)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByPrice() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	req := httptest.NewRequest("GET", "/catalog?priceLessThan=10", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	var response catalog.Response
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(s.T(), err)

	assert.Len(s.T(), response.Products, 3)
	for _, p := range response.Products {
		assert.Less(s.T(), p.Price, 10.0)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByCategoryAndPrice() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	req := httptest.NewRequest("GET", "/catalog?category=shoes&priceLessThan=10", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	var response catalog.Response
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(s.T(), err)

	assert.Len(s.T(), response.Products, 1)
	for _, p := range response.Products {
		assert.Less(s.T(), p.Price, 10.0)
		assert.NotNil(s.T(), p.Category)
		assert.Equal(s.T(), "shoes", p.Category.Code)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_ReturnsCorrectTotalCount() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	req := httptest.NewRequest("GET", "/catalog?limit=1", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	var response catalog.Response
	err := json.NewDecoder(recorder.Body).Decode(&response)
	assert.NoError(s.T(), err)

	totalCount := response.Total
	assert.Equal(s.T(), totalCount, int64(8))
}
