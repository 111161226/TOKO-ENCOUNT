package handler

import (
	"net/http"

	"github.com/cs-sysimpl/SakataKintoki/db/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
	//入力取得
	u := new(model.UserSimple)
	if err := c.Bind(u); err != nil {
		return err
	}

	//ログインチェック
	user, err := h.ui.CheckRightUser(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid name or password")
	}

	//セッション作成
	createSessionAndSetCookie(c, h, user.UserId)

	return c.JSON(http.StatusOK, &user)
}

func (h *Handler) Test(c echo.Context) error {
	var user model.UserWithoutPass
	user.UserId = "a"
	user.UserName = "b"
	user.Prefect = "c"
	user.Gender = "d"

	return c.JSON(http.StatusOK, &user)
}