package repository

import (
	"log"
	"project/domain"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserPermissionRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewUserPermissionRepository(db *gorm.DB, log *zap.Logger) *UserPermissionRepository {
	return &UserPermissionRepository{db: db, log: log}
}

func (repo UserPermissionRepository) Create(userPermission domain.UserPermission) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&userPermission).Error; err != nil {
			repo.log.Error("Error creating user notification", zap.Error(err))
			return err
		}
		return nil
	})
}

func (repo UserPermissionRepository) Update(user domain.User) error {
	log.Println(user.Permissions)
	return nil
}
