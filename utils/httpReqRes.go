package utils

import (
	"github.com/gin-gonic/gin"
)

type SwaggerRequestResponse struct {
	Error string
	Data  string
}

func ReqResHelper(c *gin.Context, status int, data any, err any) {
	c.JSON(status, gin.H{
		"error": err,
		"data":  data,
	})
}
