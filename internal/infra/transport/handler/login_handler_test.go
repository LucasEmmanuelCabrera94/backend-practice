package handler_test

import (
	"backend-practice/internal/core/entity"
	"backend-practice/internal/core/usecase"
	"backend-practice/internal/infra/transport/dto"
	"backend-practice/internal/infra/transport/handler"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// --- Mocks para el UseCase ---
type mockUserPort struct {
	mock.Mock
}

func (m *mockUserPort) GetUserByEmail(email string) (entity.User, error) {
	args := m.Called(email)
	return args.Get(0).(entity.User), args.Error(1)
}

func (m *mockUserPort) CreateUser(u entity.User) (entity.User, error) {
	args := m.Called(u)
	return args.Get(0).(entity.User), args.Error(1)
}

type mockSessionPort struct {
	mock.Mock
}

func (m *mockSessionPort) CreateSession(userId int64, token string) (entity.Session, error) {
	args := m.Called(userId, token)
	return args.Get(0).(entity.Session), args.Error(1)
}

type mockJWTService struct {
	mock.Mock
}

func (m *mockJWTService) GenerateToken(userID int64) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *mockJWTService) ValidateToken(token string) (int64, error) {
	args := m.Called(token)
	return args.Get(0).(int64), args.Error(1)
}

// --- Setup router ---
func setupLoginRouter(h *handler.LoginHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/login", h.Login)
	return r
}

// --- Tests ---
func TestLoginHandler_Success(t *testing.T) {
	userPort := new(mockUserPort)
	sessionPort := new(mockSessionPort)
	jwtService := new(mockJWTService)

	// Password hash de ejemplo
	password := "secret"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Configurar mocks
	userPort.On("GetUserByEmail", "lucas@mail.com").Return(entity.User{
		ID:           1,
		Name:         "Lucas",
		Email:        "lucas@mail.com",
		PasswordHash: string(hashed),
	}, nil)
	jwtService.On("GenerateToken", int64(1)).Return("mytoken", nil)
	sessionPort.On("CreateSession", int64(1), "mytoken").Return(entity.Session{}, nil)

	uc := usecase.NewLoginUseCase(userPort, jwtService, sessionPort)
	h := handler.NewLoginHandler(*uc)

	router := setupLoginRouter(h)

	body := dto.LoginRequest{Email: "lucas@mail.com", Password: password}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Login successful")
	assert.Contains(t, w.Body.String(), "mytoken")

	userPort.AssertExpectations(t)
	sessionPort.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	userPort := new(mockUserPort)
	sessionPort := new(mockSessionPort)
	jwtService := new(mockJWTService)

	userPort.On("GetUserByEmail", "lucas@mail.com").Return(entity.User{
		ID:           1,
		Name:         "Lucas",
		Email:        "lucas@mail.com",
		PasswordHash: "$2a$12$invalidhash", // no coincide
	}, nil)

	uc := usecase.NewLoginUseCase(userPort, jwtService, sessionPort)
	h := handler.NewLoginHandler(*uc)
	router := setupLoginRouter(h)

	body := dto.LoginRequest{Email: "lucas@mail.com", Password: "wrongpw"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
	userPort.AssertExpectations(t)
}
