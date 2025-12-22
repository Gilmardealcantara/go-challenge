package server

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB) *http.ServeMux {
	// Initialize handlers
	prodRepo := catalog.NewRepository(db)
	catHandler := catalog.NewHandler(prodRepo)

	return SetupRoutes(catHandler)
}

func SetupRoutes(catHandler *catalog.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", catHandler.HandleGet)
	mux.HandleFunc("GET /catalog/{code}", catHandler.HandleGetByCode)
	return mux
}
