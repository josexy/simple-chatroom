package user

import (
	"strings"

	"github.com/josexy/gochatroom/global"
	"github.com/josexy/gochatroom/model"
	"github.com/josexy/gochatroom/pkg/codes"
	"github.com/josexy/gochatroom/serializer"
	"gorm.io/gorm"
)

type GetService struct {
}

type LoginService struct {
	UserName string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=3,max=20"`
}

type RegisterService struct {
	Nickname string `json:"nickname" binding:"required,min=3,max=20"`
	UserName string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=3,max=20"`
	AvatarId int    `json:"avatar_id"`
}

func (service *RegisterService) check() (serializer.Response, bool) {
	var count int64
	global.DB.Model(&model.User{}).Where("username = ?", service.UserName).Count(&count)
	if count > 0 {
		return serializer.BuildResponse(codes.ErrorUserAlreadyExisted), false
	}

	count = 0
	global.DB.Model(&model.User{}).Where("nickname = ?", service.Nickname).Count(&count)
	if count > 0 {
		return serializer.BuildResponse(codes.ErrorUserNicknameAlreadyExisted), false
	}

	return serializer.Response{}, true
}

func (service *RegisterService) Register() serializer.Response {
	service.Nickname = strings.TrimSpace(service.Nickname)
	service.UserName = strings.TrimSpace(service.UserName)
	service.Password = strings.TrimSpace(service.Password)

	if resp, ok := service.check(); !ok {
		return resp
	}

	user := model.User{
		Nickname: service.Nickname,
		Username: service.UserName,
		AvatarId: service.AvatarId,
	}

	if user.EncryptPassword(service.Password) != nil {
		return serializer.BuildResponse(codes.Error)
	}

	if err := global.DB.Create(&user).Error; err != nil {
		return serializer.BuildResponse(codes.Error)
	}

	return serializer.BuildResponseWithData(
		codes.Success,
		serializer.BuildUser(user),
	)
}

func (service *LoginService) Login() *serializer.User {
	service.UserName = strings.TrimSpace(service.UserName)
	service.Password = strings.TrimSpace(service.Password)

	var user model.User

	if err := global.DB.Where("username = ?", service.UserName).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return nil
	}

	if !user.CheckPassword(service.Password) {
		return nil
	}
	u := serializer.BuildUser(user)
	return &u
}

func (service *GetService) Get(id int) serializer.Response {
	var user model.User

	if err := global.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return serializer.BuildResponse(codes.ErrorUserNotExist)
	}
	return serializer.BuildResponseWithData(codes.Success, serializer.BuildUser(user))
}
