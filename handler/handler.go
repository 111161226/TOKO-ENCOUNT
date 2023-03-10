package handler

import (
	"net/http"

	"github.com/111161226/TOKO-ENCOUNT/db/repository"
	"github.com/111161226/TOKO-ENCOUNT/infra"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	ui repository.UserRepository
	ci repository.ChatRepository
	si repository.SessionRepository
	ws *webSocketPublisher
}

//define hundler
func NewHandler(db *sqlx.DB) *Handler {
	ui := infra.NewUserInfra(db)
	ci := infra.NewChatInfra(db)
	si := infra.NewSessionInfra(db)
	ws := newWebSocketPublisher()
	return &Handler{
		ui: ui,
		ci: ci,
		si: si,
		ws: ws,
	}
}

func (h *Handler) NotImpl(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h *Handler) Ping(c echo.Context) error {
	type res struct {
		Msg string `json:"msg"`
	}
	return c.JSON(http.StatusOK, &res{Msg: "pong"})
}
