package categories

import (
	"encoding/json"
	"net/http"

	"github.com/Gilmardealcantara/go-challenge/app/api"
)

// CreateCategoryRequest is the DTO for creating a category
type CreateCategoryRequest struct {
	Code string `json:"code" example:"electronics"`
	Name string `json:"name" example:"Electronics"`
}

type CreateHandler struct {
	repo Repository
}

func NewCreateHandler(r Repository) *CreateHandler {
	return &CreateHandler{
		repo: r,
	}
}

// Handle godoc
//
//	@Summary		Create category
//	@Description	Create a new product category
//	@Tags			categories
//	@Accept			json
//	@Produce		json
//	@Param			request	body		CreateCategoryRequest	true	"Category data"
//	@Success		201		{object}	CategoryResponse
//	@Failure		400		{object}	api.ErrorDataResponse
//	@Failure		500		{object}	api.ErrorDataResponse
//	@Router			/categories [post]
func (h *CreateHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Code == "" || req.Name == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "code and name are required")
		return
	}

	category := &Category{
		Code: req.Code,
		Name: req.Name,
	}

	if err := h.repo.Create(category); err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := &CategoryResponse{
		Code: category.Code,
		Name: category.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
