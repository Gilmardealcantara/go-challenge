package categories

import (
	"net/http"

	"github.com/Gilmardealcantara/go-challenge/app/api"
)

// CategoryResponse is the DTO for a category
type CategoryResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type GetHandler struct {
	repo Repository
}

func NewGetHandler(r Repository) *GetHandler {
	return &GetHandler{
		repo: r,
	}
}

// Handle godoc
//
//	@Summary		List categories
//	@Description	Get a list of all available product categories
//	@Tags			categories
//	@Produce		json
//	@Success		200	{array}		CategoryResponse
//	@Failure		500	{object}	api.ErrorDataResponse
//	@Router			/categories [get]
func (h *GetHandler) Handle(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.GetAll()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	responses := make([]CategoryResponse, len(categories))
	for i, c := range categories {
		responses[i] = CategoryResponse{
			Code: c.Code,
			Name: c.Name,
		}
	}

	api.OKResponse(w, responses)
}
