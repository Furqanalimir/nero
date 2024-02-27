package handlers

import (
	"nero/models"
	"nero/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AddOrder		godoc
// @Summary		CreateOrder
// @Description	Save order data into databae
// @Param		order body forms.OrderSwaggerForm true "create order"
// @produce		application/json
// @Tags		Orders
// @Success		200 {object} forms.ReqResSwagger "create response"
// @Success		400	{object} forms.ReqResSwagger "error response"
// @Router		/order	[post]
// @Security	ApiKeyAuth
func HandlerAddOrder(c *gin.Context) {
	userId := c.MustGet("user_id")
	userIdStr, _ := userId.(string)

	o := &models.Order{}
	if err := c.ShouldBindJSON(o); err != nil {
		utils.LogError("handlers/orders.go", err, "func-HandlerAddOrder, line-16, bindJson")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	o.UserId = string(userIdStr)
	err := o.Validate()
	if err != nil {
		utils.LogError("handlers/orders.go", err, "func-HandlerAddOrder, line-34, userValidate")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	order, err := o.CreateOrder()
	if err != nil {
		utils.LogError("handlers/orders.go", err, "func-HandlerAddOrder, line-40, CreateUser")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, "Unable to create order")
		return
	}
	utils.ReqResHelper(c, http.StatusBadRequest, order, nil)
	return
}

// GetOrder		godoc
// @Summary		fet order from database
// @Description	fetch user order by order id
// @Param		id path string true "order id"
// @Tags		Orders
// @produce		application/json
// @Success		200 {object} models.Order "success response"
// @Success		400 {object}  forms.ReqResSwagger "error response"
// @Router		/order/:id	[get]
// @Security	ApiKeyAuth
func HandlerGetOrderById(c *gin.Context) {
	id, ok := c.Params.Get("id")
	var order *models.Order
	if !ok {
		utils.ColoredPrintln("handlers/orders.go func-HandlerGetOrderById, line-62, id-missing", utils.CRed)
		utils.ReqResHelper(c, http.StatusBadRequest, nil, "id is required and must be a string.")
		return
	}
	ord, err := order.GetById(id)
	if err != nil {
		utils.LogError("handlers/orders.go", err, " func-HandlerGetOrderById, line-68, order.getById")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, "Could not get order")
		return
	}
	utils.ReqResHelper(c, http.StatusOK, ord, nil)
	return
}
