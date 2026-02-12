package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) ExistByID(userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) FindByUserID(userID string) (*models.User, error) {
	var user models.User
	// 아이디로 유저 한 명을 찾음
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	return &user, err
}
