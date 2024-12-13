package routes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project/infra"
	"syscall"
	"time"
)

func NewRoutes(ctx infra.ServiceContext) {
	r := gin.Default()

	r.Use(ctx.Middleware.Logger())
	r.POST("/login", ctx.Ctl.AuthHandler.Login)
	r.POST("/register", ctx.Ctl.UserHandler.Registration)
	r.GET("/users", ctx.Ctl.UserHandler.All)
	r.POST("/password-reset", ctx.Ctl.PasswordResetHandler.Create)

	r.GET("/staffs", ctx.Middleware.UserCan("list-staff"), func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	r.GET("/staffs/:id", ctx.Middleware.UserCan("view-staff"), func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	r.POST("/staffs", ctx.Middleware.UserCan("create-staff"), func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	r.PUT("/staffs/:id", ctx.Middleware.UserCan("update-staff"), func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	r.DELETE("/staffs/:id", ctx.Middleware.UserCan("delete-staff"), func(c *gin.Context) {
		c.JSON(200, gin.H{"hello": "world"})
	})

	categoriesRoutes := r.Group("/categories")
	{
		categoriesRoutes.GET("/", ctx.Ctl.CategoryHandler.All)
	}


	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	gracefulShutdown(ctx, r.Handler())
}

func gracefulShutdown(ctx infra.ServiceContext, handler http.Handler) {
	const ShutdownTimeout = 5

	srv := &http.Server{
		Addr:    ctx.Cfg.ServerPort,
		Handler: handler,
	}

	go func() {
		// service connections
		log.Println("Listening and serving HTTP on", ctx.Cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	appContext, cancel := context.WithTimeout(context.Background(), ShutdownTimeout*time.Second)
	defer cancel()
	if err := srv.Shutdown(appContext); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching appContext.Done(). timeout of ShutdownTimeout seconds.
	select {
	case <-appContext.Done():
		log.Println(fmt.Sprintf("timeout of %d seconds.", ShutdownTimeout))
	}
	log.Println("Server exiting")
}
