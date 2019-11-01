package table

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

type GetTablesResponse struct {
	TableNum  int                     `json:"table_num"`
	TableList *[]model.ClassTableInfo `json:"table_list"`
}

// 获取课表
func Get(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	tables, err := model.GetTablesByUserId(userId)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	var tableInfoList []model.ClassTableInfo
	for _, table := range *tables {
		tableInfo, err := service.GetTableInfoByTableModel(&table)
		if err != nil {
			handler.SendError(c, err, nil, err.Error())
		}
		tableInfoList = append(tableInfoList, *tableInfo)
	}

	result := &GetTablesResponse{
		TableNum:  len(tableInfoList),
		TableList: &tableInfoList,
	}

	handler.SendResponse(c, nil, result)
}
