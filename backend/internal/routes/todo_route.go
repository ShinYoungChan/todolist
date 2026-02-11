package routes

import (
	"backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, h *handlers.UserHandler) {
	r.POST("/signup", h.Signup)
}
