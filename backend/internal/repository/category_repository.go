package repository

import (
	"backend/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) FindAllByUserID(userID uint) ([]models.Category, error) {
	var categoryList []models.Category

	err := r.db.Where("user_id = ?", userID).Find(&categoryList).Error

	return categoryList, err
}
