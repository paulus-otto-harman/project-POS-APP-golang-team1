package domain

import (
	"gorm.io/gorm"
)

func (v *Product) BeforeSave(tx *gorm.DB) (err error) {
	v.Availability = "In Stock"

	if v.Stock == 0 {
		v.Availability = "Out Of Stock"
		return nil
	}

	if v.Stock <= 5 {
		v.Availability = "Low Stock"
		return nil
	}

	return nil
}
