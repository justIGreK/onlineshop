package service

import (
	"onlineshop/internal/models"
	"onlineshop/internal/storage"
)

type Authorization interface {
	CreateUser(login string, password string) (int, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (int, error)
}

type UserList interface {
	GetUsersList() ([]models.User, error)
	GetUserById(id int) (models.User, error)
	ChangeBalance(id int, changeBalance float64) error
	DeleteAccount(id int, login string, password string) error
}

type Product interface {
	CreateProduct(name string, cost float64, description string, amount int) (int, error)
	GetProductList() ([]models.Product, error)
	GetProductById(id int) (models.Product, error)
	DeleteProduct(id int) error
	UpdateProduct(id int, product models.UpdateProduct) error
}

type Cart interface {
	GetCart(id int) ([]models.GetCart, error)
	AddProductToCart(user_id int, product_id int, quantity int) error
	MakeOrder(user_id int) error
	RandomDiscount() int
}

type Order interface {
}

type Service struct {
	Authorization
	UserList
	Product
	Cart
	Order
}


func NewService(store *storage.Store) *Service {
	return &Service{
		Authorization: NewAuthService(store.Authorization),
		Product:       NewProdService(store.Product),
		UserList:      NewUserService(store.UserList),
		Cart:          NewCartService(store.Cart, store.Product, store.Order, store.UserList),
		Order:         NewOrderService(store.Order),
	}
}
