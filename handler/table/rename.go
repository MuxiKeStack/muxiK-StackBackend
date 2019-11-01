package table

import (
	"errors"
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
	idStr := c.Param("id")
	if idStr == "" {
		handler.SendBadRequest(c, errno.ErrBind, nil, "Table id is expected. ")
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		handler.SendError(c, err, nil, err.Error())
	}
	var data renameBodyData
	if err := c.BindJSON(&data); err != nil {
		handler.SendBadRequest(c, errno.ErrBind, nil, "Table id is expected. ")
	}
	table := &model.ClassTableModel{
		Id:      uint32(id),
		UserId:  userId,
	}
	if !table.Existing() {
		handler.SendError(c, errors.New("Table is not existing "), nil, "Table is not existing ")
	}
	if err := table.Rename(data.NewName); err != nil {
		handler.SendError(c, err, nil, err.Error())
	}

	handler.SendResponse(c, nil, nil)
}
