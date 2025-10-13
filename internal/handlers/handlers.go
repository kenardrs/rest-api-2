package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"rest-api-2/internal/usecases"
)

type Handlers struct {
	useCases *usecases.UseCases
}

func New(useCases *usecases.UseCases) *Handlers {
	return &Handlers{
		useCases: useCases, // ← Inicializar o campo!
	}
}

func (h Handlers) Listen(port int) error {
	h.registerUserEndpoints()
	slog.Info("Listening on ", "port", port)
	return http.ListenAndServe(
		fmt.Sprintf(":%v", port),
		nil,
	)
}
