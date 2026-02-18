package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 공통 응답 구조체
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// 1. 가장 일반적인 성공 (200 OK)
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Success: true,
		Data:    data,
	})
}

// 1-1. 성공 응답 (단순 메시지만 보낼 때)
func SuccessMsg(c *gin.Context, message string) {
	c.JSON(http.StatusOK, ApiResponse{
		Success: true,
		Message: message,
	})
}

// 2. 리소스 생성 성공 (201 Created) - 회원가입, 할 일 추가 등
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, ApiResponse{
		Success: true,
		Data:    data,
	})
}

// 3. 커스텀 상태 코드가 필요한 경우 (유연성 확보)
func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, ApiResponse{
		Success: true,
		Data:    data,
	})
}

// 4. 에러 응답
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, ApiResponse{
		Success: false,
		Message: message,
	})
}
