package models

import (
	"time"
)

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`        // 유저와 연결
	Token     string    `gorm:"unique;index"` // 실제 리프레시 토큰 문자열
	ExpiresAt time.Time `gorm:"index"`        // 만료 시간
	CreatedAt time.Time
}
