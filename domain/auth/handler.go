package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
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
