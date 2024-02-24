package server

import (
	"context"
	"fmt"
	"log"
	"nero/config"
	"nero/db"
	"nero/utils"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func Init() {
	// middlewares only in dev environment
	// env := config.EnvVars("ENV")
	// if env == "dev" {
	// 	gin.SetMode(gin.DebugMode)
	// }
	// if env == "prod" {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	// // initialize gin server
	// router := gin.Default()

	// // gin middleware
	// router.Use(gin.Logger())
	// router.Use(gin.Recovery())

	// //register routes
	// router.GET("/", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "pong")
	// })
	// controllers.NewUsersHandler(&controllers.UsersController{
	// 	R: router,
	// })
	// controllers.NewOrdersHandler(&controllers.OrdersController{
	// 	R: router,
	// })
	env := config.EnvVars("ENV")
	router := InitRouter(env)
	db.Init()
	startServer(router, env)

}

func startServer(router http.Handler, env string) {
	srv := &http.Server{
		Addr:         config.EnvVars("PORT"),
		Handler:      router,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 5,
	}
	// initialize server with new go routine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to initialize server:\n%v\n", err)
		}
	}()
	utils.ColoredPrintln(fmt.Sprintf("Server listening on port [%v] in [%v] environment.", srv.Addr, strings.ToUpper(env)), utils.CGreen)

	// Wait for kill signal of channel
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// This context is used to inform the server it has 5 seconds to finish
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*5)
	defer cancle()

	// shutting down server
	utils.ColoredPrintln("\nShutting down server...", utils.CPurple)
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shut down: ", err)
	}
}
