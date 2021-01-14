package admin

import "github.com/labstack/echo/v4"

// Router in admin domain
func Router(e *echo.Echo) {
	h := NewHandler()
	g := e.Group("/admin")
	{
		g.GET("/", h.IndexPage)
		g.GET("/list", h.ListPage)
		g.GET("/customer/list", h.CustomerListPage)
		g.GET("/customer/chat/:customerID", h.CustomerChatPage)
		//g.GET("/products/:productId", h.ProductPage)
	}
}
