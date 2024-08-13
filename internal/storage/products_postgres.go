package storage

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"onlineshop/internal/models"
)

type ProductsPostgres struct {
	db *sqlx.DB
}

func NewProductsPostgres(db *sqlx.DB) *ProductsPostgres {
	return &ProductsPostgres{db: db}
}

func (p *ProductsPostgres) CreateProduct(prod models.Product) (int, error) {
	var id int
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			logger.Error("Failed to sync logger", zap.Error(syncErr))
		}
	}()
	query := fmt.Sprintf(
		"INSERT INTO %s (name, cost, description, amount)"+
			"values ($1, $2, $3, $4) RETURNING id",
		productsTable)
	row := p.db.QueryRow(query, prod.Name, prod.Cost, prod.Description, prod.Amount)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error during Scanning id fr creating product: %w", err)
	}
	logger.Info("product is created")
	return id, nil
}

func (p *ProductsPostgres) GetAllProducts() ([]models.Product, error) {
	var lists []models.Product
	query := fmt.Sprintf("SELECT * FROM %s WHERE is_active=TRUE", productsTable)
	err := p.db.Select(&lists, query)
	if err != nil {
		return lists, fmt.Errorf("error during getting all products: %w", err)
	}
	return lists, nil
}

func (p *ProductsPostgres) GetProductById(id int) (models.Product, error) {
	var product models.Product
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1 AND is_active=TRUE", productsTable)
	err := p.db.Get(&product, query, id)
	if err != nil {
		return product, fmt.Errorf("error during getting product by id: %w", err)
	}
	return product, nil
}

func (p *ProductsPostgres) CheckForExisting(id int, tableName string) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", tableName)
	err := p.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error during scanning exists:%w", err)
	}
	return exists, nil
}

func (p *ProductsPostgres) DeleteProductById(id int) error {
	var exists bool
	exists, err := p.CheckForExisting(id, productsTable)
	if exists {
		query := fmt.Sprintf("UPDATE %s SET is_active=FALSE WHERE id=$1", productsTable)
		_, err := p.db.Exec(query, id)
		if err != nil {
			return fmt.Errorf("error during deleting product: %w", err)
		}
		return nil
	} else {
		return err
	}
}

func (p *ProductsPostgres) UpdateProduct(id int, product models.UpdateProduct) error {
	var exists bool
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	exists, err = p.CheckForExisting(id, productsTable)
	if err != nil {
		return err
	}
	if exists {
		setValues := make([]string, 0)
		args := make([]interface{}, 0)
		argsId := 1

		if product.Name != nil {
			setValues = append(setValues, fmt.Sprintf("name=$%d", argsId))
			args = append(args, *product.Name)
			argsId++
		}
		if product.Cost != nil {
			setValues = append(setValues, fmt.Sprintf("cost=$%d", argsId))
			args = append(args, *product.Cost)
			argsId++
		}
		if product.Description != nil {
			setValues = append(setValues, fmt.Sprintf("description=$%d", argsId))
			args = append(args, *product.Description)
			argsId++
		}
		if product.Amount != nil {
			setValues = append(setValues, fmt.Sprintf("amount=$%d", argsId))
			args = append(args, *product.Amount)
		}

		setQuery := strings.Join(setValues, ", ")
		query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d", productsTable, setQuery, id)

		logger.Debug("updateQuery: %s", zap.String("query", query))
		_, err = p.db.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("error during updating product: %w", err)
		}
		return nil
	} else {
		return fmt.Errorf("Product with id %d was not found", id)
	}
}

func (p *ProductsPostgres) ChangeAmountOfProduct(id int, amount int) error {
	query := fmt.Sprintf("UPDATE %s SET amount = amount + $1 WHERE id = $2", productsTable)
	_, err := p.db.Exec(query, amount, id)
	if err != nil {
		return fmt.Errorf("error during changing amount of userlist: %w", err)
	}
	return nil
}
