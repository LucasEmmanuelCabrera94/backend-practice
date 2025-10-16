package usecase

import (
	"backend-practice/internal/core/port"
	"backend-practice/internal/infra/transport/dto"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase struct {
	userPort port.UserPort
    jwtService port.JWTService
}

func NewLoginUseCase(r port.UserPort, jwts port.JWTService) *LoginUseCase {
	return &LoginUseCase{userPort: r, jwtService: jwts}
}

func (uc *LoginUseCase) Login(req dto.LoginRequest) (dto.LoginResponse, error) {

	user, err := uc.userPort.GetUserByEmail(req.Email)
    if err != nil {
        return dto.LoginResponse{}, errors.New("invalid credentials")
    }

    if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
        return dto.LoginResponse{}, errors.New("invalid credentials")
    }

    token, err := uc.jwtService.GenerateToken(user.ID)
    if err != nil {
        return dto.LoginResponse{}, err
    }

    return dto.LoginResponse{
        Token: token,
        User:  user,
    }, nil
}

