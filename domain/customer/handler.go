package customer

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Handler - HTTP customer handler.
type Handler struct {
	useCase UseCase
}

// NewHandler - a factory function of customer handler.
func NewHandler(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// List customer handler.
func (h *Handler) List(c echo.Context) error {
	customers, err := h.useCase.AllCustomer()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, customers)
}

// Create customer handler.
func (h *Handler) Create(c echo.Context) error {
	payload := struct {
		Name string `json:"name"`
	}{}
	err := c.Bind(&payload)
	if err != nil {
		return err
	}
	customer := Customer{
		UUID: uuid.UUID{},
		Name: payload.Name,
	}
	err = h.useCase.CreateCustomer(&customer)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, customer)
}
