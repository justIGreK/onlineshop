package service

import "onlineshop/internal/storage"

type OrderService struct {
	store storage.Order
}

func NewOrderService(store storage.Order) *OrderService {
	return &OrderService{store: store}
}
