package storage

import (
	"fmt"
	"onlineshop/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(login, password string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (login, password) values ($1, $2) RETURNING id", userTable)
	row := a.db.QueryRow(query, login, password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	logrus.Info("user is created")
	return id, nil
}

func (a *AuthPostgres) GetUser(login, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND password=$2 AND is_active=TRUE", userTable)

	err := a.db.Get(&user, query, login, password)
	return user, err
}
