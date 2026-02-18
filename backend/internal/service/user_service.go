package service

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepository
	tokenRepo *repository.TokenRepository
	jwtSecret string
}

func NewUserService(repo *repository.UserRepository, tRepo *repository.TokenRepository, secret string) *UserService {
	return &UserService{
		repo:      repo,
		tokenRepo: tRepo,
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

func (s *UserService) GenerateToken(userID uint) (string, string, error) {

	accessToken, err := s.createAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.createAndSaveRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *UserService) createAccessToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 30).Unix(), // 짧은 수명
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *UserService) createAndSaveRefreshToken(userID uint) (string, error) {
	// 랜덤 문자열 생성
	refreshTokenStr := s.generateRandomString(64)

	// 유효 기간 설정 (7일)
	expiresAt := time.Now().Add(time.Hour * 24 * 7)

	// 레포지토리를 통해 DB 저장
	err := s.tokenRepo.SaveRefreshToken(userID, refreshTokenStr, expiresAt)
	if err != nil {
		return "", err
	}

	return refreshTokenStr, nil
}

// 랜덤 문자열 생성 함수
func (s *UserService) generateRandomString(n int) string {
	// 랜덤 시드 설정 (매번 다른 값이 나오게 하기 위함)
	rand.Seed(time.Now().UnixNano())

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *UserService) ValidateRefreshToken(token string) (string, string, error) {
	rt, err := s.tokenRepo.FindByToken(token)
	if err != nil {
		return "", "", errors.New("유효하지 않은 토큰입니다.")
	}

	// 만료시간 체크
	if time.Now().After(rt.ExpiresAt) {
		// 만료되었다면  DB에서도 삭제
		s.tokenRepo.DeleteByUserID(rt.UserID)
		return "", "", errors.New("토큰이 만료되었습니다.")
	}

	newAccessToken, err := s.createAccessToken(rt.UserID)

	newRefreshToken, err := s.createAndSaveRefreshToken(rt.UserID)

	if err != nil {
		return "", "", errors.New("토큰 생성 실패")
	}

	return newAccessToken, newRefreshToken, err
}
