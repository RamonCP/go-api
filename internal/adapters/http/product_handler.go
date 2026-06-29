package http

import (
	"encoding/json"
	"fmt"
	"go-api/internal/core/domain"
	"go-api/internal/core/ports"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type productHandler struct {
	service ports.ProductService
}

func NewProductHandler(service ports.ProductService) *productHandler {
	return &productHandler{
		service: service,
	}
}

func (h *productHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetProducts()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}
	
	response := make([]ProductResponse, len(products))
	for i, p := range products {
		response[i] = toProductResponse(p)
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *productHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Id inválido")
		return
	}

	product, err := h.service.GetProductById(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	response := toProductResponse(product)

	writeJSON(w, http.StatusOK, response)
}

func (h *productHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product domain.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	// err := ctx.BindJSON(&product)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Body inválido")
		return
	}

	product, err = h.service.CreateProduct(product)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (h *productHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Id inválido")
		return
	}

	err = h.service.DeleteProduct(id)
	fmt.Println("DeleteProduct: ", err)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Erro ao deletar produto")
		return
	}

	writeJSON(w, http.StatusOK, "Produto deletado com sucesso")
}

func (h *productHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Id inválido")
		return
	}

	var product domain.Product
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Body inválido")
		return
	}

	product, err = h.service.UpdateProduct(product, id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, "Erro ao atualizar produto")
		return
	}

	writeJSON(w, http.StatusOK, product)
}