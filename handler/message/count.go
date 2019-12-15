package message

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/gin-gonic/gin"
)

type CountModel struct {
	Count uint32 `json:"count"`
}

// @Summary 获取消息提醒的个数
// @Tags message
// @Param token header string true "token"
// @Success 200  {object} model.CountModel
// @Router /message/count/ [get]
func Count(c *gin.Context) {
	id, _ := c.Get("id")
	count, err := model.GetCount(id.(uint32))
	if err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}
	handler.SendResponse(c, nil, CountModel{Count: count})
}
