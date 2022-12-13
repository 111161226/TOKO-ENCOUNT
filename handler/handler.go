package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) NotImpl(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (h *Handler) Ping(c echo.Context) error {
	type res struct {
		msg string
	}
	return c.JSON(http.StatusOK, &res{msg: "pong"})
}
