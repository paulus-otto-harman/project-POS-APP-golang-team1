package domain

type Profile struct {
	ID       uint `gorm:"primaryKey;autoIncrement"`
	UserId   uint
	FullName string  `json:"full_name"`
	Phone    string  `json:"phone"`
	Salary   float32 `gorm:"type:money" json:"salary"`
}
