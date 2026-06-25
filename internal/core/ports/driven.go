package ports

import "go-api/internal/core/domain"

type ProductRepository interface {
	GetProducts() ([]domain.Product, error)
	GetProductById(id int) (domain.Product, error)
	CreateProduct(product domain.Product) (domain.Product, error)
	DeleteProduct(id int) error
	UpdateProduct(product domain.Product, id int) (domain.Product, error)
}