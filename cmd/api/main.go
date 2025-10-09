package main

import (
	"rest-api-2/internal/handlers"
	"rest-api-2/internal/repositories"
	"rest-api-2/internal/usecases"
)

// cadastrar e listar usuários
// handler <- usecases/services <- repositories <- models
//
//	<- external services (e.g., email, payment gateways)
func main() {
	repos := repositories.New()

	useCases := usecases.New(repos)

	h := handlers.New(useCases)

	h.Listen(3000)
}
