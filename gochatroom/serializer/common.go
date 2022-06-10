package serializer

import (
	"time"

	"github.com/josexy/gochatroom/pkg/codes"
)

type Response struct {
	Code int         `json:"code"`           // 错误码
	Msg  string      `json:"msg"`            // 错误码可读消息
	Data interface{} `json:"data,omitempty"` // 返回此次请求响应的数据
}

type Token struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}

func BuildResponse(code int) Response {
	return Response{
		Code: code,
		Msg:  codes.GetCodeMessage(code),
	}
}

func BuildResponseWithData(code int, data interface{}) Response {
	return Response{
		Code: code,
		Msg:  codes.GetCodeMessage(code),
		Data: data,
	}
}

func BuildError(code int, message string) Response {
	return Response{
		Code: code,
		Msg:  message,
	}
}
