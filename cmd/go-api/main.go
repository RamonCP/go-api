package main

import (
	httpadapter "go-api/internal/adapters/http"
	"go-api/internal/adapters/postgres"
	"go-api/internal/core/services"

	"github.com/gin-gonic/gin"
)

func main() {
	dbConnection, err := postgres.ConnectDB()
	if err != nil {
		panic(err)
	}

	if err := postgres.RunMigrations(); err != nil {
		panic(err)
	}

	repo := postgres.NewProductRepository(dbConnection)
	service := services.NewProductService(repo)
	handler := httpadapter.NewProductHandler(service)


	server := gin.Default()
	server.GET("/products", handler.GetProducts)
	server.GET("/product/:id", handler.GetProductById)
	server.POST("/product", handler.CreateProduct)
	server.DELETE("/product/:id", handler.DeleteProduct)

	server.Run(":8000")
}