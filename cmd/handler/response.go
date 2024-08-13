package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() {
		if syncErr := logger.Sync(); syncErr != nil {
			logger.Error("Failed to sync logger", zap.Error(syncErr))
		}
	}()
	logger.Info(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
