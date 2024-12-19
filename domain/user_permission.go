package domain

type UserPermission struct {
	ID           uint `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint
	PermissionID uint
}
