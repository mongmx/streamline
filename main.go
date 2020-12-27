package main

import (
	"context"
	"github.com/mongmx/streamline/domain/customer"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/mongmx/streamline/domain/admin"
	"golang.org/x/sync/errgroup"

	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mongmx/streamline/config"
	"github.com/mongmx/streamline/config/postgres"
	"github.com/mongmx/streamline/config/redis"
	_ "github.com/mongmx/streamline/docs"
	"github.com/mongmx/streamline/domain/product"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title SSE Hub API
// @version 1.0
// description This is a sample server.
// termsOfService http://swagger.io/terms/
// contact.name API Support
// contact.url http://www.swagger.io/support
// contact.email support@swagger.io
// license.name Apache 2.0
// license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config.LoadEnv()
	metricsRouter := metricsInstance()
	apiDocRouter := apiDocInstance()
	apiRouter := apiInstance(metricsRouter)
	var eg errgroup.Group
	eg.Go(func() error {
		return apiRouter.Start(":" + config.Cfg.Port)
	})
	eg.Go(func() error {
		return apiDocRouter.Start(":8081")
	})
	eg.Go(func() error {
		return metricsRouter.Start(":8082")
	})
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := metricsRouter.Shutdown(ctx); err != nil {
		metricsRouter.Logger.Fatal(err)
	}
	if err := apiDocRouter.Shutdown(ctx); err != nil {
		apiDocRouter.Logger.Fatal(err)
	}
	if err := apiRouter.Shutdown(ctx); err != nil {
		apiRouter.Logger.Fatal(err)
	}
}

func apiDocInstance() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	return e
}

func metricsInstance() *echo.Echo {
	e := echo.New()
	return e
}

func apiInstance(routerMetrics *echo.Echo) *echo.Echo {
	postgresDB, err := postgres.NewPostgres(config.Cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	redisConn, err := redis.NewRedis(config.Cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}
	e := echo.New()
	e.HideBanner = true
	e.Debug, err = strconv.ParseBool(config.Cfg.Debug)
	if err != nil {
		e.Debug = false
	}
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	p.SetMetricsPath(routerMetrics)

	productRepo := product.NewRepository(redisConn, postgresDB)
	productUseCase := product.NewUseCase(productRepo)
	productHandler := product.NewHandler(productUseCase)

	customerRepo := customer.NewRepository(postgresDB, redisConn)
	customerUseCase := customer.NewUseCase(customerRepo)
	customerHandler := customer.NewHandler(customerUseCase)

	adminHandler := admin.NewHandler()

	e.Static("/public", "public")

	e.GET("/api/customers", customerHandler.List)
	e.GET("/api/customer/:customerId", customerHandler.View)
	e.GET("/api/customer/streams/:customerId", customerHandler.Streams)
	e.POST("/api/customer", customerHandler.Create)

	e.POST("/products", productHandler.Store)
	e.PUT("/products/:productId", productHandler.Update)
	e.GET("/products/streams/:productId", productHandler.Streams)

	e.GET("/admin", adminHandler.IndexPage)
	e.GET("/admin/list", adminHandler.ListPage)
	e.GET("/admin/customer/list", adminHandler.CustomerListPage)
	e.GET("/admin/customer/chat/:customerID", adminHandler.CustomerChatPage)
	//e.GET("/admin/products/:productId", adminHandler.ProductPage)

	e.GET("/metrics", func(c echo.Context) error {
		//e.Logger.Debug(e.Debug)
		return echo.ErrNotFound
	})

	return e
}
