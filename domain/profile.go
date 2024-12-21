package domain

type Profile struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	FullName string `json:"full_name" form:"full_name" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Address  string `json:"address" form:"address" binding:"required"`
	Password string `json:"password" form:"password" binding:"omitempty,min=5"`
}
