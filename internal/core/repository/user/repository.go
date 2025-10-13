package user

import "backend-practice/internal/core/entity"

type Repository interface {
	AddUser(u entity.User) (entity.User, error)
}
