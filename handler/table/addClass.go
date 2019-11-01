package table

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

type addClassResponseData struct {
	TableId   uint32           `json:"table_id"`
	ClassInfo *model.ClassInfo `json:"class_info"`
}

// 添加课堂
func AddClass(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	tableId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, err, nil, err.Error())
		return
	}

	classId := c.DefaultQuery("classId", "")
	if classId == "" {
		handler.SendBadRequest(c, errno.ErrClassIdRequired, nil, "")
		return
	}

	table := &model.ClassTableModel{
		Id:      uint32(tableId),
		UserId:  userId,
	}
	// table是否存在
	if !table.Existing() {
		handler.SendResponse(c, errno.ErrTableExisting, nil)
		return
	}
	// 获取课表
	if err := table.GetById(); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// 验证id所属的class是否存在
	if !model.IsClassExisting(classId) {
		handler.SendBadRequest(c, errno.ErrClassExisting, nil, "")
		return
	}

	// 添加新课堂id
	var newClasses string
	if table.Classes == "" {
		newClasses = classId
	} else {
		newClasses = table.Classes + "," + classId
	}
	if err := table.UpdateClasses(newClasses); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// 获取新课堂的信息
	newClassInfo, err := service.GetClassInfoById(classId)
	if err != nil {
		handler.SendError(c, errno.ErrGetClassInfo, nil, err.Error())
		return
	}

	data := addClassResponseData{
		TableId:   uint32(tableId),
		ClassInfo: newClassInfo,
	}

	handler.SendResponse(c, nil, data)
}
