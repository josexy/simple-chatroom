package serializer

import "github.com/josexy/gochatroom/model"

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	AvatarId int    `json:"avatar_id"`
}

func BuildUser(user model.User) User {
	return User{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		AvatarId: user.AvatarId,
	}
}
