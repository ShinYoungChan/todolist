package models

import "time"

type Todo struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Content     string     `gorm:"type:text" json:"content"`
	UserID      uint       `gorm:"not null" json:"user_id"`     // FK: User.ID와 연결
	Status      bool       `gorm:"default:false" json:"status"` // false: 진행중, true: 완료
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"` // 완료 전에는 Null일 수 있으므로 포인터 사용
	StartDate   *time.Time `json:"start_date"`
	DueDate     *time.Time `json:"due_date"` // 기한이 없을 수 있으므로 포인터 사용
}
