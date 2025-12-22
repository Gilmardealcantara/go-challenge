package catalog

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
)

type Handler struct {
	service Service
}

func NewHandler(r Repository) *Handler {
	return &Handler{
		service: NewService(r),
	}
}

func NewHandlerWithService(s Service) *Handler {
	return &Handler{
		service: s,
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

	response, err := h.service.GetProducts(params)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.OKResponse(w, response)
}

func (h *Handler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	response, err := h.service.GetProductByCode(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	api.OKResponse(w, response)
}
