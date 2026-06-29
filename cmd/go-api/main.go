package main

import (
	httpadapter "go-api/internal/adapters/http"
	"go-api/internal/adapters/postgres"
	"go-api/internal/config"
	"go-api/internal/core/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	dbConnection, err := postgres.ConnectDB(cfg.Database.DSN())
	if err != nil {
		panic(err)
	}

	if err := postgres.RunMigrations(cfg.Database.URL()); err != nil {
		panic(err)
	}

	repo := postgres.NewProductRepository(dbConnection)
	service := services.NewProductService(repo)
	handler := httpadapter.NewProductHandler(service)

	healthRepo := postgres.NewHealthRepository(dbConnection)
	healthService := services.NewHealthService(healthRepo)
	healthHandler := httpadapter.NewHealthHandler(healthService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/health", healthHandler.CheckHealth)
	router.Get("/products", handler.GetProducts)
	router.Get("/product/{id}", handler.GetProductById)
	router.Post("/product", handler.CreateProduct)
	router.Delete("/product/{id}", handler.DeleteProduct)
	router.Put("/product/{id}", handler.UpdateProduct)

	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		panic(err)
	}
}
