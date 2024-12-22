package domain

type Login struct {
	Email    string `json:"email" example:"admin@mail.com" binding:"required,email"`
	Password string `json:"password" example:"admin" binding:"required,min=5"`
}
