package catalog

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/products"
)

type Response struct {
	Products []ProductResponse `json:"products"`
}

type ProductResponse struct {
	Code     string            `json:"code"`
	Price    float64           `json:"price"`
	Category *CategoryResponse `json:"category,omitempty"`
}

type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Handler struct {
	repo products.Repository
}

func NewHandler(r products.Repository) *Handler {
	return &Handler{
		repo: r,
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.GetAllProducts()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Map domain models to API response models
	catalogProducts := make([]ProductResponse, len(products))
	for i, p := range products {
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

	api.OKResponse(w, Response{Products: catalogProducts})
}
