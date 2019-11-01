package table

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/gin-gonic/gin"
	"strconv"
)

type renameBodyData struct {
	NewName string `json:"new_name"`
}

// 课表重命名
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
		Id:      uint32(id),
		UserId:  userId,
	}

	// 检测课表是否存在
	if !table.Existing() {
		handler.SendBadRequest(c, errno.ErrTableExisting, nil, "")
		return
	}

	// 更新课表
	if err := table.Rename(data.NewName); err != nil {
		handler.SendError(c, err, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, nil)
}
