package storage

import (
	"fmt"
	"onlineshop/internal/models"

	"github.com/jmoiron/sqlx"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (o *OrderPostgres) CreateOrder(userID int, cart []models.GetCart, totalPrice float64, discount int) error{
	var orderID int
	sale := float64(100-discount)/100
	query := fmt.Sprintf("INSERT INTO %s (user_id, price_before, price_after, discount) values ($1, $2, $3, $4) RETURNING id", ordersTable)
	row := o.db.QueryRow(query, userID, totalPrice, (totalPrice * sale), discount)
	if err := row.Scan(&orderID); err != nil{
		return err
	}
	if err := o.CreateOrderItems(orderID, cart); err !=nil{
		return err 
	}
	return nil
}

func (o *OrderPostgres) CreateOrderItems(orderID int, cart []models.GetCart) error{
	query := fmt.Sprintf("INSERT INTO %s (order_id, product_id, quantity, total_cost) values ($1, $2, $3, $4)", ordersItemsTable)

	for _, cartItems := range cart{
		_, err := o.db.Exec(query, orderID, cartItems.ProductId, cartItems.Quantity, cartItems.Price )	
		if err != nil{
			return err
		}
	}
	return nil 
}
