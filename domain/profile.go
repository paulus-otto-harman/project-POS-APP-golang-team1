package domain

type Profile struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" swaggerignore:"true" json:"id"`
	FullName string `json:"full_name" form:"full_name" binding:"required" example:"Update Name"`
	Email    string `json:"email" form:"email" binding:"required" example:"update@mail.com"`
	Address  string `json:"address" form:"address" binding:"required" example:"Update Address"`
	Password string `json:"password" form:"password" binding:"omitempty,min=5"`
}
