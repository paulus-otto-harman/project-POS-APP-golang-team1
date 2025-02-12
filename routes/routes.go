package routes

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project/infra"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRoutes(ctx infra.ServiceContext) {
	r := gin.Default()

	r.Static("/static", "./static")
	r.Use(cors.Default())

	r.Use(ctx.Middleware.Logger())
	r.POST("/login", ctx.Ctl.AuthHandler.Login)
	r.POST("/otp", ctx.Ctl.PasswordResetHandler.Create)
	r.PUT("/otp/:id", ctx.Ctl.PasswordResetHandler.Update)
	r.PUT("/user/:id", ctx.Ctl.UserHandler.UpdatePassword)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(ctx.Middleware.Jwt.AuthJWT())
	r.POST("/logout", ctx.Ctl.ProfileHandler.Logout)
	r.PUT("/profile", ctx.Ctl.ProfileHandler.Update)
	r.GET("/users", ctx.Middleware.OnlySuperAdmin(), ctx.Ctl.UserHandler.All)
	r.PUT("/users/:id", ctx.Middleware.OnlySuperAdmin(), ctx.Ctl.UserPermissionHandler.Update)

	staffRoutes := r.Group("/staffs", ctx.Middleware.CanAccess("Staff"))
	{
		staffRoutes.GET("/", ctx.Ctl.UserHandler.All)
		staffRoutes.GET("/:id", ctx.Ctl.UserHandler.GetByID)
		staffRoutes.POST("/", ctx.Ctl.UserHandler.Registration)
		staffRoutes.DELETE("/:id", ctx.Ctl.UserHandler.Delete)
		staffRoutes.PUT("/:id", ctx.Ctl.UserHandler.Update)
	}

	reservationsRoutes := r.Group("/reservations", ctx.Middleware.CanAccess("Reservations"))
	{
		reservationsRoutes.GET("/", ctx.Ctl.ReservationHandler.All)
		reservationsRoutes.POST("/", ctx.Ctl.ReservationHandler.Add)
		reservationsRoutes.GET("/:id", ctx.Ctl.ReservationHandler.GetByID)
		reservationsRoutes.PUT("/:id", ctx.Ctl.ReservationHandler.Update)
	}

	categoriesRoutes := r.Group("/categories", ctx.Middleware.CanAccess("Menu"))
	{
		categoriesRoutes.GET("/", ctx.Ctl.CategoryHandler.All)
		categoriesRoutes.POST("/create", ctx.Ctl.CategoryHandler.Create)
		categoriesRoutes.PUT("/:id", ctx.Ctl.CategoryHandler.Update)
	}

	r.GET("/products", ctx.Middleware.CanAccess("Menu"), ctx.Ctl.CategoryHandler.AllProducts)

	inventoryRoutes := r.Group("/inventory", ctx.Middleware.CanAccess("Inventory"))
	{
		inventoryRoutes.GET("/", ctx.Ctl.ProductHandler.All)
		inventoryRoutes.POST("/", ctx.Ctl.ProductHandler.Add)
		inventoryRoutes.PUT("/:id", ctx.Ctl.ProductHandler.Update)
		inventoryRoutes.DELETE("/:id", ctx.Ctl.ProductHandler.Delete)
	}

	dashboardRoutes := r.Group("/dashboard", ctx.Middleware.CanAccess("Dashboard"))
	{
		dashboardRoutes.GET("/", ctx.Ctl.DashboardHandler.GetDashboard)
		dashboardRoutes.GET("/export", ctx.Ctl.DashboardHandler.ExportSalesDataCSV)
		dashboardRoutes.GET("/ws", ctx.Ctl.DashboardHandler.SalesDataWebSocket)
	}

	r.GET("/tables", ctx.Middleware.CanAccess("Orders"), ctx.Ctl.OrderHandler.AllTables)
	r.GET("/payments", ctx.Middleware.CanAccess("Orders"), ctx.Ctl.OrderHandler.AllPayments)

	ordersRoutes := r.Group("/orders", ctx.Middleware.CanAccess("Orders"))
	{
		ordersRoutes.GET("/", ctx.Ctl.OrderHandler.AllOrders)
		ordersRoutes.POST("/", ctx.Ctl.OrderHandler.Create)
		ordersRoutes.PUT("/:id", ctx.Ctl.OrderHandler.Update)
		ordersRoutes.DELETE("/:id", ctx.Ctl.OrderHandler.Delete)
	}

	notificationRoutes := r.Group("/notifications")
	{
		notificationRoutes.GET("/", ctx.Ctl.NotificationHandler.All)
		notificationRoutes.PUT("/:id", ctx.Ctl.NotificationHandler.Update)
		notificationRoutes.PUT("/batch", ctx.Ctl.NotificationHandler.BatchUpdate)
		notificationRoutes.DELETE("/:id", ctx.Ctl.NotificationHandler.Delete)
	}

	revenueRoutes := r.Group("/revenue-reports", ctx.Middleware.CanAccess("Reports"))
	{
		revenueRoutes.GET("/status", ctx.Ctl.RevenueHandler.GetTotalRevenueByStatus)
		revenueRoutes.GET("/bestsellers", ctx.Ctl.RevenueHandler.GetProductRevenueDetails)
		revenueRoutes.GET("/monthly_revenue", ctx.Ctl.RevenueHandler.GetMonthlyRevenue)

	}

	gracefulShutdown(ctx, r.Handler())
}

func gracefulShutdown(ctx infra.ServiceContext, handler http.Handler) {
	srv := &http.Server{
		Addr:    ctx.Cfg.ServerPort,
		Handler: handler,
	}

	if ctx.Cfg.ShutdownTimeout == 0 {
		launchServer(srv, ctx.Cfg.ServerPort)
		return
	}

	go func() {
		launchServer(srv, ctx.Cfg.ServerPort)
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
	appContext, cancel := context.WithTimeout(context.Background(), time.Duration(ctx.Cfg.ShutdownTimeout)*time.Second)
	defer cancel()
	if err := srv.Shutdown(appContext); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching appContext.Done(). timeout of ShutdownTimeout seconds.
	select {
	case <-appContext.Done():
		log.Println(fmt.Sprintf("timeout of %d seconds.", ctx.Cfg.ShutdownTimeout))
	}
	log.Println("Server exiting")
}

func launchServer(server *http.Server, port string) {
	// service connections
	log.Println("Listening and serving HTTP on", port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
