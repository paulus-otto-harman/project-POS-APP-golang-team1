package helper

import (
	"gorm.io/gorm"
)

func Paginate(page uint, limit uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if limit > uint(0) {
			offset := (page - 1) * limit
			return db.Offset(int(offset)).Limit(int(limit))
		}
		return db
	}
}

func Sort(field, direction string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if field == "" {
			return db.Order("created_at DESC")
		}

		// Default direction to ASC if invalid
		if direction != "ASC" && direction != "DESC" {
			direction = "ASC"
		}

		return db.Order(field + " " + direction)
	}
}
