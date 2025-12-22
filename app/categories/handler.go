package categories

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetCategories()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, categories)
}
