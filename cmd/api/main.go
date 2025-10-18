package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"rest-api-2/internal/config"
	"rest-api-2/internal/database"
	"rest-api-2/internal/handlers"
	"rest-api-2/internal/repositories"
	"rest-api-2/internal/usecases"
	"rest-api-2/internal/version"
	"syscall"
)

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Conectar ao banco de dados
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	slog.Info("Starting REST API server", "version", version.Info())

	// Inicializar camadas da Clean Architecture
	repos := repositories.New(db)
	useCases := usecases.New(repos)
	h := handlers.New(useCases)

	// Canal para capturar sinais do sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	go func() {
		slog.Info("Server starting", "port", cfg.App.Port)
		if err := h.Listen(cfg.App.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Aguardar sinal de encerramento
	<-quit
	slog.Info("Server shutting down gracefully...")
}
