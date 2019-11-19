package message

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/gin-gonic/gin"
)

// @Summary 读取消息提醒
// @Tags message
// @Param token header string true "token"
// @Success 200  "OK"
// @Router /message/readall/ [post]
func ReadAll(c *gin.Context) {
	id, _ := c.Get("id")
	err := model.ReadAll(id.(uint32))
	if err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}
	handler.SendResponse(c, errno.OK, nil)
}
