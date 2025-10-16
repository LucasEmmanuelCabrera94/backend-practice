package usecase_test

import (
	"backend-practice/internal/core/entity"
	"backend-practice/internal/core/usecase"
	"backend-practice/internal/infra/transport/dto"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type inMemoryUserRepo struct{}

func (r *inMemoryUserRepo) CreateUser(u entity.User) (entity.User, error) {
    return u, nil
}

func (r *inMemoryUserRepo) GetUserByEmail(email string) (entity.User, error) {
    return entity.User{}, nil
}

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) CreateUser(u entity.User) (entity.User, error) {
	args := m.Called(u)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *mockUserRepository) GetUserByEmail(email string) (entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(entity.User), args.Error(1)
}

func TestCreateUser_Success(t *testing.T) {
	repo := &inMemoryUserRepo{} // repo que devuelve el mismo user
	uc := usecase.NewCreateUserUseCase(repo)

	input := dto.CreateUserRequest{Name: "Lucas", Email: "lucas@mail.com", Password: "password"}

	result, err := uc.CreateUser(input)

	assert.NoError(t, err)
	assert.Equal(t, "Lucas", result.Name)
	assert.Equal(t, "lucas@mail.com", result.Email)
	assert.NotEmpty(t, result.PasswordHash)
}


func TestCreateUser_InvalidData(t *testing.T) {
	repo := new(mockUserRepository)
	uc := usecase.NewCreateUserUseCase(repo)

	_, err := uc.CreateUser(dto.CreateUserRequest{Name: "", Email: "invalid-email", Password: "password"})
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid user data")
	repo.AssertNotCalled(t, "CreateUser")
}

func TestCreateUser_RepositoryError(t *testing.T) {
	repo := new(mockUserRepository)
	uc := usecase.NewCreateUserUseCase(repo)

	input := dto.CreateUserRequest{Name: "Lucas", Email: "lucas@mail.com", Password: "password"}
	repo.On("CreateUser", mock.MatchedBy(func(u entity.User) bool {
		return u.Name == input.Name && u.Email == input.Email && u.Password != ""
	})).Return(entity.User{}, errors.New("db error"))

	_, err := uc.CreateUser(input)

	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
	repo.AssertExpectations(t)
}

func TestCreateUser_PasswordIsHashed(t *testing.T) {
	repo := &inMemoryUserRepo{}
	uc := usecase.NewCreateUserUseCase(repo)

	input := dto.CreateUserRequest{Name: "Lucas", Email: "lucas@mail.com", Password: "secret"}

	result, err := uc.CreateUser(input)

	assert.NoError(t, err)
	assert.NotEqual(t, input.Password, result.PasswordHash)
	assert.Contains(t, result.PasswordHash, "$2a$")
}

