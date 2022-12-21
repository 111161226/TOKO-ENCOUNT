package handler

import (
	"net/http"

	"github.com/cs-sysimpl/SakataKintoki/db/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ChatPost(c echo.Context) error {
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}
	cid := c.Param("cid")
	did := c.QueryParam("did")
	post := &model.MessageSimple{}
	err = validatedBind(c, post)
	if err != nil {
		return err
	}

	postedMessage, err := h.ci.PostChat(cid, did, post, sess.UserId)
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

func validatedBind(c echo.Context, i interface{}) error {
	err := c.Bind(i) // リクエストボディの取り出し
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}