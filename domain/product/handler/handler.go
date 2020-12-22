package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/didiyudha/sse-redis/domain/product/model"
	"github.com/didiyudha/sse-redis/domain/product/usecase"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ProductHandler - HTTP product handler.
type ProductHandler struct {
	ProductUseCase usecase.ProductUseCase
}

// NewProductHandler - a factory function of product handler.
func NewProductHandler(productUseCase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		ProductUseCase: productUseCase,
	}
}

// Store product handler.
func (p *ProductHandler) Store(c echo.Context) error {
	payload := struct {
		ID       string `json:"id,omitempty"`
		Name     string `json:"name"`
		Category string `json:"category"`
		Qty      int    `json:"qty"`
	}{}
	if err := c.Bind(&payload); err != nil {
		return err
	}
	product := model.Product{
		Name:      payload.Name,
		Category:  payload.Category,
		Qty:       payload.Qty,
		DeletedAt: nil,
	}
	log.Println(payload)
	if payload.ID == "" {
		product.ID = uuid.New()
		product.CreatedAt = time.Now()
		product.UpdatedAt = time.Now()
	} else {
		id, err := uuid.Parse(payload.ID)
		if err != nil {
			return err
		}
		product.ID = id
		product.UpdatedAt = time.Now()
	}
	if err := p.ProductUseCase.Store(&product); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, product)
}

// Streams a product update.
func (p *ProductHandler) Streams(c echo.Context) error {
	productID := c.Param("productId")
	productUUID, err := uuid.Parse(productID)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()
	prodChan := make(chan model.Product, 1)

	go p.ProductUseCase.StreamProduct(ctx, productUUID, prodChan)

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
