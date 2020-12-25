package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mongmx/sse-redis/config"
	"github.com/mongmx/sse-redis/config/postgres"
	"github.com/mongmx/sse-redis/config/redis"
	"github.com/mongmx/sse-redis/domain/product"
	"golang.org/x/sync/errgroup"
)

var eg errgroup.Group

func main() {
	godotenv.Load(".env")
	config.LoadEnv()
	eg.Go(func() error {
		return serveAPI()
	})
	eg.Go(func() error {
		return serveAPIDoc()
	})
	eg.Go(func() error {
		return serveMetrics()
	})
	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}

func serveAPI() error {
	postgresDB, err := postgres.NewPostgres(config.Cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	redisConn, err := redis.NewRedis(config.Cfg.Redis)
	if err != nil {
		log.Fatal(err)
	}

	productRepo := product.NewRepository(redisConn, postgresDB)
	productUseCase := product.NewUseCase(productRepo)
	productHandler := product.NewHandler(productUseCase)

	e := echo.New()
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)
	e.HideBanner = true
	// e.Debug, err = strconv.ParseBool(config.Cfg.Debug)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	e.Server = &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Cfg.Port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	e.POST("/products", productHandler.Store)
	e.PUT("/products/:productId", productHandler.Update)
	e.GET("/products/streams/:productId", productHandler.Streams)
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(
			http.StatusOK,
			map[string]string{
				"message": "pong",
			},
		)
	})

	// e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Cfg.Port)))
	// e.Use(
	// 	gin.Logger(),
	// 	gin.Recovery(),
	// 	ginprom.PromMiddleware(nil),
	// 	middleware.TraceMiddleware(),
	// 	// middleware.LoggerToFile(),
	// )

	// memberRepo, err := postgres.NewMemberRepository(db)
	// must(err)
	// memberService, err := member.NewService(memberRepo)
	// must(err)

	// authRepo, err := postgres.NewMemberRepository(db)
	// must(err)
	// authService, err := auth.NewService(authRepo, nil)
	// must(err)

	// member.Routes(e, memberService, authService)
	// auth.Routes(e, authService)

	log.Printf("API server listen on :%s", config.Cfg.Port)
	return e.StartServer(e.Server)
}

func serveAPIDoc() error {
	e := echo.New()
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)
	e.HideBanner = true
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(
			http.StatusOK,
			map[string]string{
				"message": "pong",
			},
		)
	})
	log.Println("API Doc server listen on :8081")
	return e.Start(":8081")
}

func serveMetrics() error {
	e := echo.New()
	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)
	e.HideBanner = true
	// e.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(
			http.StatusOK,
			map[string]string{
				"message": "pong",
			},
		)
	})
	log.Println("Metrics server listen on :8082")
	return e.Start(":8082")
}
