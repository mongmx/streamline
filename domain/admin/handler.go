package admin

import (
	"log"
	"net/http"

	"github.com/mongmx/sse-redis/templates/t"
)

// Handler - HTTP product handler.
type Handler struct{}

// NewHandler - a factory function of product handler.
func NewHandler() *Handler {
	return &Handler{}
}

// Index product handler.
func (p *Handler) Index(w http.ResponseWriter, r *http.Request) {
	_, err := t.ViewTestPage(w)
	if err != nil {
		log.Println(err)
	}
}
