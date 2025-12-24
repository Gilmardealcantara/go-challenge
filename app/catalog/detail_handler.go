package catalog

import (
	"net/http"

	"github.com/Gilmardealcantara/go-challenge/app/api"
	"github.com/Gilmardealcantara/go-challenge/app/categories"
	"github.com/Gilmardealcantara/go-challenge/app/products"
	"github.com/Gilmardealcantara/go-challenge/app/variants"
)

// ProductDetailsResponse is the DTO for product details
type ProductDetailsResponse struct {
	Code     string            `json:"code"`
	Price    string            `json:"price"`
	Category *CategoryResponse `json:"category,omitempty"`
	Variants []VariantResponse `json:"variants"`
}

// VariantResponse is the DTO for a product variant
type VariantResponse struct {
	Name  string `json:"name"`
	SKU   string `json:"sku"`
	Price string `json:"price"`
}

type DetailHandler struct {
	repo products.Repository
}

func NewDetailHandler(r products.Repository) *DetailHandler {
	return &DetailHandler{
		repo: r,
	}
}

// Handle godoc
//
//	@Summary		Get product details
//	@Description	Get detailed information about a specific product including all variants
//	@Tags			catalog
//	@Produce		json
//	@Param			code	path		string	true	"Product code"
//	@Success		200		{object}	ProductDetailsResponse
//	@Failure		404		{object}	api.ErrorDataResponse
//	@Router			/catalog/{code} [get]
func (h *DetailHandler) Handle(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	product, err := h.repo.GetByCode(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "product not found")
		return
	}

	// Convert variants to response DTOs
	variantResponses := make([]VariantResponse, len(product.Variants))
	for i, v := range product.Variants {
		variantResponses[i] = h.toVariantResponse(v, product.Price)
	}

	response := &ProductDetailsResponse{
		Code:     product.Code,
		Price:    product.Price.String(),
		Category: h.toCategoryResponse(product.Category),
		Variants: variantResponses,
	}

	api.OKResponse(w, response)
}

// toVariantResponse converts a Variant to VariantResponse, inheriting product price if variant price is zero
func (h *DetailHandler) toVariantResponse(v variants.Variant, productPrice any) VariantResponse {
	price := v.Price.String()
	if v.Price.IsZero() {
		// Inherit product price if variant price is null
		if pd, ok := productPrice.(interface{ String() string }); ok {
			price = pd.String()
		}
	}
	return VariantResponse{
		Name:  v.Name,
		SKU:   v.SKU,
		Price: price,
	}
}

// toCategoryResponse converts a Category to CategoryResponse DTO, handling nil values
func (h *DetailHandler) toCategoryResponse(c *categories.Category) *CategoryResponse {
	if c == nil {
		return nil
	}
	return &CategoryResponse{
		Code: c.Code,
		Name: c.Name,
	}
}
