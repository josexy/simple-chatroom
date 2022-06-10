package model

import "time"

type Message struct {
	ID        uint `gorm:"primaryKey"`
	UserId    int  `gorm:"index"`
	ToUserId  int  `gorm:"index"`
	RoomId    int  `gorm:"index"`
	Content   string
	ImageUrl  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type HistoryMessage struct {
	Message
	AvatarId int
	Username string
	Nickname string
}

func BuildMessage(userId, toUserId, roomId int, content, imageUrl string) *Message {
	return &Message{
		UserId:   userId,
		ToUserId: toUserId,
		RoomId:   roomId,
		Content:  content,
		ImageUrl: imageUrl,
	}
}
