package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func validatedBind(c echo.Context, i interface{}) error {
	err := c.Bind(i) //get request body
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
