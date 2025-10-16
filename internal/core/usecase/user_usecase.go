package usecase

import (
	"backend-practice/internal/core/entity"
	"backend-practice/internal/core/port"
	"backend-practice/internal/infra/transport/dto"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	port port.UserPort
}

func NewCreateUserUseCase(p port.UserPort) *UserUseCase {
	return &UserUseCase{port: p}
}

func (uc *UserUseCase) CreateUser(req dto.CreateUserRequest) (entity.User, error) {
	u := entity.User{Name: req.Name, Email: req.Email, Password: req.Password}
	if !u.IsValid() {
		return entity.User{}, fmt.Errorf("invalid user data")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u.PasswordHash = string(hashed)
	
	return uc.port.CreateUser(u)
}

