package websocket

type Room struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	OnlineClients int    `json:"online_clients"`
}

var ChatRooms = []Room{
	{
		ID:          1,
		Title:       "测试房间1",
		Description: "测试房间1",
	},
	{
		ID:          2,
		Title:       "测试房间2",
		Description: "测试房间2",
	},
	{
		ID:          3,
		Title:       "测试房间3",
		Description: "测试房间3",
	},
}
