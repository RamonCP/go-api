package services

import (
	"go-api/internal/core/domain"
	"go-api/internal/core/ports"
)


type productService struct {
	repo ports.ProductRepository
}	

func NewProductService(repo ports.ProductRepository) *productService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) GetProducts() ([]domain.Product, error) {
	return s.repo.GetProducts()
}

func (s *productService) GetProductById(id int) (domain.Product, error) {
	return s.repo.GetProductById(id)
}

func (s *productService) CreateProduct(product domain.Product) (domain.Product, error) {
	return s.repo.CreateProduct(product)
}

func (s *productService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}

func (s *productService) UpdateProduct(product domain.Product, id int) (domain.Product, error) {
	return s.repo.UpdateProduct(product, id)
}