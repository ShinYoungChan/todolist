package models

import (
	"mime/multipart"
	"time"
)

type Todo struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Content     string     `gorm:"type:text" json:"content"`
	UserID      uint       `gorm:"not null" json:"user_id"`     // FK: User.ID와 연결
	Status      bool       `gorm:"default:false" json:"status"` // false: 진행중, true: 완료
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"` // 완료 전에는 Null일 수 있으므로 포인터 사용
	StartDate   *time.Time `json:"start_date"`
	DueDate     *time.Time `json:"due_date"`              // 기한이 없을 수 있으므로 포인터 사용
	CategoryID  *uint      `json:"category_id"`           // null 허용을 위해 포인터 사용
	Category    Category   `gorm:"foreignKey:CategoryID"` // Relation 설정
	ImageURL    *string    `json:"image_url"`             // 사진이 없을 수도 있으니 포인터(*)로 null 허용
}

// TodoCreateRequest: 할 일 생성 시 데이터를 받기 위한 구조체
type TodoCreateRequest struct {
	Title      string `form:"title" binding:"required"`
	Content    string `form:"content"`
	CategoryID *uint  `form:"category_id"`
	// 파일 자체를 받기 위한 필드 (multipart.FileHeader)
	Image     *multipart.FileHeader `form:"image"`
	StartDate *time.Time            `form:"start_date" time_format:"2006-01-02"`
	DueDate   *time.Time            `form:"due_date" time_format:"2006-01-02"`
}
