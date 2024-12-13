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
