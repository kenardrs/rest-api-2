package usecases

import (
	"errors"
	"log/slog"
	"rest-api-2/internal/models"
	"rest-api-2/internal/repositories"

	"github.com/google/uuid"
)

type UseCases struct {
	repos *repositories.Repositories
}

func New(repos *repositories.Repositories) *UseCases {
	return &UseCases{
		repos: repos,
	}
}

func (u UseCases) GetAll() []models.User {
	users := u.repos.User.GetAll()

	return users
}

func (u UseCases) Add(newUser models.User) (uuid.UUID, error) {
	exists := u.repos.User.EmailExists(newUser.Email)
	if exists {
		slog.Error("this user already exists", "email", newUser.Email)

		return uuid.Nil, errors.New("email already exists")
	}

	repoReq := models.User{
		ID:    uuid.New(),
		Name:  newUser.Name,
		Email: newUser.Email,
	}

	u.repos.User.Add(repoReq)

	return repoReq.ID, nil
}
