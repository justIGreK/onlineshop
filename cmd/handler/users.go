package handler

import (
	"net/http"
	"onlineshop/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type getUserListResponse struct {
	Data []models.User `json:"data"`
}

type changeBalance struct {
	Balance int `json:"balance"`
}

type deleteAccount struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Get Userlist
// @Security BearerAuth
// @Tags users
// @Description get list of users
// @Param X-Auth-Token header string true "Authentication Token"
// @Produce  json
// @Router /api/users/ [get]
func (h *Handler) getUserList(c *gin.Context) {
	users, err := h.services.UserList.GetUsersList()
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
// @Router /api/users/:id [get]
func (h *Handler) getUser(c *gin.Context) {
	searchId := c.Param("id")
	id, err := strconv.Atoi(searchId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.services.UserList.GetUserById(id)
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
// @Param balance body changeBalance true "NewBalance"
// @Accept  json
// @Produce  json
// @Router /api/users/:id [put]
func (h *Handler) changeBalance(c *gin.Context) {
	var changeBalance changeBalance
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

	if err := c.BindJSON(&changeBalance); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.UserList.ChangeBalance(id, changeBalance.Balance)
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
// @Param data body deleteAccount true "User data"
// @Accept  json
// @Produce  json
// @Router /api/users/:id [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	var userInfo deleteAccount
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.BindJSON(&userInfo); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.UserList.DeleteAccount(userId, userInfo.Login, userInfo.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})

}
