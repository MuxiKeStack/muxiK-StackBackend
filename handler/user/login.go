package user

import (
	. "github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/gin-gonic/gin"
)

// @Tags auth
// @Summary 学号登录
// @Description 用学号登录，返回token，如果isnew==1，就要post用户信息。
// @Accept  json
// @Produce  json
// @Param data body model.LoginModel true "sid学号，password密码"
// @Success 200  {object} model.AuthResponse
// @Router /login [post]
func Login(c *gin.Context) {
	// Binding the data with the user struct.
	var l model.LoginModel
	if err := c.ShouldBindJSON(&l); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	// FIX grpc 调 data_service
	// Compare the login password with the user password.
	// 业务逻辑异常，使用 SendResponse 发送 200 请求 + 自定义错误码
	if err := util.LoginRequest(l.Sid, l.Password); err != nil {
		SendResponse(c, errno.ErrAuthFailed, nil)
		return
	}

	// judge whether there is this user or not
	IsNewUser := service.HaveUser(l.Sid)
	if IsNewUser == 1 {
		err := model.CreateUser(l.Sid)
		if err != nil {
			SendError(c, errno.ErrCreateUser, nil, err.Error())
		}
	}
	u, err := service.GetUserBySid(l.Sid)
	if err != nil {
		SendError(c, errno.ErrUserNotFound, nil, err.Error())
		return
	}

	// Sign the json web token.
	t, err := token.Sign(c, token.Context{Id: u.Id}, "")
	if err != nil {
		SendError(c, errno.ErrToken, nil, err.Error())
		return
	}

	SendResponse(c, errno.OK, model.AuthResponse{Token: t, IsNew: IsNewUser})
}
