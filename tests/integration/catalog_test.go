//go:build integration

package integration

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/tests/helpers"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type CatalogTestSuite struct {
	IntegrationTestSuite
}

func TestCatalogSuite(t *testing.T) {
	suite.Run(t, new(CatalogTestSuite))
}

func (s *CatalogTestSuite) TestGetCatalog_ReturnsAllProductsWithDefaultPagination() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	s.Equal(200, recorder.Code)
	s.Equal(len(response.Products), 8)
	s.Equal(response.Total, int64(8))
	for _, p := range response.Products {
		s.NotZero(p.Price)
		s.NotEmpty(p.Code)
		s.NotNil(p.Category)
		s.NotEmpty(p.Category.Code)
		s.NotEmpty(p.Category.Name)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_RespectsOffsetAndLimitParameters() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog?offset=2&limit=2", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	s.Equal(200, recorder.Code)
	s.LessOrEqual(len(response.Products), 2)
	s.Equal("PROD003", response.Products[0].Code)
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByCategory() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog?category=clothing", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	s.Equal(200, recorder.Code)

	s.LessOrEqual(len(response.Products), 3)
	for _, p := range response.Products {
		s.NotNil(p.Category)
		s.Equal("clothing", p.Category.Code)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByPrice() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog?priceLessThan=10", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	s.Equal(200, recorder.Code)
	s.Len(response.Products, 3)
	for _, p := range response.Products {
		price, _ := decimal.NewFromString(p.Price)
		s.True(price.LessThan(decimal.NewFromInt(10)))
	}
}

func (s *CatalogTestSuite) TestGetCatalog_FiltersProductsByCategoryAndPrice() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog?category=shoes&priceLessThan=10", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	s.Equal(200, recorder.Code)
	s.Len(response.Products, 1)
	for _, p := range response.Products {
		price, _ := decimal.NewFromString(p.Price)
		s.True(price.LessThan(decimal.NewFromInt(10)))
		s.NotNil(p.Category)
		s.Equal("shoes", p.Category.Code)
	}
}

func (s *CatalogTestSuite) TestGetCatalog_ReturnsCorrectTotalCount() {
	recorder := helpers.MakeRequest(s.T(), s.mux, "GET", "/catalog?limit=1", nil)
	response := helpers.DecodeCatalogResponse(s.T(), recorder)

	totalCount := response.Total
	s.Equal(totalCount, int64(8))
}
