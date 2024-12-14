package repository

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project/database"
	"project/domain"
)

type AuthRepository struct {
	db     *gorm.DB
	cacher database.Cacher
	log    *zap.Logger
}

func NewAuthRepository(db *gorm.DB, cacher database.Cacher, log *zap.Logger) *AuthRepository {
	return &AuthRepository{db: db, cacher: cacher, log: log}
}

func (repo AuthRepository) Authenticate(user domain.Login) (string, bool, error) {
	if err := repo.db.Where(user).First(&user).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return "", false, errors.New("invalid username and/or password")
	}

	return "", true, nil

}
