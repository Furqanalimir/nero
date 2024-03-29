package handlers

import (
	"fmt"
	"nero/forms"
	"nero/helper"
	"nero/models"
	"nero/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Signup		godc
// @Summary		Add User
// @Description	Add user data to database
// @Param		user formData forms.UserSignUpSwagger true "add user"
// @Param 		profile formData file true "profile"
// @Tags		User
// @produce		application/json
// @Success		200 {object} forms.ReqResSwagger	"signup response"
// @Success		400	{object} forms.ReqResSwagger 	"error response"
// @Success		500	{object} forms.ReqResSwagger 	"error message"
// @Router		/users/signup 	[post]
func HandlerSignUp(c *gin.Context) {
	// var user = &models.User{}
	err := c.Request.ParseMultipartForm(1 << 20) // 1 Mb limit
	if err != nil {
		utils.LogError("handlers/user.go", err, "line-30, HandlerSignUp, parseMultipartForm")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	// create user object from form data
	user := models.CreateUserObj(c)

	err = user.Validate()
	if err != nil {
		utils.LogError("handlers/user.go", err, "line-47, HandlerSignUp, validate")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	userExists, err := models.GetByPhone(user.Phone)
	if err != nil {
		utils.LogError("handlers/user.go", err, "line-53, HandlerSignUp, userExists")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	if len(userExists.ID) > 0 || userExists.Phone > 0 {
		utils.LogError("handlers/user.go", err, "line-58, HandlerSignUp, user exists")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, "User already exits")
		return
	}

	file, err := c.FormFile("profile")
	if err == nil {
		filePath := "assets/profile/" + fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		c.SaveUploadedFile(file, filePath)
		user.ProflieUrl = filePath
	}

	_, err = user.Create()
	if err != nil {
		utils.LogError("handlers/user.go", err, "line-65, HandlerSignUp, signup")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	utils.ReqResHelper(c, http.StatusCreated, "User created!", nil)
	return
}

// userLogin 	godoc
// @Summary		login user
// @Description	validate user and get token
// @Param		credentials	body forms.Authenticate	true "login user"
// @Tags		User
// @produce		application/json
// @Success		200 {object} forms.ReqResSwagger "login response"
// @Success 	400 {object} forms.ReqResSwagger "error response"
// @Success		500	{object} forms.ReqResSwagger "error message"
// @Router		/users/login	[post]
func HandlerLogin(c *gin.Context) {
	u := &forms.Authenticate{}
	if err := c.ShouldBindJSON(u); err != nil {
		utils.LogError("handlers/user.go", err, "line-50, HandlerLogin, bindJson")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}
	user, err := models.GetByPhone(u.Phone)
	if err != nil {
		utils.LogError("handlers/user.go", err, "line-59, HandlerLogin, bindJson")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, err.Error())
		return
	}

	if len(user.ID) == 0 {
		utils.LogError("handlers/user.go", nil, "line-63, HandlerLogin, msg-userNotFound")
		utils.ReqResHelper(c, http.StatusBadRequest, nil, "invalid phone or password")
		return
	}

	pass := models.ComparePassword(user.Password, u.Password)
	if !pass {
		utils.LogError("handlers/user.go", nil, "line-69, HandlerLogin, comparePassword")
		utils.ReqResHelper(c, http.StatusForbidden, nil, "invalid phone or password")
		return
	}

	token, err := helper.GenerateToken(user.ID, "admin")
	if err != nil {
		utils.LogError("handlers/user.go", err, "line-76, HandlerLogin, generateToken")
		utils.ReqResHelper(c, http.StatusForbidden, nil, err.Error())
		return
	}
	utils.ReqResHelper(c, http.StatusOK, token, nil)
	return
}
