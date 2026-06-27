package main

import (
	httpadapter "go-api/internal/adapters/http"
	"go-api/internal/adapters/postgres"
	"go-api/internal/config"
	"go-api/internal/core/services"

	"github.com/gin-gonic/gin"
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

	server := gin.Default()
	server.GET("/health", healthHandler.CheckHealth)
	server.GET("/products", handler.GetProducts)
	server.GET("/product/:id", handler.GetProductById)
	server.POST("/product", handler.CreateProduct)
	server.DELETE("/product/:id", handler.DeleteProduct)
	server.PUT("/product/:id", handler.UpdateProduct)

	if err := server.Run(":" + cfg.Server.Port); err != nil {
		panic(err)
	}
}
