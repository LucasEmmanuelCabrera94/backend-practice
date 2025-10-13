package usecase

import (
	"backend-practice/internal/core/entity"
	"backend-practice/internal/core/repository/user"
	"fmt"
)

type CreateUserUseCase interface {
	CreateUser(name, email string) (entity.User, error)
}

type createUserUseCase struct {
	repo user.Repository
}

func NewCreateUserUseCase(r user.Repository) CreateUserUseCase {
	return &createUserUseCase{repo: r}
}

func (uc *createUserUseCase) CreateUser(name, email string) (entity.User, error) {
	u := entity.User{Name: name, Email: email}
	if !u.IsValid() {
		return entity.User{}, fmt.Errorf("invalid user data")
	}

	return uc.repo.AddUser(u)
}

