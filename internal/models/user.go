package models

type User struct {
	Id       int     `json:"id"`
	Login    string  `json:"login" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Balance  float64 `json:"balance"`
	IsActice bool    `json:"is_active" db:"is_active"`
	Role     string  `json:"role" db:"role"`
}
