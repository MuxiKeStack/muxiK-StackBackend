package user

import (
	. "github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

// @Tags user
// @Summary 上传用户信息
// @Description 带着token，上传用户的头像/avatar，名字/username.
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param data body model.UserInfoRequest true "用户信息"
// @Success 200 "OK"
// @Router /user/info/ [post]
func PostInfo(c *gin.Context) {
	var info model.UserInfoRequest
	// BindJSON 如果字段不存在会返回400
	// ShouldBindJSON 不会自动返回400
	if err := c.BindJSON(&info); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}
	if info.Username == "" || info.Avatar == "" {
		SendBadRequest(c, errno.ErrUserInfo, nil, "username or avatar cannot be empty")
		return
	}
	id, _ := c.Get("id")
	if err := service.UpdateInfoById(id.(uint32), &info); err != nil {
		SendError(c, errno.ErrUpdateUser, nil, err.Error())
		return
	}
	SendResponse(c, errno.OK, nil)
}

// @tags user
// @Summary 获取用户信息
// @Description 带着token
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Success 200 {object}  model.UserInfoResponse
// @Router /user/info/ [get]
func GetInfo(c *gin.Context) {
	id, _ := c.Get("id")
	info, err := service.GetUserInfoById(id.(uint32))
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
