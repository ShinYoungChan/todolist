package handlers

import (
	"backend/internal/models"
	"backend/internal/response"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Signup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 서비스 호출
	if err := h.service.SignupUser(user); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	responseData := models.UserResponse{
		ID:     user.ID,
		UserID: user.UserID,
		Name:   user.Name,
	}

	response.Created(c, responseData)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.LoginUser(req.UserID, req.Password)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, refreshToken, err := h.service.GenerateToken(user.ID)
	if err != nil {
		response.Error(c, 500, "토큰 생성 실패")
		return
	}

	// gin.H 부분 구조체로 빼서 보내는게 더 효율적임 / 컴파일 타임 에러 잡기도 수월
	response.Success(c, gin.H{
		"message":       "로그인 성공",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user": gin.H{
			"id":      user.ID,
			"user_id": user.UserID,
			"name":    user.Name,
		},
	})
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.service.ValidateRefreshToken(req.RefreshToken)

	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}
