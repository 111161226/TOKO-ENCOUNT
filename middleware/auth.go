package middleware

import (
	"github.com/111161226/TOKO-ENCOUNT/handler"
	"github.com/labstack/echo/v4"
)

func EnsureAuthorized(h *handler.Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Path()
			method := c.Request().Method
			if path == "/api/user" && method == "POST" {
				// アカウント作成は未ログイン状態で行われるので通す
				return next(c)
			}

			_, err := h.PickSession(c)
			if err != nil {
				return err
			}

			return next(c)
		}
	}
}
