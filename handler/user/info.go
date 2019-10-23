package user

import (
	. "github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// Update user info by token
func PostInfo(c *gin.Context) {
	log.Info("PostInfo function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var info model.UserInfoRequest
	if err := c.ShouldBindJSON(&info); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}
	id, _ := c.Get("id")
	if err := model.UpdateInfoById(id.(uint32), &info); err != nil {
		SendError(c, errno.ErrUpdateUser, nil, err.Error())
		return
	}
	SendResponse(c, errno.OK, nil)
}

//
func GetInfo(c *gin.Context) {
	log.Info("GetInfo function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	id, _ := c.Get("id")
	info, err := model.GetUserInfoById(id.(uint32))
	if err != nil {
		SendError(c, errno.ErrGetUserInfo, nil, err.Error())
		return
	}
	SendResponse(c, errno.OK, model.UserInfoResponse{
		Username: info.Username,
		Avatar:   info.Avatar,
		Sid:      info.Sid,
	})
}
