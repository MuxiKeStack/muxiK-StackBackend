package table

import (
	"errors"
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

// 新建课表
func AddTable(c *gin.Context) {
	userId := c.MustGet("id").(uint32)
	idStr := c.Param("id")

	var tableInfo *model.ClassTableInfo
	var err error

	table := model.ClassTableModel{
		UserId: userId,
	}

	if idStr == "" {
		if err := table.New(); err != nil {
			handler.SendError(c, err, nil, err.Error())
		}
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}
	table.Id = uint32(id)
	if !table.Existing() {
		handler.SendError(c, errors.New("Table is not existing "), nil, "Table is not existing ")
	}
	if err := table.GetById(); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}
	tableInfo, err = service.GetClassTableInfoById(table.Id)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, tableInfo)
}
