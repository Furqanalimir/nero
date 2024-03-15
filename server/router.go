package server

import (
	"nero/controllers"
	_ "nero/docs"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
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
	router.Use(gzip.Gzip(gzip.DefaultCompression)) // compressor middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // allow all origins
		AllowMethods:     []string{"*"}, // allow all methods
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 12 * time.Hour,
	}))

	controllers.NewUsersHandler(&controllers.UsersController{
		R: router,
	})
	controllers.NewOrdersHandler(&controllers.OrdersController{
		R: router,
	})
	return router
}
