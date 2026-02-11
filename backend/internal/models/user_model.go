package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"type:varchar(100);unique;not null" json:"user_id" binding:"required"` // 로그인용 ID (중복불가)
	Password  string    `gorm:"type:varchar(255);not null" json:"password" binding:"required"`       // 비밀번호 (JSON 응답시 보안상 제외)
	Name      string    `gorm:"type:varchar(100); not null" json:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	// 관계 설정: 한 명의 유저는 여러 개의 할 일을 가질 수 있음
	Todos []Todo `gorm:"foreignKey:UserID" json:"todos,omitempty"`
}

type UserResponse struct {
	ID     uint   `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}
