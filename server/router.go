package server

import (
	"nero/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter(env string) *gin.Engine {

	if env == "dev" {
		gin.SetMode(gin.DebugMode)
	}
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	// initialize gin server
	router := gin.Default()

	// gin middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	controllers.NewUsersHandler(&controllers.UsersController{
		R: router,
	})
	controllers.NewOrdersHandler(&controllers.OrdersController{
		R: router,
	})
	return router
}
