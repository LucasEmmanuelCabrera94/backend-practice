package db

import (
	"backend-practice/internal/core/entity"
	"backend-practice/internal/core/port"
	"database/sql"
)

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) GetUserByEmail(email string) (entity.User, error) {
	row := r.db.QueryRow("SELECT id, name, email, passwordhash FROM users WHERE email = ?", email)
	var u entity.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, nil
		}
		return entity.User{}, err
	}
	return u, nil
}

func NewUserRepository(db *sql.DB) port.UserPort {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(u entity.User) (entity.User, error) {
	result, err := r.db.Exec("INSERT INTO users (name, email, passwordhash) VALUES (?, ?, ?)", u.Name, u.Email, u.PasswordHash)
	if err != nil {
		return entity.User{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return entity.User{}, err
	}
	u.ID = id
	return u, nil
}
