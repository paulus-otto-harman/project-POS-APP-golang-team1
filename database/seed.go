package database

import (
	"fmt"
	"project/domain"
	"project/domain/seeder"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		seeds := dataSeeds()
		for _, seed := range seeds {
			err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(seed).Error
			if err != nil {
				name := reflect.TypeOf(seed).String()
				errorMessage := err.Error()
				return fmt.Errorf("%s seeder fail with %s", name, errorMessage)
			}
		}
		return nil
	})
}

func dataSeeds() []interface{} {
	return []interface{}{
		domain.ReservationSeed(),
		domain.NotificationSeed(),
		seeder.CategorySeed(),
		seeder.ProductSeed(),
		// domain.InventorySeed(),
		seeder.TableSeed(),
		seeder.PaymentMethodSeed(),
		seeder.OrderSeed(),
		seeder.Permission(),
		seeder.User(),
		seeder.PasswordResetTokenSeed(),
	}
}
