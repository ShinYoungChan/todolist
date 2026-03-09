package routes

import (
	"backend/internal/handlers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupCategoryRoutes(r *gin.Engine, h *handlers.CategoryHandler, jwtSecret string) {
	categoryGroup := r.Group("/categories")
	categoryGroup.Use(middleware.AuthMiddleware(jwtSecret))
	{
		categoryGroup.POST("", h.CreateCategory) // 카테고리 생성
		categoryGroup.GET("", h.GetCategories)   // 내 카테고리 목록 조회
	}
}
