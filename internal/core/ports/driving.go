package ports

import (
	"context"

	"go-api/internal/core/domain"
)

type ProductService interface {
	GetProducts() ([]domain.Product, error)
	GetProductById(id int) (domain.Product, error)
	CreateProduct(product domain.Product) (domain.Product, error)
	DeleteProduct(id int) error
	UpdateProduct(product domain.Product, id int) (domain.Product, error)
}

// HealthService expõe a verificação de disponibilidade da aplicação
// (readiness): responde se a app consegue atender — incluindo dependências.
type HealthService interface {
	CheckHealth(ctx context.Context) error
}