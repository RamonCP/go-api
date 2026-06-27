package http

import (
	"context"
	"net/http"
	"time"

	"go-api/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type healthHandler struct {
	service ports.HealthService
}

func NewHealthHandler(service ports.HealthService) *healthHandler {
	return &healthHandler{
		service: service,
	}
}

// CheckHealth é um readiness probe: responde 200 se a app e suas dependências
// (o banco) estão prontas, ou 503 caso contrário.
func (h *healthHandler) CheckHealth(ctx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.CheckHealth(reqCtx); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "unavailable",
			"database": "down",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": "up",
	})
}
