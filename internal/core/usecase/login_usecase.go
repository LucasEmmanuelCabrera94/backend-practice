package usecase

import (
	"backend-practice/internal/core/port"
	"backend-practice/internal/infra/transport/dto"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase struct {
	userPort port.UserPort
    sessionPort port.SessionPort
    jwtService port.JWTService
}

func NewLoginUseCase(r port.UserPort, jwts port.JWTService, sp port.SessionPort) *LoginUseCase {
	return &LoginUseCase{userPort: r, jwtService: jwts, sessionPort: sp}
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

	_, err = uc.sessionPort.CreateSession(user.ID, token)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	userDto := dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	return dto.LoginResponse{
		Token: token,
		User:  userDto,
	}, nil
}


