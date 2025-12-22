package catalog

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
)

type Response struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
}

type ProductResponse struct {
	Code     string            `json:"code"`
	Price    float64           `json:"price"`
	Category *CategoryResponse `json:"category,omitempty"`
}

type ProductDetailsResponse struct {
	Code     string            `json:"code"`
	Price    float64           `json:"price"`
	Category *CategoryResponse `json:"category,omitempty"`
	Variants []VariantResponse `json:"variants"`
}

type VariantResponse struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Handler struct {
	repo Repository
}

func NewHandler(r Repository) *Handler {
	return &Handler{
		repo: r,
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
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

	params := FilterParams{
		Offset:        offset,
		Limit:         limit,
		CategoryCode:  categoryCode,
		PriceLessThan: priceLessThan,
	}

	prods, err := h.repo.GetWithFilters(params)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	total, err := h.repo.Total()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map domain models to API response models
	catalogProducts := make([]ProductResponse, len(prods))
	for i, p := range prods {
		var category *CategoryResponse
		if p.Category != nil {
			category = &CategoryResponse{
				Code: p.Category.Code,
				Name: p.Category.Name,
			}
		}
		catalogProducts[i] = ProductResponse{
			Code:     p.Code,
			Price:    p.Price.InexactFloat64(),
			Category: category,
		}
	}

	api.OKResponse(w, Response{Products: catalogProducts, Total: total})
}

func (h *Handler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	product, err := h.repo.GetByCode(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "product not found")
		return
	}

	// Map variants, inheriting price from product if variant price is null
	variants := make([]VariantResponse, len(product.Variants))
	for i, v := range product.Variants {
		price := v.Price.InexactFloat64()
		if v.Price.IsZero() {
			price = product.Price.InexactFloat64()
		}
		variants[i] = VariantResponse{
			Name:  v.Name,
			SKU:   v.SKU,
			Price: price,
		}
	}

	var category *CategoryResponse
	if product.Category != nil {
		category = &CategoryResponse{
			Code: product.Category.Code,
			Name: product.Category.Name,
		}
	}

	response := ProductDetailsResponse{
		Code:     product.Code,
		Price:    product.Price.InexactFloat64(),
		Category: category,
		Variants: variants,
	}

	api.OKResponse(w, response)
}
