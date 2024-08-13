package storage

import (
	"github.com/jmoiron/sqlx"

	"onlineshop/internal/models"
)

type Authorization interface {
	CreateUser(login, password string) (int, error)
	GetUser(login, password string) (models.User, error)
}

type UserList interface {
	GetUsersList() ([]models.User, error)
	GetUserById(id int) (models.User, error)
	UpdateUserBalance(id int, changeBalance float64) error
	DeleteAccount(id int, login string, password string) error
}

type Cart interface {
	CreateCart(user_id int, product_id int, quantity int, price float64) error
	GetCart(id int) ([]models.GetCart, error)
	GetCartByUserAndProduct(user_id, product_id int) (models.Cart, error)
	UpdateCart(userID int, productID int, quantity int, price float64) error
	DeleteCartByProduct(userID int, productID int) error
	ClearCart(userID int) error
}

type Product interface {
	CreateProduct(prod models.Product) (int, error)
	GetAllProducts() ([]models.Product, error)
	GetProductById(id int) (models.Product, error)
	DeleteProductById(id int) error
	CheckForExisting(id int, tableName string) (bool, error)
	UpdateProduct(id int, product models.UpdateProduct) error
	ChangeAmountOfProduct(id int, amount int) error
}
type Order interface {
	CreateOrder(userID int, cart []models.GetCart, totalPrice float64, discount int, sale float64) error
	CreateOrderItems(orderID int, cart []models.GetCart) error
	GetAllOrders(userID int) ([]models.GetOrder, error)
	GetOrderDetails(userID int, orderID int) (models.GetOrder, error)
	GetOrderItems(orderID int) ([]models.OrderItems, error)
}
type Store struct {
	Authorization
	UserList
	Product
	Cart
	Order
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		Authorization: NewAuthPostgres(db),
		Product:       NewProductsPostgres(db),
		UserList:      NewUsersPostgres(db),
		Cart:          NewCartPostgres(db),
		Order:         NewOrderPostgres(db),
	}
}
