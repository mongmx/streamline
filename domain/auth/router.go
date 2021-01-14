package auth

import "github.com/labstack/echo/v4"

// Router in auth domain
func Router(e *echo.Echo, u UseCase) {
	h := NewHandler(u)
	g := e.Group("/auth")
	{
		g.POST("/register", h.Register)
		g.GET("/signin", h.GetSignin)
		g.POST("/signin", h.PostSignin)
		g.GET("/profile", h.GetProfile)
	}
}
