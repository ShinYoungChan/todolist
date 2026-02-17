package repository

import (
	"backend/internal/models"
	"errors"

	"gorm.io/gorm"
)

type TodoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// 힌트: db.Create() 사용
func (r *TodoRepository) Create(todo *models.Todo) error {
	return r.db.Create(todo).Error
}

// 힌트: 특정 유저의 ID로만 필터링 (Where("user_id = ?", userID))
func (r *TodoRepository) FindByUserID(userID uint, sortBy string) ([]models.Todo, error) {
	var todos []models.Todo

	query := r.db.Where("user_id = ?", userID)
	orderQuery := "created_at desc"

	if sortBy == "start_date" {
		orderQuery = "start_date asc"
	} else if sortBy == "due_date" {
		orderQuery = "due_date asc"
	}

	query = query.Order(orderQuery)

	err := query.Find(&todos).Error

	return todos, err
}

// 힌트: 수정/삭제 시에는 해당 Todo의 ID뿐만 아니라 작성자인 userID도 함께 체크해야 보안상 안전함
func (r *TodoRepository) Update(todoID, userID uint, data interface{}) error {
	result := r.db.Model(&models.Todo{}).
		Where("id = ? AND user_id = ?", todoID, userID).
		Select("Title", "Content", "Status", "StartDate", "DueDate").
		Updates(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("수정된 내역이 없습니다.")
	}
	return nil
}
func (r *TodoRepository) Delete(todoID, userID uint) error {
	// 보안을 위해 id와 user_id가 모두 일치하는 데이터만 삭제
	result := r.db.Where("id = ? AND user_id = ?", todoID, userID).Delete(&models.Todo{})

	if result.Error != nil {
		return result.Error
	}
	// 삭제된 행이 0개라면 내 글이 아니거나 이미 없는 글
	if result.RowsAffected == 0 {
		return errors.New("삭제할 데이터를 찾을 수 없습니다")
	}
	return nil
}

// 1. 특정 Todo 하나를 가져오는 함수
func (r *TodoRepository) FindByID(todoID, userID uint) (*models.Todo, error) {
	var todo models.Todo
	err := r.db.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error
	return &todo, err
}

// 2. 전체 구조체를 덮어씌워 저장하는 함수
func (r *TodoRepository) Save(todo *models.Todo) error {
	return r.db.Save(todo).Error
}
