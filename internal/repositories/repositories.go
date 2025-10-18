package repositories

import (
	"database/sql"
	"rest-api-2/internal/models"
	"rest-api-2/internal/repositories/users"
)

type Repositories struct {
	User interface {
		GetAll() []models.User
		Add(user models.User)
		EmailExists(email string) bool
	}
}

// New cria repositórios com conexão PostgreSQL
func New(db *sql.DB) *Repositories {
	return &Repositories{
		User: users.New(db),
	}
}
