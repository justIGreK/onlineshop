package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userIdCtx           = "userId"
	userRoleCtx         = "userRole"
	parts               = 2
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != parts {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userId, userRole, err := h.Auth.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := h.User.GetUserById(userId)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if !user.IsActice {
		newErrorResponse(c, http.StatusUnauthorized, "This token of banned or deleted user")
	}
	c.Set(userIdCtx, userId)
	c.Set(userRoleCtx, userRole)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userIdCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}
	return idInt, nil
}

func getUserRole(c *gin.Context) (string, error) {
	role, ok := c.Get(userRoleCtx)
	if !ok {
		return "", errors.New("user role not found")
	}
	userRole, ok := role.(string)
	if !ok {
		return "", errors.New("user role is of invalid type")
	}
	return userRole, nil
}
