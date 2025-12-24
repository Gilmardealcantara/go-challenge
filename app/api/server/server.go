package server

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"github.com/mytheresa/go-hiring-challenge/app/products"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) *http.ServeMux {
	// Initialize handlers
	prodRepo := products.NewRepository(db)
	listHandler := catalog.NewListHandler(prodRepo)
	detailHandler := catalog.NewDetailHandler(prodRepo)

	catRepo := categories.NewRepository(db)
	getHandler := categories.NewGetHandler(catRepo)
	createHandler := categories.NewCreateHandler(catRepo)

	return SetupRoutes(listHandler, detailHandler, getHandler, createHandler)
}

func SetupRoutes(listHandler *catalog.ListHandler, detailHandler *catalog.DetailHandler, getHandler *categories.GetHandler, createHandler *categories.CreateHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", listHandler.Handle)
	mux.HandleFunc("GET /catalog/{code}", detailHandler.Handle)
	mux.HandleFunc("GET /categories", getHandler.Handle)
	mux.HandleFunc("POST /categories", createHandler.Handle)
	return mux
}
