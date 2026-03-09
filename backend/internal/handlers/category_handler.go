package handlers

import (
	"backend/internal/middleware"
	"backend/internal/response"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	// 1. 유저 ID 가져오기
	userID := middleware.GetUserID(c)

	// 2. 바디 바인딩용 구조체
	var req struct {
		Name string `json:"name" binding:"required"`
	}

	// 3. BindJSON... 에러 처리...
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "이름을 입력하세요")
		return
	}

	// 4. 서비스 호출 (userID와 Name을 넘김)
	if err := h.service.Create(userID, req.Name); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessMsg(c, "카테고리 생성 성공")
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID := middleware.GetUserID(c)

	categoryList, err := h.service.ListCategories(userID)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, categoryList)
}
