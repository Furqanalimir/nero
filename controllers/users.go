package controllers

import (
	"nero/handlers"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	R *gin.Engine
}

func NewUsersHandler(c *UsersController) {
	g := c.R.Group("/users")
	g.POST("/signup", handlers.HandlerSignUp)
	g.POST("/login", handlers.HandlerLogin)
}
