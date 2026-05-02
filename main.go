package main

import (
	"strings"
	"os"
	"log"

	"github.com/111161226/TOKO-ENCOUNT/handler"
	mid "github.com/111161226/TOKO-ENCOUNT/middleware"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// main.go の冒頭
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Aiven接続なら "pgx" または "postgres" ドライバが必要[cite: 1]
	db, err := sqlx.Connect("pgx", dsn) 
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

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
	e.Static("/api/swagger-ui", "./docs/swagger-ui/dist")

	e.POST("/api/login", h.Login)
	e.POST("/api/check", h.CheckUser)
	e.POST("/api/restore", h.RestoreUser)

	api := e.Group("/api", mid.EnsureAuthorized(h))
	{
		api.GET("/ping", h.Ping)

		api.GET("/ws", h.GetWebSocket)
		api.POST("/logout", h.Logout)
		api.GET("/user", h.SearchUser)
		api.POST("/user", h.SignUp)
		api.GET("/user/me", h.GetMyUser)
		api.PATCH("/user/me", h.EditProfile)
		api.DELETE("/delete", h.DeleteUser)

		apiChat := api.Group("/chat")
		{
			apiChat.GET("", h.GetChatList)
			apiChat.POST("", h.CreateChat)

			apiRoomId := apiChat.Group("/:rid", mid.EnsureExistChatAndHaveAccessRight(h))
			{
				apiRoomId.POST("", h.ChatPost)
				apiRoomId.GET("", h.GetMessages)
				apiRoomId.GET("/name", h.GetroomName)
				apiRoomId.POST("/name", h.EditRoomName)
				apiRoomId.GET("/add", h.SearchUserbyUsername)
				apiRoomId.POST("/add", h.AddChatUser)
			}
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // ローカル用のデフォルト
	}
	e.Logger.Fatal(e.Start(":" + port))
}
