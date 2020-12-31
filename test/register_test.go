package test

import (
	"log"
	"net/http"
	"testing"

	"github.com/appleboy/gofight/v2"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/mongmx/streamline/domain/auth"
	"github.com/stretchr/testify/assert"
)

// EchoEngine is echo router.
func EchoEngine() *echo.Echo {
	// Echo instance
	e := echo.New()

	dsn := "host=127.0.0.1 port=5432 dbname=streamline user=mongmx sslmode=disable"
	postgresDB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	authRepo := auth.NewRepository(postgresDB, nil)
	authUseCase := auth.NewUseCase(authRepo)
	authHandler := auth.NewHandler(authUseCase)
	defer postgresDB.Close()

	// Routes
	e.POST("/api/auth/register", authHandler.Register)

	return e
}

func TestRegister(t *testing.T) {
	r := gofight.New()

	r.POST("/api/auth/register").
		SetDebug(true).
		Run(EchoEngine(), func(r gofight.HTTPResponse, rq gofight.HTTPRequest) {
			assert.Equal(t, "success", r.Body.String())
			assert.Equal(t, http.StatusOK, r.Code)
		})
}
