package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project/infra"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRoutes(ctx infra.ServiceContext) {
	r := gin.Default()

	r.Use(ctx.Middleware.Logger())
	r.POST("/login", ctx.Ctl.AuthHandler.Login)
	r.POST("/register", ctx.Ctl.UserHandler.Registration)
	r.GET("/users", ctx.Ctl.UserHandler.All)
	r.POST("/password-reset", ctx.Ctl.PasswordResetHandler.Create)

	staffRoutes := r.Group("/staffs")
	{
		staffRoutes.GET("/", ctx.Ctl.UserHandler.All)
		staffRoutes.POST("/", ctx.Ctl.UserHandler.Registration)
	}

	reservationsRoutes := r.Group("/reservations")
	{
		reservationsRoutes.GET("/", ctx.Ctl.ReservationHandler.All)
		reservationsRoutes.POST("/", ctx.Ctl.ReservationHandler.Add)
		reservationsRoutes.GET("/:id", ctx.Ctl.ReservationHandler.GetByID)
		reservationsRoutes.PUT("/:id", ctx.Ctl.ReservationHandler.Update)
	}

	categoriesRoutes := r.Group("/categories")
	{
		categoriesRoutes.GET("/", ctx.Ctl.CategoryHandler.All)
		categoriesRoutes.POST("/create", ctx.Ctl.CategoryHandler.Create)
		categoriesRoutes.PUT("/:id", ctx.Ctl.CategoryHandler.Update)
	}

	productsRoutes := r.Group("/products")
	{
		productsRoutes.GET("/", ctx.Ctl.CategoryHandler.AllProducts)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	notificationRoutes(ctx, r)

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
