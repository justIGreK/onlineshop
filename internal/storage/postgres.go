package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	userTable        = "users"
	ordersTable      = "orders"
	ordersItemsTable = "orders_items"
	cartTable        = "cart"
	productsTable    = "products"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname= %s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, fmt.Errorf("error during opening dbserver: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("server is not responding after ping: %w", err)
	}
	return db, nil
}
