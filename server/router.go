package server

import (
	"nero/controllers"
	_ "nero/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// initializing swagger for docs <host>/docs/index.html#/ (e.g: http://localhost:5050/docs/index.html#/)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
