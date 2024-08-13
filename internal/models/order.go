package models

import "time"

type Order struct {
	Id          int       `json:"id" db:"id"`
	UserId      int       `json:"user_id" db:"user_id" binding:"required"`
	PriceBefore float64   `json:"price_before" db:"price_before" binding:"required"`
	PriceAfter  float64   `json:"price_after" db:"price_after" binding:"required"`
	Discount    int       `json:"discount" db:"discount"`
	PaidAt      time.Time `json:"paid_at" db:"paid_at"`
}

type GetOrder struct {
	Id          int       `json:"id" db:"id"`
	PriceBefore float64   `json:"price_before" db:"price_before" binding:"required"`
	PriceAfter  float64   `json:"price_after" db:"price_after" binding:"required"`
	Discount    int       `json:"discount" db:"discount"`
	PaidAt      time.Time `json:"paid_at" db:"paid_at"`
}

type OrderItems struct {
	ProductId int `json:"product_id" db:"product_id" binding:"required"`
	Quantity  int `json:"quantity" db:"quantity" binding:"required"`
	TotalCost int `json:"total_cost" db:"total_cost"`
}
