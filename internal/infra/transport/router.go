package transport

import (
	"backend-practice/internal/infra/transport/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(h *handler.HealthHandler, uh *handler.UserHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/health", h.Health)
	r.POST("/users", uh.CreateUser)

	return r
}
