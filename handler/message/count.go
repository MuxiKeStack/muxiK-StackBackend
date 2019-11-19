package message

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/gin-gonic/gin"
)

// @Summary 获取消息提醒的个数
// @Tags message
// @Param token header string true "token"
// @Success 200  {object} model.CountModel
// @Router /message/count/ [get]
func Count(c *gin.Context) {
	id, _ := c.Get("id")
	count := model.GetCount(id.(uint32))
	handler.SendResponse(c, nil, model.CountModel{Count: count})
}
