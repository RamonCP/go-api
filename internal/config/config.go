package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config agrega toda a configuração da aplicação, lida do ambiente.
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// ServerConfig controla o servidor HTTP.
type ServerConfig struct {
	Port string
}

// DatabaseConfig controla a conexão com o PostgreSQL.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Load monta a Config a partir das variáveis de ambiente, caindo em defaults
// voltados ao desenvolvimento local quando a variável não está definida.
func Load() Config {
	return Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8000"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "1234"),
			Name:     getEnv("DB_NAME", "postgres"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

// DSN devolve a connection string no formato key=value, usada pelo driver lib/pq.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// URL devolve a connection string no formato postgres://, usada pelo golang-migrate.
func (d DatabaseConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode,
	)
}

// getEnv lê uma variável de ambiente ou devolve o fallback se ela não existir.
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsInt é como getEnv, mas converte para int; mantém o fallback se a
// variável não existir ou não for um número válido.
func getEnvAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(value); err == nil {
			return n
		}
	}
	return fallback
}
