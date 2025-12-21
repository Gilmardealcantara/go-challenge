//go:build integration

package integration

import (
	"net/http"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/products"
	"github.com/mytheresa/go-hiring-challenge/tests/helpers"
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

	recorder := helpers.MakeRequest(s.T(), handler.HandleGet, "GET", "/catalog")
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)
	assert.Equal(s.T(), len(response.Products), 8)
	assert.Equal(s.T(), response.Total, int64(8))
}

func (s *CatalogTestSuite) TestGetCatalog_ReturnsProductsWithCategoryInformation() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	recorder := helpers.MakeRequest(s.T(), handler.HandleGet, "GET", "/catalog")
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	for _, p := range response.Products {
		assert.NotNil(s.T(), p.Category)
		assert.NotEmpty(s.T(), p.Category.Code)
		assert.NotEmpty(s.T(), p.Category.Name)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_RespectsOffsetAndLimitParameters() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	recorder := helpers.MakeRequest(s.T(), handler.HandleGet, "GET", "/catalog?offset=2&limit=2")
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)
	assert.LessOrEqual(s.T(), len(response.Products), 2)
	assert.Equal(s.T(), "PROD003", response.Products[0].Code)
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByCategory() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	recorder := helpers.MakeRequest(s.T(), handler.HandleGet, "GET", "/catalog?category=clothing")
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)

	for _, p := range response.Products {
		assert.NotNil(s.T(), p.Category)
		assert.Equal(s.T(), "clothing", p.Category.Code)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByPrice() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	recorder := helpers.MakeRequest(s.T(), handler.HandleGet, "GET", "/catalog?priceLessThan=10")
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)
	assert.Len(s.T(), response.Products, 3)
	for _, p := range response.Products {
		assert.Less(s.T(), p.Price, 10.0)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByCategoryAndPrice() {
	repo := products.NewRepository(s.DB)
	handler := catalog.NewHandler(repo)

	recorder := helpers.MakeRequest(s.T(), handler.HandleGet, "GET", "/catalog?category=shoes&priceLessThan=10")
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), http.StatusOK, recorder.Code)
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

	recorder := helpers.MakeRequest(s.T(), handler.HandleGet, "GET", "/catalog?limit=1")
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	totalCount := response.Total
	assert.Equal(s.T(), totalCount, int64(8))
}
