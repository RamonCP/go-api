package http

import (
	"context"
	"net/http"
	"time"

	"go-api/internal/core/ports"
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
func (h *healthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	reqCtx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.service.CheckHealth(reqCtx); err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]any{
			"status":   "unavailable",
			"database": "down",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "ok",
		"database": "up",
	})
}
