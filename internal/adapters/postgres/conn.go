package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// ConnectDB abre a conexão com o PostgreSQL a partir de uma DSN e valida que o
// banco responde com um Ping. A construção da DSN fica a cargo de quem chama
// (a camada de config), mantendo este adapter sem conhecer variáveis de ambiente.
func ConnectDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrindo conexão com o banco: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("verificando conexão com o banco (ping): %w", err)
	}

	return db, nil
}
