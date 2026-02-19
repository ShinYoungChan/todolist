package routes

import (
	"backend/internal/handlers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine, h *handlers.UserHandler, jwtSecret string) {
	r.POST("/signup", h.Signup)
	r.POST("/login", h.Login)
	r.POST("/refresh", h.RefreshToken)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware(jwtSecret))
	{
		auth.POST("/logout", h.Logout)
	}
}
