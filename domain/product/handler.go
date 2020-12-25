package product

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
func (p *Handler) Store(c echo.Context) error {
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
	if err := p.useCase.Store(&product); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, product)
}

// Update product handler.
func (p *Handler) Update(c echo.Context) error {
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
	if err := p.useCase.Save(&product); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, product)
}

// Streams a product update.
func (p *Handler) Streams(c echo.Context) error {
	productID := c.Param("productId")
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()
	prodChan := make(chan Product, 1)

	go p.useCase.StreamProduct(ctx, productUUID, prodChan)

	// c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
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
			log.Printf("product p: %+v\n", p)
			if err := enc.Encode(p); err != nil {
				return err
			}
			// var w string
			// fmt.Sprintf(w, "data: %s\n\n", "xxxxxxxxx")
			b, err := json.Marshal(p)
			if err != nil {
				return err
			}
			str := fmt.Sprintf("id: 1\nevent: product update\ndata: %s\n\n", b)
			c.Response().Write([]byte(str))
			c.Response().Flush()
		}
	}
	return nil
}
