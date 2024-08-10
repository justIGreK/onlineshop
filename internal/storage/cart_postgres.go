package storage

import (
	"fmt"
	"onlineshop/internal/models"

	"github.com/jmoiron/sqlx"
)

type CartPostgres struct {
	db *sqlx.DB
}

func NewCartPostgres(db *sqlx.DB) *CartPostgres {
	return &CartPostgres{db: db}
}

func (c *CartPostgres) CreateCart(user_id int, product_id int, quantity int, price float64) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, product_id, quantity, price) values ($1, $2, $3, $4)", cartTable)
	_, err := c.db.Exec(query, user_id, product_id, quantity, price)
	if err != nil {
		return err
	}
	return err
}

func (c *CartPostgres) GetCart(user_id int) ([]models.GetCart, error) {
	var cartItems []models.GetCart
	query := fmt.Sprintf("SELECT product_id, quantity, price FROM %s WHERE user_id=$1", cartTable)

	err := c.db.Select(&cartItems, query, user_id)

	return cartItems, err
}


// func (c *CartPostgres) CheckForProductInCart(user_id, product_id int) (bool, error){
// 	query := fmt.Sprintf("SELECT user_id, product_id FROM %s WHERE user_id=$1 AND product_id=$2", cartTable)
// 	rows, err := c.db.Query(query, user_id, product_id)
// 	if err != nil {
// 		return false, err
// 	}
// 	if !rows.Next() {
// 		return false, nil
// 	}
// 	return true, nil
// }

func (c *CartPostgres) GetCartByUserAndProduct(user_id, product_id int) (models.Cart, error){
	var cart models.Cart
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1 AND product_id=$2", cartTable)
	err := c.db.Get(&cart, query, user_id, product_id)
	return cart, err
}


func (c *CartPostgres) UpdateCart(userID int, productID int, quantity int, price float64) error {
	query := fmt.Sprintf("UPDATE %s SET quantity=quantity+$1, price=price+$2 WHERE user_id=$3 AND product_id=$4", cartTable)
	_, err := c.db.Exec(query, quantity, price, userID, productID)
	return err 
}

func (c *CartPostgres) DeleteCartByProduct(userID int, productID int)(error){
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1 AND product_id=$2", cartTable)
	_, err := c.db.Exec(query, userID, productID)
	return err 
	
}

func (c *CartPostgres) ClearCart(userID int)error{
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1", cartTable)
	_, err := c.db.Exec(query, userID)
	return err
}

