package handler

import (
	"backend-practice/internal/core/usecase"
	"backend-practice/internal/infra/transport/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{
	uc usecase.CreateUserUseCase
}

func NewUserHandler(uc usecase.CreateUserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.uc.CreateUser(req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"id":      user.ID,
		"name":    user.Name,
		"email":   user.Email,
	})
}