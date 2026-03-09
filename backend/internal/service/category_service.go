package service

import (
	"backend/internal/models"
	"backend/internal/repository"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) Create(userID uint, name string) error {
	category := models.Category{
		Name:   name,
		UserID: userID,
	}
	return s.repo.Create(&category)
}

func (s *CategoryService) ListCategories(userID uint) ([]models.Category, error) {
	return s.repo.FindAllByUserID(userID)
}
