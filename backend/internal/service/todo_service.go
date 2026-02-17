package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"time"
)

// 비즈니스 로직 처리
type TodoService struct {
	repo *repository.TodoRepository
}

func NewTodoService(repo *repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

// 힌트: 컨트롤러에서 받은 데이터를 바탕으로 검증 후 저장 요청
func (s *TodoService) AddTodo(todo *models.Todo) error {
	if len(todo.Title) < 2 {
		return errors.New("제목은 최소 2자 이상이어야 합니다.")
	}

	if todo.StartDate == nil {
		now := time.Now()
		todo.StartDate = &now
	}
	return s.repo.Create(todo)
}

// 힌트: 유저별 목록 가져오기 로직
func (s *TodoService) GetUserTodos(userID uint, sortBy, filter, keyword string) ([]models.Todo, error) {
	return s.repo.FindByUserID(userID, sortBy, filter, keyword)
}

func (s *TodoService) UpdateTodo(todoID uint, userID uint, data *models.Todo) error {
	// 1. 기존 데이터 조회 (Fetch)
	todo, err := s.repo.FindByID(todoID, userID)
	if err != nil {
		return errors.New("해당 할 일을 찾을 수 없거나 권한이 없습니다")
	}

	// 2. 바뀐 데이터만 덮어쓰기 (이 과정에서 nil 체크가 핵심!)
	if data.Title != "" {
		todo.Title = data.Title
	}
	if data.Content != "" {
		todo.Content = data.Content
	}

	if todo.Status != data.Status {
		todo.Status = data.Status
		if todo.Status {
			now := time.Now()
			todo.CompletedAt = &now
		} else {
			todo.CompletedAt = nil
		}
	}

	// 포인터 타입인 날짜들은 보낸 데이터가 있을 때만 덮어씀
	if data.StartDate != nil {
		todo.StartDate = data.StartDate
	}
	if data.DueDate != nil {
		todo.DueDate = data.DueDate
	}

	if todo.DueDate != nil && todo.StartDate != nil {
		if todo.DueDate.Before(*todo.StartDate) {
			return errors.New("마감일은 시작일보다 빠를 수 없습니다")
		}
	}

	return s.repo.Save(todo)
}

func (s *TodoService) DeleteTodo(todoID uint, userID uint) error {
	return s.repo.Delete(todoID, userID)
}
