package message

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type messageListResponse struct {
	messageList *[]model.MessageSub
}

// @Summary 获取消息提醒
// @Tags message
// @Param token header string true "token"
// @Param page query integer true "页码"
// @Param limit query integer true "每页最大数"
// @Success 200 {object} message.messageListResponse
// @Router /message/ [get]
func Get(c *gin.Context) {
	uid, _ := c.Get("id")
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.ParseUint(limitStr, 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
		return
	}
	messageList, err := service.MessageList(uint32(page), uint32(limit), uid.(uint32))
	if err != nil {
		handler.SendError(c, errno.ErrGetMessage, nil, err.Error())
	}
	handler.SendResponse(c, nil, messageList)
}
