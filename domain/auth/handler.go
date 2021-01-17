package auth

import (
	"bytes"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/mongmx/streamline/templates/t"
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

func (h Handler) GetRegister(c echo.Context) error {
	b := new(bytes.Buffer)
	t.ViewSignupPage(b)
	return c.Stream(http.StatusOK, echo.MIMETextHTMLCharsetUTF8, b)
}

func (h Handler) PostRegister(c echo.Context) error {
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
	var cred Credentials
	if err := c.Bind(&cred); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	user := &User{
		model: model{
			ID:        1,
			CreatedAt: time.Now(),
		},
		Email: cred.Email,
	}
	c.Logger().Printf("signin user %v", user)
	sess, _ := session.Get("session", c)
	token, ok := sess.Values["session_token"].(string)
	if !ok {
		b := new(bytes.Buffer)
		t.ViewErrorForbiddenPage(b)
		return c.Stream(http.StatusForbidden, echo.MIMETextHTMLCharsetUTF8, b)
	}
	if err := h.useCase.Signin(token, user); err != nil {
		return c.JSON(http.StatusBadGateway, err)
	}
	return c.JSON(http.StatusOK, "success")
}

func (h Handler) GetProfile(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	c.Logger().Printf("%v", sess.Values)
	token, ok := sess.Values["session_token"].(string)
	if !ok {
		b := new(bytes.Buffer)
		t.ViewErrorForbiddenPage(b)
		return c.Stream(http.StatusForbidden, echo.MIMETextHTMLCharsetUTF8, b)
	}
	_, err = h.useCase.Profile(token)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	b := new(bytes.Buffer)
	t.ProfilePage(b)
	return c.Stream(http.StatusForbidden, echo.MIMETextHTMLCharsetUTF8, b)
}
