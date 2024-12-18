package helper

import (
	"gorm.io/gorm"
)

func Paginate(page uint, limit uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(int(offset)).Limit(int(limit))
	}
}

func Sort(field, direction string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" {
			return db.Order("full_name DESC")
		}

		// Default direction to ASC if invalid
		if direction != "ASC" && direction != "DESC" {
			direction = "ASC"
		}

		return db.Order(field + " " + direction)
	}
}
