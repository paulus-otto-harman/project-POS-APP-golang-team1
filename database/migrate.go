package database

import (
	"gorm.io/gorm"
	"project/domain"
)

func Migrate(db *gorm.DB) error {
	var err error

	if err = dropTables(db); err != nil {
		return err
	}

	if err = setupJoinTables(db); err != nil {
		return err
	}

	if err = autoMigrates(db); err != nil {
		return err
	}

	return createViews(db)
}

func autoMigrates(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Permission{},
		&domain.User{},
		&domain.Notification{},
		&domain.Category{},
		&domain.Product{},
		&domain.Profile{},
		&domain.PasswordResetToken{},
	)
}

func dropTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&domain.PasswordResetToken{},
		&domain.Profile{},
		&domain.User{},
		&domain.Notification{},
		&domain.Category{},
		&domain.Product{},
		&domain.Permission{},
		"user_permissions",
	)
}

func setupJoinTables(db *gorm.DB) error {
	var err error
	if err = db.SetupJoinTable(&domain.User{}, "Permissions", &domain.UserPermission{}); err != nil {
		return err
	}
	return err
}

func createViews(db *gorm.DB) error {
	var err error

	return err
}
