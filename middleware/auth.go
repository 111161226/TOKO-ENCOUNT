package middleware

import (
	"github.com/111161226/TOKO-ENCOUNT/handler"
	"github.com/labstack/echo/v4"
)

//check user is loginned
func EnsureAuthorized(h *handler.Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()
			method := c.Request().Method
			if path == "/api/user" && method == "POST" {
				//go through when creating account because it is done with noy loggined
				return next(c)
			}
			//get session info
			_, err := h.PickSession(c)
			if err != nil {
				return err
			}

			return next(c)
		}
	}
}
