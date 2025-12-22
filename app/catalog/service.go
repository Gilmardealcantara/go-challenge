package catalog

import (
	"fmt"

	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/app/variants"
)

type Service interface {
	GetProducts(params FilterParams) (*Response, error)
	GetProductByCode(code string) (*ProductDetailsResponse, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

// GetProducts retrieves products with filters and pagination
func (s *service) GetProducts(params FilterParams) (*Response, error) {

	prods, err := s.repo.GetWithFilters(params)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Total()
	if err != nil {
		return nil, err
	}

	catalogProducts := make([]ProductResponse, len(prods))
	for i, p := range prods {
		catalogProducts[i] = toProductResponse(p)
	}

	return &Response{Products: catalogProducts, Total: total}, nil
}

// GetProductByCode retrieves a product with its variants and category
func (s *service) GetProductByCode(code string) (*ProductDetailsResponse, error) {
	product, err := s.repo.GetByCode(code)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	variantResponses := make([]VariantResponse, len(product.Variants))
	for i, v := range product.Variants {
		variantResponses[i] = toVariantResponse(v, product.Price)
	}

	return &ProductDetailsResponse{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Category: toCategoryResponse(product.Category),
		Variants: variantResponses,
	}, nil
}

// toProductResponse converts a Product domain model to ProductResponse DTO
func toProductResponse(p Product) ProductResponse {
	return ProductResponse{
		Code:     p.Code,
		Price:    p.Price.InexactFloat64(),
		Category: toCategoryResponse(p.Category),
	}
}

// toVariantResponse converts a Variant to VariantResponse, inheriting product price if variant price is zero
func toVariantResponse(v variants.Variant, productPrice interface{}) VariantResponse {
	price := v.Price.InexactFloat64()
	if v.Price.IsZero() {
		// Inherit product price if variant price is null
		if pd, ok := productPrice.(interface{ InexactFloat64() float64 }); ok {
			price = pd.InexactFloat64()
		}
	}
	return VariantResponse{
		Name:  v.Name,
		SKU:   v.SKU,
		Price: price,
	}
}

// toCategoryResponse converts a Category to CategoryResponse DTO, handling nil values
func toCategoryResponse(c *categories.Category) *CategoryResponse {
	if c == nil {
		return nil
	}
	return &CategoryResponse{
		Code: c.Code,
		Name: c.Name,
	}
}
