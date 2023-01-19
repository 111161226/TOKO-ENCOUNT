package main

import (
	"strings"

	"github.com/cs-sysimpl/SakataKintoki/handler"
	mid "github.com/cs-sysimpl/SakataKintoki/middleware"
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

	e.POST("/api/login", h.Login)

	api := e.Group("/api", mid.EnsureAuthorized(h))
	{
		api.GET("/ping", h.Ping)

		api.GET("/ws", h.GetWebSocket)
		api.POST("/logout", h.Logout)
		api.GET("/user", h.SearchUser)
		api.POST("/user", h.SignUp)
		api.GET("/user/me", h.GetMyUser)
		api.PATCH("/user/me", h.EditProfile)

		apiChat := api.Group("/chat")
		{
			apiChat.GET("/", h.GetChatList)
			apiChat.POST("/", h.CreateChat)

			apiRoomId := apiChat.Group("/:rid", mid.EnsureExistChatAndHaveAccessRight(h))
			{
				apiRoomId.POST("/", h.ChatPost)
				apiRoomId.GET("/", h.GetMessages)
			}
		}
	}

	e.Logger.Fatal(e.Start(":80"))
}
