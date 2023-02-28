package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"fmt"

	"github.com/111161226/TOKO-ENCOUNT/db/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ChatPost(c echo.Context) error {
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}
	rid := c.Param("rid")
	did := c.QueryParam("did")
	if did == "" {
		did = "0"
	}

	post := &model.MessageSimple{}
	err = validatedBind(c, post)
	if err != nil {
		return err
	}

	//send message
	postedMessage, err := h.ci.PostChat(rid, post, sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//notify new message to talk opponent
	if(rid == "0") {
		err = h.ws.NotifyNewMessage([]string{did}, rid, postedMessage)
	} else {
		users, err := h.ui.GetRoomUsers(rid, sess.UserId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}	
		err = h.ws.NotifyNewMessage(*users, rid, postedMessage)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	//increment talk opponent not read number
	err = h.ci.IncrementNotRead(rid, sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, postedMessage)
}

//add new chat
func (h *Handler) CreateChat(c echo.Context) error {
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}

	did := c.QueryParam("did")
	if did == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "`did` is required")
	}
	//chaeck did user is valid in creating chat
	u, err := h.ui.GetUserByUserId(did)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid `did`")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if u.UserId == sess.UserId {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid `did`")
	}

	roomData, err := h.ci.CreateChat(did, sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//notify new message to talk opponent
	err = h.ws.NotifyNewMessage([]string{did}, roomData.RoomId, &roomData.LatestMessage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, roomData)
}

//rename chat name
func (h *Handler) EditRoomName(c echo.Context) error {
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}
	newName := c.QueryParam("newName")
	if newName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "`did` is required")
	}
	rid := c.Param("rid")
	chatinfo, err := h.ci.UpdateRoomName(rid, newName, sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//notify new message to talk opponent
	if(rid == "0") {
		err = h.ws.NotifyNewMessage([]string{"0"}, rid, &chatinfo.LatestMessage)
	} else {
		users, err := h.ui.GetRoomUsers(rid, sess.UserId)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}	
		err = h.ws.NotifyNewMessage(*users, rid, &chatinfo.LatestMessage)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	
	//increment talk opponent not read number
	err = h.ci.IncrementNotRead(rid, sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, chatinfo)
}

//add new chat
func (h *Handler) AddChatUser(c echo.Context) error {
	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}

	did := c.QueryParam("did")
	if did == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "`did` is required")
	}
	if did == sess.UserId {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid `did`")
	}
	rid := c.Param("rid")

	roomData, err := h.ci.AddPrivateChat(rid, did)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	//notify new message to talk opponent
	users, err := h.ui.GetRoomUsers(rid, sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}	
	err = h.ws.NotifyNewMessage(*users, rid, &roomData.LatestMessage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, roomData)
}

//get current room name
func (h *Handler) GetroomName(c echo.Context) error {
	rid := c.Param("rid")
	room, err := h.ci.GetRoomName(rid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, room)
}


//get current chat message
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

	sess, err := h.PickSession(c)
	if err != nil {
		return err
	}
	//reset not read message
	err = h.ci.ResetNotRead(rid, sess.UserId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, message)
}

//get users who join the designated room
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

//get chat list
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
