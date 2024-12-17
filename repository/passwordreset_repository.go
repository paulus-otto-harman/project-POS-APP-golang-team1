package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project/domain"
	"time"
)

type PasswordResetRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewPasswordResetRepository(db *gorm.DB, log *zap.Logger) *PasswordResetRepository {
	return &PasswordResetRepository{db: db, log: log}
}

func (repo PasswordResetRepository) Create(token *domain.PasswordResetToken) error {
	return repo.db.Create(&token).Error
}

func (repo PasswordResetRepository) GetValidToken(token *domain.PasswordResetToken) (*domain.PasswordResetToken, error) {
	result := repo.db.Where("expired_at >= ?", time.Now()).Where("validated_at IS NULL").First(&token)
	return token, result.Error
}

func (repo PasswordResetRepository) Update(token *domain.PasswordResetToken) error {
	return repo.db.Model(&token).Update("validated_at", token.ValidatedAt).Error
}
