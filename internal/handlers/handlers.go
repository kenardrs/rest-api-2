package handlers

import (
	"fmt"
	"net/http"
	"rest-api-2/internal/usecases"
)

type Handlers struct {
	usecases *usecases.UseCases
}

func New(usecases *usecases.UseCases) *Handlers {
	return &Handlers{}
}

func (h Handlers) Listen(port int) error {
	return http.ListenAndServe(
		fmt.Sprintf(":%v", port),
		nil,
	)
}
