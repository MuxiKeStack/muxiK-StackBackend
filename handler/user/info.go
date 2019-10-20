package user

import (
	. "github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Update user info by sid
func PostInfo(c *gin.Context) {
	var info model.UserInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := model.UpdateInfoById(uint32(id), &info); err != nil {
		SendBadRequest(c, errno.ErrUpdateUser, nil, err.Error())
	}
	SendResponse(c, errno.OK, nil)
}

//
func GetInfo(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	info, err := model.GetUserInfoById(uint32(id))
	if err != nil {
		SendError(c, errno.ErrGetUserInfo, nil, err.Error())
	}
	SendResponse(c, errno.OK, info)
}
