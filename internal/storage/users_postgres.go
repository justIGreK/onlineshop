package storage

import (
	"errors"
	"fmt"
	"onlineshop/internal/models"

	"github.com/jmoiron/sqlx"
)

type UsersPostgres struct {
	db *sqlx.DB
}

func NewUsersPostgres(db *sqlx.DB) *UsersPostgres {
	return &UsersPostgres{db: db}
}

func (u *UsersPostgres) GetUsersList() ([]models.User, error) {
	var lists []models.User
	query := fmt.Sprintf("SELECT * FROM %s", userTable)
	err := u.db.Select(&lists, query)
	if err != nil {
		return lists, fmt.Errorf("error during getting userlist: %w", err)
	}
	return lists, nil
}

func (u *UsersPostgres) GetUserById(id int) (models.User, error) {
	var user, empty models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", userTable)
	err := u.db.Get(&user, query, id)
	if err != nil {
		return empty, fmt.Errorf("error during getting user by id:%w", err)
	}
	return user, nil
}

func (u *UsersPostgres) UpdateUserBalance(id int, changeBalance float64) error {
	isOk, err := u.CheckForValidUpdateBalance(id, changeBalance)
	if err != nil {
		return err
	}
	if !isOk {
		return errors.New("you cant make balance below zero")
	}
	query := fmt.Sprintf("UPDATE %s SET balance=balance+$1 WHERE id=$2", userTable)

	_, err = u.db.Exec(query, changeBalance, id)
	if err != nil {
		return fmt.Errorf("error during updating user balance: %w", err)
	}
	return nil
}

func (u *UsersPostgres) CheckForValidUpdateBalance(id int, changeBalance float64) (bool, error) {
	var currentBalance float64
	query := fmt.Sprintf("SELECT balance FROM %s WHERE id=$1", userTable)
	err := u.db.Get(&currentBalance, query, id)
	if err != nil {
		return false, fmt.Errorf("it was not possible to get the balance status for the reason: %w", err)
	}
	isOk := true
	if currentBalance+changeBalance < 0 {
		isOk = false
	}
	return isOk, nil
}

func (u *UsersPostgres) DeleteAccount(id int, login string, password string) error {
	query := fmt.Sprintf("UPDATE %s SET is_active=FALSE WHERE id=$1", userTable)
	_, err := u.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("delete account has failed: %w", err)
	}
	return nil
}
