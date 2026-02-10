package utils

import (
	"qlass-be/dtos"

	"github.com/gin-gonic/gin"
)

func SendError(c *gin.Context, status int, errCode string, message string) {
    c.JSON(status, dtos.GlobalErrorResponse{
        Error:   errCode,
        Message: message,
    })
}