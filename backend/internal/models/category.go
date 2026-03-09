package models

type Category struct {
	ID     uint   `gorm:"primaryKey"`
	Name   string `gorm:"not null"`
	UserID uint   `gorm:"index"`
}
