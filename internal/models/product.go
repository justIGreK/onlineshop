package models

import "errors"

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name" binding:"required"`
	Cost        float64 `json:"cost" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Amount      int     `json:"amount" binding:"required"`
	IsActice    bool    `json:"is_active" db:"is_active"`
}

type UpdateProduct struct {
	Name        *string  `json:"name"`
	Cost        *float64 `json:"cost"`
	Description *string  `json:"description"`
	Amount      *int     `json:"amount"`
}

func (i UpdateProduct) Validate() error {
	if i.Name == nil && i.Description == nil && i.Cost == nil && i.Amount == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
