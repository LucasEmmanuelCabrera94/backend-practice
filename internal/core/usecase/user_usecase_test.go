package usecase_test

import (
	"backend-practice/internal/core/entity"
	"backend-practice/internal/core/usecase"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) AddUser(u entity.User) (entity.User, error) {
	args := m.Called(u)
	return args.Get(0).(entity.User), args.Error(1)
}

func TestCreateUser_Success(t *testing.T) {
	repo := new(mockUserRepository)
	uc := usecase.NewCreateUserUseCase(repo)

	input := entity.User{Name: "Lucas", Email: "lucas@mail.com"}
	expected := entity.User{ID: 1, Name: "Lucas", Email: "lucas@mail.com"}

	repo.On("AddUser", input).Return(expected, nil)

	result, err := uc.CreateUser(input.Name, input.Email)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}

func TestCreateUser_InvalidData(t *testing.T) {
	repo := new(mockUserRepository)
	uc := usecase.NewCreateUserUseCase(repo)

	_, err := uc.CreateUser("", "")
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid user data")
	repo.AssertNotCalled(t, "AddUser")
}

func TestCreateUser_RepositoryError(t *testing.T) {
	repo := new(mockUserRepository)
	uc := usecase.NewCreateUserUseCase(repo)

	input := entity.User{Name: "Lucas", Email: "lucas@mail.com"}
	repo.On("AddUser", input).Return(entity.User{}, errors.New("db error"))

	_, err := uc.CreateUser(input.Name, input.Email)

	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
	repo.AssertExpectations(t)
}
