package main

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/mongmx/streamline/domain/auth"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/mongmx/streamline/config"
	"github.com/mongmx/streamline/config/postgres"
	"github.com/mongmx/streamline/config/redis"
	_ "github.com/mongmx/streamline/docs"
	"github.com/mongmx/streamline/domain/admin"
	"github.com/mongmx/streamline/domain/customer"
	"github.com/mongmx/streamline/domain/product"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/sync/errgroup"
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
	appRouter := appInstance(metricsRouter)
	var eg errgroup.Group
	eg.Go(func() error {
		return appRouter.Start(":" + config.Cfg.Port)
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
	if err := appRouter.Shutdown(ctx); err != nil {
		appRouter.Logger.Fatal(err)
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

func appInstance(routerMetrics *echo.Echo) *echo.Echo {
	dsn := postgres.NewPostgres(config.Cfg.Postgres)
	postgresDB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = postgres.MigrateUp(dsn)
	if err != nil {
		log.Fatal(err)
	}
	redisPool, err := redis.NewRedis(config.Cfg.Redis)
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
		session.Middleware(sessions.NewCookieStore([]byte("secret"))),
	)
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	p.SetMetricsPath(routerMetrics)
	authRepo := auth.NewRepository(postgresDB, redisPool)
	authUseCase := auth.NewUseCase(authRepo)
	productRepo := product.NewRepository(postgresDB, redisPool)
	productUseCase := product.NewUseCase(productRepo)
	customerRepo := customer.NewRepository(postgresDB, redisPool)
	customerUseCase := customer.NewUseCase(customerRepo)

	e.Static("/public", "public")
	auth.Router(e, authUseCase)
	admin.Router(e)
	customer.Router(e, customerUseCase)
	product.Router(e, productUseCase)

	e.GET("/metrics", func(c echo.Context) error {
		//e.Logger.Debug(e.Debug)
		return echo.ErrNotFound
	})

	return e
}
