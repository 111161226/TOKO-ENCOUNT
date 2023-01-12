package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func validatedBind(c echo.Context, i interface{}) error {
	err := c.Bind(i) // リクエストボディの取り出し
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
