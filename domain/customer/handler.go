package customer

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
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

// View customer handler.
func (h *Handler) View(c echo.Context) error {
	customerID := c.Param("customerId")
	customerUUID, err := uuid.Parse(customerID)
	if err != nil {
		return err
	}
	customer, err := h.useCase.GetCustomer(customerUUID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, customer)
}

// Streams a product update.
func (h *Handler) Streams(c echo.Context) error {
	customerID := c.Param("customerId")
	customerUUID, err := uuid.Parse(customerID)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()
	chatChan := make(chan ChatMessage, 1)

	go h.useCase.StreamCustomer(ctx, customerUUID, chatChan)

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().WriteHeader(http.StatusOK)
	enc := json.NewEncoder(c.Response())

	select {
	case <-ctx.Done():
		return nil
	default:
		for cc := range chatChan {
			if err := enc.Encode(cc); err != nil {
				return err
			}
			b, err := json.Marshal(cc)
			if err != nil {
				return err
			}
			b = formatSSE("customer chat", string(b))
			_, err = c.Response().Write(b)
			if err != nil {
				return err
			}
			c.Response().Flush()
		}
	}
	return nil
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
		UUID: uuid.New(),
		Name: payload.Name,
	}
	err = h.useCase.CreateCustomer(&customer)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, customer)
}

func formatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	return []byte(eventPayload + "\n")
}
