package test

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/mongmx/streamline/domain/auth"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
	dsn := fmt.Sprintf(
		"host=127.0.0.1 port=5432 dbname=streamline user=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
	)
	postgresDB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	authRepo := auth.NewRepository(postgresDB, nil)
	authUseCase := auth.NewUseCase(authRepo)
	authHandler := auth.NewHandler(authUseCase)
	defer func() {
		err = postgresDB.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, authHandler.Register(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Regexp(t, "success", rec.Body.String())
	}
}
