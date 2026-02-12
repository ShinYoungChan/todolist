package handlers

import (
	"backend/internal/models"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 서비스 호출
	if err := h.service.SignupUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.UserResponse{
		ID:     user.ID,
		UserID: user.UserID,
		Name:   user.Name,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "회원가입 성공",
		"response": response,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.LoginUser(req.UserID, req.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.GenerateToken(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "로그인 성공",
		"token":   token,
		"user": gin.H{
			"id":      user.ID,
			"user_id": user.UserID,
			"name":    user.Name,
		},
	})
}
