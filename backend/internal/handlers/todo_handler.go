package handlers

import (
	"backend/internal/models"
	"backend/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HTTP 요청/응답 처리
type TodoHandler struct {
	service *service.TodoService
}

func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	// 힌트 1: c.Get("user_id")를 통해 미들웨어가 넣어준 유저 식별자 추출
	// 힌트 2: 추출한 user_id를 Todo 모델의 UserID 필드에 할당
	// 힌트 3: c.ShouldBindJSON으로 할 일 내용(Title 등) 받기
	val, exists := c.Get("user_id")
	userID := uint(val.(float64))

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "유저 정보를 찾을 수 없습니다."})
		return
	}

	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo.UserID = userID

	if err := h.service.AddTodo(&todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "작성 완료",
		"title":      todo.Title,
		"content":    todo.Content,
		"start_date": todo.StartDate,
		"due_date":   todo.DueDate,
	})

}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	// 힌트: 미들웨어에서 받은 user_id로 해당 유저의 목록만 요청
	val, exists := c.Get("user_id")
	userID := uint(val.(float64))

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "유저 정보를 찾을 수 없습니다."})
		return
	}

	sortBy := c.DefaultQuery("sort", "created_at")
	filter := c.DefaultQuery("filter", "all")
	keyword := c.Query("keyword")

	todos, err := h.service.GetUserTodos(userID, sortBy, filter, keyword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "리스트 확인",
		"todos":   todos,
	})
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	// 힌트: URL 파라미터에서 Todo ID 추출 (c.Param("id"))
	idStr := c.Param("id")
	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID 형식이 올바르지 않습니다."})
		return
	}

	val, exists := c.Get("user_id")
	userID := uint(val.(float64))

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "유저 정보를 찾을 수 없습니다."})
		return
	}

	var updateDate models.Todo
	if err := c.ShouldBindJSON(&updateDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "잘못된 데이터 형식입니다."})
		return
	}

	if err := h.service.UpdateTodo(uint(todoID), userID, &updateDate); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "수정 권한이 없거나 항목을 찾을 수 없습니다."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "변경 성공"})
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	// 힌트: 삭제할 권한(내 글이 맞는지) 확인 로직 고려
	idStr := c.Param("id")
	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID 형식이 올바르지 않습니다."})
		return
	}

	val, exists := c.Get("user_id")
	userID := uint(val.(float64))

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "유저 정보를 찾을 수 없습니다."})
		return
	}

	if err := h.service.DeleteTodo(uint(todoID), userID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "삭제 권한이 없거나 항목을 찾을 수 없습니다."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "삭제 성공"})
}
