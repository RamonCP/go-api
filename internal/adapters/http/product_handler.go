package http

import (
	"go-api/internal/core/ports"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productHandler struct {
	service ports.ProductService
}

func NewProductHandler(service ports.ProductService) *productHandler {
	return &productHandler{
		service: service,
	}
}

func (h *productHandler) GetProducts(ctx *gin.Context) {
	products, err := h.service.GetProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	
	response := make([]ProductResponse, len(products))
	for i, p := range products {
		response[i] = toProductResponse(p)
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *productHandler) GetProductById(ctx *gin.Context) {
	product_id := ctx.Param("id")
	id, err := strconv.Atoi(product_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Id inválido")
		return
	}

	product, err := h.service.GetProductById(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	response := toProductResponse(product)

	ctx.JSON(http.StatusOK, response)
}