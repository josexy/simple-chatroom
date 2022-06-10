package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/josexy/gochatroom/global"
	"github.com/josexy/gochatroom/logx"
	"github.com/josexy/gochatroom/model"
)

const maxSendDataCount = 256

type wsClient struct {
	Conn    *websocket.Conn
	sendMsg chan *WsMessage
	*ClientInfo
}

func newWsClient(conn *websocket.Conn) *wsClient {
	client := &wsClient{
		Conn:       conn,
		sendMsg:    make(chan *WsMessage, maxSendDataCount),
		ClientInfo: new(ClientInfo),
	}
	return client
}

func (client *wsClient) doWork() {
	go client.readData()
	go client.writeData()
}

func (client *wsClient) readMessage() (msg WsMessage, err error) {
	err = client.Conn.ReadJSON(&msg)
	return
}

func (client *wsClient) writeMessage(msg *WsMessage) error {
	return client.Conn.WriteJSON(msg)
}

func (client *wsClient) writeHeartbeat() error {
	return client.Conn.WriteJSON(&WsMessage{Status: HeartBeat, Msg: "heartbeat ok!"})
}

func (client *wsClient) writeOnlineClients() error {
	var list []*ClientInfo
	ClientsHub.GetOnlineClients(client.RoomId).Range(func(c *wsClient) {
		list = append(list, c.ClientInfo)
	})
	return client.Conn.WriteJSON(&WsMessage{Status: OnlineClients, List: list})
}

func (client *wsClient) readData() {
	defer func() {
		logx.Debug("client left room [%d] and close read connection", client.RoomId)
		_ = recover()
		// 离开房间，断开连接
		ClientsHub.LeftRoom(client.RoomId, client)
	}()

	for {
		msg, err := client.readMessage()
		if err != nil {
			break
		}

		// 定时发送心跳包
		if msg.Status == HeartBeat && msg.Msg == "heartbeat" {
			_ = client.writeHeartbeat()
			continue
		}

		if msg.Data != nil {
			// 用户首次进入房间，开始建立连接
			if msg.Status == Online {
				// 保存用户基本信息
				client.Uid = msg.Data.Uid
				client.Username = msg.Data.Username
				client.Nickname = msg.Data.Nickname
				client.RoomId = msg.Data.RoomId
				client.AvatarId = msg.Data.AvatarId
				ClientsHub.AddRoom(client.RoomId, client)
			}

			switch msg.Status {
			case Send, Online:
				if msg.Status == Send {
					global.DB.Model(model.Message{}).Create(model.BuildMessage(
						msg.Data.Uid,
						0,
						msg.Data.RoomId,
						msg.Data.Content,
						msg.Data.ImageUrl,
					))
				}
				// 广播发送
				ClientsHub.Broadcast(client.RoomId, &msg)
			default:
				if msg.Status == PrivateChat {
					global.DB.Model(model.Message{}).Create(model.BuildMessage(
						msg.Data.Uid,
						msg.Data.ToUserId,
						msg.Data.RoomId,
						msg.Data.Content,
						msg.Data.ImageUrl,
					))
				}
				// 单播
				ClientsHub.SendDirect(client, &msg)
			}
		}
	}
}

func (client *wsClient) writeData() {
	defer func() {
		_ = recover()
	}()

	for {
		select {
		case msg, ok := <-client.sendMsg:
			if !ok {
				return
			}
			switch msg.Status {
			case Send, Online, Offline:
				// 广播
				_ = client.writeMessage(msg)
			case OnlineClients:
				// 单播
				_ = client.writeOnlineClients()
			case PrivateChat:
				// 私聊
				toClient := ClientsHub.FindClient(msg.Data.RoomId, msg.Data.ToUserId)
				if toClient != nil {
					// 将消息转发给发送者和接收者
					_ = toClient.writeMessage(msg)
					_ = client.writeMessage(msg)
				}
			}
		}
	}
}
