package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/josexy/gochatroom/logx"
	"github.com/josexy/gochatroom/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	HeartBeat     int = iota // 发送心跳包
	Online                   // 进入房间，建立连接
	Offline                  // 退出房间，断开连接
	Send                     // 已经发送数据给客户端
	OnlineClients            // 客户端获取用户列表
	PrivateChat              // 私聊
)

type clientInfo struct {
	RemoteIP string `json:"remote_ip"`
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	RoomId   int    `json:"room_id"`
	AvatarId int    `json:"avatar_id"`
}

type message struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	RoomId   int    `json:"room_id,omitempty"`
	AvatarId int    `json:"avatar_id,omitempty"`
	ToUserId int    `json:"to_user_id,omitempty"`
	Content  string `json:"content,omitempty"`
	ImageUrl string `json:"image_url,omitempty"`
}

// 消息格式
type wsMessage struct {
	// 消息类型
	Status int `json:"status"`
	// 心跳包数据
	Msg string `json:"msg,omitempty"`
	// 正常通信数据
	Data *message `json:"data,omitempty"`
	// 获取在线用户列表
	List []*clientInfo `json:"list,omitempty"`
}

func readUserFromDB() []*model.User {
	dsn := "admin:12345@tcp(192.168.1.105:3306)/db_chatroom?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	var users []*model.User
	if err := db.Model(&model.User{}).Find(&users).Error; err != nil {
		logx.ErrorBy(err)
		return nil
	}
	return users
}

func openConnection(id int, avatar_id int, username string) *websocket.Conn {

	u := url.URL{Scheme: "ws", Host: "localhost:10086", Path: "/api/v1/ws"}
	fmt.Printf("connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return nil
	}

	var msg message
	msg.Uid = id
	msg.AvatarId = rand.Intn(5)
	msg.RoomId = 1
	msg.Username = username
	msg.Nickname = username

	// 加入房间
	_ = c.WriteJSON(&wsMessage{
		Status: Online,
		Data:   &msg,
	})

	// 发送心跳包
	heartbeat := time.NewTicker(time.Second * 5)
	go func() {
		for {
			select {
			case <-heartbeat.C:
				_ = c.WriteJSON(&wsMessage{
					Status: HeartBeat,
					Msg:    "heartbeat",
				})
			}
		}
	}()

	// 读取消息
	go func() {
		for {
			var msg wsMessage
			err = c.ReadJSON(&msg)
			if err != nil {
				logx.Error("read error: %v", err)
				return
			}

			if msg.Status == HeartBeat {
				// logx.Error("心跳包, uid: %d, username: %s, %s", msg.Data.Uid, msg.Data.Username, msg.Data.Content)
			} else if msg.Status == Online {
				logx.Info("new client coming, uid: %d, username: %s", msg.Data.Uid, msg.Data.Username)
			} else if msg.Status == Offline {
				logx.Info("client offline, uid: %d, username: %s", msg.Data.Uid, msg.Data.Username)
			} else if msg.Status == OnlineClients {
				data, err := json.MarshalIndent(msg.List, "", "  ")
				if err == nil {
					fmt.Println(string(data))
				}
			} else {
				data, err := json.MarshalIndent(msg.Data, "", "  ")
				if err == nil {
					fmt.Println(string(data))
				}
			}
		}
	}()

	// 写入消息
	for {
		var typ, uid int
		var content string
		fmt.Printf("please input type and content: [0:exit, 1:send, 2:privatechat, 3:getonlineclients] ")
		_, err = fmt.Scanf("%d %s", &typ, &content)
		if err != nil {
			if strings.TrimSpace(content) == "" {
				continue
			}
			logx.Error("scan error: %v", err)
			break
		}
		var status int
		switch typ {
		case 1:
			status = Send
		case 2:
			status = PrivateChat
			fmt.Printf("please input another user id: ")
			fmt.Scanf("%d", &uid)
		case 3:
			status = OnlineClients
		default:
			return nil
		}
		_ = c.WriteJSON(&wsMessage{
			Status: status,
			Data: &message{
				Uid:      msg.Uid,
				Username: msg.Username,
				Nickname: msg.Nickname,
				RoomId:   msg.RoomId,
				AvatarId: 2,
				ToUserId: uid,
				Content:  content,
			},
		})
	}
	return c
}

func main() {

	users := readUserFromDB()
	for _, user := range users {
		logx.Info("(%d, %s, %s, %d)", user.ID, user.Username, user.Nickname, user.AvatarId)
	}

	// var conns []*websocket.Conn
	// for i := 4; i < len(users); i++ {
	// 	conn := openConnection(int(users[i].ID), users[i].AvatarId, users[i].Username)
	// 	conns = append(conns, conn)
	// }

	conn := openConnection(1, 1, "admin")
	defer conn.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT)
	<-interrupt

	// for _, conn := range conns {
	// 	conn.Close()
	// }
}
