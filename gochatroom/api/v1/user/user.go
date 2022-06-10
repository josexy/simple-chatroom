package user

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/josexy/gochatroom/api/v1/common"
	"github.com/josexy/gochatroom/service/user"
)

func Register(c *gin.Context) {
	var service user.RegisterService
	if err := c.ShouldBind(&service); err == nil {
		res := service.Register()
		common.ResponseJson(c, res)
	} else {
		common.ResponseJsonError(c, err)
	}
}

func Login(c *gin.Context) (interface{}, error) {
	var service user.LoginService
	if err := c.ShouldBind(&service); err == nil {
		if u := service.Login(); u != nil {
			return u, nil
		} else {
			return nil, jwt.ErrFailedAuthentication
		}
	} else {
		return nil, jwt.ErrMissingLoginValues
	}
}

func Get(c *gin.Context) {
	var service user.GetService
	common.ResponseJson(c, service.Get(common.GetUserId(c)))
}
