package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"onlineshop/internal/models"
	"onlineshop/pkg/util/logger"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (o *OrderPostgres) CreateOrder(
	userID int,
	cart []models.GetCart,
	totalPrice float64,
	discount int,
	sale float64) error {
	var orderID int
	query := fmt.Sprintf("INSERT INTO %s (user_id, price_before, price_after, discount)"+
		" values ($1, $2, $3, $4) RETURNING id", ordersTable)
	row := o.db.QueryRow(query, userID, totalPrice, (totalPrice * sale), discount)
	if err := row.Scan(&orderID); err != nil {
		return fmt.Errorf("error during scanning orderId: %w", err)
	}
	logger.Logger.Info("order is created")
	if err := o.CreateOrderItems(orderID, cart); err != nil {
		return err
	}

	return nil
}

func (o *OrderPostgres) CreateOrderItems(orderID int, cart []models.GetCart) error {
	query := fmt.Sprintf("INSERT INTO %s (order_id, product_id, quantity, total_cost)"+
		" values ($1, $2, $3, $4)", ordersItemsTable)
	for _, cartItems := range cart {
		_, err := o.db.Exec(query, orderID, cartItems.ProductId, cartItems.Quantity, cartItems.Price)
		if err != nil {
			return fmt.Errorf("error during creating orderitems: %w", err)
		}
	}
	logger.Logger.Info("order items is created")
	return nil
}
func (o *OrderPostgres) GetAllOrders(userID int) ([]models.GetOrder, error) {
	var orders []models.GetOrder
	query := fmt.Sprintf("SELECT id, price_before, price_after, discount, paid_at FROM %s WHERE user_id=$1", ordersTable)
	err := o.db.Select(&orders, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error while selecting from db: %w", err)
	}
	return orders, nil
}

func (o *OrderPostgres) GetOrderDetails(userID int, orderID int) (models.GetOrder, error) {
	var order, empty models.GetOrder
	query := fmt.Sprintf(
		"SELECT id, price_before, price_after, discount, paid_at"+
			" FROM %s WHERE user_id=$1 AND id=$2",
		ordersTable)
	err := o.db.Get(&order, query, userID, orderID)
	if err != nil {
		return empty, fmt.Errorf("error while selecting from db: %w", err)
	}
	return order, nil
}

func (o *OrderPostgres) GetOrderItems(orderID int) ([]models.OrderItems, error) {
	var orderItems []models.OrderItems
	query := fmt.Sprintf("SELECT product_id, quantity, total_cost FROM %s WHERE order_id=$1", ordersItemsTable)
	err := o.db.Select(&orderItems, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("error while selecting from db: %w", err)
	}
	return orderItems, nil
}
