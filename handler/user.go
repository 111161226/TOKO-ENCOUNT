package handler

import (
	"net/http"

	"github.com/cs-sysimpl/SakataKintoki/db/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
	//入力取得
	u := new(model.UserSimple)
	err := validatedBind(c, u)
	if err != nil {
		return err
	}

	//ログインチェック
	user, err := h.ui.CheckRightUser(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid name or password")
	}

	//セッション作成
	err = createSessionAndSetCookie(c, h, user.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &user)
}

func (h *Handler) SignUp(c echo.Context) error {
	//入力取得
	u := new(model.User)
	err := validatedBind(c, u)
	if err != nil {
		return err
	}

	//ユーザ名、パスワード確認
	if u.UserName == "" || u.Password == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid name or password")
	}

	//重複チェック
	userdup, err := h.ui.CheckUsedUserName(u.UserName)
	if err != nil { // DBエラーの場合
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if userdup != nil { // 重複している場合
		return echo.NewHTTPError(http.StatusUnauthorized, "Username is already taken")
	}

	//登録
	user, err := h.ui.CreateUser(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &user)
}

func (h *Handler) EditProfile(c echo.Context) error {
	//入力取得
	u := new(model.UserUpdate)
	if err := c.Bind(u); err != nil{
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//セッション取得
	sess, err := h.PickSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//プロフィール更新
	newprofile, err := h.ui.EditUser(sess.UserId, u)
	if err != nil { //DBエラー
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if newprofile == nil { //パスワード不一致
		return echo.NewHTTPError(http.StatusUnauthorized, "incorrect password")
	}

	return c.JSON(http.StatusOK, newprofile)
}

func (h *Handler) GetMyUser(c echo.Context) error{
	//セッション取得
	sess, err := h.PickSession(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//自身のユーザー情報取得
	user, err := h.ui.GetUser(sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}
