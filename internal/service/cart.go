package service

import (
	"errors"
	"fmt"
	"math/rand"
	"onlineshop/internal/models"
	"onlineshop/internal/storage"
)

type CartService struct {
	cartStore    storage.Cart
	productStore storage.Product
	orderStore   storage.Order
	userStore    storage.UserList
}

func NewCartService(store storage.Cart, product storage.Product) *CartService {
	return &CartService{
		cartStore:    store,
		productStore: product,
	}
}

// func (c *CartService) CreateCart()(error){
// 	err := c.store.CreateCart()
// 	if err != nil{
// 		fmt.Errorf("problem during creating cart:%w", err )
// 	}
// 	return nil
// }

func (c *CartService) GetCart(id int) ([]models.GetCart, error) {
	var cartItems []models.GetCart
	cartItems, err := c.cartStore.GetCart(id)
	return cartItems, err

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
		return fmt.Errorf("problem duting making order: %w", err)
	}
	if cart == nil {
		return errors.New("your cart is empty")
	}

	var totalPrice float64
	for _, cartItem := range cart {
		totalPrice = totalPrice + cartItem.Price
		product, newerr := c.productStore.GetProductById(cartItem.ProductId)
		if newerr != nil {
			return fmt.Errorf("cant check for amount of produnt in storage:%w", err)
		}
		if cartItem.Quantity > product.Amount{
			return errors.New("the requested quantity of products is greater than the quantity of products in stock")
		}
	}

	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAaaa")
	user, err := c.userStore.GetUserById(userID)
	fmt.Println("SSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSs")
	if err != nil {
		return fmt.Errorf("we cant check yout balance because of this problem: %w", err)
	}
	if totalPrice > float64(user.Balance) {
		newerr := fmt.Sprintf("you cannot place this order: your balance: %d, orders price: %f", user.Balance, totalPrice)
		return errors.New(newerr)
	}

	err = c.orderStore.CreateOrder(userID, cart, totalPrice, c.RandomDiscount())
	if err !=nil {
		return fmt.Errorf("error during making order: %w", err)
	}
	return nil
}

func (c *CartService) RandomDiscount() int {
	return rand.Intn(15)	
}
