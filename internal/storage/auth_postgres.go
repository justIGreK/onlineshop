package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"onlineshop/internal/models"
	"onlineshop/pkg/util/logger"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(login, password string) (int, error) {
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
	query := fmt.Sprintf("INSERT INTO %s (login, password) values ($1, $2) RETURNING id", userTable)
	row := a.db.QueryRow(query, login, password)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error during inserting into db:%w", err)
	}
	logger.Logger.Info("user is created")
	return id, nil
}

func (a *AuthPostgres) GetUser(login, password string) (models.User, error) {
	var user, empty models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE login=$1 AND password=$2 AND is_active=TRUE", userTable)
	err := a.db.Get(&user, query, login, password)
	if err != nil {
		return empty, fmt.Errorf("error during selecting from db: %w", err)
	}
	return user, nil
}
