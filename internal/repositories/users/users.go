package users

// Simula um repositório de usuários em memória

import "rest-api-2/internal/models"

type Users struct {
	users []models.User
}

func New() *Users {
	return &Users{users: make([]models.User, 0)}
}

func (u *Users) GetAll() []models.User {
	return []models.User{}
}

func (u *Users) EmailExists(email string) bool {
	for _, user := range u.users {
		if user.Email == email {
			return true
		}
	}
	return false
}

func (u *Users) Add(newUser models.User) {
	u.users = append(u.users, newUser)
}
