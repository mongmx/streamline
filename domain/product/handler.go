package product

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler - HTTP product handler.
type Handler struct {
	useCase UseCase
}

// NewHandler - a factory function of product handler.
func NewHandler(useCase UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// Store product handler.
func (h *Handler) Store(c echo.Context) error {
	payload := struct {
		Name     string `json:"name"`
		Category string `json:"category"`
		Qty      int    `json:"qty"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return err
	}
	product := Product{
		ID:        uuid.New(),
		Name:      payload.Name,
		Category:  payload.Category,
		Qty:       payload.Qty,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
	if err := h.useCase.Store(&product); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, product)
}

// Update product handler.
func (h *Handler) Update(c echo.Context) error {
	payload := struct {
		Name     string `json:"name"`
		Category string `json:"category"`
		Qty      int    `json:"qty"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return err
	}
	productID := c.Param("productId")
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return err
	}
	product := Product{
		ID:        productUUID,
		Name:      payload.Name,
		Category:  payload.Category,
		Qty:       payload.Qty,
		UpdatedAt: time.Now(),
	}
	if err := h.useCase.Save(&product); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, product)
}

// Streams a product update.
func (h *Handler) Streams(c echo.Context) error {
	productID := c.Param("productId")
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()
	prodChan := make(chan Product, 1)

	go h.useCase.StreamProduct(ctx, productUUID, prodChan)

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
		for p := range prodChan {
			log.Printf("product h: %+v\n", p)
			if err := enc.Encode(p); err != nil {
				return err
			}
			b, err := json.Marshal(p)
			if err != nil {
				return err
			}
			b = formatSSE("product update", string(b))
			_, err = c.Response().Write(b)
			if err != nil {
				return err
			}
			c.Response().Flush()
		}
	}
	return nil
}

func formatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	return []byte(eventPayload + "\n")
}
