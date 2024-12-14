package domain

type Permission struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"unique" json:"name"`
}

func PermissionSeed() []Permission {
	return []Permission{
		{Name: "Dashboard"},
		{Name: "Reports"},
		{Name: "Inventory"},
		{Name: "Orders"},
		{Name: "Customers"},
		{Name: "Settings"},
	}
}
