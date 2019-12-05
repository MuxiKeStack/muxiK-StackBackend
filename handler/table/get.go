package table

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

type getTablesResponse struct {
	TableNum  int                     `json:"table_num"`
	TableList *[]model.ClassTableInfo `json:"table_list"`
}

// 获取课表
// @Summary 获取课表
// @Tags table
// @Param token header string true "token"
// @Success 200 {object} table.getTablesResponse
// @Router /table/ [get]
func Get(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	// 获取该用户所属的所有课表
	tables, err := model.GetTablesByUserId(userId)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	// If there is no table, create a default table for him.
	if len(*tables) == 0 {
		originalTable := &model.ClassTableModel{
			UserId: userId,
			Name:   "默认课表",
		}

		if err := originalTable.New(); err != nil {
			handler.SendError(c, errno.ErrDatabase, nil, err.Error())
			return
		}

		*tables = append(*tables, *originalTable)
	}

	// 将课表解析为要返回的课表详情
	var tableInfoList []model.ClassTableInfo
	for _, table := range *tables {
		tableInfo, err := service.GetTableInfoByTableModel(&table)
		if err != nil {
			handler.SendError(c, err, nil, err.Error())
			return
		}

		tableInfoList = append(tableInfoList, *tableInfo)
	}

	data := &getTablesResponse{
		TableNum:  len(tableInfoList),
		TableList: &tableInfoList,
	}

	handler.SendResponse(c, nil, data)
}
