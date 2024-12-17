package database

import (
	"project/domain"

	"gorm.io/gorm"
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
		&domain.User{},
		&domain.Reservation{},
		&domain.Notification{},
		&domain.Category{},
		&domain.Product{},
		&domain.Inventory{},
	)
}

func dropTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&domain.User{},
		&domain.Reservation{},
		&domain.Notification{},
		&domain.Category{},
		&domain.Product{},
		&domain.UserNotification{},
		&domain.Inventory{},
	)
}

func setupJoinTables(db *gorm.DB) error {
	err := db.SetupJoinTable(&domain.User{}, "Notifications", &domain.UserNotification{})
	if err != nil {
		return err
	}
	return nil
}

func createViews(db *gorm.DB) error {
	var err error

	return err
}
