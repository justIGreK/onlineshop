package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"onlineshop/internal/models"
	"onlineshop/internal/storage"
)

const (
	PercentageBase = 100
	DiscountMax    = 15
)

type CartService struct {
	cartStore    storage.Cart
	productStore storage.Product
	orderStore   storage.Order
	userStore    storage.UserList
}

func NewCartService(store storage.Cart,
	product storage.Product,
	order storage.Order,
	user storage.UserList) *CartService {
	return &CartService{
		cartStore:    store,
		orderStore:   order,
		productStore: product,
		userStore:    user,
	}
}

func (c *CartService) GetCart(id int) ([]models.GetCart, error) {
	var cartItems, empty []models.GetCart
	cartItems, err := c.cartStore.GetCart(id)
	if err != nil {
		return empty, fmt.Errorf("error during getting cart: %w", err)
	}
	return cartItems, nil
}

func (c *CartService) AddProductToCart(user_id int, product_id int, quantity int) error {
	product, err := c.productStore.GetProductById(product_id)
	if err != nil {
		return fmt.Errorf("didnt find such product: %w", err)
	}
	price := product.Cost * float64(quantity)
	cart, err := c.cartStore.GetCartByUserAndProduct(user_id, product_id)
	if err != nil {
		if quantity < 0 {
			return errors.New("by adding new product to cart quantity cant be below zero")
		}
		err = c.cartStore.CreateCart(user_id, product_id, quantity, price)
		if err != nil {
			return fmt.Errorf("adding to cart got problem:%w", err)
		}
	}
	dif := cart.Quantity + quantity
	if dif < 0 {
		return errors.New("you cant make quantity of product below zero")
	}
	if dif == 0 {
		newerr := c.cartStore.DeleteCartByProduct(user_id, product_id)
		if newerr != nil {
			return fmt.Errorf("problem during deleting cart item: %w", err)
		}
		return nil
	}
	if dif > 0 {
		newerr := c.cartStore.UpdateCart(user_id, product_id, quantity, price)
		if newerr != nil {
			return fmt.Errorf("problem during updating cart item: %w", err)
		}
		return nil
	}
	return nil
}

func (c *CartService) MakeOrder(userID int) error {
	cart, err := c.GetCart(userID)
	if err != nil && cart != nil {
		return fmt.Errorf("problem during making order: %w", err)
	}
	if cart == nil {
		return errors.New("your cart is empty")
	}
	var totalPrice float64
	for _, cartItem := range cart {
		totalPrice += cartItem.Price
		product, newerr := c.productStore.GetProductById(cartItem.ProductId)
		if newerr != nil {
			return fmt.Errorf("cant check for amount of produnt in storage:%w", err)
		}
		if cartItem.Quantity > product.Amount {
			return errors.New("the requested quantity of products is greater than the quantity of products in stock")
		}
	}
	var user models.User
	user, err = c.userStore.GetUserById(userID)
	if err != nil {
		return fmt.Errorf("we cant check yout balance because of this problem: %w", err)
	}
	if totalPrice > float64(user.Balance) {
		newerr := fmt.Sprintf("you cannot place this order: your balance: %f, orders price: %f", user.Balance, totalPrice)
		return errors.New(newerr)
	}
	discount, err := c.randomDiscount(DiscountMax)
	if err != nil {
		return fmt.Errorf("cant get a discount: %w", err)
	}
	sale := float64(PercentageBase-discount) / PercentageBase
	err = c.orderStore.CreateOrder(userID, cart, totalPrice, discount, sale)
	if err != nil {
		return fmt.Errorf("error during making order: %w", err)
	}
	orderCost := -(totalPrice * sale)
	err = c.userStore.UpdateUserBalance(user.Id, orderCost)
	if err != nil {
		return fmt.Errorf("cant reduce balance after creating order:%w", err)
	}
	for _, cartItem := range cart {
		newerr := c.productStore.ChangeAmountOfProduct(cartItem.ProductId, -cartItem.Quantity)
		if newerr != nil {
			return fmt.Errorf("cant change amount of product in storage:%w", err)
		}
	}
	err = c.cartStore.ClearCart(user.Id)
	if err != nil {
		return fmt.Errorf("problem during clearing cart:%w", err)
	}
	return nil
}

func (c *CartService) randomDiscount(max int64) (int, error) {
	if max < 0 {
		return 0, errors.New("max cant be below zero")
	}
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		panic(err)
	}
	return int(n.Int64()), nil
}
