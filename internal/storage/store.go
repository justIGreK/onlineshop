package storage

import "github.com/jmoiron/sqlx"

type Store struct {
	Authorization *AuthPostgres
	UserList *UsersPostgres
	Product *ProductsPostgres
	Cart *CartPostgres
	Order *OrderPostgres
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
