package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Gilmardealcantara/go-challenge/app/catalog"
	"github.com/Gilmardealcantara/go-challenge/app/categories"
	"github.com/Gilmardealcantara/go-challenge/app/products"
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", listHandler.Handle)
	mux.HandleFunc("GET /catalog/{code}", detailHandler.Handle)
	mux.HandleFunc("GET /categories", getHandler.Handle)
	mux.HandleFunc("POST /categories", createHandler.Handle)

	// Swagger documentation
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /swagger/*", httpSwagger.WrapHandler)

	return mux
}
