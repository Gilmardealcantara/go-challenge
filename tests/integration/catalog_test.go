//go:build integration

package integration

import (
	"testing"

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
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), 200, recorder.Code)
	assert.Equal(s.T(), len(response.Products), 8)
	assert.Equal(s.T(), response.Total, int64(8))
}

func (s *CatalogTestSuite) TestGetCatalog_ReturnsProductsWithCategoryInformation() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), 200, recorder.Code)

	for _, p := range response.Products {
		assert.NotNil(s.T(), p.Category)
		assert.NotEmpty(s.T(), p.Category.Code)
		assert.NotEmpty(s.T(), p.Category.Name)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_RespectsOffsetAndLimitParameters() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog?offset=2&limit=2", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), 200, recorder.Code)
	assert.LessOrEqual(s.T(), len(response.Products), 2)
	assert.Equal(s.T(), "PROD003", response.Products[0].Code)
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByCategory() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog?category=clothing", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), 200, recorder.Code)

	for _, p := range response.Products {
		assert.NotNil(s.T(), p.Category)
		assert.Equal(s.T(), "clothing", p.Category.Code)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByPrice() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog?priceLessThan=10", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), 200, recorder.Code)
	assert.Len(s.T(), response.Products, 3)
	for _, p := range response.Products {
		assert.Less(s.T(), p.Price, 10.0)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByCategoryAndPrice() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog?category=shoes&priceLessThan=10", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	assert.Equal(s.T(), 200, recorder.Code)
	assert.Len(s.T(), response.Products, 1)
	for _, p := range response.Products {
		assert.Less(s.T(), p.Price, 10.0)
		assert.NotNil(s.T(), p.Category)
		assert.Equal(s.T(), "shoes", p.Category.Code)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_ReturnsCorrectTotalCount() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog?limit=1", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	totalCount := response.Total
	assert.Equal(s.T(), totalCount, int64(8))
}

func (s *CatalogTestSuite) TestGetProductByCode_ReturnsProductWithCategory() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog/PROD002", nil)
	response := helpers.DecodeProductDetailsResponse(s.T(), recorder)

	assert.Equal(s.T(), 200, recorder.Code)
	assert.Equal(s.T(), "PROD002", response.Code)
	assert.NotNil(s.T(), response.Category)
	assert.Equal(s.T(), "shoes", response.Category.Code)
	assert.Equal(s.T(), "Shoes", response.Category.Name)
}

func (s *CatalogTestSuite) TestGetProductByCode_VariantsInheritProductPrice() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog/PROD001", nil)
	response := helpers.DecodeProductDetailsResponse(s.T(), recorder)

	assert.Equal(s.T(), 200, recorder.Code)
	assert.NotEmpty(s.T(), response.Variants)

	// Verify that variants have prices (either their own or inherited from product)
	for _, variant := range response.Variants {
		assert.Greater(s.T(), variant.Price, 0.0)
	}
}

func (s *CatalogTestSuite) TestGetProductByCode_ReturnsNotFoundForInvalidCode() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog/INVALID_CODE", nil)
	errorResponse := helpers.DecodeErrorResponse(s.T(), recorder)

	assert.Equal(s.T(), 404, recorder.Code)
	assert.Equal(s.T(), "product not found", errorResponse.Error)
}

func (s *CatalogTestSuite) TestGetProductByCode_ReturnsAllProductVariants() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog/PROD001", nil)
	response := helpers.DecodeProductDetailsResponse(s.T(), recorder)

	assert.Equal(s.T(), 200, recorder.Code)
	assert.NotEmpty(s.T(), response.Variants)

	// Verify variant structure
	for _, variant := range response.Variants {
		assert.NotEmpty(s.T(), variant.Name)
		assert.NotEmpty(s.T(), variant.SKU)
		assert.Greater(s.T(), variant.Price, 0.0)
	}
}
