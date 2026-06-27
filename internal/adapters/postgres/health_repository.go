package postgres

import (
	"context"
	"database/sql"
)

type HealthRepository struct {
	connection *sql.DB
}

func NewHealthRepository(connection *sql.DB) *HealthRepository {
	return &HealthRepository{
		connection: connection,
	}
}

// CheckHealth verifica se o banco está acessível. PingContext respeita o
// deadline do contexto, evitando que o health check fique pendurado.
func (hr *HealthRepository) CheckHealth(ctx context.Context) error {
	return hr.connection.PingContext(ctx)
}
