package room

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/josexy/gochatroom/api/v1/common"
	"github.com/josexy/gochatroom/pkg/codes"
	"github.com/josexy/gochatroom/pkg/util"
	"github.com/josexy/gochatroom/serializer"
	"github.com/josexy/gochatroom/service/room"
	"github.com/josexy/gochatroom/websocket"
)

func Rooms(ctx *gin.Context) {
	for i := 0; i < len(websocket.ChatRooms); i++ {
		websocket.ChatRooms[i].OnlineClients = websocket.ClientsHub.GetOnlineClientCount(websocket.ChatRooms[i].ID)
	}
	common.ResponseJson(ctx, serializer.BuildResponseWithData(codes.Success, websocket.ChatRooms))
}

func Messages(ctx *gin.Context) {
	var service room.MessagesService
	if roomId, err := strconv.Atoi(ctx.Param("id")); err != nil {
		common.ResponseJsonError(ctx, err)
	} else {
		common.ResponseJson(ctx, service.List(roomId))
	}
}

func PrivateChatMessages(ctx *gin.Context) {
	roomId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		common.ResponseJsonError(ctx, err)
		return
	}
	toUserId, err := strconv.Atoi(ctx.Param("uid"))
	if err != nil {
		common.ResponseJsonError(ctx, err)
		return
	}
	var service room.PrivateChatMessagesService
	if err = ctx.ShouldBind(&service); err == nil {
		common.ResponseJson(ctx, service.List(roomId, toUserId))
	} else {
		common.ResponseJsonError(ctx, err)
	}
}

const tmpFileDir = "tmpfile"

func UploadFile(ctx *gin.Context) {
	file, err := ctx.FormFile("uploadImage")
	if err != nil {
		ctx.String(http.StatusBadRequest, "upload file failed")
		return
	}
	if _, err := os.Stat(tmpFileDir); os.IsNotExist(err) {
		_ = os.Mkdir(tmpFileDir, 0666)
	}
	tmpFilePath := tmpFileDir + "/" + file.Filename
	if err = ctx.SaveUploadedFile(file, tmpFilePath); err != nil {
		ctx.String(http.StatusBadRequest, "save upload file failed")
		return
	}
	imageUrl := util.UploadFile(tmpFilePath)

	_ = os.Remove(tmpFilePath)

	common.ResponseJson(ctx, serializer.BuildResponseWithData(codes.Success, gin.H{
		"image_url": imageUrl,
	}))
}
