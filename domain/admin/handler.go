package admin

import (
	"github.com/gobuffalo/plush"
	"github.com/labstack/echo/v4"
	"github.com/mongmx/sse-redis/domain/admin/template"
	"io/ioutil"
	"log"
	"net/http"
)

// Handler - HTTP product handler.
type Handler struct{}

// NewHandler - a factory function of product handler.
func NewHandler() *Handler {
	return &Handler{}
}

func html() string {
	b, _ := ioutil.ReadFile("templates/admin/index.html")
	return string(b)
}

// Index product handler.
func (p *Handler) Index(c echo.Context) error {
	ctx := plush.NewContext()
	s, _ := plush.Render(html(), ctx)
	return c.HTML(http.StatusOK, s)
}

func (p *Handler) Index2(w http.ResponseWriter, r *http.Request) {
	//var userList = []string{
	//	"Alice",
	//	"Bob",
	//	"Tom",
	//}
	//_, err := template.ViewUserListToWriter(userList, w)
	_, err := template.ViewTestPage(w)
	if err != nil {
		log.Println(err)
	}
}
