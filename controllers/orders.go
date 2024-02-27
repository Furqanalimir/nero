package controllers

import (
	"nero/handlers"
	"nero/middleware"

	"github.com/gin-gonic/gin"
)

type OrdersController struct {
	R *gin.Engine
}

func NewOrdersHandler(c *OrdersController) {
	// grouping routest (e.g: localhost:<basePath>/orders/*)
	g := c.R.Group("/orders")

	// authentication middleware
	g.Use(middleware.AuthenticateToken)
	// protected routes
	g.POST("/", handlers.HandlerAddOrder)
	g.GET("/:id", handlers.HandlerGetOrderById)
}
