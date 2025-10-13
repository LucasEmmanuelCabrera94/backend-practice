package db

import (
	"backend-practice/internal/core/entity"
	"backend-practice/internal/core/repository/user"
	"database/sql"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) AddUser(u entity.User) (entity.User, error) {
	result, err := r.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", u.Name, u.Email)
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
