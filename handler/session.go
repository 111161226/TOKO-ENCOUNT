package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/111161226/TOKO-ENCOUNT/db/model"
	"github.com/labstack/echo/v4"
)

//get session by session id
func (h *Handler) PickSession(c echo.Context) (*model.Session, error) {
	cookie, err := c.Cookie("session_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	sess, err := h.si.CheckSession(cookie.Value)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if sess == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	return sess, nil
}

//create and set Cookie
func createSessionAndSetCookie(c echo.Context, h *Handler, userId string) error {
	sess, err := h.si.CreateSession(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_id",
		Value:    sess.SessionId,
		Path:     "/",
		Expires:  time.Date(2030, 12, 31, 23, 59, 59, 99, time.Local),
		HttpOnly: true,
	})

	return nil
}
