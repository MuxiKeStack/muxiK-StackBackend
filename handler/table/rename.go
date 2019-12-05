package table

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/gin-gonic/gin"
)

type renameBodyData struct {
	NewName string `json:"new_name"`
}

// 课表重命名
// @Summary 课表重命名
// @Tags table
// @Param token header string true "token"
// @Param id path string true "课表id"
// @Param data body table.renameBodyData true "data"
// @Success 200 "OK"
// @Router /table/{id}/rename/ [put]
func Rename(c *gin.Context) {
	userId := c.MustGet("id").(uint32)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	var data renameBodyData
	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, "Table id is expected. ")
		return
	}

	table := &model.ClassTableModel{
		Id:     uint32(id),
		UserId: userId,
	}

	// 检测课表是否存在
	if !table.Existing() {
		handler.SendBadRequest(c, errno.ErrTableExisting, nil, "")
		return
	}

	// 更新课表
	if err := table.Rename(data.NewName); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, nil)
}
