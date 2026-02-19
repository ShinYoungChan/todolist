package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Header에서 Authorization 키의 값을 가져옴 (Bearer <token>)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "인증 헤더가 없습니다."})
			c.Abort() // 중요: 이후 핸들러 실행을 중단함
			return
		}

		// 2. "Bearer " 문자열 제거 후 토큰만 추출
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. 토큰 파싱 및 검증
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		// 4. 검증 실패 처리
		if err != nil {
			// 에러의 종류를 확인합니다.
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":  "TOKEN_EXPIRED", // 프론트가 읽을 코드
					"error": "토큰이 만료되었습니다.",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":  "INVALID_TOKEN",
					"error": "유효하지 않은 토큰입니다.",
				})
			}
			c.Abort()
			return
		}

		// 5. 토큰에서 user_id 꺼내기 (로그인할 때 넣었던 값)
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			// 꺼낸 user_id를 Context에 저장 (나중에 핸들러에서 사용)
			c.Set("user_id", claims["user_id"])
		}

		c.Next() // 다음 단계(실제 핸들러)로 진행
	}
}

func GetUserID(c *gin.Context) uint {
	val, exists := c.Get("user_id")

	if !exists {
		return 0
	}
	// 어떤 숫자 타입이든 uint로 안전하게 변환하는 스킬
	switch v := val.(type) {
	case uint:
		return v
	case float64:
		return uint(v)
	case int:
		return uint(v)
	default:
		return 0
	}
}
