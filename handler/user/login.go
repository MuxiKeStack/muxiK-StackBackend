package user

import (
	. "github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/token"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/gin-gonic/gin"
)

// Login generates the authentication token
// if the password was matched with the specified account.
func Login(c *gin.Context) {
	// Binding the data with the user struct.
	var u model.LoginModel
	if err := c.ShouldBindJSON(&u); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	// Compare the login password with the user password.
	// 业务逻辑异常，使用 SendResponse 发送 200 请求 + 自定义错误码
	if err := util.LoginRequest(u.Sid, u.Password); err != nil {
		SendResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	// Sign the json web token.
	t, err := token.Sign(c, token.Context{Sid: u.Sid}, "")
	if err != nil {
		SendError(c, errno.ErrToken, nil, err.Error())
		return
	}

	// judge whether there is this user or not
	IsNewUser, err := model.HaveUser(u.Sid)
	if IsNewUser == 1 {
		err := model.CreateUser(u.Sid)
		if err != nil {
			SendResponse(c, errno.ErrCreateUser, nil)
		}
	}
	SendResponse(c, errno.OK, model.AuthResponse{Token: t, IsNew: IsNewUser})
}
