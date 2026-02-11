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
