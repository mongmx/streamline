package customer

import "github.com/labstack/echo/v4"

// Router in customer domain
func Router(e *echo.Echo, u UseCase) {
	h := NewHandler(u)
	g := e.Group("/customer")
	{
		g.GET("/", h.List)
		g.GET("/:customerId", h.View)
		g.GET("/streams/:customerId", h.Streams)
		g.POST("/", h.Create)
	}
}
