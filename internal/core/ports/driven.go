package ports

import (
	"context"

	"go-api/internal/core/domain"
)

type ProductRepository interface {
	GetProducts() ([]domain.Product, error)
	GetProductById(id int) (domain.Product, error)
	CreateProduct(product domain.Product) (domain.Product, error)
	DeleteProduct(id int) error
	UpdateProduct(product domain.Product, id int) (domain.Product, error)
}

// HealthRepository representa a capacidade de verificar a saúde de uma
// dependência externa (ex.: o banco de dados respondendo a um ping).
type HealthRepository interface {
	CheckHealth(ctx context.Context) error
}