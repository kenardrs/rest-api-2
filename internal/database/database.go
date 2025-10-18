package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"rest-api-2/internal/config"
	"time"

	_ "github.com/lib/pq"
)

// Connect estabelece conexão com PostgreSQL com pool de conexões
func Connect(cfg *config.DatabaseConfig) (*sql.DB, error) {
	slog.Info("Connecting to database", "host", cfg.Host, "port", cfg.Port, "database", cfg.Name)

	// Abrir conexão
	db, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configurar pool de conexões
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Testar conexão
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	slog.Info("Database connection established successfully",
		"max_open_conns", cfg.MaxOpenConns,
		"max_idle_conns", cfg.MaxIdleConns,
		"conn_max_lifetime", cfg.ConnMaxLifetime,
	)

	return db, nil
}

// HealthCheck verifica se a conexão com o banco está saudável
func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// Close fecha a conexão com o banco de dados
func Close(db *sql.DB) error {
	slog.Info("Closing database connection")
	return db.Close()
}
