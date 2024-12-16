package repository

import (
	"errors"
	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewUserRepository(db *gorm.DB, log *zap.Logger) *UserRepository {
	return &UserRepository{db: db, log: log}
}

func (repo UserRepository) Create(user *domain.User) error {
	return repo.db.Create(&user).Error
}

func (repo UserRepository) All(user domain.User) ([]domain.User, error) {
	var users []domain.User
	result := repo.db.Where(user).Find(&users)
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}
	return users, nil
}

func (repo UserRepository) GetByRole(role string) ([]domain.User, error) {
	var users []domain.User
	result := repo.db.Where("role =?", role).Find(&users)
	if result.RowsAffected == 0 {
		return nil, errors.New("users not found")
	}
	return users, nil
}
