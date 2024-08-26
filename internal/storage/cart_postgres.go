package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"onlineshop/internal/models"
	"onlineshop/pkg/util/logger"
)

type CartPostgres struct {
	db *sqlx.DB
}

func NewCartPostgres(db *sqlx.DB) *CartPostgres {
	return &CartPostgres{db: db}
}

func (c *CartPostgres) CreateCart(user_id int, product_id int, quantity int, price float64) error {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			logger.Error("Failed to sync logger", zap.Error(syncErr))
		}
	}()
	query := fmt.Sprintf("INSERT INTO %s (user_id, product_id, quantity, price) values ($1, $2, $3, $4)", cartTable)
	_, err = c.db.Exec(query, user_id, product_id, quantity, price)
	if err != nil {
		return fmt.Errorf("error during inserting to db: %w", err)
	}
	logger.Logger.Info("cart is created")
	return nil
}

func (c *CartPostgres) GetCart(user_id int) ([]models.GetCart, error) {
	var cartItems []models.GetCart
	query := fmt.Sprintf("SELECT product_id, quantity, price FROM %s WHERE user_id=$1", cartTable)

	err := c.db.Select(&cartItems, query, user_id)
	if err != nil {
		return nil, fmt.Errorf("error during getting cart: %w", err)
	}
	return cartItems, nil
}

func (c *CartPostgres) GetCartByUserAndProduct(user_id, product_id int) (models.Cart, error) {
	var cart, empty models.Cart
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1 AND product_id=$2", cartTable)
	err := c.db.Get(&cart, query, user_id, product_id)
	if err != nil {
		return empty, fmt.Errorf("error during selecting from db: %w", err)
	}
	return cart, nil
}

func (c *CartPostgres) UpdateCart(userID int, productID int, quantity int, price float64) error {
	query := fmt.Sprintf("UPDATE %s SET quantity=quantity+$1, price=price+$2"+
		" WHERE user_id=$3 AND product_id=$4", cartTable)
	_, err := c.db.Exec(query, quantity, price, userID, productID)
	if err != nil {
		return fmt.Errorf("error during updating: %w", err)
	}
	return nil
}

func (c *CartPostgres) DeleteCartByProduct(userID int, productID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1 AND product_id=$2", cartTable)
	_, err := c.db.Exec(query, userID, productID)
	if err != nil {
		return fmt.Errorf("error during deliting cart by product: %w", err)
	}
	return nil
}

func (c *CartPostgres) ClearCart(userID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", cartTable)
	_, err := c.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error during clearing cart: %w", err)
	}
	return nil
}
