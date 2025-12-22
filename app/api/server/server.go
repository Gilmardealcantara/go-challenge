package server

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/categories"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) *http.ServeMux {
	// Initialize handlers
	prodRepo := catalog.NewRepository(db)
	catHandler := catalog.NewHandler(catalog.NewService(prodRepo))

	catRepo := categories.NewRepository(db)
	categoriesHandler := categories.NewHandler(categories.NewService(catRepo))

	return SetupRoutes(catHandler, categoriesHandler)
}

func SetupRoutes(catHandler *catalog.Handler, categoriesHandler *categories.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", catHandler.HandleGet)
	mux.HandleFunc("GET /catalog/{code}", catHandler.HandleGetByCode)
	mux.HandleFunc("GET /categories", categoriesHandler.HandleGet)
	return mux
}
