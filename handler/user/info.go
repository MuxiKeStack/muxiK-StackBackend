package user

import (
	. "github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/gin-gonic/gin"
)

// Update user info by sid
func PostInfo(c *gin.Context) {
	var info model.UserInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}
	sid := c.Param("sid")
	u, err := model.GetUserBySid(sid)
	if err != nil {
		SendBadRequest(c, errno.ErrUserNotFound, nil, err.Error())
		return
	}
	if err := u.UpdateInfo(info); err != nil {
		SendBadRequest(c, errno.ErrUpdateUser, nil, err.Error())
		return
	}
	SendResponse(c, errno.OK, nil)
}

//
func GetInfo(c *gin.Context) {
	sid := c.Param("sid")
	u, err := model.GetUserBySid(sid)
	if err != nil {
		SendBadRequest(c, errno.ErrUserNotFound, nil, err.Error())
		return
	}
	info := u.GetInfo()
	SendResponse(c, errno.OK, info)
}
