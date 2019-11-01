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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	table := &model.ClassTableModel{
		Id:      uint32(id),
		UserId:  userId,
	}

	if !table.Existing() {
		handler.SendResponse(c, errno.ErrTableExisting, nil)
		return
	}

	if err := table.Delete(); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, nil)
}
