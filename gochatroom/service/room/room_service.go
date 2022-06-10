package room

import (
	"github.com/josexy/gochatroom/global"
	"github.com/josexy/gochatroom/model"
	"github.com/josexy/gochatroom/pkg/codes"
	"github.com/josexy/gochatroom/serializer"
	"gorm.io/gorm"
)

type MessagesService struct{}

type PrivateChatMessagesService struct {
	Uid int `json:"uid" form:"uid"`
}

func getHistoryMessagesBy(roomId int, toUserId int) ([]*model.HistoryMessage, error) {
	/*
		select *from
		(select messages.*,username,nickname,users.avatar_id
			from messages
			inner join users on users.id=messages.user_id
			where messages.room_id=2 and messages.to_user_id=0
			order by messages.id desc
			limit 10
		) as T
		order by T.id asc;
	*/

	var messages []*model.HistoryMessage
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		subQuery := tx.Model(&model.Message{}).
			Select("messages.*, users.username, users.nickname, users.avatar_id").
			Joins("inner join users on users.id = messages.user_id").
			Where("room_id = ? and to_user_id = ?", roomId, toUserId).
			Order("messages.id desc").
			Limit(20)
		return tx.Table("(?) as T", subQuery).Order("T.id asc").Find(&messages).Error
	})
	return messages, err
}

func getPrivateChatHistoryMessagesBy(roomId, toUserId, Uid int) ([]*model.HistoryMessage, error) {
	/*
		select *from
		(select messages.*,username,nickname,users.avatar_id
			from messages
			inner join users on users.id=messages.user_id
			where user_id in (1,2) and to_user_id in (1,2) and room_id = 1
			order by messages.id desc
			limit 10
		) as T
		order by T.id asc;
	*/
	var messages []*model.HistoryMessage
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		subQuery := tx.Model(&model.Message{}).
			Select("messages.*, users.username, users.nickname, users.avatar_id").
			Joins("inner join users on users.id = messages.user_id").
			Where("room_id = ? and user_id in (?, ?) and to_user_id in (?, ?)",
				roomId, Uid, toUserId, toUserId, Uid).
			Order("messages.id desc").
			Limit(20)
		return tx.Table("(?) as T", subQuery).Order("T.id asc").Find(&messages).Error
	})
	return messages, err
}

func (service *MessagesService) List(roomId int) serializer.Response {
	messages, err := getHistoryMessagesBy(roomId, 0)
	if err != nil {
		return serializer.BuildResponse(codes.Error)
	}
	return serializer.BuildResponseWithData(codes.Success, serializer.BuildMessageList(messages))
}

func (service *PrivateChatMessagesService) List(roomId int, toUserId int) serializer.Response {
	messages, err := getPrivateChatHistoryMessagesBy(roomId, toUserId, service.Uid)
	if err != nil {
		return serializer.BuildResponse(codes.Error)
	}
	return serializer.BuildResponseWithData(codes.Success, serializer.BuildMessageList(messages))
}
