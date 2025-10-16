package handler_test

import (
	"backend-practice/internal/core/entity"
	"backend-practice/internal/infra/transport/dto"
	"backend-practice/internal/infra/transport/handler"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserUseCase struct {
	mock.Mock
}

func (m *mockUserUseCase) CreateUser(req dto.CreateUserRequest) (entity.User, error) {
	args := m.Called(req)
	return args.Get(0).(entity.User), args.Error(1)
}

func setupRouter(h *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/users", h.CreateUser)
	return r
}

func TestCreateUser_Success(t *testing.T) {
	mockUC := new(mockUserUseCase)
	h := handler.NewUserHandler(mockUC)
	router := setupRouter(h)

	body := map[string]string{"name": "Lucas", "email": "lucas@mail.com", "password": "secretpw"}
	jsonBody, _ := json.Marshal(body)

	expectedUser := entity.User{ID: 1, Name: "Lucas", Email: "lucas@mail.com", Password: "secretpw"}
	mockUC.On("CreateUser", mock.AnythingOfType("dto.CreateUserRequest")).Return(expectedUser, nil)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")
	mockUC.AssertExpectations(t)
}

func TestCreateUser_InvalidRequest(t *testing.T) {
	mockUC := new(mockUserUseCase)
	h := handler.NewUserHandler(mockUC)
	router := setupRouter(h)

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(`{invalid-json}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request")
	mockUC.AssertNotCalled(t, "CreateUser")
}

func TestCreateUser_UseCaseError(t *testing.T) {
	mockUC := new(mockUserUseCase)
	h := handler.NewUserHandler(mockUC)
	router := setupRouter(h)

	body := map[string]string{"name": "Lucas", "email": "lucas@mail.com", "password": "secretpw"}
	jsonBody, _ := json.Marshal(body)

	mockUC.On("CreateUser", mock.AnythingOfType("dto.CreateUserRequest")).Return(entity.User{}, errors.New("invalid user data"))

	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid user data")
	mockUC.AssertExpectations(t)
}
