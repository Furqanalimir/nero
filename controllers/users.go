package controllers

import (
	"nero/forms"
	"nero/helper"
	"nero/models"
	"nero/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UsersController struct {
	R *gin.Engine
}

func NewUsersHandler(c *UsersController) {
	// Create an fruits group
	g := c.R.Group("/users")
	g.POST("/signup", handlerSignUp)
	g.POST("/login", handlerLogin)
}

func handlerSignUp(c *gin.Context) {
	var user = &models.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		utils.LogError("controllers/user.go", err, "line-26, handlerSignUp, bindJson")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	err := user.Validate()
	if err != nil {
		utils.LogError("controllers/user.go", err, "line-33, handlerSignUp")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	_, err = user.Create()
	if err != nil {
		utils.LogError("controllers/user.go", err, "line-39, handlerSignUp, signup")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	utils.ReqResHelper(c, http.StatusCreated, "User created!", nil)
	return
}

func handlerLogin(c *gin.Context) {
	u := &forms.Authenticate{}
	if err := c.ShouldBindJSON(u); err != nil {
		utils.LogError("controllers/user.go", err, "line-50, handlerLogin, bindJson")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	user, err := models.GetByPhone(u.Phone)
	if err != nil {
		utils.LogError("controllers/user.go", err, "line-59, handlerLogin, bindJson")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	if len(user.ID) == 0 {
		utils.LogError("controllers/user.go", nil, "line-63, handlerLogin, msg-userNotFound")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, "invalid phone or password")
		return
	}
	pass := models.ComparePassword(user.Password, u.Password)
	if !pass {
		utils.LogError("controllers/user.go", nil, "line-69, handlerLogin, comparePassword")
		utils.ReqResHelper(c, http.StatusForbidden, nil, "invalid phone or password")
		return
	}

	token, err := helper.GenerateToken(user.ID, "admin")
	if err != nil {
		utils.LogError("controllers/user.go", err, "line-76, handlerLogin, generateToken")
		utils.ReqResHelper(c, http.StatusForbidden, nil, err.Error())
		return
	}
	utils.ReqResHelper(c, http.StatusOK, token, nil)
	return
}
