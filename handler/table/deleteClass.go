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
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	classId := c.DefaultQuery("classId", "")
	if classId == "" {
		handler.SendBadRequest(c,errno.ErrGetQuery, nil, "Table id is expected. ")
		return
	}

	table := &model.ClassTableModel{
		Id:      uint32(tableId),
		UserId:  userId,
	}

	// 检测课表是否已经存在
	if !table.Existing() {
		handler.SendResponse(c, errno.ErrTableExisting, nil)
		return
	}

	// 获取课表
	if err := table.GetById(); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// 移除目标课堂的id
	omitStart := strings.Index(table.Classes, classId)
	if omitStart == -1 {
		handler.SendError(c, errno.ErrClassExisting, nil, "")
		return
	}
	omitEnd := strings.Index(table.Classes[omitStart:], ",")

	var newClasses string
	if omitEnd != -1 {
		newClasses = table.Classes[:omitStart] + table.Classes[omitEnd+1:]
	} else {
		newClasses = table.Classes[:omitStart]
	}

	// 更新课表信息
	if err := table.UpdateClasses(newClasses); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, nil)
}
