package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/josexy/gochatroom/logx"
)

const (
	maxReadMessageSize = 8 * 1024
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	ClientsHub = NewWsHub()
)

func WsServe(ctx *gin.Context) {
	// 协议升级
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logx.Debug(err.Error())
		return
	}

	conn.SetReadLimit(maxReadMessageSize)

	client := newWsClient(conn)
	client.doWork()
}
