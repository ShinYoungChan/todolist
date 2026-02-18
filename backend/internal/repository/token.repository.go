package repository

import (
	"backend/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// 토큰 저장
func (r *TokenRepository) SaveRefreshToken(userID uint, token string, expiresAt time.Time) error {
	// 기존 토큰이 있다면 지우고 새로 생성 (혹은 업데이트)
	// 여기서는 간단하게 기존 유저의 토큰을 모두 지우고 새로 생성하는 방식을 예로 듭니다.
	r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{})

	userToken := models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	return r.db.Create(&userToken).Error
}

// 토큰으로 실제 있는 값인지 찾는 함수
func (r *TokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken

	query := r.db.Where("token = ?", token)

	// Find 쓰면 객체가 없더라도 빈 객체를 리턴할 수 있어서 First로 써야함
	err := query.First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, err
}

// 삭제
func (r *TokenRepository) DeleteByUserID(userID uint) error {
	result := r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("삭제할 데이터를 찾을 수 없습니다")
	}
	return nil
}
