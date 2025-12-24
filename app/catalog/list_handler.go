package catalog

import (
	"net/http"

	"github.com/Gilmardealcantara/go-challenge/app/api"
	"github.com/Gilmardealcantara/go-challenge/app/categories"
	"github.com/Gilmardealcantara/go-challenge/app/products"
)

// Response is the DTO for listing products
type Response struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
}

// ProductResponse is the DTO for a product in the list
type ProductResponse struct {
	Code     string            `json:"code"`
	Price    string            `json:"price"`
	Category *CategoryResponse `json:"category,omitempty"`
}

// CategoryResponse is the DTO for a category
type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ListHandler struct {
	repo products.Repository
}

func NewListHandler(r products.Repository) *ListHandler {
	return &ListHandler{
		repo: r,
	}
}

// Handle godoc
//
//	@Summary		List products
//	@Description	Get a paginated list of products with optional filtering by category and price
//	@Tags			catalog
//	@Produce		json
//	@Param			offset			query		integer	false	"Pagination offset"	default(0)
//	@Param			limit			query		integer	false	"Pagination limit"	default(10)
//	@Param			category		query		string	false	"Filter by category code"
//	@Param			priceLessThan	query		number	false	"Filter by price less than"
//	@Success		200				{object}	Response
//	@Failure		400				{object}	api.ErrorDataResponse
//	@Failure		500				{object}	api.ErrorDataResponse
//	@Router			/catalog [get]
func (h *ListHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters using helper
	offset := api.QueryInt(r, "offset", 0)
	limit := api.QueryInt(r, "limit", 10)

	if limit < 0 || offset < 0 {
		errMsg := "limit and offset need to be positive values"
		api.ErrorResponse(w, http.StatusBadRequest, errMsg)
		return
	}

	if limit > 100 {
		errMsg := "limit cannot be greater than 100"
		api.ErrorResponse(w, http.StatusBadRequest, errMsg)
		return
	}

	// Parse filter parameters
	categoryCode := r.URL.Query().Get("category")
	priceLessThan := api.QueryFloat(r, "priceLessThan")

	params := products.FilterParams{
		Offset:        offset,
		Limit:         limit,
		CategoryCode:  categoryCode,
		PriceLessThan: priceLessThan,
	}

	// Get products with filters
	prods, err := h.repo.GetWithFilters(params)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Get total count
	total, err := h.repo.Total()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert to response DTOs
	catalogProducts := make([]ProductResponse, len(prods))
	for i, p := range prods {
		catalogProducts[i] = h.toProductResponse(p)
	}

	response := &Response{Products: catalogProducts, Total: total}
	api.OKResponse(w, response)
}

// toProductResponse converts a Product domain model to ProductResponse DTO
func (h *ListHandler) toProductResponse(p products.Product) ProductResponse {
	return ProductResponse{
		Code:     p.Code,
		Price:    p.Price.String(),
		Category: h.toCategoryResponse(p.Category),
	}
}

// toCategoryResponse converts a Category to CategoryResponse DTO, handling nil values
func (h *ListHandler) toCategoryResponse(c *categories.Category) *CategoryResponse {
	if c == nil {
		return nil
	}
	return &CategoryResponse{
		Code: c.Code,
		Name: c.Name,
	}
}
