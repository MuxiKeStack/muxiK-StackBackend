package collection

import (
	"strconv"

	"github.com/MuxiKeStack/muxiK-StackBackend/handler"
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/pkg/errno"
	"github.com/MuxiKeStack/muxiK-StackBackend/service"

	"github.com/gin-gonic/gin"
)

type CollectionsInfo struct {
	Sum        int                              `json:"sum"`
	CourseList *[]model.CourseInfoInCollections `json:"course_list"`
}

// @Summary 课表界面获取课程清单
// @Tags collection
// @Param token header string true "token"
// @Param id path string true "课表id"
// @Success 200 {object} collection.CollectionsInfo
// @Router /collection/table/{id}/ [get]
func CollectionsForTable(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	tableId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrGetParam, nil, err.Error())
		return
	}

	// Check whether the tableId is valid or whether the table exists
	if !model.TableIsExisting(uint32(tableId), userId) {
		handler.SendResponse(c, errno.ErrTableExisting, nil)
		return
	}

	data, err := service.GetCollectionsList(userId, uint32(tableId))
	if err != nil {
		handler.SendError(c, errno.ErrGetCollections, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, &CollectionsInfo{
		Sum:        len(*data),
		CourseList: data,
	})
}
