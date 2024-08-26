package handler

import (
	"github.com/gin-gonic/gin"

	"onlineshop/pkg/util/logger"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logger.Logger.Info(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
