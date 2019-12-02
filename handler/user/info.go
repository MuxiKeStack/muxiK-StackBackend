package user

import (
	. "github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
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
	log.Info("PostInfo function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var info model.UserInfoRequest
	if err := c.ShouldBindJSON(&info); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
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
	log.Info("GetInfo function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
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
