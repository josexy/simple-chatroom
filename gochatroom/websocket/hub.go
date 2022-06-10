package websocket

import (
	"sync"
)

type WsHub struct {
	mu      sync.RWMutex
	clients map[*wsClient]bool
	rooms   map[int]*ClientList
}

func NewWsHub() *WsHub {
	return &WsHub{
		clients: map[*wsClient]bool{},
		rooms:   map[int]*ClientList{},
	}
}

func (hub *WsHub) GetOnlineClientCount(roomId int) int {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	cs := hub.rooms[roomId]
	if cs == nil {
		return 0
	}
	return cs.Size()
}

func (hub *WsHub) GetAllOnlineClients() (count int) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	for _, cs := range hub.rooms {
		count += cs.Size()
	}
	return
}

func (hub *WsHub) GetOnlineClients(roomId int) *ClientList {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	return hub.rooms[roomId]
}

func (hub *WsHub) BroadcastAllRooms(message *WsMessage) {
	if message == nil {
		return
	}
	hub.mu.Lock()
	defer hub.mu.Unlock()
	for roomId := range hub.rooms {
		hub.broadcastToRoom(roomId, message, nil)
	}
}

func (hub *WsHub) broadcastToRoom(roomId int, message *WsMessage, exclude *wsClient) {
	if message == nil {
		return
	}
	if clients, ok := hub.rooms[roomId]; ok {
		clients.Range(func(client *wsClient) {
			if client != exclude {
				select {
				case client.sendMsg <- message:
				default:
				}
			}
		})
	}
}

func (hub *WsHub) Broadcast(roomId int, message *WsMessage) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	hub.broadcastToRoom(roomId, message, nil)
}

func (hub *WsHub) BroadcastExclude(roomId int, message *WsMessage, exclude *wsClient) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	hub.broadcastToRoom(roomId, message, exclude)
}

func (hub *WsHub) SendDirect(client *wsClient, message *WsMessage) {
	if message == nil {
		return
	}
	hub.mu.Lock()
	defer hub.mu.Unlock()
	client.sendMsg <- message
}

func (hub *WsHub) FindClient(roomId, uid int) *wsClient {
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	for _, client := range *hub.rooms[roomId] {
		if client.Uid == uid {
			return client
		}
	}
	return nil
}

func (hub *WsHub) AddRoom(roomId int, client *wsClient) bool {
	if client == nil {
		return false
	}

	hub.mu.Lock()
	defer hub.mu.Unlock()

	if hub.rooms[roomId] == nil {
		hub.rooms[roomId] = newClientList()
	}

	if oldClient, ok := hub.clientExistRoom(roomId, client); ok {
		hub.leftRoom(roomId, oldClient)
	}

	hub.clients[client] = true
	hub.rooms[roomId].PushBack(client)
	return true
}

func (hub *WsHub) clientExistRoom(roomId int, client *wsClient) (*wsClient, bool) {
	if cs, ok := hub.rooms[roomId]; ok {
		for _, c := range *cs {
			if c == client || c.Uid == client.Uid || c.Username == client.Username {
				return c, true
			}
		}
	}
	return nil, false
}

func (hub *WsHub) LeftRoom(roomId int, client *wsClient) bool {
	if client == nil {
		return false
	}

	hub.mu.Lock()
	defer hub.mu.Unlock()
	return hub.leftRoom(roomId, client)
}

func (hub *WsHub) leftRoom(roomId int, client *wsClient) bool {
	if cs, ok := hub.rooms[roomId]; !ok {
		return false
	} else if !cs.Contain(client) {
		return false
	}
	delete(hub.clients, client)
	hub.rooms[roomId].Remove(client)
	close(client.sendMsg)
	_ = client.Conn.Close()

	hub.broadcastToRoom(roomId, &WsMessage{
		Status: Offline,
		Data: &Message{
			Username: client.Username, // 下线的用户名
			Nickname: client.Nickname, // 下线的用户昵称
			Uid:      client.Uid,      // 下线的用户ID
		},
	}, nil)

	return true
}
