package handler

import (
	"net/http"
	"strconv"

	"github.com/111161226/TOKO-ENCOUNT/db/model"
	"github.com/labstack/echo/v4"
)

//withdraw user by logic
func (h *Handler) DeleteUser(c echo.Context) error {
	//get userid by session
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}
	//delete user
	err = h.ui.DeleteUser(sess.SessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

//func for logout
func (h *Handler) Logout(c echo.Context) error {
	//get userid by session
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}
	//delete session
	err = h.si.DeleteSessionBySessionId(sess.SessionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

//func for login
func (h *Handler) Login(c echo.Context) error {
	//get user input
	u := new(model.UserSimple)
	err := validatedBind(c, u)
	if err != nil {
		return err
	}

	//check if user can login
	user, err := h.ui.CheckRightUser(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid name or password")
	}

	//create session
	err = createSessionAndSetCookie(c, h, user.UserId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &user)
}

func (h *Handler) SignUp(c echo.Context) error {
	//get new user input
	u := new(model.User)
	err := validatedBind(c, u)
	if err != nil {
		return err
	}

	//check username and password
	if u.UserName == "" || u.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid name or password")
	}

	//check if username is duplicated
	userdup, err := h.ui.CheckUsedUserName(u.UserName)
	if err != nil { // error for db
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if userdup != nil { // case in duplicated
		return echo.NewHTTPError(http.StatusBadRequest, "Username is already taken")
	}

	//register
	user, err := h.ui.CreateUser(u)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//add new user to open chat , create session and cookie
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
	//get updated user input
	u := new(model.UserUpdate)
	err := validatedBind(c, u)
	if err != nil {
		return err
	}

	//get user session
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}

	//update profile
	newprofile, err := h.ui.EditUser(sess.UserId, u)
	if err != nil { //error for db
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if newprofile == nil { //incorrespond of password
		return echo.NewHTTPError(http.StatusUnauthorized, "Incorrect password")
	}

	return c.JSON(http.StatusOK, newprofile)
}

func (h *Handler) GetMyUser(c echo.Context) error {
	//get session
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}

	//get my user info
	user, err := h.ui.GetUser(sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) SearchUser(c echo.Context) error {
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}
	//get input
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

	//get users who match conditions
	userlist, err := h.ui.GetUserList(limit, offset, name, gender, prefect, sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, userlist)
}
