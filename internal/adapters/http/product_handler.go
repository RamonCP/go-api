package http

import (
	"fmt"
	"go-api/internal/core/domain"
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

func (h *productHandler) CreateProduct(ctx *gin.Context) {
	var product domain.Product
	err := ctx.BindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Body inválido")
		return
	}

	product, err = h.service.CreateProduct(product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (h *productHandler) DeleteProduct(ctx *gin.Context) {
	id_param := ctx.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Id inválido")
		return
	}

	err = h.service.DeleteProduct(id)
	fmt.Println("DeleteProduct: ", err)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Erro ao deletar produto")
		return
	}

	ctx.JSON(http.StatusOK, "Produto deletado com sucesso")
}

func (h *productHandler) UpdateProduct(ctx *gin.Context) {
	id_param := ctx.Param("id")
	id, err := strconv.Atoi(id_param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Id inválido")
		return
	}

	var product domain.Product
	err = ctx.BindJSON(&product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Body inválido")
		return
	}

	product, err = h.service.UpdateProduct(product, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Erro ao atualizar produto")
		return
	}

	ctx.JSON(http.StatusOK, product)
}