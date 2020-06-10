package user

import (
	"encoding/json"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// @Tags user
// @Summary 加入成绩共享计划
// @Param token header string true "token"
// @Param data body model.LoginModel true "学号密码"
// @Success 200 "OK"
// @Router /user/licence/ [post]
func JoinPro(c *gin.Context) {
	var l model.LoginModel
	if err := c.ShouldBindJSON(&l); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	userId := c.MustGet("id").(uint32)

	// 检查该用户是否有查看成绩的许可
	if ok, err := model.UserHasLicence(userId); err != nil {
		handler.SendError(c, errno.ErrAddLicence, nil, err.Error())
		return
	} else if ok {
		// 已加入计划
		handler.SendResponse(c, nil, "Has already joined.")
		return
	}

	u, err := service.GetUserById(userId)
	if err != nil {
		handler.SendError(c, errno.ErrGetUserInfo, nil, err.Error())
		return
	}

	u.Licence = true
	if err := u.UpdateLicence(); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, nil)

	/* ------ 成绩服务 ------ */

	gMsg := &service.AsynGradeMsgModel{
		LoginModel: model.LoginModel{
			Sid:      l.Sid,
			Password: l.Password,
		},
		UserId: userId,
		New:    true,
	}
	msg, err := json.Marshal(gMsg)
	if err != nil {
		log.Errorf(err, "marshal asyn-grade-msg error for (userId=%d, sid=%s, psw=%s)", userId, gMsg.Sid, gMsg.Password)
		return
	}

	if err := model.PublishMsg(msg, model.GradeChan); err != nil {
		log.Errorf(err, "asyn-grade-msg publish error for (%s)", string(msg))
		return
	}
	log.Info("publish msg OK")
}
