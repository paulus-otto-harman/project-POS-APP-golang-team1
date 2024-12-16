package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project/domain"
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

func (repo PasswordResetRepository) Update(token *domain.PasswordResetToken) error {
	return repo.db.Save(&token).Error
}
