package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"onlineshop/internal/models"
)

type UserList interface {
	GetUsersList() ([]models.User, error)
	GetUserById(id int) (models.User, error)
	ChangeBalance(id int, changeBalance float64) error
	DeleteAccount(id int, login string, password string) error
	LinkAccount(id int, login, password, srv string) error
}

type getUserListResponse struct {
	Data []models.User `json:"data"`
}

type deleteAccount struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Get Userlist
// @Security BearerAuth
// @Tags users
// @Description get list of users
// @Produce  json
// @Router /api/users/ [get]
func (h *Handler) getUserList(c *gin.Context) {
	userRole, err := getUserRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if userRole != AdministratorRole {
		newErrorResponse(c, http.StatusUnauthorized, "You have no permission to do this")
		return
	}
	users, err := h.User.GetUsersList()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getUserListResponse{
		Data: users,
	})
}

// @Summary Get User
// @Security BearerAuth
// @Tags users
// @Description get user by id
// @Param id path int  true  "Account ID"
// @Produce  json
// @Router /api/users/{id} [get]
func (h *Handler) getUser(c *gin.Context) {
	searchId := c.Param("id")
	fmt.Println(searchId)
	id, err := strconv.Atoi(searchId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.User.GetUserById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, models.User{
		Id:       user.Id,
		Login:    user.Login,
		Password: user.Password,
		Balance:  user.Balance,
	})
}

// @Summary Update user
// @Security BearerAuth
// @Tags users
// @Description update balance of user
// @Param id path int  true  "Account ID"
// @Param balance query float64 true "how many money do you want add/receive(type with -)"
// @Accept  json
// @Produce  json
// @Router /api/users/{id} [put]
func (h *Handler) changeBalance(c *gin.Context) {
	userid, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	searchId := c.Param("id")
	id, err := strconv.Atoi(searchId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userid != id {
		newErrorResponse(c, http.StatusBadRequest, "Your token id is not equel with requested id")
		return
	}

	strBalance := c.Query("balance")
	balance, err := strconv.ParseFloat(strBalance, 64)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid balance")
		return
	}

	err = h.User.ChangeBalance(id, balance)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary Delete user
// @Security BearerAuth
// @Tags users
// @Description delete user or change user acc to inactive
// @Param login query string true "your login"
// @Param password query string true "your password"
// @Accept  json
// @Produce  json
// @Router /api/users/ [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	userInfo := deleteAccount{
		Login:    c.Query("login"),
		Password: c.Query("password"),
	}
	err = h.User.DeleteAccount(userId, userInfo.Login, userInfo.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

// @Summary link account
// @Security BearerAuth
// @Tags users
// @Description link your account with another service
// @Param service path string true "service"
// @Param login query string true "your login"
// @Param password query string true "your password"
// @Accept  json
// @Produce  json
// @Router /api/users/link/{service} [post]
func (h *Handler) linkAcc(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	srv := c.Param("service")
	userInfo := SignIn{
		Login:    c.Query("login"),
		Password: c.Query("password"),
	}
	if err := c.ShouldBind(&userInfo); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.User.LinkAccount(userId, userInfo.Login, userInfo.Password, srv)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{
		Status: "service was connected",
	})
}
