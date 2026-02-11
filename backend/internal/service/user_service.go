package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) SignupUser(user models.User) error {
	// 1. 중복 체크 로직
	exist, err := s.repo.ExistByID(user.UserID)
	if err != nil {
		return err
	}

	if exist {
		return errors.New("이미 존재하는 아이디입니다.")
	}

	// 2. 비밀번호 암호화 (bcrypt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	// user data 삽입
	userData := models.User{
		UserID:   user.UserID,
		Password: string(hashedPassword),
		Name:     user.Name,
	}
	// 3. repo.CreateUser 호출
	return s.repo.CreateUser(&userData)
}
