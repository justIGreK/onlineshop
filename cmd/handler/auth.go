package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Authorization interface {
	CreateUser(login string, password string, email string) (int, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (int, string, error)
}

type SignIn struct {
	Login    string `json:"login" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,min=6"`
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @Accept  json
// @Produce  json
// @Param login query string true "your login"
// @Param password query string true "your password"
// @Param email query string true "your email"
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	input := SignIn{
		Login:    c.Query("login"),
		Password: c.Query("password"),
		Email:    c.Query("email"),
	}

	if err := c.ShouldBind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.Auth.CreateUser(input.Login, input.Password, input.Email)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary SignIn
// @Tags auth
// @Description login
// @Accept  json
// @Produce  json
// @Param login query string true "your login"
// @Param password query string true "your password"
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	input := SignIn{
		Login:    c.Query("login"),
		Password: c.Query("password"),
	}
	fmt.Println(input)

	token, err := h.Auth.GenerateToken(input.Login, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
