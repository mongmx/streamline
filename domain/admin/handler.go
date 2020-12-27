package admin

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/mongmx/streamline/templates/t"
	"net/http"
)

// Handler - HTTP product handler.
type Handler struct{}

// NewHandler - a factory function of product handler.
func NewHandler() *Handler {
	return &Handler{}
}

// IndexPageHandler product ui handler.
func (p *Handler) IndexPage(c echo.Context) error {
	b := new(bytes.Buffer)
	t.ViewTestPage(b)
	return c.Stream(http.StatusOK, echo.MIMETextHTMLCharsetUTF8, b)
}

// ListPageHandler product ui handler.
func (p *Handler) ListPage(c echo.Context) error {
	b := new(bytes.Buffer)
	t.ViewListPage(b)
	return c.Stream(http.StatusOK, echo.MIMETextHTMLCharsetUTF8, b)
}
