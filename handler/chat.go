package handler

import (
	"net/http"
	"strconv"

	"github.com/cs-sysimpl/SakataKintoki/db/model"
	"github.com/labstack/echo/v4"
	"fmt"
)

func (h *Handler) ChatPost(c echo.Context) error {
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}
	rid := c.Param("rid")
	did := c.QueryParam("did")
	post := &model.MessageSimple{}
	err = validatedBind(c, post)
	if err != nil {
		return err
	}

	postedMessage, err := h.ci.PostChat(rid, did, post, sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, postedMessage)
}

func (h *Handler) CreateChat(c echo.Context) error {
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}
	did := c.QueryParam("did")
	message, err := h.ci.CreateChat(did, sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, message)
}

func (h *Handler) GetMessages(c echo.Context) error {
	rid := c.Param("rid")
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
	message, err := h.ci.GetMessages(rid, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, message)
}

func (h *Handler) PickChatByRoomId(roomId string) (*model.ChatUserList, error) {
	users, err := h.ci.GetChatByRoomId(roomId)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if users == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("no such room `%s`", roomId))
	}
	return users, nil
}

func (h *Handler) GetChatList(c echo.Context) error {
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}

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

	res, err := h.ci.GetChatList(sess.UserId, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func validatedBind(c echo.Context, i interface{}) error {
	err := c.Bind(i) // リクエストボディの取り出し
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
