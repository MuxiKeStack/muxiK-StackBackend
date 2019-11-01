package table

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

// 新建课表
func AddTable(c *gin.Context) {
	userId := c.MustGet("id").(uint32)
	idStr := c.Param("id")

	newTable := model.ClassTableModel{
		UserId: userId,
	}

	// id为空，新建空白课表
	if idStr == "" {
		if err := newTable.New(); err != nil {
			handler.SendError(c, err, nil, err.Error())
			return
		}

	} else {
		// 创建课表副本

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			handler.SendBadRequest(c, errno.ErrGetQuery, nil, err.Error())
			return
		}

		table := model.ClassTableModel{
			Id:      uint32(id),
			UserId:  userId,
		}

		// 检验父课表存在
		if !table.Existing() {
			handler.SendBadRequest(c, errno.ErrTableExisting, nil, "")
			return
		}

		// 根据id获取父课表
		if err := table.GetById(); err != nil {
			handler.SendError(c, err, nil, err.Error())
			return
		}

		newTable.Name = table.Name
		newTable.Classes = table.Classes

		// 创建副本课表
		if err := newTable.New(); err != nil {
			handler.SendError(c, err, nil, err.Error())
			return
		}
	}

	// 获取返回的课表信息
	tableInfo, err := service.GetTableInfoById(newTable.Id)
	if err != nil {
		handler.SendError(c, errno.ErrGetTableInfo, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, tableInfo)
}
