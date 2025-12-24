//go:build integration

package integration

import (
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/products"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RepositoryTestSuite struct {
	IntegrationTestSuite
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

func (s *RepositoryTestSuite) TestGetAll_ReturnsAllProductsFromDatabase() {
	repo := products.NewRepository(s.DB)

	products, err := repo.GetAll()

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), len(products), 8)

	hasVariants := false
	for _, p := range products {
		assert.NotEmpty(s.T(), p.Code)
		assert.True(s.T(), p.Price.GreaterThan(decimal.Zero))
		assert.NotNil(s.T(), p.Category, "product "+p.Code+" have no category")
		if len(p.Variants) > 0 {
			hasVariants = true
		}
	}
	assert.True(s.T(), hasVariants, "expected at least one product with variants")
}

func (s *RepositoryTestSuite) TestGetWithFilters_FiltersByCategoryCode() {
	repo := products.NewRepository(s.DB)

	params := products.FilterParams{Offset: 0, Limit: 100, CategoryCode: "clothing", PriceLessThan: nil}
	prods, err := repo.GetWithFilters(params)

	assert.NoError(s.T(), err)
	assert.Greater(s.T(), len(prods), 0)

	// Verify all products belong to clothing category
	for _, p := range prods {
		assert.NotNil(s.T(), p.Category)
		assert.Equal(s.T(), "clothing", p.Category.Code)
	}
}

func (s *RepositoryTestSuite) TestGetWithFilters_FiltersByPriceLessThan() {
	repo := products.NewRepository(s.DB)

	price := 10.0
	params := products.FilterParams{Offset: 0, Limit: 100, CategoryCode: "", PriceLessThan: &price}
	prods, err := repo.GetWithFilters(params)

	assert.NoError(s.T(), err)

	assert.Len(s.T(), prods, 3)
	for _, p := range prods {
		assert.Less(s.T(), p.Price.InexactFloat64(), 10.0)
	}
}

func (s *RepositoryTestSuite) TestGetWithFilters_FiltersByBothCategoryAndPrice() {
	repo := products.NewRepository(s.DB)

	price := 100.0
	params := products.FilterParams{Offset: 0, Limit: 100, CategoryCode: "shoes", PriceLessThan: &price}
	prods, err := repo.GetWithFilters(params)

	assert.NoError(s.T(), err)

	for _, p := range prods {
		assert.Less(s.T(), p.Price.InexactFloat64(), 100.0)
		assert.NotNil(s.T(), p.Category)
		assert.Equal(s.T(), "shoes", p.Category.Code)
	}
}

func (s *RepositoryTestSuite) TestGetWithFilters_ReturnsEmptyWhenNoProductsMatchFilters() {
	repo := products.NewRepository(s.DB)

	price := 1.0 // Very low price
	params := products.FilterParams{Offset: 0, Limit: 100, CategoryCode: "", PriceLessThan: &price}
	prods, err := repo.GetWithFilters(params)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), prods)
}

func (s *RepositoryTestSuite) TestGetWithFilters_ReturnsAllProductsWhenNoFiltersApplied() {
	repo := products.NewRepository(s.DB)

	params := products.FilterParams{Offset: 0, Limit: 100, CategoryCode: "", PriceLessThan: nil}
	prods, err := repo.GetWithFilters(params)

	assert.NoError(s.T(), err)
	assert.Greater(s.T(), len(prods), 0)
}

func (s *RepositoryTestSuite) TestGetWithFilters_RespectsPaginationWithFilters() {
	repo := products.NewRepository(s.DB)

	// Get first page
	params1 := products.FilterParams{Offset: 0, Limit: 2, CategoryCode: "clothing", PriceLessThan: nil}
	firstPage, err := repo.GetWithFilters(params1)
	assert.NoError(s.T(), err)

	// Get second page
	params2 := products.FilterParams{Offset: 2, Limit: 2, CategoryCode: "clothing", PriceLessThan: nil}
	secondPage, err := repo.GetWithFilters(params2)
	assert.NoError(s.T(), err)

	// Verify pages are different (if second page has items)
	if len(firstPage) > 0 && len(secondPage) > 0 {
		assert.NotEqual(s.T(), firstPage[0].Code, secondPage[0].Code)
	}
}

func (s *RepositoryTestSuite) TestTotal_ReturnsTotalProductCount() {
	repo := products.NewRepository(s.DB)

	total, err := repo.Total()

	assert.NoError(s.T(), err)
	assert.Greater(s.T(), total, int64(0))
}
