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
	categoriesHandler := categories.NewHandler(catRepo)

	return SetupRoutes(listHandler, detailHandler, categoriesHandler)
}

func SetupRoutes(listHandler *catalog.ListHandler, detailHandler *catalog.DetailHandler, categoriesHandler *categories.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", listHandler.Handle)
	mux.HandleFunc("GET /catalog/{code}", detailHandler.Handle)
	mux.HandleFunc("GET /categories", categoriesHandler.HandleGet)
	mux.HandleFunc("POST /categories", categoriesHandler.HandleCreate)
	return mux
}
