package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Router in auth domain
func Router(e *echo.Echo, u UseCase) {
	h := NewHandler(u)
	g := e.Group("/auth")
	{
		g.GET("/register", h.GetRegister)
		g.POST("/register", h.PostRegister)
		g.GET("/signin", h.GetSignin)
		g.POST("/signin", h.PostSignin)
		g.GET("/profile", h.GetProfile, mustLoginMiddleware())
	}
}

func mustLoginMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID, ok := c.Get("user_id").(int64)
			if !ok || userID <= 0 {
				//c.Logger().Error(errors.New("get user info error"))
				//b := new(bytes.Buffer)
				//t.ViewErrorForbiddenPage(b)
				//return c.Stream(http.StatusForbidden, echo.MIMETextHTMLCharsetUTF8, b)
				return c.Redirect(http.StatusFound, "/auth/signin")
			}
			//if {
			//	c.Logger().Error(errors.New("user not login"))
			//	//b := new(bytes.Buffer)
			//	//t.ViewErrorForbiddenPage(b)
			//	//return c.Stream(http.StatusForbidden, echo.MIMETextHTMLCharsetUTF8, b)
			//	return c.Redirect(http.StatusFound, "/auth/signin")
			//}
			return next(c)
		}
	}
}
