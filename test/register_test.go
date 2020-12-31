package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/mongmx/streamline/domain/auth"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	dsn := "host=127.0.0.1 port=5432 dbname=streamline user= sslmode=disable"
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
