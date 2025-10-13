package port

import "backend-practice/internal/core/entity"

type CreateUserPort interface {
	CreateUser(name, email string) (entity.User, error)
}
