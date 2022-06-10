package websocket

type msgType = int

const (
	HeartBeat     msgType = iota // 发送心跳包（单播）
	Online                       // 进入房间，建立连接（广播）
	Offline                      // 退出房间，断开连接（广播）
	Send                         // 已发送数据给客户端（广播）
	OnlineClients                // 客户端获取用户列表（单播）
	PrivateChat                  // 私聊（单播）
)

type ClientInfo struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	RoomId   int    `json:"room_id"`
	AvatarId int    `json:"avatar_id"`
}

type Message struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	RoomId   int    `json:"room_id,omitempty"`
	AvatarId int    `json:"avatar_id,omitempty"`
	ToUserId int    `json:"to_user_id,omitempty"`
	Content  string `json:"content,omitempty"`
	ImageUrl string `json:"image_url,omitempty"`
}

type WsMessage struct {
	// 消息类型
	Status msgType `json:"status"`
	// 心跳包数据
	Msg string `json:"msg,omitempty"`
	// 正常通信数据
	Data *Message `json:"data,omitempty"`
	// 获取在线用户列表
	List []*ClientInfo `json:"list,omitempty"`
}
