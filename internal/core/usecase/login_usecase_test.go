package usecase_test

import (
	"backend-practice/internal/core/entity"
	"backend-practice/internal/core/usecase"
	"backend-practice/internal/infra/transport/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type mockUserPort struct {
	user entity.User
}

func (f *mockUserPort) GetUserByEmail(email string) (entity.User, error) {
	return f.user, nil
}

func (f *mockUserPort) CreateUser(u entity.User) (entity.User, error) {
	return u, nil
}

type mockUserPortError struct{}

func (f *mockUserPortError) GetUserByEmail(email string) (entity.User, error) {
	return entity.User{}, assert.AnError
}

func (f *mockUserPortError) CreateUser(u entity.User) (entity.User, error) {
	return u, nil
}

type mockJWTService struct {
	token string
}

func (f *mockJWTService) GenerateToken(userID int64) (string, error) {
	return f.token, nil
}

func (f *mockJWTService) ValidateToken(token string) (int64, error) {
	return 1, nil 
}

type mockSessionPort struct{}

func (f *mockSessionPort) CreateSession(userID int64, token string) (entity.Session, error) {
	return entity.Session{UserID: userID, Token: token}, nil
}

func TestLogin_Success(t *testing.T) {
	password := "secret"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := entity.User{
		ID:           1,
		Name:         "Lucas",
		Email:        "lucas@mail.com",
		PasswordHash: string(hash),
	}

	userPort := &mockUserPort{user: user}
	jwtService := &mockJWTService{token: "fake-jwt-token"}
	sessionPort := &mockSessionPort{}

	uc := usecase.NewLoginUseCase(userPort, jwtService, sessionPort)

	input := dto.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	result, err := uc.Login(input)

	assert.NoError(t, err)
	assert.Equal(t, "fake-jwt-token", result.Token)
	assert.Equal(t, user.ID, result.User.ID)
	assert.Equal(t, user.Name, result.User.Name)
	assert.Equal(t, user.Email, result.User.Email)
}

func TestLogin_InvalidPassword(t *testing.T) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)

	user := entity.User{
		ID:           1,
		Name:         "Lucas",
		Email:        "lucas@mail.com",
		PasswordHash: string(hash),
	}

	userPort := &mockUserPort{user: user}
	jwtService := &mockJWTService{token: "fake-jwt-token"}
	sessionPort := &mockSessionPort{}

	uc := usecase.NewLoginUseCase(userPort, jwtService, sessionPort)

	input := dto.LoginRequest{
		Email:    user.Email,
		Password: "wrongpassword",
	}

	_, err := uc.Login(input)
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestLogin_UserNotFound(t *testing.T) {
	userPort := &mockUserPortError{}
	jwtService := &mockJWTService{token: "fake-jwt-token"}
	sessionPort := &mockSessionPort{}

	uc := usecase.NewLoginUseCase(userPort, jwtService, sessionPort)

	input := dto.LoginRequest{
		Email:    "noone@mail.com",
		Password: "whatever",
	}

	_, err := uc.Login(input)
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}
