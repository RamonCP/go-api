package ports

import "go-api/internal/core/domain"

type ProductService interface {
	GetProducts() ([]domain.Product, error)
	GetProductById(id int) (domain.Product, error)
}