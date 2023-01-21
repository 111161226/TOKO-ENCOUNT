package handler

import (
	"net/http"
	"strconv"

	"github.com/cs-sysimpl/SakataKintoki/db/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Logout(c echo.Context) error {
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}

	err = h.si.DeleteSessionBySessionId(sess.SessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

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
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid name or password")
	}

	//セッション作成
	err = createSessionAndSetCookie(c, h, user.UserId)
	if err != nil {
		return err
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
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid name or password")
	}

	//重複チェック
	userdup, err := h.ui.CheckUsedUserName(u.UserName)
	if err != nil { // DBエラーの場合
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if userdup != nil { // 重複している場合
		return echo.NewHTTPError(http.StatusBadRequest, "Username is already taken")
	}

	//登録
	user, err := h.ui.CreateUser(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.ci.AddOpenChat(user.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = createSessionAndSetCookie(c, h, user.UserId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &user)
}

func (h *Handler) EditProfile(c echo.Context) error {
	//入力取得
	u := new(model.UserUpdate)
	err := validatedBind(c, u)
	if err != nil {
		return err
	}

	//セッション取得
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}

	//プロフィール更新
	newprofile, err := h.ui.EditUser(sess.UserId, u)
	if err != nil { //DBエラー
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if newprofile == nil { //パスワード不一致
		return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect password")
	}

	return c.JSON(http.StatusOK, newprofile)
}

func (h *Handler) GetMyUser(c echo.Context) error {
	//セッション取得
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}

	//自身のユーザー情報取得
	user, err := h.ui.GetUser(sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) SearchUser(c echo.Context) error {
	//入力取得
	l := c.QueryParam("limit")
	if l == "" {
		l = "20"
	}
	limit, err := strconv.Atoi(l)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	o := c.QueryParam("offset")
	if o == "" {
		o = "0"
	}
	offset, err := strconv.Atoi(o)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	name := c.QueryParam("name")
	gender := c.QueryParam("gender")
	prefect := c.QueryParam("prefect")

	//対象となるユーザ取得
	userlist, err := h.ui.GetUserList(limit, offset, name, gender, prefect)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, userlist)
}
