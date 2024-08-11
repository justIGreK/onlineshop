package service

import (
	"fmt"
	"onlineshop/internal/models"
	"onlineshop/internal/storage"
)

type OrderService struct {
	store storage.Order
}

func NewOrderService(store storage.Order) *OrderService {
	return &OrderService{store: store}
}

func (o *OrderService) GetOrderList(id int) ([]models.GetOrder, error) {
	var orders, empty []models.GetOrder
	orders, err := o.store.GetAllOrders(id)
	if err != nil {
		return empty, fmt.Errorf("error during getting cart: %w", err)
	}
	return orders, err
}

func (o *OrderService) GetOrderDetails(userID int, orderID int) (models.GetOrder, []models.OrderItems, error) {
	var order, empty models.GetOrder
	var orderItems, emptyDetails []models.OrderItems
	order, err := o.store.GetOrderDetails(userID, orderID)
	if err != nil {
		return empty, emptyDetails ,fmt.Errorf("problem during getting order by id: %w", err)
	}
	orderItems, err = o.store.GetOrderItems(orderID)
	if err != nil {
		return empty, emptyDetails, fmt.Errorf("problem during getting orderitems: %w", err)
	}

	return order, orderItems, nil
}
