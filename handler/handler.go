package handler

import (
	"net/http"

	"github.com/cs-sysimpl/SakataKintoki/db/repository"
	"github.com/cs-sysimpl/SakataKintoki/infra"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	ui repository.UserRepository
	ci repository.ChatRepository
	si repository.SessionRepository
}

func NewHandler(db *sqlx.DB) *Handler {
	ui := infra.NewUserInfra(db)
	ci := infra.NewChatInfra(db)
	si := infra.NewSessionInfra(db)
	return &Handler{
		ui: ui,
		ci: ci,
		si: si,
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
