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

    input := entity.User{Name: "Lucas", Email: "lucas@mail.com", Password: "password"}
    expected := entity.User{ID: 1, Name: "Lucas", Email: "lucas@mail.com"}

    repo.On("AddUser", mock.MatchedBy(func(u entity.User) bool {
        return u.Name == input.Name && u.Email == input.Email && u.Password != ""
    })).Return(expected, nil)

    result, err := uc.CreateUser(dto.CreateUserRequest{Name: input.Name, Email: input.Email, Password: input.Password})

    assert.NoError(t, err)
    assert.Equal(t, expected.ID, result.ID)
    assert.NotNil(t, result.PasswordHash)
    repo.AssertExpectations(t)
}


func TestCreateUser_InvalidData(t *testing.T) {
	repo := new(mockUserRepository)
	uc := usecase.NewCreateUserUseCase(repo)

	_, err := uc.CreateUser(dto.CreateUserRequest{Name: "", Email: "invalid-email", Password: "password"})
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid user data")
	repo.AssertNotCalled(t, "AddUser")
}

func TestCreateUser_RepositoryError(t *testing.T) {
	repo := new(mockUserRepository)
	uc := usecase.NewCreateUserUseCase(repo)

    input := entity.User{Name: "Lucas", Email: "lucas@mail.com", Password: "password"}
	repo.On("AddUser", mock.MatchedBy(func(u entity.User) bool {
        return u.Name == input.Name && u.Email == input.Email && u.Password != ""
    })).Return(entity.User{}, errors.New("db error"))

	_, err := uc.CreateUser(dto.CreateUserRequest{Name: input.Name, Email: input.Email, Password: input.Password})

	assert.Error(t, err)
	assert.EqualError(t, err, "db error")
	repo.AssertExpectations(t)
}
