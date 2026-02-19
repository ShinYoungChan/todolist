package handlers

import (
	"backend/internal/middleware"
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

func (h *UserHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, "인증 정보가 유효하지 않습니다.")
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "리프레시 토큰이 필요합니다.")
		return
	}
	// 유저 ID 전체 로그아웃
	if err := h.service.RevokeToken(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "로그아웃 처리 중 오류 발생")
		return
	}

	/* 유저 ID 계정 1개만 로그아웃
		err := h.service.RevokeSpecificToken(userID, req.RefreshToken)
	    if err != nil {
	        response.Error(c, http.StatusInternalServerError, "로그아웃 실패")
	        return
	    }
	*/

	response.SuccessMsg(c, "로그아웃 성공")
}
