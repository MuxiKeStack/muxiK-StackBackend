package table

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"

	"github.com/gin-gonic/gin"
)

// 删除课表
func DeleteTable(c *gin.Context) {
	userId := c.MustGet("id").(uint32)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}
	table := &model.ClassTableModel{
		Id:      uint32(id),
		UserId:  userId,
	}
	if !table.Existing() {
		handler.SendResponse(c, errno.ErrDelete, nil)
	}
	if err := table.Delete(); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, nil)
}
