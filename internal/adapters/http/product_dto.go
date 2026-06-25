package http

import "go-api/internal/core/domain"

type ProductResponse struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
}

func toProductResponse(p domain.Product) ProductResponse {
	return ProductResponse{ID: p.ID, Name: p.Name, Price: p.Price}
}