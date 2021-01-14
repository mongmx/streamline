package auth

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mongmx/streamline/templates/t"
	"log"
	"net/http"
	"time"
)

// Handler - HTTP auth handler.
type Handler struct {
	useCase UseCase
}

// NewHandler - a factory function of auth handler.
func NewHandler(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h Handler) Register(c echo.Context) error {
	user := &User{
		Email:  "test@mail.com",
		PlanID: 1,
		Auth: &Auth{
			Type:   "email",
			Secret: "123456",
		},
		Topic: &Topic{
			Title: "default",
		},
	}
	err := h.useCase.Register(user)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "success")
}

func (h Handler) GetSignin(c echo.Context) error {
	b := new(bytes.Buffer)
	t.ViewSigninPage(b)
	return c.Stream(http.StatusOK, echo.MIMETextHTMLCharsetUTF8, b)
}

func (h Handler) PostSignin(c echo.Context) error {
	var creds Credentials
	if err := c.Bind(&creds); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	user := &User{
		model: model{
			CreatedAt: time.Now(),
		},
		Email: creds.Username,
		Auth: &Auth{
			UserID: 0,
			Type:   "",
			Secret: creds.Password,
		},
	}
	log.Println(user)
	if err := h.useCase.Signin(user); err != nil {
		return c.JSON(http.StatusBadGateway, err)
	}
	sessionToken := uuid.New().String()
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   120,
		HttpOnly: true,
	}
	sess.Values["session_token"] = sessionToken
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, "success")
}
