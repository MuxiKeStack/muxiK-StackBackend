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
// @Summary 新建课表
// @Tags table
// @Param token header string true "token"
// @Param id query string false "若是创建副本，则为课表副本id，若是添加新课表，则为0"
// @Success 200 {object} model.ClassTableInfo
// @Router /table/ [post]
func AddTable(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	// Check whether table amount is more than three
	if num := model.GetTableAmount(userId); num >= 3 {
		handler.SendBadRequest(c, errno.ErrNewTable, nil, "Table amount is limited to three.")
		return
	}

	newTable := &model.ClassTableModel{
		UserId:  userId,
		Name:    "新课表",
		Classes: "",
	}

	idStr, ok := c.GetQuery("id")

	// id为空，新建空白课表
	if !ok || idStr == "" || idStr == "0" {
		if err := newTable.New(); err != nil {
			handler.SendError(c, errno.ErrDatabase, nil, err.Error())
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
			Id:     uint32(id),
			UserId: userId,
		}

		// 检验父课表存在
		if !table.Existing() {
			handler.SendBadRequest(c, errno.ErrTableExisting, nil, "")
			return
		}

		// 根据id获取父课表
		if err := table.GetById(); err != nil {
			handler.SendError(c, errno.ErrDatabase, nil, err.Error())
			return
		}

		newTable.Name = table.Name
		newTable.Classes = table.Classes

		// 创建副本课表
		if err := newTable.New(); err != nil {
			handler.SendError(c, errno.ErrDatabase, nil, err.Error())
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
