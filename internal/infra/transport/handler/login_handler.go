package handler

import (
	"backend-practice/internal/core/usecase"
	"backend-practice/internal/infra/transport/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct{
	uc usecase.LoginUseCase
}

func NewLoginHandler(uc usecase.LoginUseCase) *LoginHandler {
	return &LoginHandler{uc: uc}
}

func (h *LoginHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	r, err := h.uc.Login(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Login successful",
		"token":   r.Token,
		"user":    r.User,
	})
}