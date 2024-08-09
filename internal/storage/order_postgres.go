package storage

import (
	"onlineshop/internal/models"

	"github.com/jmoiron/sqlx"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (o *OrderPostgres) MakeOrder(userID int, cart []models.GetCart, totalPrice float64, discount int) error{
	// query := "INSERT INTO %w (user_id, price_before, price_after) values ($1, $2, $3)"
	return nil
}

