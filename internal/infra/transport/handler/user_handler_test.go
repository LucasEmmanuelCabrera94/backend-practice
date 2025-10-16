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
)

type inMemoryUserRepo struct {
	users  []entity.User
	nextID int64
}

func (r *inMemoryUserRepo) CreateUser(u entity.User) (entity.User, error) {
	r.nextID++
	u.ID = r.nextID
	r.users = append(r.users, u)
	return u, nil
}

func (r *inMemoryUserRepo) GetUserByEmail(email string) (entity.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return entity.User{}, nil
}

func setupUserRouter(h *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/users", h.CreateUser)
	return r
}

func TestCreateUser_Handler(t *testing.T) {
	repo := &inMemoryUserRepo{}
	uc := usecase.NewCreateUserUseCase(repo)
	h := handler.NewUserHandler(*uc) 

	router := setupUserRouter(h)

	body := dto.CreateUserRequest{Name: "Lucas", Email: "lucas@mail.com", Password: "secretpw"}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User created successfully")
}
