package middleware

import (
	"net/http"

	"github.com/111161226/TOKO-ENCOUNT/handler"
	"github.com/labstack/echo/v4"
)

//check the user have the chat being entered 
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
			//open chat is gone through whoever is loginned
			if rid == "0" {
				correct = true
				validDid = true
			} else { //private chat 
				did := c.QueryParam("did")
				method := c.Request().Method
				//did is needed when method is POST
				if did == "" && method != "GET" {
					return echo.NewHTTPError(http.StatusBadRequest, "`did` is required")
				}

				if method == "GET" {
					validDid = true
				}

				//check did and userid is invalid in that chat room
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
