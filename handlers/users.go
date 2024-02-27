package handlers

import (
	"nero/forms"
	"nero/helper"
	"nero/models"
	"nero/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Signup		godc
// @Summary		Add User
// @Description	Add user data to database
// @Param		user body models.User true "user info"
// @Tags		User
// @produce		application/json
// @Success		200 {object} models.User "signup response"
// @Success		400	{object} forms.ReqResSwagger "error response"
// @Router		/user/signup 	[post]
// @Security	ApiKeyAuth
func HandlerSignUp(c *gin.Context) {
	var user = &models.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		utils.LogError("controllers/user.go", err, "line-26, HandlerSignUp, bindJson")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	err := user.Validate()
	if err != nil {
		utils.LogError("controllers/user.go", err, "line-33, HandlerSignUp")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	_, err = user.Create()
	if err != nil {
		utils.LogError("controllers/user.go", err, "line-39, HandlerSignUp, signup")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	utils.ReqResHelper(c, http.StatusCreated, "User created!", nil)
	return
}

// userLogin godoc
// @Summary		login user
// @Description	validate user and get token
// @Param	credentials	body forms.Authenticate	true "login user"
// @Tags	User
// @produce	application/json
// @Success	200 {object} forms.ReqResSwagger "login response"
// @Success 400 {object} forms.ReqResSwagger "error response"
// @Router	/user/login	[post]
// @Securit	ApiKeyAuth
func HandlerLogin(c *gin.Context) {
	u := &forms.Authenticate{}
	if err := c.ShouldBindJSON(u); err != nil {
		utils.LogError("controllers/user.go", err, "line-50, HandlerLogin, bindJson")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	user, err := models.GetByPhone(u.Phone)
	if err != nil {
		utils.LogError("controllers/user.go", err, "line-59, HandlerLogin, bindJson")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	if len(user.ID) == 0 {
		utils.LogError("controllers/user.go", nil, "line-63, HandlerLogin, msg-userNotFound")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, "invalid phone or password")
		return
	}
	pass := models.ComparePassword(user.Password, u.Password)
	if !pass {
		utils.LogError("controllers/user.go", nil, "line-69, HandlerLogin, comparePassword")
		utils.ReqResHelper(c, http.StatusForbidden, nil, "invalid phone or password")
		return
	}

	token, err := helper.GenerateToken(user.ID, "admin")
	if err != nil {
		utils.LogError("controllers/user.go", err, "line-76, HandlerLogin, generateToken")
		utils.ReqResHelper(c, http.StatusForbidden, nil, err.Error())
		return
	}
	utils.ReqResHelper(c, http.StatusOK, token, nil)
	return
}
