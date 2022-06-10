package serializer

import (
	"github.com/josexy/gochatroom/model"
)

type Message struct {
	ID       uint   `json:"id"`
	Uid      int    `json:"uid"`
	ToUserId int    `json:"to_user_id"`
	AvatarId int    `json:"avatar_id"`
	RoomId   int    `json:"room_id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Content  string `json:"content"`
	ImageUrl string `json:"image_url"`
}

func BuildMessageList(messages []*model.HistoryMessage) []*Message {
	msgList := []*Message{}
	for i := 0; i < len(messages); i++ {
		msgList = append(msgList, &Message{
			ID:       messages[i].ID,
			Uid:      messages[i].UserId,
			ToUserId: messages[i].ToUserId,
			AvatarId: messages[i].AvatarId,
			RoomId:   messages[i].RoomId,
			Content:  messages[i].Content,
			ImageUrl: messages[i].ImageUrl,
			Username: messages[i].Username,
			Nickname: messages[i].Nickname,
		})
	}
	return msgList
}
