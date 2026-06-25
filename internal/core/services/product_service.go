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