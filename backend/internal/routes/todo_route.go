package routes

import (
	"backend/internal/handlers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupTodoRoutes(r *gin.Engine, h *handlers.TodoHandler, jwtSecret string) {
	todoGroup := r.Group("/todos")
	todoGroup.Use(middleware.AuthMiddleware(jwtSecret))
	{
		todoGroup.POST("", h.CreateTodo)       // 할 일 생성
		todoGroup.GET("", h.GetTodos)          // 내 할 일 목록 조회
		todoGroup.PUT("/:id", h.UpdateTodo)    // 특정 할 일 수정
		todoGroup.DELETE("/:id", h.DeleteTodo) // 특정 할 일 삭제
	}
}
