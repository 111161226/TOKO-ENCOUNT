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
			validDid := false
			if rid == "0" {
				correct = true
				validDid = true
			} else {
				did := c.QueryParam("did")
				if did == "" {
					return echo.NewHTTPError(http.StatusBadRequest, "`did` is required")
				}

				users, err := h.PickChatByRoomId(rid)
				if err != nil {
					return err
				}
				for _, user := range *(users.ChatUsers) {
					if sess.UserId == user.UserId {
						correct = true
					}

					if did == user.UserId && did != sess.UserId {
						validDid = true
					}
				}
			}
			if !correct {
				return echo.NewHTTPError(http.StatusForbidden, "cannot access other's chat")
			}
			if !validDid {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid `did`")
			}

			return next(c)
		}
	}
}
