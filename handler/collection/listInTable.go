package collection

import (
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
// @Success 200 {object} collection.CollectionsInfo
// @Router /collection/table/ [get]
func CollectionsForTable(c *gin.Context) {
	userId := c.MustGet("id").(uint32)

	data, err := service.GetCollectionsList(userId)
	if err != nil {
		handler.SendError(c, errno.ErrGetCollections, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, &CollectionsInfo{
		Sum:        len(*data),
		CourseList: data,
	})
}
