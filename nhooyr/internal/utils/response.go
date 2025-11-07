package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func RespondJSON(c *gin.Context, statusCode int, message string, data interface{}) {
	status := "success"
	if statusCode >= http.StatusBadRequest {
		status = "error"
	}
	c.JSON(statusCode, APIResponse{
		Status:  status,
		Message: message,
		Data:    data,
	})
}
