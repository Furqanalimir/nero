package middleware

import (
	"nero/helper"
	"nero/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticateToken(c *gin.Context) {
	id, err := helper.ExtractAuthToken(c.Request)
	if err != nil {
		utils.ReqResHelper(c, http.StatusNonAuthoritativeInfo, nil, err.Error())
		c.Abort()
		return
	}
	c.Set("user_id", id)
	c.Next()
}
