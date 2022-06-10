package router

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gochatroom/api/v1/common"
	"github.com/josexy/gochatroom/api/v1/room"
	"github.com/josexy/gochatroom/api/v1/user"
	"github.com/josexy/gochatroom/global"
	"github.com/josexy/gochatroom/logx"
	"github.com/josexy/gochatroom/middleware"
	"github.com/josexy/gochatroom/websocket"
)

func NewRouter() *gin.Engine {

	if global.AppConfig.Server.Mode == "release" {
		logx.DebugMode = false
		gin.DefaultWriter = ioutil.Discard
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(middleware.RateLimit(), middleware.Cors(), middleware.Logger())
	middleware.NewJWTAuthMiddleWare(&middleware.AllUserAuthorizator{})

	v1 := r.Group("/api/v1")
	{
		v1.POST("user/register", user.Register)
		v1.POST("user/login", middleware.AuthMiddleware.LoginHandler)

		v1.GET("/online", common.Get)
		v1.GET("/ws", websocket.WsServe)

		auth := v1.Group("/auth")
		auth.GET("/refresh_token", middleware.AuthMiddleware.RefreshHandler)

		auth.Use(middleware.AuthMiddleware.MiddlewareFunc())
		{
			u := auth.Group("/user")
			{
				u.GET("", user.Get)
				u.DELETE("", middleware.AuthMiddleware.LogoutHandler)
			}

			auth.GET("/rooms", room.Rooms)
			auth.GET("/room/:id", room.Messages)
			auth.GET("/room/p/:id/:uid", room.PrivateChatMessages)
			auth.POST("/img_upload", room.UploadFile)
		}
	}
	return r
}
