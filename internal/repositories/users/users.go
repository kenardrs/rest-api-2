package users

import (
	"database/sql"
	"log/slog"
	"rest-api-2/internal/models"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Users implementa repository para PostgreSQL
type Users struct {
	db *sql.DB
}

// New cria nova instância do repository PostgreSQL
func New(db *sql.DB) *Users {
	return &Users{
		db: db,
	}
}

// GetAll retorna todos os usuários do banco
func (u *Users) GetAll() []models.User {
	query := `
		SELECT id, name, email 
		FROM users 
		ORDER BY created_at DESC
	`

	rows, err := u.db.Query(query)
	if err != nil {
		slog.Error("Error querying users", "error", err)
		return []models.User{}
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			slog.Error("Error scanning user", "error", err)
			continue
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Error iterating user rows", "error", err)
	}

	return users
}

// EmailExists verifica se email já existe no banco
func (u *Users) EmailExists(email string) bool {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := u.db.QueryRow(query, email).Scan(&exists)
	if err != nil {
		slog.Error("Error checking email existence", "error", err, "email", email)
		return false
	}

	return exists
}

// Add adiciona novo usuário no banco
func (u *Users) Add(newUser models.User) {
	query := `
		INSERT INTO users (id, name, email) 
		VALUES ($1, $2, $3)
	`

	_, err := u.db.Exec(query, newUser.ID, newUser.Name, newUser.Email)
	if err != nil {
		slog.Error("Error inserting user", "error", err, "user", newUser)
		return
	}

	slog.Info("User created successfully", "id", newUser.ID, "email", newUser.Email)
}
