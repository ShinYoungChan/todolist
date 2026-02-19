package handlers

import (
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/response"
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
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "인증 정보가 유효하지 않습니다.")
		return
	}

	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	todo.UserID = userID

	if err := h.service.AddTodo(&todo); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(c, gin.H{
		"message":    "작성 완료",
		"title":      todo.Title,
		"content":    todo.Content,
		"start_date": todo.StartDate,
		"due_date":   todo.DueDate,
	})
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
	// 힌트: 미들웨어에서 받은 user_id로 해당 유저의 목록만 요청
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "인증 정보가 유효하지 않습니다.")
		return
	}

	sortBy := c.DefaultQuery("sort", "created_at")
	filter := c.DefaultQuery("filter", "all")
	keyword := c.Query("keyword")

	todos, err := h.service.GetUserTodos(userID, sortBy, filter, keyword)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 객체로 감싸서 보내주면 유연성 증가
	response.Success(c, gin.H{"todos": todos})
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	// 힌트: URL 파라미터에서 Todo ID 추출 (c.Param("id"))
	idStr := c.Param("id")
	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID 형식이 올바르지 않습니다.")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "인증 정보가 유효하지 않습니다.")
		return
	}

	var updateDate models.Todo
	if err := c.ShouldBindJSON(&updateDate); err != nil {
		response.Error(c, http.StatusBadRequest, "잘못된 데이터 형식입니다.")
		return
	}

	if err := h.service.UpdateTodo(uint(todoID), userID, &updateDate); err != nil {
		response.Error(c, http.StatusForbidden, "수정 권한이 없거나 항목을 찾을 수 없습니다.")
		return
	}

	response.SuccessMsg(c, "변경 성공")
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	// 힌트: 삭제할 권한(내 글이 맞는지) 확인 로직 고려
	idStr := c.Param("id")
	todoID, err := strconv.Atoi(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID 형식이 올바르지 않습니다.")
		return
	}

	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "인증 정보가 유효하지 않습니다.")
		return
	}

	if err := h.service.DeleteTodo(uint(todoID), userID); err != nil {
		response.Error(c, http.StatusForbidden, "삭제 권한이 없거나 항목을 찾을 수 없습니다.")
		return
	}

	response.SuccessMsg(c, "삭제 성공")
}
