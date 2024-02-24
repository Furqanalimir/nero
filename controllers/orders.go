package controllers

import (
	"nero/middleware"
	"nero/models"
	"nero/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrdersController struct {
	R *gin.Engine
}

func NewOrdersHandler(c *OrdersController) {
	g := c.R.Group("/orders")

	g.Use(middleware.AuthenticateToken)
	g.POST("/", handlerAddOrder)
	g.POST("/:id", handlerGetOrderById)
}

func handlerAddOrder(c *gin.Context) {
	userId := c.MustGet("user_id")
	userIdStr, _ := userId.(string)

	o := &models.Order{}
	if err := c.ShouldBindJSON(o); err != nil {
		utils.LogError("controllers/orders.go", err, "func-handlerAddOrder, line-24, bindJson")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	o.UserId = string(userIdStr)
	err := o.Validate()
	if err != nil {
		utils.LogError("controllers/orders.go", err, "func-handlerAddOrder, line-31, userValidate")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	order, err := o.CreateOrder()
	if err != nil {
		utils.LogError("controllers/orders.go", err, "func-handlerAddOrder, line-37, CreateUser")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, "Unable to create order")
		return
	}
	utils.ReqResHelper(c, http.StatusBadRequest, order, nil)
	return
}

func handlerGetOrderById(c *gin.Context) {
	id, ok := c.Params.Get("id")
	var order *models.Order
	if !ok {
		utils.ColoredPrintln("controllers/orders.go func-handlerGetOrderById, line-49, id-missing", utils.CRed)
		utils.ReqResHelper(c, http.StatusBadRequest, nil, "id is required and must be a string.")
		return
	}
	ord, err := order.GetById(id)
	if err != nil {
		utils.LogError("controllers/orders.go", err, " func-handlerGetOrderById, line-55, order.getById")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, "Could not get order")
		return
	}
	utils.ReqResHelper(c, http.StatusOK, ord, nil)
	return
}
