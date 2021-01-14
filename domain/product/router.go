package product

import "github.com/labstack/echo/v4"

// Router in product domain
func Router(e *echo.Echo, u UseCase) {
	h := NewHandler(u)
	g := e.Group("/product")
	{
		g.POST("/", h.Store)
		g.PUT("/:productId", h.Update)
		g.GET("/streams/:productId", h.Streams)
	}
}
