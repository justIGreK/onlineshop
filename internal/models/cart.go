package models

type Cart struct {
	Id        int     `json:"id"`
	UserId    int     `json:"user_id" db:"user_id"`
	ProductId int     `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type GetCart struct {
	ProductId int     `json:"product_id" db:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
