package port

import (
	"backend-practice/internal/core/entity"
)

type UserPort interface {
	CreateUser(req entity.User) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
}
