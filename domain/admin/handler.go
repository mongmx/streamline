package admin

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/mongmx/streamline/templates/t"
	"net/http"
)

// Handler - HTTP admin handler.
type Handler struct{}

// NewHandler - a factory function of admin handler.
func NewHandler() *Handler {
	return &Handler{}
}

// IndexPage admin ui handler.
func (h *Handler) IndexPage(c echo.Context) error {
	b := new(bytes.Buffer)
	t.ViewTestPage(b)
	return c.Stream(http.StatusOK, echo.MIMETextHTMLCharsetUTF8, b)
}

// ListPage admin ui handler.
func (h *Handler) ListPage(c echo.Context) error {
	b := new(bytes.Buffer)
	t.ViewListPage(b)
	return c.Stream(http.StatusOK, echo.MIMETextHTMLCharsetUTF8, b)
}

// CustomerListPage admin ui handler.
func (h *Handler) CustomerListPage(c echo.Context) error {
	b := new(bytes.Buffer)
	t.ViewCustomerListPage(b)
	return c.Stream(http.StatusOK, echo.MIMETextHTMLCharsetUTF8, b)
}
