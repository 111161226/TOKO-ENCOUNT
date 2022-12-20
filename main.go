package main

import (
	"strings"

	"github.com/cs-sysimpl/SakataKintoki/handler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := sqlx.MustConnect("mysql", "root:password@tcp(mysql:3306)/chatdb?parseTime=true")

	h := handler.NewHandler(db)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			return strings.HasPrefix(c.Request().URL.Path, "/api")
		},
		Root:  "web/dist",
		HTML5: true,
	}))

	e.File("/api/swagger.yaml", "./docs/swagger.yaml")
	e.Static("/api", "./docs/swagger-ui/dist")

	e.POST("/api/login", h.NotImpl)

	api := e.Group("/api")
	{
		api.GET("/ping", h.Ping)
	}

	e.Logger.Fatal(e.Start(":80"))
}
