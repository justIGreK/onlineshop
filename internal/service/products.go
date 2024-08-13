package service

import (
	"fmt"

	"onlineshop/internal/models"
	"onlineshop/internal/storage"
)

type ProdService struct {
	store storage.Product
}

func NewProdService(store storage.Product) *ProdService {
	return &ProdService{store: store}
}

func (p *ProdService) CreateProduct(name string, cost float64, description string, amount int) (int, error) {
	prod := models.Product{
		Name:        name,
		Cost:        cost,
		Description: description,
		Amount:      amount,
	}
	resp, err := p.store.CreateProduct(prod)
	if err != nil {
		return 0, fmt.Errorf("create product err: %w", err)
	}
	return resp, nil
}

func (p *ProdService) GetProductList() ([]models.Product, error) {
	resp, err := p.store.GetAllProducts()
	if err != nil {
		return nil, fmt.Errorf("get product list err: %w", err)
	}
	return resp, nil
}

func (p *ProdService) GetProductById(id int) (models.Product, error) {
	var product, empty_product models.Product
	product, err := p.store.GetProductById(id)
	if err != nil {
		return empty_product, fmt.Errorf("get product by id problem: %w", err)
	}
	return product, nil
}

func (p *ProdService) DeleteProduct(id int) error {
	err := p.store.DeleteProductById(id)
	if err != nil {
		return fmt.Errorf("error during deleting product:%w", err)
	}
	return nil
}

func (p *ProdService) UpdateProduct(id int, product models.UpdateProduct) error {
	err := product.Validate()
	if err != nil {
		return fmt.Errorf("error during validate:%w", err)
	}
	err = p.store.UpdateProduct(id, product)
	if err != nil {
		return fmt.Errorf("update product got problem:%w", err)
	}
	return nil
}
