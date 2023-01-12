package middleware

import (
	"net/http"

	"github.com/cs-sysimpl/SakataKintoki/handler"
	"github.com/labstack/echo/v4"
)

func EnsureExistChatAndHaveAccessRight(h *handler.Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := h.PickSession(c)
			if err != nil {
				return err
			}

			rid := c.Param("rid")
			if rid == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "`rid` is required")
			}

			correct := false
			if rid == "0" {
				correct = true
			} else {
				users, err := h.PickChatByRoomId(rid)
				if err != nil {
					return err
				}
				for _, user := range *(users.ChatUsers) {
					if sess.UserId == user.UserId {
						correct = true
						break
					}
				}
			}
			if !correct {
				return echo.NewHTTPError(http.StatusForbidden, "cannot access other's chat")
			}
			return next(c)
		}
	}
}
