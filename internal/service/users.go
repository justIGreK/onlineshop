package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"onlineshop/internal/models"
	"onlineshop/internal/storage"
)

type UserList interface {
	GetUsersList() ([]models.User, error)
	GetUserById(id int) (models.User, error)
	UpdateUserBalance(id int, changeBalance float64) error
	DeleteAccount(id int, login string, password string) error
	AddConnection(userID, serviceID int, service string) error
	GetConnections(userID int) ([]storage.Connection, error)
}
type UserService struct {
	store UserList
}

const (
	joker = "http://localhost:8000/user/check-user"
)

func NewUserService(store UserList) *UserService {
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

type UserData struct {
	Login    string
	Password string
}

func (u *UserService) LinkAccount(id int, login, password, srv string) error {
	var url string
	switch srv {
	case "joker":
		url = joker
	default:
		return fmt.Errorf("no such service")
	}
	conns, err := u.store.GetConnections(id)
	fmt.Println(conns)
	if err != nil {
		return fmt.Errorf("error during gettin connections:%w", err)
	}
	for _, conn := range conns {
		if conn.ServiceName == srv {
			return errors.New("you are already connected with this service")
		}
	}
	jsonData, err := json.Marshal(UserData{Login: login, Password: password})
	if err != nil {
		return fmt.Errorf("failed to marshal account data: %w", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response: %s", resp.Status)
	}
	var result map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	err = u.store.AddConnection(id, result["id"], srv)
	if err != nil {
		return fmt.Errorf("error during linking account:%w", err)
	}

	return nil
}
