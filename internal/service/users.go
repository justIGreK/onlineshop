package service

import (
	"errors"
	"fmt"
	"onlineshop/internal/models"
	"onlineshop/internal/storage"
)

type UserService struct {
	store storage.UserList
}

func NewUserService(store storage.UserList) *UserService {
	return &UserService{store: store}
}

func (u *UserService) GetUsersList() ([]models.User, error) {
	var emptylist, userlist []models.User
	userlist, err := u.store.GetUsersList()
	if err != nil {
		return emptylist, fmt.Errorf("problem with getting userlist:%w", err)
	}
	return userlist, nil
}

func (u *UserService) GetUserById(id int) (models.User, error) {
	var emptyuser, user models.User
	user, err := u.store.GetUserById(id)
	if err != nil {
		return emptyuser, fmt.Errorf("error during getting user: %w", err)
	}
	return user, nil
}

func (u *UserService) ChangeBalance(id int, changeBalance float64) error {
	err := u.store.UpdateUserBalance(id, changeBalance)
	if err != nil {
		return fmt.Errorf("the balance has not been updated due to this error: %w", err)
	}
	return nil
}

func (u *UserService) DeleteAccount(id int, login string, password string) error {
	password = generatePasswordHash(password)
	user, err := u.GetUserById(id)
	if err != nil {
		return fmt.Errorf("error during getting user: %w", err)
	}
	if user.Login != login || user.Password != password {
		return errors.New("your password or login is not correct")
	}
	err = u.store.DeleteAccount(id, login, password)
	if err != nil {
		return fmt.Errorf("account cant be deleted due this error: %w", err)
	}
	return nil
}
