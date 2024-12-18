package repository

import (
	"errors"
	"project/domain"
	"project/helper"

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
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			repo.log.Error("Error creating user", zap.Error(err))
			return err
		}
		// Return nil to indicate success
		return nil
	})

}

func (repo UserRepository) All(sortField, sortDirection string, page, limit uint) ([]domain.User, int64, error) {
	var users []domain.User
	var count int64

	// Count total users
	err := repo.db.Model(&domain.User{}).Count(&count).Error
	if err != nil {
		repo.log.Error("Error counting users", zap.Error(err))
		return nil, 0, err
	}

	// Query with dynamic age calculation
	err = repo.db.Scopes(
		helper.Paginate(page, limit),
		helper.Sort(sortField, sortDirection),
	).Select("*, EXTRACT(YEAR FROM AGE(CURRENT_DATE, birth_date)) AS Age").
		Find(&users).Error

	if err != nil {
		repo.log.Error("Error fetching users", zap.Error(err))
		return nil, count, err
	}

	return users, count, nil
}

func (repo UserRepository) Get(criteria domain.User) (*domain.User, error) {
	var user domain.User
	err := repo.db.Preload("Permissions").Where(criteria).First(&user).Error
	return &user, err
}

func (repo UserRepository) GetByRole(role string) ([]domain.User, error) {
	var users []domain.User
	result := repo.db.Where("role =?", role).Find(&users)
	if result.RowsAffected == 0 {
		return nil, errors.New("users not found")
	}
	return users, nil
}

func (repo UserRepository) GetByEmail(email string) *domain.User {
	var user domain.User
	result := repo.db.Where("email =?", email).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		return nil
	}
	return &user
}

func (repo UserRepository) Update(user *domain.User) error {
	return repo.db.Save(user).Error
}
