package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, secret string) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: secret,
	}
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

func (s *UserService) LoginUser(userID string, password string) (*models.User, error) {
	user, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("아이디 또는 비밀번호가 틀렸습니다.")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("아이디 또는 비밀번호가 틀렸습니다.")
	}

	return user, nil
}

func (s *UserService) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 주입받은 secret 사용
	return token.SignedString([]byte(s.jwtSecret))
}
