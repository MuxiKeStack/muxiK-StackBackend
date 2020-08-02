package table

import (
	"fmt"
	"strconv"
	"strings"

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
// @Summary 收藏的课堂加入课表
// @Tags table
// @Param token header string true "token"
// @Param id path string true "课表id"
// @Param course_id query string true "课程hash id"
// @Param class_id query string true "课堂教学班编号"
// @Success 200 {object} table.addClassResponseData
// @Router /table/{id}/class/ [post]
func AddClass(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	tableId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	classId := c.DefaultQuery("class_id", "")
	if classId == "" {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, "The class_id is required.")
		return
	}

	courseId := c.DefaultQuery("course_id", "")
	if classId == "" {
		handler.SendBadRequest(c, errno.ErrGetQuery, nil, "The course_id is required.")
		return
	}

	table := &model.ClassTableModel{
		Id:     uint32(tableId),
		UserId: userId,
	}
	// table是否存在
	if !table.Existing() {
		handler.SendResponse(c, errno.ErrTableExisting, nil)
		return
	}
	// 获取课表
	if err := table.Get(); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// 验证id所属的class是否存在
	if !model.IsClassExisting(courseId, classId) {
		handler.SendBadRequest(c, errno.ErrClassExisting, nil, "")
		return
	}

	// 验证该class在table中是否已存在
	if ok := strings.Contains(table.Classes, courseId); ok {
		handler.SendBadRequest(c, errno.ErrCourseHasExisted, nil, "")
		return
	}

	// 用于表中存储的字符串
	courseStr := fmt.Sprintf("%s#%s", courseId, classId)

	// 添加新课堂id
	var newClasses string
	if table.Classes == "" {
		newClasses = courseStr
	} else {
		newClasses = table.Classes + "," + courseStr
	}
	if err := table.UpdateClasses(newClasses); err != nil {
		handler.SendError(c, errno.ErrDatabase, nil, err.Error())
		return
	}

	// 获取新课堂的信息
	newClassInfo, err := service.GetClassInfoForTableById(courseId, classId)
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
