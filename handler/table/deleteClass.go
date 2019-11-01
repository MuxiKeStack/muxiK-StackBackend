package table

import (
	"strconv"
	"strings"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"

	"github.com/gin-gonic/gin"
)

// 删除课堂
func DeleteClass(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	tableId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, err, nil, err.Error())
	}

	classId := c.DefaultQuery("classId", "")
	if classId == "" {
		handler.SendBadRequest(c,nil, nil, "Table id is expected. ")
	}

	// 获取课表
	table := &model.ClassTableModel{
		Id:      uint32(tableId),
		UserId:  userId,
	}
	if !table.Existing() {
		handler.SendResponse(c, errno.ErrTableExisting, nil)
	}
	if err := table.GetById(); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
	}

	// 移除目标课堂的id
	omitStart := strings.Index(table.Classes, classId)
	if omitStart == -1 {
		handler.SendError(c, errno.ErrClassExisting, nil, "")
	}
	omitEnd := strings.Index(table.Classes[omitStart:], ",")

	var newClasses string
	if omitEnd != -1 {
		newClasses = table.Classes[:omitStart] + table.Classes[omitEnd+1:]
	} else {
		newClasses = table.Classes[:omitStart]
	}

	if err := table.UpdateClasses(newClasses); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
	}

	handler.SendResponse(c, nil, nil)
}
